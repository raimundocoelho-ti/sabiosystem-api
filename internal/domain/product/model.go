/*
|------------------------------------------------
| File: internal/domain/product/model.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package product

import "time"

// Product representa o produto no banco de dados.
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AgentID     uint      `gorm:"not null" json:"agent_id"`
	CategoryID  uint      `gorm:"not null" json:"category_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductDTO - dados para criar produto
type CreateProductDTO struct {
	AgentID     uint    `json:"agent_id"`
	CategoryID  uint    `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active"`
}

// UpdateProductDTO - dados para atualizar produto
type UpdateProductDTO struct {
	CategoryID  uint    `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active"`
}

// PaginatedProducts - resposta paginada
type PaginatedProducts struct {
	Data       []Product `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	TotalPages int       `json:"total_pages"`
}
