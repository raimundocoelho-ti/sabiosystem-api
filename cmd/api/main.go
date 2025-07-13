// cmd/api/main.go
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/raimundocoelho-ti/sabiosystem-api/config"                   // <-- ATUALIZE O CAMINHO
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/database"        // <-- ATUALIZE O CAMINHO
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/category" // <-- ATUALIZE O CAMINHO
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/router"          // <-- ATUALIZE O CAMINHO
)

func main() {
	// 1. Carregar configuração do .env
	config.LoadConfig()

	// 2. Conectar ao banco de dados
	database.ConnectDB()

	// 3. Injeção de Dependências (DI)
	categoryRepo := category.NewRepository(database.DB)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	// 4. Iniciar o Fiber
	app := fiber.New()

	// 5. Configurar as rotas
	router.SetupRoutes(app, categoryHandler)

	// 6. Iniciar o servidor
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // Porta padrão
	}
	log.Fatal(app.Listen(":" + port))
}
