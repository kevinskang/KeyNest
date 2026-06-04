package service_test

import (
	"database/sql"
	"os"
	"testing"

	"KeyNest/internal/database"
	"KeyNest/internal/repository"
	"KeyNest/internal/service"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	f, err := os.CreateTemp("", "keynest_test_*.db")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	f.Close()

	db, err := database.Open(f.Name())
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	return db, func() {
		db.Close()
		os.Remove(f.Name())
	}
}

func TestRegister_Success(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	svc := service.NewAuthService(repository.NewUserRepository(db))
	if err := svc.Register("user@example.com", "password123"); err != nil {
		t.Fatalf("Register: %v", err)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	svc := service.NewAuthService(repository.NewUserRepository(db))
	_ = svc.Register("user@example.com", "password123")
	if err := svc.Register("user@example.com", "other-pass"); err == nil {
		t.Error("중복 이메일 등록이 허용되어서는 안 됩니다")
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	svc := service.NewAuthService(repository.NewUserRepository(db))
	if err := svc.Register("not-an-email", "password123"); err == nil {
		t.Error("잘못된 이메일 형식이 허용되어서는 안 됩니다")
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	svc := service.NewAuthService(repository.NewUserRepository(db))
	if err := svc.Register("user@example.com", "1234567"); err == nil {
		t.Error("7자 이하 비밀번호가 허용되어서는 안 됩니다")
	}
}

func TestLogin_Success(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	svc := service.NewAuthService(repository.NewUserRepository(db))
	_ = svc.Register("user@example.com", "password123")

	result, err := svc.Login("user@example.com", "password123")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if result.UserID == 0 {
		t.Error("UserID가 0이어서는 안 됩니다")
	}
	if len(result.EncryptionKey) != 32 {
		t.Errorf("EncryptionKey 길이: got %d, want 32", len(result.EncryptionKey))
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	svc := service.NewAuthService(repository.NewUserRepository(db))
	_ = svc.Register("user@example.com", "password123")

	_, err := svc.Login("user@example.com", "wrong-password")
	if err == nil {
		t.Error("잘못된 비밀번호로 로그인이 성공해서는 안 됩니다")
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	svc := service.NewAuthService(repository.NewUserRepository(db))
	_, err := svc.Login("nobody@example.com", "password123")
	if err == nil {
		t.Error("존재하지 않는 이메일로 로그인이 성공해서는 안 됩니다")
	}
}

func TestLogin_DifferentUsersDifferentKeys(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	svc := service.NewAuthService(repository.NewUserRepository(db))
	_ = svc.Register("user1@example.com", "password123")
	_ = svc.Register("user2@example.com", "password123")

	r1, _ := svc.Login("user1@example.com", "password123")
	r2, _ := svc.Login("user2@example.com", "password123")

	same := true
	for i := range r1.EncryptionKey {
		if r1.EncryptionKey[i] != r2.EncryptionKey[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("다른 사용자는 다른 암호화 키를 가져야 합니다 (salt가 다르므로)")
	}
}
