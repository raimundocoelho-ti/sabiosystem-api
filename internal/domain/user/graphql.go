/*
|------------------------------------------------
| File: internal/domain/user/graphql.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package user

import "github.com/graphql-go/graphql"

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":         &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"agent_id":   &graphql.Field{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
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

func GetQueryFields(service Service) graphql.Fields {
	return graphql.Fields{
		"users": &graphql.Field{
			Type:        paginatedUsersType,
			Description: "Obtém usuários de um agente específico.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				page, _ := p.Args["page"].(int)
				return service.GetAllUsers(uint(agentId), page)
			},
		},
		"user": &graphql.Field{
			Type:        userType,
			Description: "Obtém um usuário pelo ID, dentro de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				id, _ := p.Args["id"].(int)
				return service.GetUserByID(uint(agentId), uint(id))
			},
		},
		"searchUsers": &graphql.Field{
			Type:        graphql.NewList(userType),
			Description: "Busca usuários por nome/email, dentro de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"name":    &graphql.ArgumentConfig{Type: graphql.String},
				"email":   &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				name, _ := p.Args["name"].(string)
				email, _ := p.Args["email"].(string)
				return service.SearchUsers(uint(agentId), name, email)
			},
		},
	}
}

func GetMutationFields(service Service) graphql.Fields {
	return graphql.Fields{
		"createUser": &graphql.Field{
			Type:        userType,
			Description: "Cria um novo usuário para um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"name":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				dto := CreateUserDTO{
					AgentID:  uint(p.Args["agentId"].(int)),
					Name:     p.Args["name"].(string),
					Email:    p.Args["email"].(string),
					Password: p.Args["password"].(string),
				}
				return service.CreateUser(dto)
			},
		},
		"updateUser": &graphql.Field{
			Type:        userType,
			Description: "Atualiza um usuário de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				id, _ := p.Args["id"].(int)
				dto := UpdateUserDTO{
					Name:  p.Args["name"].(string),
					Email: p.Args["email"].(string),
				}
				return service.UpdateUser(uint(agentId), uint(id), dto)
			},
		},
		"deleteUser": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name:   "DeleteUserPayload",
				Fields: graphql.Fields{"deletedId": &graphql.Field{Type: graphql.Int}, "success": &graphql.Field{Type: graphql.Boolean}},
			}),
			Description: "Deleta um usuário de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				id, _ := p.Args["id"].(int)
				err := service.DeleteUser(uint(agentId), uint(id))
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"deletedId": id, "success": true}, nil
			},
		},
	}
}
