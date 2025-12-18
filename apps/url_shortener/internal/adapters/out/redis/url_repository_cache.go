package redis

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/lucasbrito3001/url_shortner/internal/app/ports"
	"github.com/lucasbrito3001/url_shortner/internal/domain"
	"github.com/redis/go-redis/v9"
)

type CachedUrlRepository struct {
	inner ports.UrlRepository
	redis *redis.Client
	ttl   time.Duration
}

func NewCachedUrlRepository(inner ports.UrlRepository, redis *redis.Client, ttl time.Duration) ports.UrlRepository {
	return &CachedUrlRepository{
		inner: inner,
		redis: redis,
		ttl:   ttl,
	}
}

func (r *CachedUrlRepository) Save(ctx context.Context, url *domain.ShortenedUrl) error {
	return r.inner.Save(ctx, url)
}

func (r *CachedUrlRepository) FindByCode(ctx context.Context, code domain.Code) (*domain.ShortenedUrl, error) {
	val, err := r.redis.Get(ctx, "code:"+string(code)).Result()
	if err != nil {
		return r.inner.FindByCode(ctx, code)
	}

	var model shortenedUrlCacheModel
	json.Unmarshal([]byte(val), &model)

	jitter := time.Duration(rand.Intn(30)) * time.Minute
	r.redis.Expire(ctx, "code:"+string(code), r.ttl+jitter)

	domain, err := model.toDomain()
	if err != nil {
		return nil, err
	}

	return domain, nil
}
