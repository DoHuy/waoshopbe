package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data,omitempty"`
}

func Success[T any](w http.ResponseWriter, data T) {
	body := Body[T]{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	httpx.OkJson(w, body)
}

func Error(w http.ResponseWriter, code int, msg string) {
	body := Body[any]{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	// Chúng ta vẫn trả về HTTP 200 hoặc mã lỗi tương ứng tùy chiến lược của bạn
	// Ở đây tôi để OkJson để Frontend luôn nhận được body JSON để xử lý
	httpx.WriteJson(w, code, body)
}
