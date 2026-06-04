package models

// User represents the users table.
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	EncSalt      string
	CreatedAt    string
	UpdatedAt    string
}
