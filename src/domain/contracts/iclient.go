package contracts

import (
	"context"
	"net/http"
)

type IRequest interface {
	GetPath() string
	GetBody() []byte
}

type IHTTPClient interface {
	Get(ctx context.Context, request IRequest) (*http.Response, error)
	Post(ctx context.Context, request IRequest) (*http.Response, error)
	Delete(ctx context.Context, request IRequest) (*http.Response, error)
	SetToken(token string)
}
