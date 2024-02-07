package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type URL string

type URLS []URL

func (urls URLS) IsValid() bool {
	for _, u := range urls {
		_, err := url.ParseRequestURI(string(u))
		if err != nil {
			panic(fmt.Sprintf("invalid urls: %s", string(u)))
		}
	}
	return true
}

func NewURLS(urls ...string) URLS {
	var result URLS
	for _, url := range urls {
		result = append(result, URL(url))
	}
	return result
}

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
	Staging     Environment = "staging"
)

type Config struct {
	PORT                   int
	RUN_SEEDS              bool
	ENVIRONMENT            Environment
	NATS_URL               string
	JAEGER_ENDPOINT        string
	SERVICE_NAME           string
	EXCHANGE_SERVICE_URL_1 string
	EXCHANGE_SERVICE_URL_2 string
	API_TOKEN              string
}

func Get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("%s: %s", key, err)
			return fallback
		}
		return i
	}
	return fallback
}

func GetEnvironment() Environment {
	if env := Get("ENV", ""); env == "" {
		return Development
	} else {
		return Environment(env)
	}
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading env file")
	}

	return &Config{
		PORT:                   GetInt("PORT", 3000),
		ENVIRONMENT:            GetEnvironment(),
		SERVICE_NAME:           Get("SERVICE_NAME", "auth"),
		EXCHANGE_SERVICE_URL_1: Get("EXCHANGE_SERVICE_URL_1", "http://localhost:3001/api/v1/rates"),
		EXCHANGE_SERVICE_URL_2: Get("EXCHANGE_SERVICE_URL_2", "http://localhost:3002/api/v1/rates"),
		API_TOKEN:              Get("API_TOKEN", "​8a395ccb-7f3e-4a5a-b35c-4fea034d24f2"),
	}
}
