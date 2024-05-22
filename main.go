package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	// Configuração do Redis
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	redisClient := NewRedisClient(redisHost, redisPort)

	// Configuração do Rate Limiter
	rateLimitIP := getEnvAsInt("RATE_LIMIT_IP", 10)
	rateLimitToken := getEnvAsInt("RATE_LIMIT_TOKEN", 100)
	blockDurationIP := getEnvAsDuration("BLOCK_DURATION_IP", 300*time.Second)
	blockDurationToken := getEnvAsDuration("BLOCK_DURATION_TOKEN", 300*time.Second)

	rateLimiter := NewRateLimiter(redisClient, rateLimitIP, rateLimitToken, blockDurationIP, blockDurationToken)

	http.Handle("/", rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})))

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getEnvAsInt(name string, defaultValue int) int {
	value, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Valor inválido para %s: %v", name, err)
	}
	return intValue
}

func getEnvAsDuration(name string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}
	durationValue, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("Valor inválido para %s: %v", name, err)
	}
	return durationValue
}
