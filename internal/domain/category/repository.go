/*
|------------------------------------------------
| File: internal/domain/category/repository.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package category

import "gorm.io/gorm"

type Repository interface {
	FindAll(agentID uint, page, perPage int) ([]Category, int64, error)
	FindByID(agentID, id uint) (Category, error)
	Search(agentID uint, name string) ([]Category, error)
	Create(category Category) (Category, error)
	Update(category Category) (Category, error)
	Delete(agentID, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(agentID uint, page, perPage int) ([]Category, int64, error) {
	var categories []Category
	var total int64
	// <-- MUDANÇA: Adicionado Where("agent_id = ?", agentID)
	if err := r.db.Model(&Category{}).Where("agent_id = ?", agentID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	// <-- MUDANÇA: Adicionado Where("agent_id = ?", agentID)
	err := r.db.Where("agent_id = ?", agentID).Limit(perPage).Offset(offset).Order("id asc").Find(&categories).Error
	return categories, total, err
}

func (r *repository) FindByID(agentID, id uint) (Category, error) {
	var category Category
	// <-- MUDANÇA: Adicionado Where("agent_id = ?", agentID) para segurança
	err := r.db.Where("agent_id = ?", agentID).First(&category, id).Error
	return category, err
}

func (r *repository) Search(agentID uint, name string) ([]Category, error) {
	var categories []Category
	// <-- MUDANÇA: Adicionado Where("agent_id = ?", agentID)
	err := r.db.Where("agent_id = ? AND name ILIKE ?", agentID, "%"+name+"%").Find(&categories).Error
	return categories, err
}

func (r *repository) Create(category Category) (Category, error) {
	err := r.db.Create(&category).Error
	return category, err
}

func (r *repository) Update(category Category) (Category, error) {
	err := r.db.Save(&category).Error
	return category, err
}

func (r *repository) Delete(agentID, id uint) error {
	// <-- MUDANÇA: Adicionado Where("agent_id = ?", agentID) para segurança
	return r.db.Where("agent_id = ?", agentID).Delete(&Category{}, id).Error
}
