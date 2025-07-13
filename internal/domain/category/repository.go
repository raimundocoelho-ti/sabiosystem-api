// internal/domain/category/repository.go
package category

import "gorm.io/gorm"

// Repository define a interface para as operações de banco de dados para Categoria.
type Repository interface {
	FindAll(page, perPage int) ([]Category, int64, error)
	FindByID(id uint) (Category, error)
	Create(category Category) (Category, error)
	Update(category Category) (Category, error)
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository cria uma nova instância do repositório.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(page, perPage int) ([]Category, int64, error) {
	var categories []Category
	var total int64

	// Contar o total de registros
	if err := r.db.Model(&Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calcular o offset
	offset := (page - 1) * perPage

	// Buscar os registros paginados
	err := r.db.Limit(perPage).Offset(offset).Order("id asc").Find(&categories).Error
	return categories, total, err
}

func (r *repository) FindByID(id uint) (Category, error) {
	var category Category
	// Esta é a linha 43 (ou próxima a ela)
	err := r.db.First(&category, id).Error
	return category, err
}

func (r *repository) Create(category Category) (Category, error) {
	err := r.db.Create(&category).Error
	return category, err
}

func (r *repository) Update(category Category) (Category, error) {
	err := r.db.Save(&category).Error
	return category, err
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Category{}, id).Error
}
