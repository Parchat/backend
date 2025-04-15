package config

import (
	"os"
	"path/filepath"
)

// Config contiene la configuración de la aplicación
type Config struct {
	Port            string
	FirebaseCredFile string
	Environment     string
}

// NewConfig crea una nueva instancia de Config
func NewConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		FirebaseCredFile: getEnv("FIREBASE_CREDENTIALS", "./firebase-credentials.json"),
		Environment:     getEnv("ENVIRONMENT", "development"),
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
