/*
|------------------------------------------------
| File: internal/domain/agent/model.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package agent

import "time"

// Agent representa a entidade no banco de dados.
type Agent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	Domain    string    `gorm:"unique;not null" json:"domain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateAgentDTO é o Data Transfer Object para a criação de um agent.
type CreateAgentDTO struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

// UpdateAgentDTO é o Data Transfer Object para a atualização de um agent.
type UpdateAgentDTO struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

// PaginatedAgents é a estrutura de resposta para a lista paginada de agents.
type PaginatedAgents struct {
	Data       []Agent `json:"data"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PerPage    int     `json:"per_page"`
	TotalPages int     `json:"total_pages"`
}
