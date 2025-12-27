package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func saveBook(ctx context.Context, datasource *redis.Client, book *Book) error {
	fmt.Println("adding book:", book.SKU)
	fmt.Println(book)

	bookJson, err := json.Marshal(book)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	err = datasource.Set(ctx, book.SKU, bookJson, 0).Err()
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	return nil
}

func removeBook(ctx context.Context, datasource *redis.Client, sku string) error {
	fmt.Println("removing book:", sku)

	err := datasource.Del(ctx, sku).Err()
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	return nil
}

func notify(ctx context.Context, action string, name string) {
	notification := struct {
		Identifier string `json:"identifier"`
		Type       string `json:"type"`
	}{
		Identifier: name,
		Type:       action,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if _, err := http.Post(os.Getenv("NOTIFICATION_URL"), "application/json", bytes.NewReader(body)); err != nil {
		fmt.Println("error:", err)
		return
	}

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

	router.POST("/", func(c *gin.Context) {
		var book Book
		c.ShouldBindJSON(&book)

		err := saveBook(c, redisClient, &book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		notify(c, "BOOK_CREATED", book.Name)

		c.JSON(http.StatusCreated, gin.H{
			"message": "created successfully",
		})
	})

	router.DELETE("/:sku", func(c *gin.Context) {
		sku := c.Param("sku")

		err := removeBook(c, redisClient, sku)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "deleted successfully",
		})
	})

	router.Run()
}
