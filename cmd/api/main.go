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
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/category"
	gql "github.com/raimundocoelho-ti/sabiosystem-api/internal/graphql"
)

func main() {
	// 1. Carregar configuração
	config.LoadConfig()

	// 2. Conectar ao banco de dados
	database.ConnectDB()

	// 3. Injeção de Dependências
	categoryRepo := category.NewRepository(database.DB)
	categoryService := category.NewService(categoryRepo)

	// 4. Injetar o serviço nos resolvers GraphQL
	resolver := &gql.Resolver{
		CategoryService: categoryService,
	}

	// 5. Criar o schema GraphQL
	schema, err := gql.NewSchema(resolver)
	if err != nil {
		log.Fatalf("Falha ao criar o schema GraphQL: %v", err)
	}

	// 6. Criar um handler HTTP para o schema GraphQL
	gqlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		Playground: true, // Habilita o GraphQL Playground
	})

	// 7. Iniciar o Fiber
	app := fiber.New()
	app.Use(logger.New())

	// 8. Configurar a ROTA ÚNICA para GraphQL
	app.All("/graphql", adaptor.HTTPHandler(gqlHandler))

	// 9. Iniciar o servidor
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // Porta padrão
	}
	log.Printf("Servidor GraphQL rodando em http://localhost:%s/graphql", port)
	log.Fatal(app.Listen(":" + port))
}
