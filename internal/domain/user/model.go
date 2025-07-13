/*
|------------------------------------------------
| File: internal/domain/user/model.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User representa a entidade no banco de dados.
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // O 'json:"-"' impede que a senha seja exposta em respostas JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Antes de salvar, cria um hash da senha
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return
}

// CreateUserDTO é o DTO para a criação de um usuário.
type CreateUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserDTO é o DTO para a atualização de um usuário.
type UpdateUserDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// PaginatedUsers é a estrutura de resposta para a lista paginada de usuários.
type PaginatedUsers struct {
	Data       []User `json:"data"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	TotalPages int    `json:"total_pages"`
}
