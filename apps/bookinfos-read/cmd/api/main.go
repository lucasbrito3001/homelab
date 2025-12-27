package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

type Book struct {
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	ReleaseDate time.Time `json:"release_date"`
	Price       int64     `json:"price"`
	SKU         string    `json:"sku"`
}

func configureRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	err := checkRedisConn(client)
	if err != nil {
		return nil
	}

	return client
}

func checkRedisConn(client *redis.Client) error {
	ctx := context.Background()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	fmt.Println("redis connected:", pong)

	return nil
}

func getBySKU(ctx context.Context, datasource *redis.Client, sku string) (*Book, error) {
	val, err := datasource.Get(ctx, sku).Result()
	if err != nil {
		fmt.Println("error getting book by sku:", err)
		return nil, err
	}

	var book Book
	err = json.Unmarshal([]byte(val), &book)
	if err != nil {
		log.Fatalf("error unmarshaling book: %v", err)
	}

	return &book, nil
}

func main() {
	redisClient := configureRedis()
	if redisClient == nil {
		panic("error connecting to redis")
	}

	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		err := checkRedisConn(redisClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "redis connection is down",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "healthy",
		})
	})

	router.GET("/:sku", func(c *gin.Context) {
		sku := c.Param("sku")

		book, err := getBySKU(c, redisClient, sku)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, book)
	})

	router.Run()
}
