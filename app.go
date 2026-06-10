package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"KeyNest/internal/database"
	"KeyNest/internal/models"
	"KeyNest/internal/repository"
	"KeyNest/internal/service"
)

// App is the Wails application struct. All public methods are bound to the frontend.
type App struct {
	ctx context.Context

	// Session state — protected by mu. encKey is held in memory only, never persisted.
	mu     sync.RWMutex
	userID int64
	encKey []byte

	authService *service.AuthService
	keyService  *service.APIKeyService
}

func NewApp() *App {
	return &App{}
}

// startup is called by Wails when the application initialises.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	exePath, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("cannot locate executable: %v", err))
	}
	dbPath := filepath.Join(filepath.Dir(exePath), "keynest.db")

	db, err := database.Open(dbPath)
	if err != nil {
		panic(fmt.Sprintf("cannot open database at %s: %v", dbPath, err))
	}

	userRepo := repository.NewUserRepository(db)
	keyRepo := repository.NewAPIKeyRepository(db)
	a.authService = service.NewAuthService(userRepo)
	a.keyService = service.NewAPIKeyService(keyRepo)
}

// shutdown is called by Wails when the application is closing.
func (a *App) shutdown(ctx context.Context) {
	a.clearSession()
}

// ─── Auth ────────────────────────────────────────────────────────────────────

// Register creates a new user account.
func (a *App) Register(email, password string) error {
	return a.authService.Register(email, password)
}

// Login authenticates the user and stores the derived encryption key in memory.
func (a *App) Login(email, password string) error {
	result, err := a.authService.Login(email, password)
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.userID = result.UserID
	a.encKey = result.EncryptionKey
	a.mu.Unlock()
	return nil
}

// Logout clears the in-memory session. The encryption key is immediately zeroed.
func (a *App) Logout() {
	a.clearSession()
}

// IsLoggedIn reports whether a user session is active.
func (a *App) IsLoggedIn() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.encKey != nil
}

// ─── Keys ─────────────────────────────────────────────────────────────────────

// GetKeys returns the decrypted key list matching the filter.
// All filter fields empty → returns all keys for the current user.
func (a *App) GetKeys(filter models.KeyFilter) ([]models.APIKeyDTO, error) {
	uid, key, err := a.requireSession()
	if err != nil {
		return nil, err
	}
	return a.keyService.GetKeys(uid, key, filter)
}

// CreateKey encrypts and persists a new key entry.
func (a *App) CreateKey(req models.CreateKeyRequest) error {
	uid, key, err := a.requireSession()
	if err != nil {
		return err
	}
	return a.keyService.CreateKey(uid, key, req)
}

// UpdateKey re-encrypts the key value and updates the row.
func (a *App) UpdateKey(req models.UpdateKeyRequest) error {
	uid, key, err := a.requireSession()
	if err != nil {
		return err
	}
	return a.keyService.UpdateKey(uid, key, req)
}

// DeleteKey removes the key with the given id.
func (a *App) DeleteKey(id int64) error {
	uid, _, err := a.requireSession()
	if err != nil {
		return err
	}
	return a.keyService.DeleteKey(uid, id)
}

// GetExpiringKeys returns decrypted keys expiring within the next `days` days,
// including already-expired keys.
func (a *App) GetExpiringKeys(days int) ([]models.APIKeyDTO, error) {
	uid, key, err := a.requireSession()
	if err != nil {
		return nil, err
	}
	return a.keyService.GetExpiringKeys(uid, key, days)
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

func (a *App) requireSession() (int64, []byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.encKey == nil {
		return 0, nil, fmt.Errorf("인증이 필요합니다")
	}
	return a.userID, a.encKey, nil
}

// clearSession zeros the encryption key before releasing it to the GC.
func (a *App) clearSession() {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := range a.encKey {
		a.encKey[i] = 0
	}
	a.encKey = nil
	a.userID = 0
}
