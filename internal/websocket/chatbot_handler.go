package websocket

import (
	"context"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"

	"dropshipbe/common/utils"
	"dropshipbe/dropshipbe"
	"dropshipbe/dropshipbeclient"
	"dropshipbe/gateway/config"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	writeWait      = 10 * time.Second
	maxMessageSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:3000",
			"https://yourdomain.com",
			"null",
		}

		if slices.Contains(allowedOrigins, origin) {
			return true
		}

		if origin == "" {
			return true
		}

		logx.Errorf("Invalid Origin: %s", origin)
		return false
	},
}

var (
	ipConnections sync.Map
	maxConnPerIP  = 5 // Maximum 5 connection for each IP
)

func tryConnectIP(ip string) bool {
	count, _ := ipConnections.LoadOrStore(ip, 0)
	if count.(int) >= maxConnPerIP {
		return false
	}
	ipConnections.Store(ip, count.(int)+1)
	return true
}

func disconnectIP(ip string) {
	count, ok := ipConnections.Load(ip)
	if ok && count.(int) > 0 {
		ipConnections.Store(ip, count.(int)-1)
	}
}

func writeWSMessage(conn *websocket.Conn, messageType int, payload []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.WriteMessage(messageType, payload)
}

func ServeChatbotWS(rpcClient dropshipbeclient.Dropshipbe, rds *redis.Redis, jwt config.JwtConfig) http.HandlerFunc {
	// Init rate limiter for messages: 1 message per 3 seconds per guest_id
	msgLimiter := limit.NewPeriodLimit(1, 3, rds, "ws:msg_limit:")

	return func(w http.ResponseWriter, r *http.Request) {

		tokenStr := r.URL.Query().Get("token")
		guestID, err := utils.ValidateGuestToken(tokenStr, jwt)
		if err != nil {
			logx.Errorf("Unauthorized WebSocket attempt: %v", err)
			http.Error(w, "Unauthorized: Invalid or missing token", http.StatusUnauthorized)
			return
		}

		// Limit number of connections from the same IP to prevent abuse
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			if !tryConnectIP(ip) {
				logx.Errorf("IP %s reached max WS connections", ip)
				http.Error(w, "Too many connections from this IP", http.StatusTooManyRequests)
				return
			}
			defer disconnectIP(ip) // rollback connection count when function exits
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Error("Error upgrading WebSocket:", err)
			return
		}
		defer conn.Close()

		conn.SetReadLimit(maxMessageSize)
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		go func() {
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()
			for range ticker.C {
				if err := writeWSMessage(conn, websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}()

		welcomeMsg := "Hello " + guestID[:10] + ", I'm your assistant! How can I help you today?"
		writeWSMessage(conn, websocket.TextMessage, []byte(welcomeMsg))

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logx.Errorf("WS Disconnected unexpectedly: %v", err)
				}
				break
			}

			takeRes, err := msgLimiter.Take(guestID)
			if err != nil {
				logx.Error("Redis Limiter error:", err)
				continue
			}
			if takeRes == limit.OverQuota {
				logx.Infof("User %s is spamming. Blocked a message.", guestID)
				writeWSMessage(conn, messageType, []byte("⚠️ Please slow down. You are sending messages too quickly."))
				continue
			}

			cleanMessage := utils.SanitizeMessage(string(p), 500)
			if cleanMessage == "" {
				continue
			}

			logx.Infof("Received clean message from %s: %s", guestID, cleanMessage)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			resp, err := rpcClient.ChatbotMessage(ctx, &dropshipbe.ChatbotRequest{
				Message: cleanMessage,
			})
			cancel()

			if err != nil {
				logx.Error("Error calling gRPC Chatbot:", err)
				writeWSMessage(conn, messageType, []byte("Sorry, there was an error processing your message. Please try again later."))
				continue
			}

			err = writeWSMessage(conn, messageType, []byte(resp.Reply))
			if err != nil {
				logx.Error("Error sending WS message:", err)
				break
			}
		}
	}
}
