/*
|------------------------------------------------
| File: internal/domain/category/service.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package category

import "math"

const PageSize = 8

type Service interface {
	GetAllCategories(agentID uint, page int) (PaginatedCategories, error)
	GetCategoryByID(agentID, id uint) (Category, error)
	SearchCategories(agentID uint, name string) ([]Category, error)
	CreateCategory(dto CreateCategoryDTO) (Category, error)
	UpdateCategory(agentID, id uint, dto UpdateCategoryDTO) (Category, error)
	DeleteCategory(agentID, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllCategories(agentID uint, page int) (PaginatedCategories, error) {
	categories, total, err := s.repo.FindAll(agentID, page, PageSize)
	if err != nil {
		return PaginatedCategories{}, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(PageSize)))
	return PaginatedCategories{
		Data:       categories,
		Total:      total,
		Page:       page,
		PerPage:    PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *service) GetCategoryByID(agentID, id uint) (Category, error) {
	return s.repo.FindByID(agentID, id)
}

func (s *service) SearchCategories(agentID uint, name string) ([]Category, error) {
	return s.repo.Search(agentID, name)
}

func (s *service) CreateCategory(dto CreateCategoryDTO) (Category, error) {
	category := Category{
		AgentID: dto.AgentID,
		Name:    dto.Name,
	}
	return s.repo.Create(category)
}

func (s *service) UpdateCategory(agentID, id uint, dto UpdateCategoryDTO) (Category, error) {
	categoryToUpdate, err := s.repo.FindByID(agentID, id)
	if err != nil {
		return Category{}, err
	}
	categoryToUpdate.Name = dto.Name
	return s.repo.Update(categoryToUpdate)
}

func (s *service) DeleteCategory(agentID, id uint) error {
	// A verificação de existência já acontece no repositório, mas podemos manter aqui também.
	_, err := s.repo.FindByID(agentID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(agentID, id)
}
