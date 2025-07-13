// internal/domain/category/handler.go
package category

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler contém o serviço de categoria.
type Handler struct {
	service Service
}

// NewHandler cria uma nova instância do handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetAllCategories manipula a requisição para listar categorias paginadas.
func (h *Handler) GetAllCategories(c *fiber.Ctx) error {
	pageQuery := c.Query("page", "1")
	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		page = 1
	}

	paginatedResult, err := h.service.GetAllCategories(page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve categories"})
	}

	return c.Status(fiber.StatusOK).JSON(paginatedResult)
}

// GetCategoryByID manipula a requisição para buscar uma categoria por ID.
func (h *Handler) GetCategoryByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	category, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

// CreateCategory manipula a requisição para criar uma nova categoria.
func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	var dto CreateCategoryDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Validação simples
	if dto.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	}

	category, err := h.service.CreateCategory(dto)
	if err != nil {
		// Pode ser um erro de 'unique constraint'
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create category"})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// UpdateCategory manipula a requisição para atualizar uma categoria existente.
func (h *Handler) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var dto UpdateCategoryDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if dto.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	}

	category, err := h.service.UpdateCategory(uint(id), dto)
	if err != nil {
		// Pode ser porque a categoria não foi encontrada
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

// DeleteCategory manipula a requisição para deletar uma categoria.
func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
