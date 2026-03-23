package activity

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {

	return &Service{repo: repo}

}

const ContextActivityKey = "job_activity_payload"

func (s *Service) Log(payload Payload) error {

	return s.repo.Create(payload)
}
