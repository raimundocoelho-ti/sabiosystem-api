// Cole este código em internal/graphql/schema.go
package graphql

import (
	"github.com/graphql-go/graphql"
)

// NewSchema cria e retorna o schema GraphQL completo para a aplicação.
func NewSchema(resolver *Resolver) (graphql.Schema, error) {
	// RootQuery define os pontos de entrada para consultas (leitura de dados)
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"categories": &graphql.Field{
				Type:        PaginatedCategoriesType,
				Description: "Obtém uma lista paginada de categorias.",
				Args: graphql.FieldConfigArgument{
					"page": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 1,
					},
				},
				Resolve: resolver.resolveCategories,
			},
			"category": &graphql.Field{
				Type:        CategoryType,
				Description: "Obtém uma única categoria pelo seu ID.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolver.resolveCategory,
			},
		},
	})

	// RootMutation define os pontos de entrada para modificações (escrita de dados)
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createCategory": &graphql.Field{
				Type:        CategoryType,
				Description: "Cria uma nova categoria.",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolver.resolveCreateCategory,
			},
			"updateCategory": &graphql.Field{
				Type:        CategoryType,
				Description: "Atualiza uma categoria existente.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolver.resolveUpdateCategory,
			},
			"deleteCategory": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "DeleteCategoryPayload",
					Fields: graphql.Fields{
						"deletedId": &graphql.Field{Type: graphql.Int},
						"success":   &graphql.Field{Type: graphql.Boolean},
					},
				}),
				Description: "Deleta uma categoria pelo seu ID.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolver.resolveDeleteCategory,
			},
		},
	})

	// O Schema final que une Queries e Mutations
	schemaConfig := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	return graphql.NewSchema(schemaConfig)
}
