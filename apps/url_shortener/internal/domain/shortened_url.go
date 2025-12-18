package domain

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/lucasbrito3001/url_shortner/pkg/validator"
)

type (
	OriginalUrl string
	Code        string

	ShortenedUrl struct {
		ID        string
		Original  OriginalUrl
		Code      Code
		CreatedAt time.Time
	}
)

func NewOriginalUrl(raw string) (OriginalUrl, error) {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return "", err
	}

	if u.Scheme != "https" {
		return "", ErrNotSecureOriginalUrl
	}

	return OriginalUrl(raw), nil
}

func NewCode(raw string) (Code, error) {
	if len(raw) > 7 {
		return "", ErrTooLargeCode
	}

	if !validator.IsAlphanumeric(raw) {
		return "", ErrNonAlphanumericCode
	}

	return Code(raw), nil
}

func NewShortenedUrl(original OriginalUrl, code Code) *ShortenedUrl {
	uuid := uuid.New().String()
	now := time.Now().UTC()

	return &ShortenedUrl{
		ID:        uuid,
		Original:  original,
		Code:      code,
		CreatedAt: now,
	}
}
