/*
|------------------------------------------------
| File: internal/auth/service.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5" // <-- MUDANÇA: Import adicionado
)

type Service interface {
	StoreRefreshToken(tokenString string, userID uint) (*time.Time, error)
	ValidateRefreshToken(tokenString string) (*RefreshToken, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// hashToken cria um hash SHA-256 de uma string. É determinístico.
func hashToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}

// StoreRefreshToken cria o hash de um refresh token e o salva no banco.
func (s *service) StoreRefreshToken(tokenString string, userID uint) (*time.Time, error) {
	tokenHash := hashToken(tokenString)

	// A duração do token de atualização é de 7 dias (168 horas).
	expiresAt := time.Now().Add(168 * time.Hour)

	token := RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	err := s.repo.Store(token)
	if err != nil {
		return nil, err
	}
	return &expiresAt, nil
}

// ValidateRefreshToken verifica se um refresh token é válido.
func (s *service) ValidateRefreshToken(tokenString string) (*RefreshToken, error) {
	tokenHash := hashToken(tokenString)
	storedToken, err := s.repo.FindByTokenHash(tokenHash)
	if err != nil {
		return nil, err // Token não encontrado
	}

	if time.Now().After(storedToken.ExpiresAt) {
		// Usamos a variável de erro do pacote jwt importado
		return nil, jwt.ErrTokenExpired
	}

	return storedToken, nil
}
