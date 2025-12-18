package internal

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucasbrito3001/url_shortner/internal/config"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type application struct {
	engine      *gin.Engine
	env         *config.Environment
	mongoClient *mongo.Client
	redisClient *redis.Client
	httpServer  *http.Server
}

func NewApplication(env *config.Environment) (Application, error) {
	mongoClient, redisClient, err := bootstrapConnections(env)
	if err != nil {
		return nil, err
	}

	shortener, err := bootstrapDomain(env)
	if err != nil {
		return nil, err
	}

	urlRepository, counterRepository := bootstrapRepositories(mongoClient, redisClient, env.CacheTTL)
	shortenUrlUseCase, redirectUrlUseCase := bootstrapUseCases(urlRepository, counterRepository, shortener, env.CounterOffset)
	router := bootstrapRouter(shortenUrlUseCase, redirectUrlUseCase)

	return &application{
		engine:      router,
		env:         env,
		mongoClient: mongoClient,
		redisClient: redisClient,
	}, nil
}
func (a *application) Run() error {
	a.httpServer = &http.Server{
		Addr:    ":" + a.env.ServerPort,
		Handler: a.engine,
	}

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *application) Shutdown(ctx context.Context) error {
	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Printf("error closing http server: %v", err)
	}

	if err := a.redisClient.Close(); err != nil {
		log.Printf("error closing redis: %v", err)
	}

	if err := a.mongoClient.Disconnect(ctx); err != nil {
		log.Printf("error closing mongodb: %v", err)
	}

	log.Println("shutdown completed successfully")

	return nil
}
