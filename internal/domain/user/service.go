/*
|------------------------------------------------
| File: internal/domain/user/service.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/

package user

import "math"

const PageSize = 8

type Service interface {
	GetAllUsers(page int) (PaginatedUsers, error)
	GetUserByID(id uint) (User, error)
	SearchUsers(name, email string) ([]User, error)
	CreateUser(dto CreateUserDTO) (User, error)
	UpdateUser(id uint, dto UpdateUserDTO) (User, error)
	DeleteUser(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllUsers(page int) (PaginatedUsers, error) {
	users, total, err := s.repo.FindAll(page, PageSize)
	if err != nil {
		return PaginatedUsers{}, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(PageSize)))
	return PaginatedUsers{
		Data:       users,
		Total:      total,
		Page:       page,
		PerPage:    PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *service) GetUserByID(id uint) (User, error) {
	return s.repo.FindByID(id)
}

func (s *service) SearchUsers(name, email string) ([]User, error) {
	return s.repo.Search(name, email)
}

func (s *service) CreateUser(dto CreateUserDTO) (User, error) {
	user := User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password, // A senha será hasheada pelo hook BeforeSave no model
	}
	return s.repo.Create(user)
}

func (s *service) UpdateUser(id uint, dto UpdateUserDTO) (User, error) {
	userToUpdate, err := s.repo.FindByID(id)
	if err != nil {
		return User{}, err
	}
	userToUpdate.Name = dto.Name
	userToUpdate.Email = dto.Email
	// Nota: Não permitimos a atualização da senha aqui por segurança.
	// Isso geralmente é feito em um fluxo separado de "reset de senha".
	return s.repo.Update(userToUpdate)
}

func (s *service) DeleteUser(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
