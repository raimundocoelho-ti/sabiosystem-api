/*
|------------------------------------------------
| File: internal/domain/user/repository.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package user

import "gorm.io/gorm"

type Repository interface {
	FindAll(agentID uint, page, perPage int) ([]User, int64, error)
	FindByID(agentID, id uint) (User, error)
	Search(agentID uint, name, email string) ([]User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(agentID, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(agentID uint, page, perPage int) ([]User, int64, error) {
	var users []User
	var total int64
	// <-- MUDANÇA: Adicionado Where("agent_id = ?", agentID)
	if err := r.db.Model(&User{}).Where("agent_id = ?", agentID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	// <-- MUDANÇA: Adicionado Where("agent_id = ?", agentID)
	err := r.db.Where("agent_id = ?", agentID).Limit(perPage).Offset(offset).Order("id asc").Find(&users).Error
	return users, total, err
}

func (r *repository) FindByID(agentID, id uint) (User, error) {
	var user User
	// <-- MUDANÇA: Adicionado Where("agent_id = ?", agentID)
	err := r.db.Where("agent_id = ?", agentID).First(&user, id).Error
	return user, err
}

func (r *repository) Search(agentID uint, name, email string) ([]User, error) {
	var users []User
	// <-- MUDANÇA: Adicionado Where("agent_id = ?", agentID)
	query := r.db.Where("agent_id = ?", agentID)

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}

	err := query.Find(&users).Error
	return users, err
}

func (r *repository) Create(user User) (User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *repository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error
	return user, err
}

func (r *repository) Delete(agentID, id uint) error {
	// <-- MUDANÇA: Adicionado Where("agent_id = ?", agentID)
	return r.db.Where("agent_id = ?", agentID).Delete(&User{}, id).Error
}
