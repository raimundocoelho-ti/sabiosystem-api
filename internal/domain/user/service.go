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
	GetAllUsers(agentID uint, page int) (PaginatedUsers, error)
	GetUserByID(agentID, id uint) (User, error)
	SearchUsers(agentID uint, name, email string) ([]User, error)
	CreateUser(dto CreateUserDTO) (User, error)
	UpdateUser(agentID, id uint, dto UpdateUserDTO) (User, error)
	DeleteUser(agentID, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllUsers(agentID uint, page int) (PaginatedUsers, error) {
	users, total, err := s.repo.FindAll(agentID, page, PageSize)
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

func (s *service) GetUserByID(agentID, id uint) (User, error) {
	return s.repo.FindByID(agentID, id)
}

func (s *service) SearchUsers(agentID uint, name, email string) ([]User, error) {
	return s.repo.Search(agentID, name, email)
}

func (s *service) CreateUser(dto CreateUserDTO) (User, error) {
	user := User{
		AgentID:  dto.AgentID, // <-- MUDANÇA
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
	return s.repo.Create(user)
}

func (s *service) UpdateUser(agentID, id uint, dto UpdateUserDTO) (User, error) {
	userToUpdate, err := s.repo.FindByID(agentID, id)
	if err != nil {
		return User{}, err
	}
	userToUpdate.Name = dto.Name
	userToUpdate.Email = dto.Email
	return s.repo.Update(userToUpdate)
}

func (s *service) DeleteUser(agentID, id uint) error {
	_, err := s.repo.FindByID(agentID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(agentID, id)
}
