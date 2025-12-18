package redis

import (
	"time"

	"github.com/lucasbrito3001/url_shortner/internal/domain"
)

type shortenedUrlCacheModel struct {
	ID        string    `json:"id"`
	Original  string    `json:"original"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *shortenedUrlCacheModel) toDomain() (*domain.ShortenedUrl, error) {
	originalUrl, err := domain.NewOriginalUrl(c.Original)
	if err != nil {
		return nil, err
	}

	code, err := domain.NewCode(c.Code)
	if err != nil {
		return nil, err
	}

	return &domain.ShortenedUrl{
		ID:        c.ID,
		Original:  originalUrl,
		Code:      code,
		CreatedAt: c.CreatedAt,
	}, nil
}
