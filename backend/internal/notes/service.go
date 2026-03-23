package notes

import "errors"

type Service struct {
	repo *Repository
}

var ErrEmptyContent = errors.New("content cannot by empty")
var ErrNoteNotFound = errors.New("note not found")

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

func (s *Service) Update(noteID, content string) (*Note, error) {
	if content == "" {
		return nil, ErrEmptyContent
	}

	note, err := s.repo.UpdateByID(noteID, content)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, ErrNoteNotFound
	}

	return note, nil
}
