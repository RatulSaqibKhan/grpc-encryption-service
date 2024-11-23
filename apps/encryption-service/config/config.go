package config

import "os"

type Config struct {
	MySQLDSN     string
	RedisAddr    string
	RedisPassword string
	GRPCPort     string
}

func LoadConfig() Config {
	return Config{
		MySQLDSN:      os.Getenv("MYSQL_DSN"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		GRPCPort:      os.Getenv("GRPC_PORT"),
	}
}
