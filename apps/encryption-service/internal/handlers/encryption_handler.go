package handlers

import (
	"context"
	"log"
	"runtime"
	"time"

	"encryption-service/internal/encryption"
	repository "encryption-service/internal/repositories"
	pb "encryption-service/proto"
)

type EncryptionHandler struct {
	RedisRepo repository.RedisRepository
	MySQLRepo repository.MySQLRepository
}

// logSystemStats logs the CPU and memory utilization.
func logSystemStats(operation string) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Log memory usage
	log.Printf("[%s] Memory Usage: Alloc = %v KB, TotalAlloc = %v KB, Sys = %v KB, NumGC = %v",
		operation, memStats.Alloc/1024, memStats.TotalAlloc/1024, memStats.Sys/1024, memStats.NumGC)

	// CPU stats (simple tracking with Goroutines)
	numCPU := runtime.NumCPU()
	log.Printf("[%s] CPU Utilization: NumCPU = %d", operation, numCPU)
}

// Encrypt handles encryption of an array of plaintexts.
func (h *EncryptionHandler) Encrypt(ctx context.Context, req *pb.EncryptRequest) (*pb.EncryptResponse, error) {
	startTime := time.Now() // Track the start time
	defer func() {
		duration := time.Since(startTime)
		log.Printf("[Encrypt] Total Execution Time: %v", duration)
		logSystemStats("Encrypt")
	}()

	var encryptedTexts []string
	var redisKey string

	for _, plaintext := range req.Plaintexts {
		log.Printf("Encryption request of '%s'", plaintext)

		// Check if the encrypted data exists in MySQL
		encrypted, err := h.MySQLRepo.GetEncrypted(plaintext)
		if err == nil {
			encryptedTexts = append(encryptedTexts, encrypted)
			continue
		}

		// Encrypt the plaintext
		encrypted, err = encryption.Encrypt(plaintext)
		if err != nil {
			log.Printf("Error during encryption of '%s': %v", plaintext, err)
			return nil, err
		}

		// Save to MySQL
		if err := h.MySQLRepo.SaveEncrypted(plaintext, encrypted); err != nil {
			log.Printf("Error saving to MySQL: %v", err)
		}

		redisKey = "go_encryption:decrypted:" + encrypted
		// Save to Redis
		if err := h.RedisRepo.Set(ctx, redisKey, plaintext); err != nil {
			log.Printf("Error saving to Redis: %v", err)
		}

		encryptedTexts = append(encryptedTexts, encrypted)
	}

	return &pb.EncryptResponse{EncryptedTexts: encryptedTexts}, nil
}

// Decrypt handles decryption of an array of encrypted texts.
func (h *EncryptionHandler) Decrypt(ctx context.Context, req *pb.DecryptRequest) (*pb.DecryptResponse, error) {
	startTime := time.Now() // Track the start time
	defer func() {
		duration := time.Since(startTime)
		log.Printf("[Decrypt] Total Execution Time: %v", duration)
		logSystemStats("Decrypt")
	}()

	var plaintexts []string
	var redisKey string

	for _, encryptedText := range req.EncryptedTexts {
		log.Printf("Decryption request of '%s'", encryptedText)
		
		if encryptedText == "" {
			log.Printf("Error: Empty encrypted text in decryption request")
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
				log.Printf("Error caching to Redis: %v", cacheErr)
			}
			plaintexts = append(plaintexts, plaintext)
			continue
		}

		// Step 3: Decrypt using the decryption script
		plaintext, err = encryption.Decrypt(encryptedText)
		if err != nil {
			log.Printf("Skipping decryption of '%s' due to error: %v", encryptedText, err)
			continue
		}

		// Save the decrypted plaintext to MySQL
		if dbErr := h.MySQLRepo.SaveEncrypted(encryptedText, plaintext); dbErr != nil {
			log.Printf("Error saving to MySQL: %v", dbErr)
		}

		// Save the decrypted plaintext to Redis
		if cacheErr := h.RedisRepo.Set(ctx, encryptedText, plaintext); cacheErr != nil {
			log.Printf("Error caching to Redis: %v", cacheErr)
		}

		plaintexts = append(plaintexts, plaintext)
	}

	return &pb.DecryptResponse{Plaintexts: plaintexts}, nil
}
