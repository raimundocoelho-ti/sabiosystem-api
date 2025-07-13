/*
|------------------------------------------------
| File: internal/domain/category/graphql.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package category

import "github.com/graphql-go/graphql"

// --- TYPES ---

var categoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Category",
		Fields: graphql.Fields{
			"id":         &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
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

// --- QUERIES ---

func GetQueryFields(service Service) graphql.Fields {
	return graphql.Fields{
		// Query paginada original
		"categories": &graphql.Field{
			Type:        paginatedCategoriesType,
			Description: "Obtém uma lista paginada de categorias.",
			Args: graphql.FieldConfigArgument{
				"page": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				page, _ := p.Args["page"].(int)
				return service.GetAllCategories(page)
			},
		},
		// Query por ID original
		"category": &graphql.Field{
			Type:        categoryType,
			Description: "Obtém uma única categoria pelo seu ID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				return service.GetCategoryByID(uint(id))
			},
		},
		// Nova query de busca
		"searchCategories": &graphql.Field{
			Type:        graphql.NewList(categoryType),
			Description: "Busca categorias por nome.",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String), // Argumento obrigatório
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				return service.SearchCategories(name)
			},
		},
	}
}

// --- MUTATIONS ---

// GetMutationFields retorna os campos de mutation para o módulo de categoria.
func GetMutationFields(service Service) graphql.Fields {
	return graphql.Fields{
		"createCategory": &graphql.Field{
			Type:        categoryType,
			Description: "Cria uma nova categoria.",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				return service.CreateCategory(CreateCategoryDTO{Name: name})
			},
		},
		"updateCategory": &graphql.Field{
			Type:        categoryType,
			Description: "Atualiza uma categoria existente.",
			Args: graphql.FieldConfigArgument{
				"id":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				name, _ := p.Args["name"].(string)
				return service.UpdateCategory(uint(id), UpdateCategoryDTO{Name: name})
			},
		},
		"deleteCategory": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name:   "DeleteCategoryPayload",
				Fields: graphql.Fields{"deletedId": &graphql.Field{Type: graphql.Int}, "success": &graphql.Field{Type: graphql.Boolean}},
			}),
			Description: "Deleta uma categoria pelo seu ID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				err := service.DeleteCategory(uint(id))
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"deletedId": id, "success": true}, nil
			},
		},
	}
}
