package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Servicio
	AppEnv   string
	GRPCHost string
	GRPCPort string

	// PostgreSQL
	PSQLHost     string
	PSQLPort     string
	PSQLUser     string
	PSQLPassword string
	PSQLDB       string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Security
	EncryptionSecret   string
	BindexSecret       string
	PasswordBcryptCost int
	JWTSecret          string
	SessionTTL         time.Duration
}

// LoadEnv carga y valida la configuración del sistema.
func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("ADVERTENCIA: No se encontró el archivo .env. Se intentarán cargar las variables de entorno desde el sistema.")
	}

	cfg := &Config{
		// Servicio
		AppEnv:   GetEnvOrDefault("APP_ENV", "prod"),
		GRPCHost: GetEnvOrDefault("GRPC_HOST", "0.0.0.0"),
		GRPCPort: GetEnvOrDefault("GRPC_PORT", "50051"),

		// PostgreSQL
		PSQLHost:     GetEnvOrCrash("PSQL_HOST"),
		PSQLPort:     GetEnvOrDefault("PSQL_PORT", "5432"),
		PSQLUser:     GetEnvOrCrash("PSQL_USER"),
		PSQLPassword: GetEnvOrCrash("PSQL_PASSWORD"),
		PSQLDB:       GetEnvOrCrash("PSQL_DB"),

		// Redis
		RedisHost:     GetEnvOrCrash("REDIS_HOST"),
		RedisPort:     GetEnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword: GetEnvOrCrash("REDIS_PASSWORD"),
		RedisDB:       GetEnvAsIntOrDefault("REDIS_DB", 0),

		// Security
		EncryptionSecret:   GetEnvOrCrash("ENCRYPTION_SECRET"),
		BindexSecret:       GetEnvOrCrash("BINDEX_SECRET"),
		PasswordBcryptCost: GetEnvAsIntOrDefault("PASSWORD_BCRYPT_COST", 12),
		JWTSecret:          GetEnvOrCrash("JWT_SECRET"),
		SessionTTL:         GetEnvAsDurationOrDefault("SESSION_TTL", time.Hour),
	}

	return cfg
}

// GetEnvOrDefault busca una variable de entorno (string); si no existe, devuelve un valor por defecto.
func GetEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

// GetEnvOrCrash busca una variable de entorno (string); si no existe, aborta la ejecución.
func GetEnvOrCrash(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("ERROR CRÍTICO: La variable de entorno '%s' es obligatoria y no está definida.", key)
	}
	return value
}

// GetEnvAsIntOrDefault busca una variable de entorno (int); si no existe, devuelve un valor por defecto.
func GetEnvAsIntOrDefault(key string, fallback int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists || valueStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("ERROR CRÍTICO: La variable de entorno '%s' debe ser un número entero válido.", key)
	}
	return val
}

// GetEnvAsDurationOrDefault busca una variable de entorno y la convierte a time.Duration; si no existe, devuelve el valor por defecto.
func GetEnvAsDurationOrDefault(key string, fallback time.Duration) time.Duration {
	valueStr, exists := os.LookupEnv(key)
	if !exists || valueStr == "" {
		return fallback
	}

	d, err := time.ParseDuration(valueStr)
	if err == nil {
		return d
	}

	if seconds, err := strconv.Atoi(valueStr); err == nil {
		return time.Duration(seconds) * time.Second
	}

	log.Fatalf("ERROR CRÍTICO: La variable de entorno '%s' debe tener un formato de duración válido (ej. 1h, 60m, 3600s).", key)
	return fallback
}