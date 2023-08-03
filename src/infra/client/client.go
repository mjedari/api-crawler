package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mjedari/vgang-project/domain/contracts"
	"net/http"
)

type Client struct {
	BaseURL string
	Token   string
	Headers Headers
}
type Headers []map[string]string

func NewClient(baseURL string) *Client {
	headers := Headers{map[string]string{"Content-Type": "application/json"}}
	return &Client{BaseURL: baseURL, Headers: headers}
}

func (c *Client) SetToken(token string) {
	c.Token = token
}

func (c *Client) Get(ctx context.Context, request contracts.IRequest) (*http.Response, error) {
	fmt.Printf("get request to: %v\n", c.BaseURL)
	url := fmt.Sprintf("%v%v", c.BaseURL, request.GetPath())
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

func (c *Client) Post(ctx context.Context, request contracts.IRequest) (*http.Response, error) {
	url := fmt.Sprintf("%v%v", c.BaseURL, request.GetPath())
	req, err := http.NewRequest("POST", url, bytes.NewReader(request.GetBody()))
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

func (c *Client) Delete(ctx context.Context, request contracts.IRequest) (*http.Response, error) {
	return nil, nil
}
