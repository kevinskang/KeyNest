package service

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"KeyNest/internal/crypto"
	"KeyNest/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// LoginResult carries the session data established by a successful login.
// The EncryptionKey is held in memory only — never persisted to disk.
type LoginResult struct {
	UserID        int64
	EncryptionKey []byte
}

// Register creates a new user account.
func (s *AuthService) Register(email, password string) error {
	email = normalizeEmail(email)
	if err := validateAuthInput(email, password); err != nil {
		return err
	}

	existing, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("계정 생성 중 오류가 발생했습니다")
	}
	if existing != nil {
		return errors.New("이미 등록된 이메일입니다")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return fmt.Errorf("계정 생성 중 오류가 발생했습니다")
	}

	salt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("계정 생성 중 오류가 발생했습니다")
	}

	if _, err = s.userRepo.Create(email, string(hash), salt); err != nil {
		return fmt.Errorf("계정 생성 중 오류가 발생했습니다")
	}
	return nil
}

// Login authenticates the user and derives the AES encryption key from their password.
func (s *AuthService) Login(email, password string) (*LoginResult, error) {
	email = normalizeEmail(email)
	if email == "" || password == "" {
		return nil, errors.New("이메일과 비밀번호를 입력해주세요")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("로그인 중 오류가 발생했습니다")
	}
	if user == nil {
		return nil, errors.New("이메일 또는 비밀번호가 올바르지 않습니다")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("이메일 또는 비밀번호가 올바르지 않습니다")
	}

	encKey, err := crypto.DeriveKey(password, user.EncSalt)
	if err != nil {
		return nil, fmt.Errorf("로그인 중 오류가 발생했습니다")
	}

	return &LoginResult{UserID: user.ID, EncryptionKey: encKey}, nil
}

func validateAuthInput(email, password string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("유효하지 않은 이메일 형식입니다")
	}
	if len(password) < 8 {
		return errors.New("비밀번호는 최소 8자 이상이어야 합니다")
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
