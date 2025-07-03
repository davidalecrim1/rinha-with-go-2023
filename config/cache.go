package config

import (
	"fmt"
	"go-rinha-de-backend-2023/config/env"
	"log"

	"github.com/redis/rueidis"
)

func InitializeCache() rueidis.Client {
	cache_host := env.GetEnvOrSetDefault("CACHE_HOST", "localhost")
	cache_port := env.GetEnvOrSetDefault("CACHE_PORT", "6379")
	cache_url := fmt.Sprintf("%s:%s", cache_host, cache_port)

	cache, err := rueidis.NewClient(
		rueidis.ClientOption{
			InitAddress:      []string{cache_url},
			AlwaysPipelining: true,
		})

	if err != nil {
		log.Fatalf("error loading cache configuration: %v", err)
	}

	return cache
}
