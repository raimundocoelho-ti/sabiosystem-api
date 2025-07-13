/*
|------------------------------------------------
| File: internal/graphql/schema.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/category"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/user"
)

// SchemaServices contém todos os serviços necessários para construir o schema.
type SchemaServices struct {
	CategorySvc category.Service
	UserSvc     user.Service
}

// NewSchema cria e retorna o schema GraphQL completo, montado a partir dos módulos.
func NewSchema(services SchemaServices) (graphql.Schema, error) {
	// Juntando os campos de Query de todos os módulos
	queryFields := mergeFields(
		category.GetQueryFields(services.CategorySvc),
		user.GetQueryFields(services.UserSvc),
	)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: queryFields,
	})

	// Juntando os campos de Mutation de todos os módulos
	mutationFields := mergeFields(
		category.GetMutationFields(services.CategorySvc),
		user.GetMutationFields(services.UserSvc),
	)

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: mutationFields,
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}

	return graphql.NewSchema(schemaConfig)
}

// mergeFields é uma função utilitária para juntar múltiplos mapas de campos.
func mergeFields(fieldMaps ...graphql.Fields) graphql.Fields {
	result := graphql.Fields{}
	for _, fieldMap := range fieldMaps {
		for key, value := range fieldMap {
			result[key] = value
		}
	}
	return result
}
