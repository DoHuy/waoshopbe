package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CreatePayPalOrder(paypalBaseURL, clientID, secret, mode string, amount float64) (string, string, error) {

	token, err := getPayPalAccessToken(paypalBaseURL, clientID, secret)
	if err != nil {
		return "", "", fmt.Errorf("error getting PayPal access token: %v", err)
	}

	payload := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]string{
					"currency_code": "GBP",
					"value":         fmt.Sprintf("%.2f", amount),
				},
			},
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", paypalBaseURL+"/v2/checkout/orders", bytes.NewBuffer(payloadBytes))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("error creating PayPal order, HTTP status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	orderID := result["id"].(string)

	approveLink := ""
	links := result["links"].([]interface{})
	for _, l := range links {
		link := l.(map[string]interface{})
		if link["rel"] == "approve" {
			approveLink = link["href"].(string)
			break
		}
	}

	return orderID, approveLink, nil
}

func getPayPalAccessToken(baseURL, clientID, secret string) (string, error) {
	req, _ := http.NewRequest("POST", baseURL+"/v1/oauth2/token", bytes.NewBufferString("grant_type=client_credentials"))

	// Basic Auth
	auth := clientID + ":" + secret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+basicAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cannot get token, HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return result["access_token"].(string), nil
}

func CapturePayPalOrder(paypalBaseURL, clientID, secret, orderID string) (string, map[string]interface{}, error) {

	token, err := getPayPalAccessToken(paypalBaseURL, clientID, secret)
	if err != nil {
		return "", nil, fmt.Errorf("error getting PayPal access token: %v", err)
	}

	url := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", paypalBaseURL, orderID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return "", nil, fmt.Errorf("error creating capture request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("error connecting to PayPal: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", result, fmt.Errorf("paypal rejected the transaction, HTTP Code: %d", resp.StatusCode)
	}

	status := "UNKNOWN"
	if val, ok := result["status"].(string); ok {
		status = val
	}

	return status, result, nil
}

type PayPalWebhookEvent struct {
	Id        string `json:"id"`
	EventType string `json:"event_type"`
	Resource  struct {
		Id       string `json:"id"`
		CustomId string `json:"custom_id"`
	} `json:"resource"`
}

func VerifyPayPalSignature(paypalBaseURL string, headers map[string]string, rawBody []byte, clientID, secret, webhookID string) error {
	accessToken, err := getPayPalAccessToken(paypalBaseURL, clientID, secret)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	verifyReqBody := map[string]interface{}{
		"auth_algo":         headers["paypal-auth-algo"],
		"cert_url":          headers["paypal-cert-url"],
		"transmission_id":   headers["paypal-transmission-id"],
		"transmission_sig":  headers["paypal-transmission-sig"],
		"transmission_time": headers["paypal-transmission-time"],
		"webhook_id":        webhookID,
		"webhook_event":     json.RawMessage(rawBody),
	}

	jsonReq, _ := json.Marshal(verifyReqBody)

	req, err := http.NewRequest("POST", paypalBaseURL+"/v1/notifications/verify-webhook-signature", bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		VerificationStatus string `json:"verification_status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.VerificationStatus != "SUCCESS" {
		return fmt.Errorf("signature verification failed: %s", result.VerificationStatus)
	}

	return nil
}
