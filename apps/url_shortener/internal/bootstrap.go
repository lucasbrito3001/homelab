package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/lucasbrito3001/url_shortner/internal/adapters/in/http"
	"github.com/lucasbrito3001/url_shortner/internal/adapters/out/mongodb"
	redisinternal "github.com/lucasbrito3001/url_shortner/internal/adapters/out/redis"
	"github.com/lucasbrito3001/url_shortner/internal/app/ports"
	"github.com/lucasbrito3001/url_shortner/internal/config"
	"github.com/lucasbrito3001/url_shortner/internal/domain"
	"github.com/lucasbrito3001/url_shortner/internal/usecases"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func bootstrapConnections(env *config.Environment) (*mongo.Client, *redis.Client, error) {
	mongoClient, err := mongodb.NewClient(env.MongoURI)
	if err != nil {
		return nil, nil, err
	}

	redisConfig := redisinternal.RedisConfig{
		Host:     env.RedisHost,
		Port:     env.RedisPort,
		Password: env.RedisPass,
		DB:       0,
	}
	redisClient, err := redisinternal.NewClient(redisConfig)
	if err != nil {
		return nil, nil, err
	}

	return mongoClient, redisClient, nil
}

func bootstrapRepositories(mongoClient *mongo.Client, redisClient *redis.Client, cacheTtl int) (ports.UrlRepository, ports.CounterRepository) {
	mongoDbUrlRepository := mongodb.NewMongoDbUrlRepository(mongoClient.Database("url_shortner"))
	// redisUrlRepository := redisinternal.NewCachedUrlRepository(mongoDbUrlRepository, redisClient, time.Second*time.Duration(cacheTtl))
	redisCounterRepository := redisinternal.NewRedisCounterRepository(redisClient)

	return mongoDbUrlRepository, redisCounterRepository
}

func bootstrapDomain(env *config.Environment) (*domain.Shortener, error) {
	shortener, err := domain.NewShortener(env.Base62Alphabet)
	if err != nil {
		return nil, err
	}

	return shortener, nil
}

func bootstrapUseCases(urlRepository ports.UrlRepository, counterRepository ports.CounterRepository, shortener *domain.Shortener, counterOffset int64) (ports.ShortenUrlUseCase, ports.RedirectUrlUseCase) {
	shortenUrlUseCase := usecases.NewShortenUrl(urlRepository, counterRepository, shortener, counterOffset)
	redirectUrlUseCase := usecases.NewRedirectUrl(urlRepository)

	return shortenUrlUseCase, redirectUrlUseCase
}

func bootstrapRouter(shortenUrlUseCase ports.ShortenUrlUseCase, redirectUrlUseCase ports.RedirectUrlUseCase) *gin.Engine {
	router := gin.Default()

	shortenedUrlController := http.NewShortenedUrlController(shortenUrlUseCase, redirectUrlUseCase)
	config.SetRoutes(router, *shortenedUrlController)

	return router
}
