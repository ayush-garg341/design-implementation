package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/linkvault/config"
	"github.com/linkvault/internal/analytics"
	db_pkg "github.com/linkvault/internal/db"
	"github.com/linkvault/internal/events"
	"github.com/linkvault/internal/handler"
	"github.com/linkvault/internal/middleware"
	redisInternal "github.com/linkvault/internal/redis"
	"github.com/linkvault/internal/service"
	"github.com/linkvault/internal/store"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load the config
	cfg := config.Config()

	// Create the db connection
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Run the migrations
	db_pkg.RunMigrations(db)

	// Create the redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	// Create the kafka client
	kafkaClient := events.NewKafkaProducer(cfg.KafkaBrokers, "analytics", 1000, 4)

	// Create auth store and handler
	authStore := store.NewPostgresUserStore(db)
	authService := service.NewUserService(authStore)
	authHandler := handler.NewAuthHandler(authService)

	// Create link store and handler
	linkStore := store.NewPostgresLinkStore(db)
	linkService := service.NewLinkService(linkStore)
	linkHandler := handler.NewLinkHandler(linkService)

	// Create link stats handler
	analyticsClient, err := analytics.NewClient("localhost:8085")
	if err != nil {
		log.Fatalf("Something went wrong")
	}
	linkStatsHandler := handler.NewLinkStatsHandler(analyticsClient)

	// Redirect handler
	cache := redisInternal.NewRedisCache(redisClient)
	redirectService := service.NewRedirectService(linkStore, cache, kafkaClient)
	redirectHandler := handler.NewRedirectHandler(redirectService)

	r := gin.New()

	auth := r.Group("/auth")
	auth.Use(middleware.Timeout(500 * time.Millisecond))
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.Timeout(100 * time.Millisecond))
	api.Use(middleware.Auth())
	{
		api.POST("/link", linkHandler.Create)
		api.GET("/links", linkHandler.List)
		api.GET("/stats/link/:code", linkStatsHandler.GetLinkStats)
	}

	r.GET("/long/:shortcode", middleware.Timeout(100*time.Millisecond), redirectHandler.Redirect)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Println("server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit
	log.Println("shutdown signal received")

	// Timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	// Close infrastructure
	log.Println("closing kafka producer")
	kafkaClient.Close()

	log.Println("closing redis")
	redisClient.Close()

	log.Println("closing database")
	db.Close()

	log.Println("server exiting")

}
