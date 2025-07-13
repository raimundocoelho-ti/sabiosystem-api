/*
|------------------------------------------------
| File: internal/domain/product/service.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package product

import "math"

const PageSize = 8

type Service interface {
	GetAllProducts(agentID uint, page int) (PaginatedProducts, error)
	GetProductByID(agentID, id uint) (Product, error)
	SearchProducts(agentID uint, name string) ([]Product, error)
	SearchByCategory(agentID, categoryID uint) ([]Product, error) // NOVO
	CreateProduct(dto CreateProductDTO) (Product, error)
	UpdateProduct(agentID, id uint, dto UpdateProductDTO) (Product, error)
	DeleteProduct(agentID, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllProducts(agentID uint, page int) (PaginatedProducts, error) {
	products, total, err := s.repo.FindAll(agentID, page, PageSize)
	if err != nil {
		return PaginatedProducts{}, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(PageSize)))
	return PaginatedProducts{
		Data:       products,
		Total:      total,
		Page:       page,
		PerPage:    PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *service) GetProductByID(agentID, id uint) (Product, error) {
	return s.repo.FindByID(agentID, id)
}

func (s *service) SearchProducts(agentID uint, name string) ([]Product, error) {
	return s.repo.Search(agentID, name)
}

func (s *service) SearchByCategory(agentID, categoryID uint) ([]Product, error) {
	return s.repo.SearchByCategory(agentID, categoryID)
}

func (s *service) CreateProduct(dto CreateProductDTO) (Product, error) {
	product := Product{
		AgentID:     dto.AgentID,
		CategoryID:  dto.CategoryID,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		ImageURL:    dto.ImageURL,
		IsActive:    dto.IsActive,
	}
	return s.repo.Create(product)
}

func (s *service) UpdateProduct(agentID, id uint, dto UpdateProductDTO) (Product, error) {
	productToUpdate, err := s.repo.FindByID(agentID, id)
	if err != nil {
		return Product{}, err
	}
	productToUpdate.CategoryID = dto.CategoryID
	productToUpdate.Name = dto.Name
	productToUpdate.Description = dto.Description
	productToUpdate.Price = dto.Price
	productToUpdate.ImageURL = dto.ImageURL
	productToUpdate.IsActive = dto.IsActive
	return s.repo.Update(productToUpdate)
}

func (s *service) DeleteProduct(agentID, id uint) error {
	_, err := s.repo.FindByID(agentID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(agentID, id)
}
