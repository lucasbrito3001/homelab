package ports

import (
	"context"
)

type ShortenUrlUseCase interface {
	Execute(ctx context.Context, request *ShortenUrlRequest) (*ShortenUrlResponse, error)
}

type ShortenUrlRequest struct {
	OriginalUrl string `json:"original_url"`
}

type ShortenUrlResponse struct {
	Code string `json:"code"`
}
