/*
|------------------------------------------------
| File: internal/auth/model.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package auth

import "time"

// RefreshToken representa a entidade no banco de dados.
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	TokenHash string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
