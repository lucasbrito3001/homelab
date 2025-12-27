package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Notification struct {
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
}

func notify(ctx context.Context, notification *Notification) error {
	switch notification.Type {
	case "BOOK_CREATED":
		fmt.Printf("The book %s has been created", notification.Identifier)
		return nil
	case "BOOK_DELETED":
		fmt.Printf("The book %s has been deleted", notification.Identifier)
		return nil
	}

	return errors.New("the event is not supported")
}

func main() {
	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "healthy",
		})
	})

	router.POST("/", func(c *gin.Context) {
		var notification Notification
		c.ShouldBindJSON(&notification)

		err := notify(c, &notification)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "created successfully",
		})
	})

	router.Run("0.0.0.0:8080")
}
