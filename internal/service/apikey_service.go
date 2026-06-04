package service

import (
	"errors"
	"fmt"

	"KeyNest/internal/crypto"
	"KeyNest/internal/models"
	"KeyNest/internal/repository"
)

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
	if req.KeyName == "" {
		return errors.New("Key Name은 필수입니다")
	}
	if req.KeyValue == "" {
		return errors.New("Key Value는 필수입니다")
	}

	encrypted, err := crypto.Encrypt(encKey, req.KeyValue)
	if err != nil {
		return fmt.Errorf("키 암호화 중 오류가 발생했습니다")
	}

	return s.keyRepo.Create(&models.APIKey{
		UserID:         userID,
		KeyName:        req.KeyName,
		KeyValue:       encrypted,
		URL:            req.URL,
		ExpiryDate:     req.ExpiryDate,
		RegisteredDate: req.RegisteredDate,
		Memo:           req.Memo,
	})
}

// UpdateKey re-encrypts the key value and persists the updated row.
func (s *APIKeyService) UpdateKey(userID int64, encKey []byte, req models.UpdateKeyRequest) error {
	if req.KeyName == "" {
		return errors.New("Key Name은 필수입니다")
	}
	if req.KeyValue == "" {
		return errors.New("Key Value는 필수입니다")
	}

	encrypted, err := crypto.Encrypt(encKey, req.KeyValue)
	if err != nil {
		return fmt.Errorf("키 암호화 중 오류가 발생했습니다")
	}

	return s.keyRepo.Update(&models.APIKey{
		ID:             req.ID,
		UserID:         userID,
		KeyName:        req.KeyName,
		KeyValue:       encrypted,
		URL:            req.URL,
		ExpiryDate:     req.ExpiryDate,
		RegisteredDate: req.RegisteredDate,
		Memo:           req.Memo,
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
