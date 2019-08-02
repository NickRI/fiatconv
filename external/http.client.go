package external

import (
	"context"
	"net/http"

	"golang.org/x/xerrors"
)

type CtxHttpClient struct {
	*http.Client
}

func NewCtxHttpClient() *CtxHttpClient {
	return &CtxHttpClient{&http.Client{
		Transport: http.DefaultTransport,
	}}
}

func (b *CtxHttpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := b.Client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return resp, xerrors.Errorf("error from context: %w", ctx.Err())
		default:
		}
		return resp, xerrors.Errorf("error during making request: %w", err)
	}
	return resp, nil
}
