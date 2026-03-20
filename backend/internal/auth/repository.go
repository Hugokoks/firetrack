package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var user User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateSession(userID string, expiresAt time.Time) (*Session, error) {
	session := &Session{
		ID:        uuid.NewString(),
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO sessions (id, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, session.ID, session.UserID, session.ExpiresAt, session.CreatedAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}


func (r *Repository) GetUserBySessionID(sessionID string) (*User, error) {
	query := `
		SELECT 
			u.id::text,
			u.name,
			u.email,
			u.password_hash,
			u.role,
			u.is_active,
			u.created_at,
			u.updated_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = $1
		  AND s.expires_at > NOW()
		LIMIT 1
	`

	var user User
	err := r.db.QueryRow(query, sessionID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}