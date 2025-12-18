package mongodb

import (
	"time"

	"github.com/lucasbrito3001/url_shortner/internal/domain"
)

type shortenedUrlDocument struct {
	ID        string    `bson:"_id"`
	Original  string    `bson:"original_url"`
	Code      string    `bson:"code"`
	CreatedAt time.Time `bson:"created_at"`
}

func (ud *shortenedUrlDocument) toDomain() *domain.ShortenedUrl {
	originalUrl, err := domain.NewOriginalUrl(ud.Original)
	if err != nil {
		return nil
	}

	code, err := domain.NewCode(ud.Code)
	if err != nil {
		return nil
	}

	return &domain.ShortenedUrl{
		ID:        ud.ID,
		Original:  originalUrl,
		Code:      code,
		CreatedAt: ud.CreatedAt,
	}
}

func fromDomain(url *domain.ShortenedUrl) *shortenedUrlDocument {
	return &shortenedUrlDocument{
		ID:        url.ID,
		Original:  string(url.Original),
		Code:      string(url.Code),
		CreatedAt: url.CreatedAt,
	}
}
