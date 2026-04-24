package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

func RegisterGatewayMiddlewares(gw *rest.Server, rds *redis.Redis) {
	rateLimitMiddleware := NewOrderRateLimitMiddleware(rds)

	gw.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next(w, r)
				return
			}

			if r.URL.Path == "/api/v1/order/create" || r.URL.Path == "/api/v1/order/capture" {
				rateLimitMiddleware.Handle(next).ServeHTTP(w, r)
				return
			}
			next(w, r)
		}
	})

	gw.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next(w, r)
				return
			}

			if r.URL.Path == "/api/v1/webhooks/paypal" || r.URL.Path == "/ws/chat" {
				next(w, r)
				return
			}

			BuildCommonResponse(next)(w, r)
		}
	})
}
