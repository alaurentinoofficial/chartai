package infrastructure_config

import "os"

type Config struct {
	Port               string
	DatabaseConnection string
	TokenSignature     string
	RedisAddress       string
	RedisPassword      string
}

func NewConfig() *Config {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	return &Config{
		Port:               port,
		DatabaseConnection: os.Getenv("DATABASE_CONNECTION"),
		TokenSignature:     os.Getenv("TOKEN_SIGNATURE"),
		RedisAddress:       os.Getenv("REDIS_ADDRESS"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
	}
}
