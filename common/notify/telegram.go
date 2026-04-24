package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SendTelegramMessage gửi tin nhắn qua bot Telegram
func SendTelegramMessage(botToken, chatID, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := map[string]string{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "HTML", // Hỗ trợ in đậm <b>, in nghiêng <i>
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("lỗi gọi API Telegram: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram từ chối tin nhắn, status: %d", resp.StatusCode)
	}
	return nil
}
