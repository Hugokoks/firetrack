package jobs

import "errors"

type Service struct {
	repo *Repository
}

var ErrJobNotFound = errors.New("job not found")

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateJobInput, createdBy string) (*Job, error) {
	if input.Country == "" {
		input.Country = "CZ"
	}

	if input.Priority == "" {
		input.Priority = "normal"
	}

	return s.repo.Create(input, createdBy)
}

func (s *Service) GetAll() ([]Job, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id string) (*Job, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateJobInput) (*Job, error) {
	job, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrJobNotFound
	}

	applyJobUpdates(job, input)

	return s.repo.Update(job)
}
