/*
|------------------------------------------------
| File: internal/auth/graphql.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package auth

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/agent"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// AuthPayload é o tipo de resposta para as mutations de login e refreshToken.
var AuthPayload = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AuthPayload",
		Fields: graphql.Fields{
			"access_token":  &graphql.Field{Type: graphql.String},
			"refresh_token": &graphql.Field{Type: graphql.String},
			"user":          &graphql.Field{Type: user.UserType}, // Usando o UserType que já definimos
		},
	},
)

// GetMutationFields retorna as mutations relacionadas à autenticação.
func GetMutationFields(authSvc Service, userSvc user.Service, agentSvc agent.Service) graphql.Fields {
	return graphql.Fields{
		"login": &graphql.Field{
			Type:        AuthPayload,
			Description: "Autentica um usuário e retorna um par de tokens.",
			Args: graphql.FieldConfigArgument{
				"agentDomain": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"password":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentDomain, _ := p.Args["agentDomain"].(string)
				email, _ := p.Args["email"].(string)
				password, _ := p.Args["password"].(string)

				// 1. Encontrar o Agente pelo domínio
				agents, err := agentSvc.SearchAgents("", agentDomain)
				if err != nil || len(agents) != 1 {
					return nil, errors.New("agente inválido ou não encontrado")
				}
				targetAgent := agents[0]

				// 2. Encontrar o Usuário pelo email DENTRO daquele agente
				users, err := userSvc.SearchUsers(targetAgent.ID, "", email)
				if err != nil || len(users) != 1 {
					return nil, errors.New("email ou senha inválidos")
				}
				targetUser := users[0]

				// 3. Verificar a senha
				err = bcrypt.CompareHashAndPassword([]byte(targetUser.Password), []byte(password))
				if err != nil {
					return nil, errors.New("email ou senha inválidos")
				}

				// 4. Gerar os tokens
				accessToken, err := GenerateAccessToken(targetUser)
				if err != nil {
					return nil, errors.New("falha ao gerar token de acesso")
				}
				refreshToken, err := GenerateRefreshToken(targetUser)
				if err != nil {
					return nil, errors.New("falha ao gerar refresh token")
				}

				// 5. Salvar o refresh token no banco
				_, err = authSvc.StoreRefreshToken(refreshToken, targetUser.ID)
				if err != nil {
					return nil, errors.New("falha ao salvar sessão")
				}

				// 6. Retornar o payload
				return map[string]interface{}{
					"access_token":  accessToken,
					"refresh_token": refreshToken,
					"user":          targetUser,
				}, nil
			},
		},
		"refreshToken": &graphql.Field{
			Type:        AuthPayload,
			Description: "Gera um novo access token usando um refresh token válido.",
			Args: graphql.FieldConfigArgument{
				"refreshToken": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tokenString, _ := p.Args["refreshToken"].(string)

				// 1. Validar o refresh token
				storedToken, err := authSvc.ValidateRefreshToken(tokenString)
				if err != nil {
					return nil, errors.New("refresh token inválido ou expirado")
				}

				// 2. Encontrar o usuário associado ao token
				// Precisamos passar o agentId do token para o serviço de usuário
				user, err := userSvc.GetUserByID(storedToken.UserID, storedToken.UserID)
				if err != nil {
					return nil, errors.New("usuário do token não encontrado")
				}

				// 3. Gerar um novo access token
				newAccessToken, err := GenerateAccessToken(user)
				if err != nil {
					return nil, errors.New("falha ao gerar novo token de acesso")
				}

				// Nota: Para segurança máxima, poderíamos gerar um novo refresh token também
				// e invalidar o antigo (rotação de tokens), mas manteremos simples por agora.

				return map[string]interface{}{
					"access_token":  newAccessToken,
					"refresh_token": nil, // Não retornamos um novo refresh token aqui
					"user":          nil, // Não retornamos o usuário aqui
				}, nil
			},
		},
	}
}
