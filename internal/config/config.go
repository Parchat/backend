package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config contiene la configuración de la aplicación
type Config struct {
	Port             string
	FirebaseCredFile string
	Environment      string
	ServerURL        string
}

// NewConfig crea una nueva instancia de Config
func NewConfig() *Config {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:             getEnv("PORT", "8080"),
		FirebaseCredFile: getEnv("FIREBASE_CREDENTIALS", "./firebase-credentials.json"),
		Environment:      getEnv("ENVIRONMENT", "development"),
		ServerURL:        getEnv("SERVER_URL", "http://localhost:8080"),
	}
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetFirebaseCredentialsPath devuelve la ruta absoluta al archivo de credenciales de Firebase
func (c *Config) GetFirebaseCredentialsPath() string {
	if filepath.IsAbs(c.FirebaseCredFile) {
		return c.FirebaseCredFile
	}

	// Si es una ruta relativa, convertirla a absoluta
	absPath, err := filepath.Abs(c.FirebaseCredFile)
	if err != nil {
		return c.FirebaseCredFile
	}
	return absPath
}
