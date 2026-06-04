package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"KeyNest/internal/models"
)

type APIKeyRepository struct {
	db *sql.DB
}

func NewAPIKeyRepository(db *sql.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// FindByFilter returns keys matching the filter with OR logic.
// If all filter fields are empty, all keys for the user are returned.
func (r *APIKeyRepository) FindByFilter(userID int64, filter models.KeyFilter) ([]models.APIKey, error) {
	query := `
		SELECT id, key_name, key_value, url, expiry_date, registered_date, memo,
		       created_at, updated_at,
		       CASE
		           WHEN expiry_date IS NULL                          THEN 0
		           WHEN DATE(expiry_date) < DATE('now')              THEN 3
		           WHEN DATE(expiry_date) <= DATE('now', '+30 days') THEN 2
		           ELSE 1
		       END AS expiry_status
		FROM api_keys
		WHERE user_id = ?`

	args := []any{userID}

	hasName := filter.KeyName != ""
	hasFrom := filter.DateFrom != ""
	hasTo := filter.DateTo != ""

	if hasName || hasFrom || hasTo {
		var conds []string
		if hasName {
			conds = append(conds, "key_name LIKE ?")
			args = append(args, "%"+filter.KeyName+"%")
		}
		if hasFrom {
			conds = append(conds, "DATE(registered_date) >= ?")
			args = append(args, filter.DateFrom)
		}
		if hasTo {
			conds = append(conds, "DATE(registered_date) <= ?")
			args = append(args, filter.DateTo)
		}
		query += " AND (" + strings.Join(conds, " OR ") + ")"
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("FindByFilter: %w", err)
	}
	defer rows.Close()
	return scanRows(rows)
}

// FindByID returns a single key owned by userID, or nil if not found.
func (r *APIKeyRepository) FindByID(userID, id int64) (*models.APIKey, error) {
	var k models.APIKey
	var url, expiry, registered, memo sql.NullString
	err := r.db.QueryRow(`
		SELECT id, user_id, key_name, key_value, url, expiry_date, registered_date,
		       memo, created_at, updated_at
		FROM api_keys WHERE id = ? AND user_id = ?`, id, userID,
	).Scan(&k.ID, &k.UserID, &k.KeyName, &k.KeyValue,
		&url, &expiry, &registered, &memo,
		&k.CreatedAt, &k.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	k.URL = url.String
	k.ExpiryDate = expiry.String
	k.RegisteredDate = registered.String
	k.Memo = memo.String
	return &k, nil
}

// Create inserts a new api_key row.
func (r *APIKeyRepository) Create(k *models.APIKey) error {
	_, err := r.db.Exec(`
		INSERT INTO api_keys (user_id, key_name, key_value, url, expiry_date, registered_date, memo)
		VALUES (?, ?, ?, NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''))`,
		k.UserID, k.KeyName, k.KeyValue, k.URL, k.ExpiryDate, k.RegisteredDate, k.Memo,
	)
	if err != nil {
		return fmt.Errorf("Create api_key: %w", err)
	}
	return nil
}

// Update modifies an existing api_key row owned by k.UserID.
func (r *APIKeyRepository) Update(k *models.APIKey) error {
	_, err := r.db.Exec(`
		UPDATE api_keys
		SET key_name = ?, key_value = ?,
		    url             = NULLIF(?, ''),
		    expiry_date     = NULLIF(?, ''),
		    registered_date = NULLIF(?, ''),
		    memo            = NULLIF(?, ''),
		    updated_at      = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?`,
		k.KeyName, k.KeyValue, k.URL, k.ExpiryDate, k.RegisteredDate, k.Memo,
		k.ID, k.UserID,
	)
	if err != nil {
		return fmt.Errorf("Update api_key: %w", err)
	}
	return nil
}

// Delete removes the key with the given id that belongs to userID.
func (r *APIKeyRepository) Delete(userID, id int64) error {
	_, err := r.db.Exec(
		`DELETE FROM api_keys WHERE id = ? AND user_id = ?`, id, userID,
	)
	if err != nil {
		return fmt.Errorf("Delete api_key: %w", err)
	}
	return nil
}

// FindExpiring returns keys whose expiry_date is within the next `days` days (including already expired).
func (r *APIKeyRepository) FindExpiring(userID int64, days int) ([]models.APIKey, error) {
	rows, err := r.db.Query(`
		SELECT id, key_name, key_value, url, expiry_date, registered_date, memo,
		       created_at, updated_at,
		       CASE
		           WHEN DATE(expiry_date) < DATE('now')              THEN 3
		           WHEN DATE(expiry_date) <= DATE('now', '+30 days') THEN 2
		           ELSE 1
		       END AS expiry_status
		FROM api_keys
		WHERE user_id = ?
		  AND expiry_date IS NOT NULL
		  AND DATE(expiry_date) <= DATE('now', '+'||?||' days')
		ORDER BY expiry_date ASC`, userID, days,
	)
	if err != nil {
		return nil, fmt.Errorf("FindExpiring: %w", err)
	}
	defer rows.Close()
	return scanRows(rows)
}

// scanRows reads all rows from a query result into a slice.
// Expects columns: id, key_name, key_value, url, expiry_date, registered_date, memo,
//
//	created_at, updated_at, expiry_status
func scanRows(rows *sql.Rows) ([]models.APIKey, error) {
	var keys []models.APIKey
	for rows.Next() {
		var k models.APIKey
		var url, expiry, registered, memo sql.NullString
		if err := rows.Scan(
			&k.ID, &k.KeyName, &k.KeyValue,
			&url, &expiry, &registered, &memo,
			&k.CreatedAt, &k.UpdatedAt, &k.ExpiryStatus,
		); err != nil {
			return nil, fmt.Errorf("scan api_key row: %w", err)
		}
		k.URL = url.String
		k.ExpiryDate = expiry.String
		k.RegisteredDate = registered.String
		k.Memo = memo.String
		keys = append(keys, k)
	}
	return keys, rows.Err()
}
