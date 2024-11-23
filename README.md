# Encryption Service

## Overview

The **Encryption Service** is a high-performance microservice for encrypting and decrypting strings using the AES-256-CTR algorithm. It includes support for Redis caching and MySQL persistence to improve performance and scalability.

---

## Features

- **Encryption & Decryption**: Handles secure encryption and decryption of plaintexts.
- **Caching**: Uses Redis for caching results to reduce latency.
- **Database Integration**: MySQL for persistence of encrypted and decrypted values.
- **Resource Monitoring**: Logs CPU and memory utilization for each request.
- **Protobuf-based API**: gRPC APIs for efficient communication.

---

## Requirements

- **Go**: Version 1.23 or higher
- **Redis**: For caching
- **MySQL**: For persistence
- **Environment Variables**: Configured for encryption settings and database connections

---

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/RatulSaqibKhan/grpc-encryption-service.git
cd grpc-encryption-service
```
### 2. Run using docker
  - cd docker
  - create env files in `docker/envs` using example file
  - create `docker/.env` file from `docker/.env.example`
  - create `docker/docker-compose.override.yml` from `docker/docker-compose.override.yml.example`
  - Run `docker compose build encryption-service`
  - Run `docker compose up -d encryption-service` to run the service

