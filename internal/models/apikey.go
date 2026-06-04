package models

// APIKey is the database model. KeyValue is AES-256-GCM encrypted.
type APIKey struct {
	ID             int64
	UserID         int64
	KeyName        string
	KeyValue       string // encrypted: base64(nonce || ciphertext)
	URL            string
	ExpiryDate     string
	RegisteredDate string
	Memo           string
	CreatedAt      string
	UpdatedAt      string
	ExpiryStatus   int // 0=none, 1=normal, 2=expiring soon, 3=expired
}

// APIKeyDTO is the frontend-facing model. KeyValue is the decrypted plain text.
type APIKeyDTO struct {
	ID             int64  `json:"id"`
	KeyName        string `json:"keyName"`
	KeyValue       string `json:"keyValue"` // decrypted plain text
	URL            string `json:"url"`
	ExpiryDate     string `json:"expiryDate"`
	RegisteredDate string `json:"registeredDate"`
	Memo           string `json:"memo"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	ExpiryStatus   int    `json:"expiryStatus"`
}

// KeyFilter defines search conditions (OR logic between non-empty fields).
type KeyFilter struct {
	KeyName  string `json:"keyName"`
	DateFrom string `json:"dateFrom"` // YYYY-MM-DD, matches registered_date
	DateTo   string `json:"dateTo"`   // YYYY-MM-DD, matches registered_date
}

// CreateKeyRequest is the payload for key registration.
type CreateKeyRequest struct {
	KeyName        string `json:"keyName"`
	KeyValue       string `json:"keyValue"`
	URL            string `json:"url"`
	ExpiryDate     string `json:"expiryDate"`
	RegisteredDate string `json:"registeredDate"`
	Memo           string `json:"memo"`
}

// UpdateKeyRequest is the payload for key modification.
type UpdateKeyRequest struct {
	ID             int64  `json:"id"`
	KeyName        string `json:"keyName"`
	KeyValue       string `json:"keyValue"`
	URL            string `json:"url"`
	ExpiryDate     string `json:"expiryDate"`
	RegisteredDate string `json:"registeredDate"`
	Memo           string `json:"memo"`
}
