package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-url-shortener/internal/handler"
	"go-url-shortener/internal/repository"
	"go-url-shortener/internal/router"
	"go-url-shortener/internal/service"
)

var ctx = context.Background()

func initializePostgres() *gorm.DB {

	dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
		log.Fatalf("Fatl Error: DATABASE_URL environment variable not set.")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fatal Error: Failed to connect to PostgreSQL: %v", err)
	}
	
	log.Println("Database connection established: PostgreSQL")
	return db
}

func initializeRedis() *redis.Client {
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379"
    }

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Fatal Error: Failed to connect to Redis: %v", err)
	}
	
	log.Println("Cache connection established: Redis")
	return rdb
}

func main() {
	log.Printf("Starting URL Shortener Server")

	db := initializePostgres()
	rdb := initializeRedis()

	pgRepo := repository.NewPostgresRepo(db)
	redisRepo := repository.NewRedisRepo(rdb)

	urlService := service.NewURLService(pgRepo, redisRepo)

	h := handler.NewShortenerHandler(urlService)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.APIRouter(r, h)
	router.RedirectRouter(r, h)

	srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

	go func() {
        log.Printf("Server listening on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

	stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    <-stop
    log.Println("Shutting down server...")

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Closing database connections...")
    
    sqlDB, _ := db.DB()
    sqlDB.Close()
    
    rdb.Close()

    log.Println("Server exited gracefully")
}