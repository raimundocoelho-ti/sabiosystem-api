// Cole este código em internal/graphql/types.go
package graphql

import (
	"github.com/graphql-go/graphql"
)

// CategoryType é a representação GraphQL da nossa Categoria
var CategoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Category",
		Description: "Representa uma categoria de produtos ou serviços.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "O ID único da categoria.",
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "O nome da categoria.",
			},
			"created_at": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String), // Exposto como string para simplicidade
				Description: "Data de criação.",
			},
			"updated_at": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Data da última atualização.",
			},
		},
	},
)

// PaginatedCategoriesType é a representação para a resposta paginada
var PaginatedCategoriesType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PaginatedCategories",
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type: graphql.NewList(CategoryType),
			},
			"total": &graphql.Field{
				Type: graphql.Int,
			},
			"page": &graphql.Field{
				Type: graphql.Int,
			},
			"per_page": &graphql.Field{
				Type: graphql.Int,
			},
			"total_pages": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
