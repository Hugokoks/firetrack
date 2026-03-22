package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("user is inactive")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type LoginResult struct {
	User    *User
	Session *Session
}

func (s *Service) Login(email, password string) (*LoginResult, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrInactiveUser
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	session, err := s.repo.CreateSession(user.ID, expiresAt)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:    user,
		Session: session,
	}, nil
}

func (s *Service) Logout(sessionID string) error {

	if sessionID == "" {

		return nil
	}

	return s.repo.DeleteSessionById(sessionID)
}
