package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	_ = initializePostgres()
	_ = initializeRedis()

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	log.Printf("Server listenting on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}