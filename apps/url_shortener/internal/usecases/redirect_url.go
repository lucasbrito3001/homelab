package usecases

import (
	"context"

	"github.com/lucasbrito3001/url_shortner/internal/app/ports"
	"github.com/lucasbrito3001/url_shortner/internal/domain"
)

type redirectUrl struct {
	urlRepository ports.UrlRepository
}

func NewRedirectUrl(urlRepository ports.UrlRepository) ports.RedirectUrlUseCase {
	return &redirectUrl{
		urlRepository: urlRepository,
	}
}

func (u *redirectUrl) Execute(ctx context.Context, req *ports.RedirectUrlRequest) (*ports.RedirectUrlResponse, error) {
	code, err := domain.NewCode(req.Code)
	if err != nil {
		return nil, err
	}

	shortened, err := u.urlRepository.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return u.mapDomainToResponse(shortened), nil
}

func (u *redirectUrl) mapDomainToResponse(url *domain.ShortenedUrl) *ports.RedirectUrlResponse {
	return &ports.RedirectUrlResponse{
		OriginalUrl: string(url.Original),
	}
}
