/*
|------------------------------------------------
| File: internal/auth/jwt.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/user"
)

// JWTClaims são os dados que vamos armazenar dentro do nosso token de acesso.
type JWTClaims struct {
	UserID  uint   `json:"user_id"`
	AgentID uint   `json:"agent_id"`
	Name    string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateAccessToken cria um novo token de acesso JWT de curta duração.
func GenerateAccessToken(appUser user.User) (string, error) {
	// O token de acesso expira em 1 hora.
	// Em produção, um tempo menor como 15 minutos é mais seguro.
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &JWTClaims{
		UserID:  appUser.ID,
		AgentID: appUser.AgentID,
		Name:    appUser.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sabiosystem-api",
		},
	}

	// Lê o segredo do .env
	jwtSecret := os.Getenv("JWT_ACCESS_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}

// GenerateRefreshToken cria um novo token de atualização de longa duração.
func GenerateRefreshToken(appUser user.User) (string, error) {
	// O token de atualização expira em 7 dias (168 horas).
	expirationTime := time.Now().Add(168 * time.Hour)

	claims := &JWTClaims{
		UserID:  appUser.ID,
		AgentID: appUser.AgentID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sabiosystem-api",
		},
	}

	// Usa um segredo DIFERENTE para o refresh token, por segurança.
	jwtSecret := os.Getenv("JWT_REFRESH_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
