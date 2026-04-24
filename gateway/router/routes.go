package router

import (
	"io"
	"net/http"
	"strings"

	"dropshipbe/dropshipbe"
	"dropshipbe/dropshipbeclient"
	"dropshipbe/gateway/config"
	ws "dropshipbe/internal/websocket"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

func RegisterRoutes(gw *rest.Server, svc dropshipbeclient.Dropshipbe, rds *redis.Redis, c config.GwConfig) {
	gw.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/ws/chat",
				Handler: ws.ServeChatbotWS(svc, rds, c.Jwt),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/webhooks/paypal",
				Handler: handlePayPalWebhook(svc),
			},
		},
	)
}

func handlePayPalWebhook(svc dropshipbeclient.Dropshipbe) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Cannot read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		headers := make(map[string]string)
		for k, v := range r.Header {
			if strings.HasPrefix(strings.ToLower(k), "paypal-") {
				headers[k] = v[0]
			}
		}

		_, err = svc.HandlePaypalWebhook(r.Context(), &dropshipbe.PayPalWebhookRequest{
			RawBody: bodyBytes,
			Headers: headers,
		})

		if err != nil {
			logx.Errorf("PayPal Webhook RPC Error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
