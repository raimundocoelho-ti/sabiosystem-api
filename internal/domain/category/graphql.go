/*
|------------------------------------------------
| File: internal/domain/category/graphql.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package category

import "github.com/graphql-go/graphql"

var categoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Category",
		Fields: graphql.Fields{
			"id":         &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"agent_id":   &graphql.Field{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA: Expondo o agent_id
			"name":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"created_at": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"updated_at": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	},
)

var paginatedCategoriesType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PaginatedCategories",
		Fields: graphql.Fields{
			"data":        &graphql.Field{Type: graphql.NewList(categoryType)},
			"total":       &graphql.Field{Type: graphql.Int},
			"page":        &graphql.Field{Type: graphql.Int},
			"per_page":    &graphql.Field{Type: graphql.Int},
			"total_pages": &graphql.Field{Type: graphql.Int},
		},
	},
)

func GetQueryFields(service Service) graphql.Fields {
	return graphql.Fields{
		"categories": &graphql.Field{
			Type:        paginatedCategoriesType,
			Description: "Obtém categorias para um agente específico.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				page, _ := p.Args["page"].(int)
				return service.GetAllCategories(uint(agentId), page)
			},
		},
		"category": &graphql.Field{
			Type:        categoryType,
			Description: "Obtém uma categoria pelo seu ID, dentro de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				id, _ := p.Args["id"].(int)
				return service.GetCategoryByID(uint(agentId), uint(id))
			},
		},
		"searchCategories": &graphql.Field{
			Type:        graphql.NewList(categoryType),
			Description: "Busca categorias por nome, dentro de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"name":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				name, _ := p.Args["name"].(string)
				return service.SearchCategories(uint(agentId), name)
			},
		},
	}
}

func GetMutationFields(service Service) graphql.Fields {
	return graphql.Fields{
		"createCategory": &graphql.Field{
			Type:        categoryType,
			Description: "Cria uma nova categoria para um agente.",
			Args: graphql.FieldConfigArgument{
				// ↓↓ MUDANÇA PRINCIPAL AQUI ↓↓
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// O agentId é capturado...
				agentId, _ := p.Args["agentId"].(int)
				name, _ := p.Args["name"].(string)
				// ...e passado para o serviço através do DTO.
				return service.CreateCategory(CreateCategoryDTO{AgentID: uint(agentId), Name: name})
			},
		},
		"updateCategory": &graphql.Field{
			Type:        categoryType,
			Description: "Atualiza uma categoria de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				id, _ := p.Args["id"].(int)
				name, _ := p.Args["name"].(string)
				return service.UpdateCategory(uint(agentId), uint(id), UpdateCategoryDTO{Name: name})
			},
		},
		"deleteCategory": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name:   "DeleteCategoryPayload",
				Fields: graphql.Fields{"deletedId": &graphql.Field{Type: graphql.Int}, "success": &graphql.Field{Type: graphql.Boolean}},
			}),
			Description: "Deleta uma categoria de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)}, // <-- MUDANÇA
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				id, _ := p.Args["id"].(int)
				err := service.DeleteCategory(uint(agentId), uint(id))
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"deletedId": id, "success": true}, nil
			},
		},
	}
}
