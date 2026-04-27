package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Hostname   string
	Username   string
	Password   string
	HTTPClient *http.Client
}

func NewClient(hostname, username, password string) (*Client, error) {
	if hostname == "" {
		return nil, fmt.Errorf("hostname must not be empty")
	}
	if username == "" {
		return nil, fmt.Errorf("username must not be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password must not be empty")
	}

	return &Client{
		Hostname:   hostname,
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{},
	}, nil
}

type graphQLRequest struct {
	Query string `json:"query"`
}

type graphQLError struct {
	Message string `json:"message"`
}

type graphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []graphQLError  `json:"errors"`
}

func (c *Client) DoGraphQL(ctx context.Context, query string, out interface{}) error {
	body, err := json.Marshal(graphQLRequest{Query: query})
	if err != nil {
		return fmt.Errorf("marshal graphql request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Hostname, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create graphql request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute graphql request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read graphql response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("graphql request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	var envelope graphQLResponse
	if err := json.Unmarshal(responseBody, &envelope); err != nil {
		return fmt.Errorf("decode graphql response: %w", err)
	}

	if len(envelope.Errors) > 0 {
		return fmt.Errorf("graphql error: %s", envelope.Errors[0].Message)
	}

	if out == nil || len(envelope.Data) == 0 {
		return nil
	}

	if err := json.Unmarshal(envelope.Data, out); err != nil {
		return fmt.Errorf("decode graphql data: %w", err)
	}

	return nil
}
