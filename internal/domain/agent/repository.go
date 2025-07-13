/*
|------------------------------------------------
| File: internal/domain/agent/repository.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package agent

import "gorm.io/gorm"

// Repository define a interface para as operações de banco de dados.
type Repository interface {
	FindAll(page, perPage int) ([]Agent, int64, error)
	FindByID(id uint) (Agent, error)
	Search(name, domain string) ([]Agent, error)
	Create(agent Agent) (Agent, error)
	Update(agent Agent) (Agent, error)
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(page, perPage int) ([]Agent, int64, error) {
	var agents []Agent
	var total int64
	if err := r.db.Model(&Agent{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	err := r.db.Limit(perPage).Offset(offset).Order("id asc").Find(&agents).Error
	return agents, total, err
}

func (r *repository) FindByID(id uint) (Agent, error) {
	var agent Agent
	err := r.db.First(&agent, id).Error
	return agent, err
}

func (r *repository) Search(name, domain string) ([]Agent, error) {
	var agents []Agent
	query := r.db
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if domain != "" {
		query = query.Where("domain ILIKE ?", "%"+domain+"%")
	}
	err := query.Find(&agents).Error
	return agents, err
}

func (r *repository) Create(agent Agent) (Agent, error) {
	err := r.db.Create(&agent).Error
	return agent, err
}

func (r *repository) Update(agent Agent) (Agent, error) {
	err := r.db.Save(&agent).Error
	return agent, err
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Agent{}, id).Error
}
