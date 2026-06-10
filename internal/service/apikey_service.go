package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"KeyNest/internal/crypto"
	"KeyNest/internal/models"
	"KeyNest/internal/repository"
)

const isoDate = "2006-01-02"

type APIKeyService struct {
	keyRepo *repository.APIKeyRepository
}

func NewAPIKeyService(keyRepo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{keyRepo: keyRepo}
}

// GetKeys returns decrypted keys matching the filter.
func (s *APIKeyService) GetKeys(userID int64, encKey []byte, filter models.KeyFilter) ([]models.APIKeyDTO, error) {
	keys, err := s.keyRepo.FindByFilter(userID, filter)
	if err != nil {
		return nil, err
	}
	return s.toDTO(encKey, keys)
}

// CreateKey encrypts the key value and persists it.
func (s *APIKeyService) CreateKey(userID int64, encKey []byte, req models.CreateKeyRequest) error {
	keyName := strings.TrimSpace(req.KeyName)
	if err := validateKeyPayload(keyName, req.KeyValue, req.URL, req.ExpiryDate, req.RegisteredDate); err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(encKey, req.KeyValue)
	if err != nil {
		return fmt.Errorf("키 암호화 중 오류가 발생했습니다")
	}

	return s.keyRepo.Create(&models.APIKey{
		UserID:         userID,
		KeyName:        keyName,
		KeyValue:       encrypted,
		URL:            strings.TrimSpace(req.URL),
		ExpiryDate:     strings.TrimSpace(req.ExpiryDate),
		RegisteredDate: strings.TrimSpace(req.RegisteredDate),
		Memo:           strings.TrimSpace(req.Memo),
	})
}

// UpdateKey re-encrypts the key value and persists the updated row.
func (s *APIKeyService) UpdateKey(userID int64, encKey []byte, req models.UpdateKeyRequest) error {
	keyName := strings.TrimSpace(req.KeyName)
	if err := validateKeyPayload(keyName, req.KeyValue, req.URL, req.ExpiryDate, req.RegisteredDate); err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(encKey, req.KeyValue)
	if err != nil {
		return fmt.Errorf("키 암호화 중 오류가 발생했습니다")
	}

	return s.keyRepo.Update(&models.APIKey{
		ID:             req.ID,
		UserID:         userID,
		KeyName:        keyName,
		KeyValue:       encrypted,
		URL:            strings.TrimSpace(req.URL),
		ExpiryDate:     strings.TrimSpace(req.ExpiryDate),
		RegisteredDate: strings.TrimSpace(req.RegisteredDate),
		Memo:           strings.TrimSpace(req.Memo),
	})
}

// DeleteKey removes the key that belongs to userID.
func (s *APIKeyService) DeleteKey(userID, keyID int64) error {
	return s.keyRepo.Delete(userID, keyID)
}

// GetExpiringKeys returns decrypted keys expiring within the given number of days
// (including already-expired keys).
func (s *APIKeyService) GetExpiringKeys(userID int64, encKey []byte, days int) ([]models.APIKeyDTO, error) {
	keys, err := s.keyRepo.FindExpiring(userID, days)
	if err != nil {
		return nil, err
	}
	return s.toDTO(encKey, keys)
}

// toDTO decrypts key values and converts models to DTOs.
func (s *APIKeyService) toDTO(encKey []byte, keys []models.APIKey) ([]models.APIKeyDTO, error) {
	dtos := make([]models.APIKeyDTO, 0, len(keys))
	for _, k := range keys {
		plain, err := crypto.Decrypt(encKey, k.KeyValue)
		if err != nil {
			return nil, fmt.Errorf("키 복호화 중 오류가 발생했습니다")
		}
		dtos = append(dtos, models.APIKeyDTO{
			ID:             k.ID,
			KeyName:        k.KeyName,
			KeyValue:       plain,
			URL:            k.URL,
			ExpiryDate:     k.ExpiryDate,
			RegisteredDate: k.RegisteredDate,
			Memo:           k.Memo,
			CreatedAt:      k.CreatedAt,
			UpdatedAt:      k.UpdatedAt,
			ExpiryStatus:   k.ExpiryStatus,
		})
	}
	return dtos, nil
}

func validateKeyPayload(keyName, keyValue, rawURL, expiryDate, registeredDate string) error {
	if keyName == "" {
		return errors.New("Key Name은 필수입니다")
	}
	if keyValue == "" {
		return errors.New("Key Value는 필수입니다")
	}

	if rawURL = strings.TrimSpace(rawURL); rawURL != "" {
		if _, err := url.ParseRequestURI(rawURL); err != nil {
			return errors.New("URL 형식이 올바르지 않습니다")
		}
	}

	if err := validateISODate(expiryDate, "만료예정일"); err != nil {
		return err
	}
	if err := validateISODate(registeredDate, "등록일"); err != nil {
		return err
	}
	return nil
}

func validateISODate(value, label string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	if _, err := time.Parse(isoDate, value); err != nil {
		return fmt.Errorf("%s은 YYYY-MM-DD 형식이어야 합니다", label)
	}
	return nil
}
