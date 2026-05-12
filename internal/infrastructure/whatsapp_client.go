package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WhatsAppClient struct {
	endpoint string
	client  *http.Client
}

type WhatsAppRequest struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type WhatsAppResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewWhatsAppClient(endpoint string) *WhatsAppClient {
	return &WhatsAppClient{
		endpoint: endpoint,
	client: &http.Client{
		Timeout: 30 * time.Second,
	},
	}
}

func (c *WhatsAppClient) SendWhatsApp(ctx context.Context, phone, message string) error {
	if c.endpoint == "" {
		return fmt.Errorf("whatsapp endpoint not configured")
	}

	reqBody := WhatsAppRequest{
		Phone:   phone,
		Message: message,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal whatsapp request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create whatsapp request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send whatsapp request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("whatsapp API returned status %d", resp.StatusCode)
	}

	var waResp WhatsAppResponse
	if err := json.NewDecoder(resp.Body).Decode(&waResp); err != nil {
		return fmt.Errorf("failed to decode whatsapp response: %w", err)
	}

	if !waResp.Success {
		return fmt.Errorf("whatsapp sending failed: %s", waResp.Message)
	}

	return nil
}