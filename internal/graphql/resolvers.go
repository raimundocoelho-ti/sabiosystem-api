// Cole este código em internal/graphql/resolvers.go
package graphql

import (
	"github.com/graphql-go/graphql"
	// ↓↓ ATUALIZE SEU CAMINHO AQUI SE FOR DIFERENTE ↓↓
	"github.com/raimundocoelho-ti/sabiosystem-api/internal/domain/category"
)

// Resolver encapsula os serviços para injeção de dependência.
type Resolver struct {
	CategoryService category.Service
}

// resolveCategories busca todas as categorias com paginação.
func (r *Resolver) resolveCategories(p graphql.ResolveParams) (interface{}, error) {
	page, ok := p.Args["page"].(int)
	if !ok || page < 1 {
		page = 1
	}
	return r.CategoryService.GetAllCategories(page)
}

// resolveCategory busca uma única categoria pelo ID.
func (r *Resolver) resolveCategory(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, nil
	}
	return r.CategoryService.GetCategoryByID(uint(id))
}

// resolveCreateCategory cria uma nova categoria.
func (r *Resolver) resolveCreateCategory(p graphql.ResolveParams) (interface{}, error) {
	name, ok := p.Args["name"].(string)
	if !ok || name == "" {
		return nil, nil
	}
	dto := category.CreateCategoryDTO{Name: name}
	return r.CategoryService.CreateCategory(dto)
}

// resolveUpdateCategory atualiza uma categoria existente.
func (r *Resolver) resolveUpdateCategory(p graphql.ResolveParams) (interface{}, error) {
	id, idOk := p.Args["id"].(int)
	name, nameOk := p.Args["name"].(string)
	if !idOk || !nameOk {
		return nil, nil
	}
	dto := category.UpdateCategoryDTO{Name: name}
	return r.CategoryService.UpdateCategory(uint(id), dto)
}

// resolveDeleteCategory deleta uma categoria.
func (r *Resolver) resolveDeleteCategory(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, nil
	}
	err := r.CategoryService.DeleteCategory(uint(id))
	// Em GraphQL, é comum retornar o ID do item deletado ou um status
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"deletedId": id,
		"success":   true,
	}, nil
}
