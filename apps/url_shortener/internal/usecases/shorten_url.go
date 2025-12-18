package usecases

import (
	"context"

	"github.com/lucasbrito3001/url_shortner/internal/app/ports"
	"github.com/lucasbrito3001/url_shortner/internal/domain"
)

type shortenUrl struct {
	urlRepository     ports.UrlRepository
	counterRepository ports.CounterRepository
	shortener         *domain.Shortener
	idOffset          int64
}

func NewShortenUrl(urlRepository ports.UrlRepository, counterRepository ports.CounterRepository, shortener *domain.Shortener, idOffset int64) ports.ShortenUrlUseCase {
	return &shortenUrl{
		urlRepository:     urlRepository,
		counterRepository: counterRepository,
		shortener:         shortener,
		idOffset:          idOffset,
	}
}

func (u *shortenUrl) Execute(ctx context.Context, req *ports.ShortenUrlRequest) (*ports.ShortenUrlResponse, error) {
	seq, err := u.counterRepository.Increment(ctx)
	if err != nil {
		return nil, err
	}

	code := u.shortener.Encode(u.idOffset + seq)

	url, err := u.mapRequestToDomain(req, code)
	if err != nil {
		return nil, err
	}

	err = u.urlRepository.Save(ctx, url)
	if err != nil {
		return nil, err
	}

	return u.mapDomainToResponse(url), nil
}

func (u *shortenUrl) mapRequestToDomain(req *ports.ShortenUrlRequest, rawCode string) (*domain.ShortenedUrl, error) {
	originalUrl, err := domain.NewOriginalUrl(req.OriginalUrl)
	if err != nil {
		return nil, err
	}

	code, err := domain.NewCode(rawCode)
	if err != nil {
		return nil, err
	}

	return domain.NewShortenedUrl(originalUrl, code), nil
}

func (u *shortenUrl) mapDomainToResponse(url *domain.ShortenedUrl) *ports.ShortenUrlResponse {
	return &ports.ShortenUrlResponse{
		Code: string(url.Code),
	}
}
