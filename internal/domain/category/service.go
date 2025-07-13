// Cole este código completo em internal/domain/category/service.go

package category

import (
	"math"
)

const PageSize = 8 // Itens por página

// Service define a interface para a lógica de negócio de Categoria.
type Service interface {
	GetAllCategories(page int) (PaginatedCategories, error)
	GetCategoryByID(id uint) (Category, error)
	SearchCategories(name string) ([]Category, error)
	CreateCategory(dto CreateCategoryDTO) (Category, error)
	UpdateCategory(id uint, dto UpdateCategoryDTO) (Category, error)
	DeleteCategory(id uint) error
}

type service struct {
	repo Repository
}

// NewService cria uma nova instância do serviço.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllCategories(page int) (PaginatedCategories, error) {
	categories, total, err := s.repo.FindAll(page, PageSize)
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

func (s *service) GetCategoryByID(id uint) (Category, error) {
	return s.repo.FindByID(id)
}

func (s *service) SearchCategories(name string) ([]Category, error) {
	return s.repo.Search(name)
}

func (s *service) CreateCategory(dto CreateCategoryDTO) (Category, error) {
	category := Category{Name: dto.Name}
	return s.repo.Create(category)
}

// UpdateCategory contém a lógica de atualização corrigida.
func (s *service) UpdateCategory(id uint, dto UpdateCategoryDTO) (Category, error) {
	// 1. Busca a categoria para garantir que ela existe.
	categoryToUpdate, err := s.repo.FindByID(id)
	if err != nil {
		// 2. Se não encontrou (erro != nil), retorna o erro imediatamente.
		//    Isso impede que o código continue.
		return Category{}, err
	}

	// 3. Se encontrou, atualiza o nome e salva.
	categoryToUpdate.Name = dto.Name
	return s.repo.Update(categoryToUpdate)
}

func (s *service) DeleteCategory(id uint) error {
	// Opcional: verificar se a categoria existe antes de deletar
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
