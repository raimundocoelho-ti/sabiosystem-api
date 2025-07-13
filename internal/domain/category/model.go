/*
|------------------------------------------------
| File: internal/domain/category/model.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package category

import "time"

// Category representa a entidade no banco de dados.
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `gorm:"not null" json:"agent_id"` // <-- MUDANÇA: Adicionado AgentID
	Name      string    `gorm:"not null" json:"name"`     // Removi o 'unique' daqui. Nomes podem se repetir entre agentes diferentes.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCategoryDTO é o Data Transfer Object para a criação de uma categoria.
type CreateCategoryDTO struct {
	AgentID uint   `json:"agent_id"` // <-- MUDANÇA: Adicionado AgentID
	Name    string `json:"name"`
}

// UpdateCategoryDTO é o Data Transfer Object para a atualização de uma categoria.
type UpdateCategoryDTO struct {
	Name string `json:"name"`
}

// PaginatedCategories é a estrutura de resposta para a lista paginada de categorias.
type PaginatedCategories struct {
	Data       []Category `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
}
