package middleware

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type OrderRateLimitMiddleware struct {
	limiter *limit.PeriodLimit
}

func NewOrderRateLimitMiddleware(rds *redis.Redis) *OrderRateLimitMiddleware {

	// - seconds (60)
	// - quota (5)
	limiter := limit.NewPeriodLimit(60, 5, rds, "rate_limit:create_order:")

	return &OrderRateLimitMiddleware{
		limiter: limiter,
	}
}

func (m *OrderRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = strings.Split(r.RemoteAddr, ":")[0]
		}

		code, err := m.limiter.Take(ip)
		if err != nil {
			logx.Errorf("Error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if code == limit.OverQuota {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "Too Many Requests. Please wait a minute before trying again."}`))
			return
		}

		next(w, r)
	}
}
