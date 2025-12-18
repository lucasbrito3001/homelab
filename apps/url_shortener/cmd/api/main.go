package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucasbrito3001/url_shortner/internal"
	"github.com/lucasbrito3001/url_shortner/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	env, err := config.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	app, err := internal.NewApplication(env)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("error during application run: %v", err)
		}
	}()

	log.Println("app started successfully!")

	<-ctx.Done()

	timeoutDuration := 10 * time.Second
	log.Printf("SIGTERM received, starting graceful shutdown by %v", timeoutDuration)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error closing connections: %v", err)
	}

	log.Println("application shutdown successfully.")
}
