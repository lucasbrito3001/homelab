package ports

import (
	"context"

	"github.com/lucasbrito3001/url_shortner/internal/domain"
)

type UrlRepository interface {
	Save(ctx context.Context, url *domain.ShortenedUrl) error
	FindByCode(ctx context.Context, code domain.Code) (*domain.ShortenedUrl, error)
}
