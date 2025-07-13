/*
|------------------------------------------------
| File: internal/domain/product/graphql.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package product

import "github.com/graphql-go/graphql"

// Type principal do produto
var productType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"agent_id":    &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"category_id": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"name":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"description": &graphql.Field{Type: graphql.String},
			"price":       &graphql.Field{Type: graphql.NewNonNull(graphql.Float)},
			"image_url":   &graphql.Field{Type: graphql.String},
			"is_active":   &graphql.Field{Type: graphql.Boolean},
			"created_at":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"updated_at":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	},
)

var paginatedProductsType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PaginatedProducts",
		Fields: graphql.Fields{
			"data":        &graphql.Field{Type: graphql.NewList(productType)},
			"total":       &graphql.Field{Type: graphql.Int},
			"page":        &graphql.Field{Type: graphql.Int},
			"per_page":    &graphql.Field{Type: graphql.Int},
			"total_pages": &graphql.Field{Type: graphql.Int},
		},
	},
)

func GetQueryFields(service Service) graphql.Fields {
	return graphql.Fields{
		"products": &graphql.Field{
			Type:        paginatedProductsType,
			Description: "Obtém produtos para um agente específico.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"page":    &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				page, _ := p.Args["page"].(int)
				return service.GetAllProducts(uint(agentId), page)
			},
		},
		"product": &graphql.Field{
			Type:        productType,
			Description: "Obtém um produto pelo seu ID, dentro de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				id, _ := p.Args["id"].(int)
				return service.GetProductByID(uint(agentId), uint(id))
			},
		},
		"searchProducts": &graphql.Field{
			Type:        graphql.NewList(productType),
			Description: "Busca produtos por nome, dentro de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				name, _ := p.Args["name"].(string)
				return service.SearchProducts(uint(agentId), name)
			},
		},
		"productsByCategory": &graphql.Field{
			Type:        graphql.NewList(productType),
			Description: "Busca produtos por category_id e agent_id.",
			Args: graphql.FieldConfigArgument{
				"agentId":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"categoryId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId := uint(p.Args["agentId"].(int))
				categoryId := uint(p.Args["categoryId"].(int))
				return service.SearchByCategory(agentId, categoryId)
			},
		},
	}
}

func GetMutationFields(service Service) graphql.Fields {
	return graphql.Fields{
		"createProduct": &graphql.Field{
			Type:        productType,
			Description: "Cria um novo produto para um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"categoryId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"description": &graphql.ArgumentConfig{Type: graphql.String},
				"price":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				"imageUrl":    &graphql.ArgumentConfig{Type: graphql.String},
				"isActive":    &graphql.ArgumentConfig{Type: graphql.Boolean, DefaultValue: true},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				dto := CreateProductDTO{
					AgentID:     uint(p.Args["agentId"].(int)),
					CategoryID:  uint(p.Args["categoryId"].(int)),
					Name:        p.Args["name"].(string),
					Description: "",
					Price:       p.Args["price"].(float64),
					ImageURL:    "",
					IsActive:    true,
				}
				if v, ok := p.Args["description"]; ok && v != nil {
					dto.Description = v.(string)
				}
				if v, ok := p.Args["imageUrl"]; ok && v != nil {
					dto.ImageURL = v.(string)
				}
				if v, ok := p.Args["isActive"]; ok && v != nil {
					dto.IsActive = v.(bool)
				}
				return service.CreateProduct(dto)
			},
		},
		"updateProduct": &graphql.Field{
			Type:        productType,
			Description: "Atualiza um produto de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"id":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"categoryId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":        &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"description": &graphql.ArgumentConfig{Type: graphql.String},
				"price":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				"imageUrl":    &graphql.ArgumentConfig{Type: graphql.String},
				"isActive":    &graphql.ArgumentConfig{Type: graphql.Boolean},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				dto := UpdateProductDTO{
					CategoryID:  uint(p.Args["categoryId"].(int)),
					Name:        p.Args["name"].(string),
					Description: "",
					Price:       p.Args["price"].(float64),
					ImageURL:    "",
					IsActive:    true,
				}
				if v, ok := p.Args["description"]; ok && v != nil {
					dto.Description = v.(string)
				}
				if v, ok := p.Args["imageUrl"]; ok && v != nil {
					dto.ImageURL = v.(string)
				}
				if v, ok := p.Args["isActive"]; ok && v != nil {
					dto.IsActive = v.(bool)
				}
				agentId := uint(p.Args["agentId"].(int))
				id := uint(p.Args["id"].(int))
				return service.UpdateProduct(agentId, id, dto)
			},
		},
		"deleteProduct": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "DeleteProductPayload",
				Fields: graphql.Fields{
					"deletedId": &graphql.Field{Type: graphql.Int},
					"success":   &graphql.Field{Type: graphql.Boolean},
				},
			}),
			Description: "Deleta um produto de um agente.",
			Args: graphql.FieldConfigArgument{
				"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				agentId, _ := p.Args["agentId"].(int)
				id, _ := p.Args["id"].(int)
				err := service.DeleteProduct(uint(agentId), uint(id))
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"deletedId": id, "success": true}, nil
			},
		},
	}
}
