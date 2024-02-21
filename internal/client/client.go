package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"
)

var defaultTimeout = 10 * time.Second

type Client struct {
	accessToken string
	client      *http.Client
}

func NewClient(accessToken string, timeout *time.Duration) *Client {
	if timeout == nil {
		timeout = &defaultTimeout
	}

	return &Client{
		accessToken: accessToken,
		client: &http.Client{
			Timeout: *timeout,
		},
	}
}

func (client *Client) Do(req *http.Request) (*http.Response, error) {
	err := signRequest(req, client.accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed sending request to acs: %w", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("acs rejected request: HTTP %v", resp.Status)
	}

	return resp, nil
}

func signRequest(r *http.Request, token string) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(body)
	r.Body = io.NopCloser(buf)

	contentHash, err := hashBody(body)
	if err != nil {
		return err
	}

	ts := time.Now().UTC().String()
	stringToSign := fmt.Sprintf("%s\n%s\n%s;%s;%s",
		r.Method,
		r.URL.Path+"?"+r.URL.RawQuery,
		ts,
		r.URL.Host,
		contentHash)
	decodedSecret, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return err
	}

	hmac := hmac.New(sha256.New, decodedSecret)
	_, err = hmac.Write([]byte(stringToSign))
	if err != nil {
		return err
	}
	hash := hmac.Sum(nil)
	hashStr := base64.StdEncoding.EncodeToString(hash)
	r.Header.Set("x-ms-date", ts)
	r.Header.Set("x-ms-content-sha256", contentHash)
	r.Header.Set("Authorization",
		fmt.Sprintf("HMAC-SHA256 SignedHeaders=x-ms-date;host;x-ms-content-sha256&Signature=%s", hashStr))

	return nil
}

func hashBody(body []byte) (string, error) {
	sha := sha256.New()
	_, err := sha.Write(body)
	if err != nil {
		return "", err
	}
	sum := sha.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum), nil
}
