// config/config.go
package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadConfig carrega as variáveis de ambiente do arquivo .env.
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables from OS")
	}
}
