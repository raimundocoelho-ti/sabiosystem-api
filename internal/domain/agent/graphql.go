/*
|------------------------------------------------
| File: internal/domain/agent/graphql.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package agent

import "github.com/graphql-go/graphql"

var agentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Agent",
		Fields: graphql.Fields{
			"id":         &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"name":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"domain":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"created_at": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"updated_at": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	},
)

var paginatedAgentsType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PaginatedAgents",
		Fields: graphql.Fields{
			"data":        &graphql.Field{Type: graphql.NewList(agentType)},
			"total":       &graphql.Field{Type: graphql.Int},
			"page":        &graphql.Field{Type: graphql.Int},
			"per_page":    &graphql.Field{Type: graphql.Int},
			"total_pages": &graphql.Field{Type: graphql.Int},
		},
	},
)

func GetQueryFields(service Service) graphql.Fields {
	return graphql.Fields{
		"agents": &graphql.Field{
			Type:        paginatedAgentsType,
			Description: "Obtém uma lista paginada de agentes.",
			Args: graphql.FieldConfigArgument{
				"page": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				page, _ := p.Args["page"].(int)
				return service.GetAllAgents(page)
			},
		},
		"agent": &graphql.Field{
			Type:        agentType,
			Description: "Obtém um único agente pelo seu ID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				return service.GetAgentByID(uint(id))
			},
		},
		"searchAgents": &graphql.Field{
			Type:        graphql.NewList(agentType),
			Description: "Busca agentes por nome e/ou domínio.",
			Args: graphql.FieldConfigArgument{
				"name":   &graphql.ArgumentConfig{Type: graphql.String},
				"domain": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				domain, _ := p.Args["domain"].(string)
				return service.SearchAgents(name, domain)
			},
		},
	}
}

func GetMutationFields(service Service) graphql.Fields {
	return graphql.Fields{
		"createAgent": &graphql.Field{
			Type:        agentType,
			Description: "Cria um novo agente.",
			Args: graphql.FieldConfigArgument{
				"name":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"domain": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				dto := CreateAgentDTO{Name: p.Args["name"].(string), Domain: p.Args["domain"].(string)}
				return service.CreateAgent(dto)
			},
		},
		"updateAgent": &graphql.Field{
			Type:        agentType,
			Description: "Atualiza um agente existente.",
			Args: graphql.FieldConfigArgument{
				"id":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"domain": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				dto := UpdateAgentDTO{Name: p.Args["name"].(string), Domain: p.Args["domain"].(string)}
				return service.UpdateAgent(uint(id), dto)
			},
		},
		"deleteAgent": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name:   "DeleteAgentPayload",
				Fields: graphql.Fields{"deletedId": &graphql.Field{Type: graphql.Int}, "success": &graphql.Field{Type: graphql.Boolean}},
			}),
			Description: "Deleta um agente pelo seu ID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				err := service.DeleteAgent(uint(id))
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"deletedId": id, "success": true}, nil
			},
		},
	}
}
