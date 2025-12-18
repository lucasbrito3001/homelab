package ports

import (
	"context"
)

type RedirectUrlUseCase interface {
	Execute(ctx context.Context, req *RedirectUrlRequest) (*RedirectUrlResponse, error)
}

type RedirectUrlRequest struct {
	Code string `uri:"code"`
}

type RedirectUrlResponse struct {
	OriginalUrl string `json:"original_url"`
}
