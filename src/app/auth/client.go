package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
	Token   string
	// todo: Headers
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

type GetRequest struct {
	Path string
}

func (c *Client) Get(ctx context.Context, request GetRequest) (*http.Response, error) {
	fmt.Printf("get request to: %v\n", c.BaseURL)
	url := fmt.Sprintf("%v%v", c.BaseURL, request.Path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", c.Token))

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type PostRequest struct {
	Path  string
	Body  []byte
	Token string
	// todo: header
}

func (c *Client) Post(ctx context.Context, request PostRequest) (*http.Response, error) {
	url := fmt.Sprintf("%v%v", c.BaseURL, request.Path)
	req, err := http.NewRequest("POST", url, bytes.NewReader(request.Body))
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
