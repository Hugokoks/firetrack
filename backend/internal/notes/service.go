package notes

import "errors"

type Service struct {
	repo *Repository
}

var ErrEmptyContent = errors.New("content cannot by empty")

func NewService(repo *Repository) *Service {

	return &Service{repo: repo}

}

func (s *Service) Create(jobID, authorID, content string) (*Note, error) {

	if content == "" {

		return nil, ErrEmptyContent
	}

	return s.repo.Create(jobID, authorID, content)

}

func (s *Service) Delete(noteID string) error {
	return s.repo.DeleteByID(noteID)
}

func (s *Service) GetByJobID(jobID string) ([]Note, error) {
	return s.repo.GetByJobID(jobID)
}
