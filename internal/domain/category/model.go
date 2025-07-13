// internal/domain/category/model.go
package category

import "time"

// Category representa a entidade no banco de dados.
// As tags `gorm` e `json` mapeiam os campos da struct para a tabela e para as respostas JSON.
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCategoryDTO é o Data Transfer Object para a criação de uma categoria.
type CreateCategoryDTO struct {
	Name string `json:"name" binding:"required"`
}

// UpdateCategoryDTO é o Data Transfer Object para a atualização de uma categoria.
type UpdateCategoryDTO struct {
	Name string `json:"name" binding:"required"`
}

// PaginatedCategories é a estrutura de resposta para a lista paginada de categorias.
type PaginatedCategories struct {
	Data       []Category `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
}
