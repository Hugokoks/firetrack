package jobs

type Service struct {
	repo *Repository
}

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
