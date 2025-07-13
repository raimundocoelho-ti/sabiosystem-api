/*
|------------------------------------------------
| File: internal/auth/repository.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package auth

import "gorm.io/gorm"

type Repository interface {
	Store(token RefreshToken) error
	FindByTokenHash(hash string) (*RefreshToken, error)
	DeleteByUserID(userID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Store salva um novo refresh token no banco.
func (r *repository) Store(token RefreshToken) error {
	return r.db.Create(&token).Error
}

// FindByTokenHash busca um token pelo seu hash.
func (r *repository) FindByTokenHash(hash string) (*RefreshToken, error) {
	var token RefreshToken
	err := r.db.Where("token_hash = ?", hash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// DeleteByUserID deleta todos os refresh tokens de um usuário. Útil para "deslogar de todos os dispositivos".
func (r *repository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}
