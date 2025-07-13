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
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/auth"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/database"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/agent"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/category"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/product" // <-- ADICIONADO
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/user"
	gql "github.com/raimundocoelho-ti/sabiosystem-api/internal/graphql"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()

	// 1. Instanciar todos os repositórios e serviços
	categoryRepo := category.NewRepository(database.DB)
	categoryService := category.NewService(categoryRepo)
	productRepo := product.NewRepository(database.DB) // <-- ADICIONADO
	productService := product.NewService(productRepo) // <-- ADICIONADO
	userRepo := user.NewRepository(database.DB)
	userService := user.NewService(userRepo)
	agentRepo := agent.NewRepository(database.DB)
	agentService := agent.NewService(agentRepo)
	authRepo := auth.NewRepository(database.DB)
	authService := auth.NewService(authRepo)

	// 2. Agrupar todos os serviços para o montador de schema
	schemaServices := gql.SchemaServices{
		CategorySvc: categoryService,
		ProductSvc:  productService, // <-- ADICIONADO
		UserSvc:     userService,
		AgentSvc:    agentService,
		AuthSvc:     authService,
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

	app.Post("/upload", product.UploadImageHandler)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor GraphQL rodando em http://localhost:%s/graphql", port)
	log.Fatal(app.Listen(":" + port))
}
