/*
|------------------------------------------------
| File: internal/domain/product/repository.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package product

import "gorm.io/gorm"

type Repository interface {
	FindAll(agentID uint, page, perPage int) ([]Product, int64, error)
	FindByID(agentID, id uint) (Product, error)
	Search(agentID uint, name string) ([]Product, error)
	SearchByCategory(agentID, categoryID uint) ([]Product, error) // NOVO
	Create(product Product) (Product, error)
	Update(product Product) (Product, error)
	Delete(agentID, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(agentID uint, page, perPage int) ([]Product, int64, error) {
	var products []Product
	var total int64
	if err := r.db.Model(&Product{}).Where("agent_id = ?", agentID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	err := r.db.Where("agent_id = ?", agentID).Limit(perPage).Offset(offset).Order("id asc").Find(&products).Error
	return products, total, err
}

func (r *repository) FindByID(agentID, id uint) (Product, error) {
	var product Product
	err := r.db.Where("agent_id = ?", agentID).First(&product, id).Error
	return product, err
}

func (r *repository) Search(agentID uint, name string) ([]Product, error) {
	var products []Product
	err := r.db.Where("agent_id = ? AND name ILIKE ?", agentID, "%"+name+"%").Find(&products).Error
	return products, err
}

func (r *repository) SearchByCategory(agentID, categoryID uint) ([]Product, error) {
	var products []Product
	err := r.db.Where("agent_id = ? AND category_id = ?", agentID, categoryID).Find(&products).Error
	return products, err
}

func (r *repository) Create(product Product) (Product, error) {
	err := r.db.Create(&product).Error
	return product, err
}

func (r *repository) Update(product Product) (Product, error) {
	err := r.db.Save(&product).Error
	return product, err
}

func (r *repository) Delete(agentID, id uint) error {
	return r.db.Where("agent_id = ?", agentID).Delete(&Product{}, id).Error
}
