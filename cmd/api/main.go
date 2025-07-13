/*
|------------------------------------------------
| File: cmd/api/main.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/

package main

import (
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"log"
	"os"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/graphql-go/handler"
	"github.com/raimundocoelho-ti/sabiosystem-api/config"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/database"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/category" // Caminho corrigido
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/user"
	gql "github.com/raimundocoelho-ti/sabiosystem-api/internal/graphql"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()

	// 1. Instanciar repositórios e serviços para TODOS os módulos
	categoryRepo := category.NewRepository(database.DB)
	categoryService := category.NewService(categoryRepo)
	userRepo := user.NewRepository(database.DB)
	userService := user.NewService(userRepo)

	// 2. Agrupar serviços para passar ao montador de schema
	schemaServices := gql.SchemaServices{
		CategorySvc: categoryService,
		UserSvc:     userService,
	}

	// 3. Criar o schema a partir dos serviços agrupados
	schema, err := gql.NewSchema(schemaServices)
	if err != nil {
		log.Fatalf("Falha ao criar o schema GraphQL: %v", err)
	}

	// 4. Criar o handler GraphQL
	gqlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		Playground: true,
	})

	// 5. Iniciar e configurar o Fiber
	app := fiber.New()
	app.Use(logger.New())
	app.All("/graphql", adaptor.HTTPHandler(gqlHandler))

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor GraphQL rodando em http://localhost:%s/graphql", port)
	log.Fatal(app.Listen(":" + port))
}
