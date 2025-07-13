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
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/auth"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/agent"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/category"
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/product" // <-- ADICIONADO
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/user"
)

// SchemaServices contém todos os serviços necessários para construir o schema.
type SchemaServices struct {
	CategorySvc category.Service
	ProductSvc  product.Service // <-- ADICIONADO
	UserSvc     user.Service
	AgentSvc    agent.Service
	AuthSvc     auth.Service
}

func NewSchema(services SchemaServices) (graphql.Schema, error) {
	// Juntando os campos de Query de todos os módulos (auth não tem queries)
	queryFields := mergeFields(
		category.GetQueryFields(services.CategorySvc),
		product.GetQueryFields(services.ProductSvc), // <-- ADICIONADO
		user.GetQueryFields(services.UserSvc),
		agent.GetQueryFields(services.AgentSvc),
	)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: queryFields,
	})

	// Juntando os campos de Mutation de todos os módulos
	mutationFields := mergeFields(
		category.GetMutationFields(services.CategorySvc),
		product.GetMutationFields(services.ProductSvc), // <-- ADICIONADO
		user.GetMutationFields(services.UserSvc),
		agent.GetMutationFields(services.AgentSvc),
		auth.GetMutationFields(services.AuthSvc, services.UserSvc, services.AgentSvc),
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
