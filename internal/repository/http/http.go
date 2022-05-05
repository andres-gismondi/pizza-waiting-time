package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{client: &http.Client{}}
}

func (c Client) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	encodedBody, err := encode(body)
	req, err := http.NewRequest("POST", path, encodedBody)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func encode(value interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(value)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
