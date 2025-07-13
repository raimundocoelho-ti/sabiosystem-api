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
	FindAll(page, perPage int) ([]User, int64, error)
	FindByID(id uint) (User, error)
	Search(name, email string) ([]User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(page, perPage int) ([]User, int64, error) {
	var users []User
	var total int64
	if err := r.db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	err := r.db.Limit(perPage).Offset(offset).Order("id asc").Find(&users).Error
	return users, total, err
}

func (r *repository) FindByID(id uint) (User, error) {
	var user User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *repository) Search(name, email string) ([]User, error) {
	var users []User
	query := r.db

	if name != "" {
		// ILIKE faz a busca parcial e case-insensitive (só funciona bem no PostgreSQL)
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

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
