package middleware

import (
	"bytes"
	"dropshipbe/common/response"
	"encoding/json"
	"net/http"
)

type customResponseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *customResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func BuildCommonResponse(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/ws/chat" {
			next(w, r)
			return
		}

		cw := &customResponseWriter{
			ResponseWriter: w,
			body:           bytes.NewBuffer(nil),
			statusCode:     http.StatusOK,
		}

		next(cw, r)

		grpcData := cw.body.Bytes()

		if cw.statusCode >= 200 && cw.statusCode < 300 {

			var data json.RawMessage
			if len(grpcData) > 0 {
				data = json.RawMessage(grpcData)
			}

			response.Success(w, data)
		} else {
			msg := string(grpcData)
			if msg == "" {
				msg = "Error occurred with status code: " + http.StatusText(cw.statusCode)
			}

			response.Error(w, cw.statusCode, msg)
		}
	}
}
