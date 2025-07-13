/*
|------------------------------------------------
| File: internal/domain/user/graphql.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package user

import "github.com/graphql-go/graphql"

// --- TYPES ---

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":         &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"name":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"email":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"created_at": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"updated_at": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	},
)

var paginatedUsersType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PaginatedUsers",
		Fields: graphql.Fields{
			"data":        &graphql.Field{Type: graphql.NewList(userType)},
			"total":       &graphql.Field{Type: graphql.Int},
			"page":        &graphql.Field{Type: graphql.Int},
			"per_page":    &graphql.Field{Type: graphql.Int},
			"total_pages": &graphql.Field{Type: graphql.Int},
		},
	},
)

// --- QUERIES ---

func GetQueryFields(service Service) graphql.Fields {
	// A versão corrigida retorna todos os campos de query para o usuário
	return graphql.Fields{
		// 1. A query paginada (que tinha sumido)
		"users": &graphql.Field{
			Type:        paginatedUsersType,
			Description: "Obtém uma lista paginada de usuários.",
			Args: graphql.FieldConfigArgument{
				"page": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				page, _ := p.Args["page"].(int)
				return service.GetAllUsers(page)
			},
		},
		// 2. A query por ID (que tinha sumido)
		"user": &graphql.Field{
			Type:        userType,
			Description: "Obtém um único usuário pelo seu ID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				return service.GetUserByID(uint(id))
			},
		},
		// 3. A nova query de busca (que já tínhamos adicionado)
		"searchUsers": &graphql.Field{
			Type:        graphql.NewList(userType), // Retorna uma lista de usuários
			Description: "Busca usuários por nome e/ou email.",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				email, _ := p.Args["email"].(string)
				return service.SearchUsers(name, email)
			},
		},
	}
}

// --- MUTATIONS ---

func GetMutationFields(service Service) graphql.Fields {
	return graphql.Fields{
		"createUser": &graphql.Field{
			Type:        userType,
			Description: "Cria um novo usuário.",
			Args: graphql.FieldConfigArgument{
				"name":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				dto := CreateUserDTO{
					Name:     p.Args["name"].(string),
					Email:    p.Args["email"].(string),
					Password: p.Args["password"].(string),
				}
				return service.CreateUser(dto)
			},
		},
		"updateUser": &graphql.Field{
			Type:        userType,
			Description: "Atualiza um usuário existente.",
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				dto := UpdateUserDTO{
					Name:  p.Args["name"].(string),
					Email: p.Args["email"].(string),
				}
				return service.UpdateUser(uint(id), dto)
			},
		},
		"deleteUser": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name:   "DeleteUserPayload",
				Fields: graphql.Fields{"deletedId": &graphql.Field{Type: graphql.Int}, "success": &graphql.Field{Type: graphql.Boolean}},
			}),
			Description: "Deleta um usuário pelo seu ID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				err := service.DeleteUser(uint(id))
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"deletedId": id, "success": true}, nil
			},
		},
	}
}
