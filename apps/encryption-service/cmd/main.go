package main

import (
	"database/sql"
	"log"

	"encryption-service/config"
	"encryption-service/internal/handlers"
	repository "encryption-service/internal/repositories"
	"encryption-service/internal/service"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize MySQL
	db, err := sql.Open("mysql", cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Initialize Redis with password
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	defer redisClient.Close()

	// Create repositories
	mysqlRepo := repository.MySQLRepository{DB: db}
	redisRepo := repository.RedisRepository{Client: redisClient}

	// Create handler
	handler := handlers.EncryptionHandler{
		RedisRepo: redisRepo,
		MySQLRepo: mysqlRepo,
	}

	// Start gRPC server
	if err := service.StartGRPCServer(&handler, cfg.GRPCPort); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
