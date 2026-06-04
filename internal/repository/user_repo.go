package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"KeyNest/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByEmail returns the user with the given email, or nil if not found.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		`SELECT id, email, password_hash, enc_salt, created_at, updated_at
		 FROM users WHERE email = ?`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.EncSalt, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindByEmail: %w", err)
	}
	return &u, nil
}

// Create inserts a new user and returns the new row ID.
func (r *UserRepository) Create(email, passwordHash, encSalt string) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO users (email, password_hash, enc_salt) VALUES (?, ?, ?)`,
		email, passwordHash, encSalt,
	)
	if err != nil {
		return 0, fmt.Errorf("Create user: %w", err)
	}
	return res.LastInsertId()
}
