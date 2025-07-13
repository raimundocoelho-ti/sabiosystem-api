/*
|------------------------------------------------
| File: internal/domain/agent/service.go
| Developer: Raimundo Coelho
| GitHub: https://github.com/raimundocoelho-ti
| ------------------------------------------------
*/
package agent

import "math"

const PageSize = 8

type Service interface {
	GetAllAgents(page int) (PaginatedAgents, error)
	GetAgentByID(id uint) (Agent, error)
	SearchAgents(name, domain string) ([]Agent, error)
	CreateAgent(dto CreateAgentDTO) (Agent, error)
	UpdateAgent(id uint, dto UpdateAgentDTO) (Agent, error)
	DeleteAgent(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllAgents(page int) (PaginatedAgents, error) {
	agents, total, err := s.repo.FindAll(page, PageSize)
	if err != nil {
		return PaginatedAgents{}, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(PageSize)))
	return PaginatedAgents{
		Data:       agents,
		Total:      total,
		Page:       page,
		PerPage:    PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *service) GetAgentByID(id uint) (Agent, error) {
	return s.repo.FindByID(id)
}

func (s *service) SearchAgents(name, domain string) ([]Agent, error) {
	return s.repo.Search(name, domain)
}

func (s *service) CreateAgent(dto CreateAgentDTO) (Agent, error) {
	agent := Agent{
		Name:   dto.Name,
		Domain: dto.Domain,
	}
	return s.repo.Create(agent)
}

func (s *service) UpdateAgent(id uint, dto UpdateAgentDTO) (Agent, error) {
	agentToUpdate, err := s.repo.FindByID(id)
	if err != nil {
		return Agent{}, err
	}
	agentToUpdate.Name = dto.Name
	agentToUpdate.Domain = dto.Domain
	return s.repo.Update(agentToUpdate)
}

func (s *service) DeleteAgent(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
