package handlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
	"os"

	"encryption-service/internal/encryption"
	repository "encryption-service/internal/repositories"
	pb "encryption-service/proto"
)

type EncryptionHandler struct {
	RedisRepo repository.RedisRepository
	MySQLRepo repository.MySQLRepository
}

var logger = logrus.New()

// Initialize logger to use JSON format and set log level based on environment
func init() {
	// Set JSON formatter
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Get the environment variable
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // Default to 'dev' if no environment variable is set
	}

	// Set log level based on environment
	switch env {
	case "prod":
		logger.SetLevel(logrus.WarnLevel) // In production, log at 'warn' or higher
	case "stage":
		logger.SetLevel(logrus.InfoLevel) // In staging, log at 'info' level
	default:
		logger.SetLevel(logrus.DebugLevel) // In dev, log everything (debug, info, warn, error, fatal)
	}
}

// logSystemStats logs the CPU and memory utilization.
func logSystemStats(operation string, requestID string) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Log memory usage
	logger.WithFields(logrus.Fields{
		"operation": operation,
		"request_id": requestID,
		"memory_alloc": memStats.Alloc / 1024,
		"memory_total_alloc": memStats.TotalAlloc / 1024,
		"memory_sys": memStats.Sys / 1024,
		"gc_count": memStats.NumGC,
	}).Infof("System stats")

	// CPU stats (simple tracking with Goroutines)
	numCPU := runtime.NumCPU()
	logger.WithFields(logrus.Fields{
		"operation": operation,
		"request_id": requestID,
		"cpu_count": numCPU,
	}).Infof("CPU stats")
}

// Encrypt handles encryption of an array of plaintexts.
func (h *EncryptionHandler) Encrypt(ctx context.Context, req *pb.EncryptRequest) (*pb.EncryptResponse, error) {
	startTime := time.Now() // Track the start time
	requestID := req.RequestId
	plainTexts := req.Plaintexts
	
	logger.WithFields(logrus.Fields{
		"operation": "[Encryption]",
		"request_id": requestID,
		"texts": plainTexts,
		"total_texts": len(plainTexts),
	}).Warnf("Encryption Request")

	defer func() {
		duration := time.Since(startTime)
		logger.WithFields(logrus.Fields{
			"operation": "Encrypt",
			"request_id": requestID,
			"duration": duration,
		}).Info("Total Execution Time")
		logSystemStats("Encrypt", requestID)
	}()

	var encryptedTexts []string
	var redisKey string

	for _, plaintext := range plainTexts {
		// Check if the encrypted data exists in MySQL
		encrypted, err := h.MySQLRepo.GetEncrypted(plaintext)
		if err == nil {
			encryptedTexts = append(encryptedTexts, encrypted)
			continue
		}

		// Encrypt the plaintext
		encrypted, err = encryption.Encrypt(plaintext)
		if err != nil {
			logger.Errorf("Error during encryption of '%s': %v", plaintext, err)
			return nil, err
		}

		// Save to MySQL
		if err := h.MySQLRepo.SaveEncrypted(plaintext, encrypted); err != nil {
			logger.Errorf("Error saving to MySQL: %v", err)
		}

		redisKey = "go_encryption:decrypted:" + encrypted
		// Save to Redis
		if err := h.RedisRepo.Set(ctx, redisKey, plaintext); err != nil {
			logger.Errorf("Error saving to Redis: %v", err)
		}

		encryptedTexts = append(encryptedTexts, encrypted)
	}

	return &pb.EncryptResponse{EncryptedTexts: encryptedTexts}, nil
}

// Decrypt handles decryption of an array of encrypted texts.
func (h *EncryptionHandler) Decrypt(ctx context.Context, req *pb.DecryptRequest) (*pb.DecryptResponse, error) {
	startTime := time.Now() // Track the start time
	requestID := req.RequestId
	encryptedTexts := req.EncryptedTexts
	
	logger.WithFields(logrus.Fields{
		"operation": "[Decryption]",
		"request_id": requestID,
		"encrypted_texts": encryptedTexts,
		"total_encrypted_texts": len(encryptedTexts),
	}).Warnf("Decryption Request")

	defer func() {
		duration := time.Since(startTime)
		logger.WithFields(logrus.Fields{
			"operation": "Decrypt",
			"request_id": requestID,
			"duration": duration,
		}).Info("Total Execution Time")
		logSystemStats("Decrypt", requestID)
	}()

	var plaintexts []string
	var redisKey string

	for _, encryptedText := range encryptedTexts {
		if encryptedText == "" {
			logger.Warn("Error: Empty encrypted text in decryption request")
			continue // Skip empty encrypted texts
		}

		var plaintext string
		var err error
		redisKey = "go_encryption:decrypted:" + encryptedText

		// Step 1: Check Redis for the plaintext
		plaintext, err = h.RedisRepo.Get(ctx, redisKey)
		if err == nil {
			plaintexts = append(plaintexts, plaintext)
			continue
		}

		// Step 2: Check MySQL for the plaintext
		plaintext, err = h.MySQLRepo.GetDecrypted(encryptedText)
		if err == nil {
			// Store in Redis for future use
			if cacheErr := h.RedisRepo.Set(ctx, encryptedText, plaintext); cacheErr != nil {
				logger.Errorf("Error caching to Redis: %v", cacheErr)
			}
			plaintexts = append(plaintexts, plaintext)
			continue
		}

		// Step 3: Decrypt using the decryption script
		plaintext, err = encryption.Decrypt(encryptedText)
		if err != nil {
			logger.Warnf("Skipping decryption of '%s' due to error: %v", encryptedText, err)
			continue
		}

		// Save the decrypted plaintext to MySQL
		if dbErr := h.MySQLRepo.SaveEncrypted(encryptedText, plaintext); dbErr != nil {
			logger.Errorf("Error saving to MySQL: %v", dbErr)
		}

		// Save the decrypted plaintext to Redis
		if cacheErr := h.RedisRepo.Set(ctx, encryptedText, plaintext); cacheErr != nil {
			logger.Errorf("Error caching to Redis: %v", cacheErr)
		}

		plaintexts = append(plaintexts, plaintext)
	}

	return &pb.DecryptResponse{Plaintexts: plaintexts}, nil
}
