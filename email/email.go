package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wellsync/acs-go/internal/client"
)

type EmailClient struct {
	client   *client.Client
	endpoint string
}

func NewEmailClient(endpoint string, accessToken string, timeout *time.Duration) *EmailClient {
	return &EmailClient{
		client:   client.NewClient(accessToken, timeout),
		endpoint: endpoint,
	}
}

func (ec *EmailClient) Send(ctx context.Context, msg Message) (*SendResult, error) {
	url := fmt.Sprintf("%s/emails:send?api-version=2023-03-31", ec.endpoint)

	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	err := e.Encode(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := ec.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send message to acs: %w", err)
	}

	d := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	var response SendResult
	err = d.Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal send result: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return &response, fmt.Errorf("acs rejected request: HTTP %v, %+v", resp.Status, response.Error)
	}

	return &response, nil
}
