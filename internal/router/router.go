// internal/router/router.go
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/category" // <-- ATUALIZE O CAMINHO
)

// SetupRoutes configura as rotas da aplicação.
func SetupRoutes(app *fiber.App, categoryHandler *category.Handler) {
	// Middleware
	app.Use(logger.New())

	// Agrupar rotas da API
	api := app.Group("/api/v1")

	// Rotas de Categoria
	categoryRoutes := api.Group("/categories")
	categoryRoutes.Get("/", categoryHandler.GetAllCategories)
	categoryRoutes.Get("/:id", categoryHandler.GetCategoryByID)
	categoryRoutes.Post("/", categoryHandler.CreateCategory)
	categoryRoutes.Put("/:id", categoryHandler.UpdateCategory)
	categoryRoutes.Delete("/:id", categoryHandler.DeleteCategory)
}
