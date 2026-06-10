package service_test

import (
	"database/sql"
	"os"
	"testing"

	"KeyNest/internal/database"
	"KeyNest/internal/models"
	"KeyNest/internal/repository"
	"KeyNest/internal/service"
)

func setupAPIKeyTestDB(t *testing.T) (*sql.DB, func()) {
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

func TestCreateKey_InvalidURL(t *testing.T) {
	db, teardown := setupAPIKeyTestDB(t)
	defer teardown()

	svc := service.NewAPIKeyService(repository.NewAPIKeyRepository(db))
	err := svc.CreateKey(1, make([]byte, 32), models.CreateKeyRequest{
		KeyName:        "My Key",
		KeyValue:       "secret",
		URL:            "ht!tp://bad-url",
		ExpiryDate:     "",
		RegisteredDate: "",
	})
	if err == nil {
		t.Fatal("잘못된 URL을 허용해서는 안 됩니다")
	}
}

func TestCreateKey_InvalidDate(t *testing.T) {
	db, teardown := setupAPIKeyTestDB(t)
	defer teardown()

	svc := service.NewAPIKeyService(repository.NewAPIKeyRepository(db))
	err := svc.CreateKey(1, make([]byte, 32), models.CreateKeyRequest{
		KeyName:        "My Key",
		KeyValue:       "secret",
		URL:            "https://example.com",
		ExpiryDate:     "2024-13-01",
		RegisteredDate: "2024-02-30",
	})
	if err == nil {
		t.Fatal("잘못된 날짜 형식을 허용해서는 안 됩니다")
	}
}

func TestCreateKey_Success(t *testing.T) {
	db, teardown := setupAPIKeyTestDB(t)
	defer teardown()

	userRepo := repository.NewUserRepository(db)
	if _, err := userRepo.Create("user@example.com", "$2a$12$abcdefghijklmnopqrstuv", "0123456789abcdef0123456789abcdef"); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	svc := service.NewAPIKeyService(repository.NewAPIKeyRepository(db))
	err := svc.CreateKey(1, make([]byte, 32), models.CreateKeyRequest{
		KeyName:        "My Key",
		KeyValue:       "secret",
		URL:            "https://example.com",
		ExpiryDate:     "2026-12-31",
		RegisteredDate: "2026-06-05",
		Memo:           "test key",
	})
	if err != nil {
		t.Fatalf("CreateKey: %v", err)
	}
}
