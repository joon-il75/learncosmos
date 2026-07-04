package systemsettings

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	appcrypto "github.com/learnweaver/backend/internal/pkg/crypto"
	secretmanager "github.com/learnweaver/backend/internal/pkg/secretmanager"
)

type ProviderCredentials struct {
	Provider    string
	APIKey      string
	APIKey2     string
	EndpointURL *string
}

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) GetProviderCredentials(ctx context.Context, provider string) (ProviderCredentials, error) {
	provider = strings.TrimSpace(provider)
	creds := ProviderCredentials{Provider: provider}
	if provider == "" {
		return creds, nil
	}

	var apiKeyEncrypted *string
	var apiKey2Encrypted *string
	var endpointURL *string

	err := s.db.QueryRow(ctx, `
		SELECT api_key_encrypted, api_key_secondary_encrypted, endpoint_url
		FROM system_api_keys
		WHERE provider = $1
	`, provider).Scan(&apiKeyEncrypted, &apiKey2Encrypted, &endpointURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return creds, nil
		}
		return creds, err
	}

	creds.EndpointURL = endpointURL

	if apiKeyEncrypted != nil && strings.TrimSpace(*apiKeyEncrypted) != "" {
		decrypted, err := secretmanager.Decrypt(*apiKeyEncrypted)
		if err != nil {
			// Backward compatibility: legacy rows may still contain plaintext.
			creds.APIKey = *apiKeyEncrypted
		} else {
			creds.APIKey = decrypted
		}
	}

	if apiKey2Encrypted != nil && strings.TrimSpace(*apiKey2Encrypted) != "" {
		decrypted, err := secretmanager.Decrypt(*apiKey2Encrypted)
		if err != nil {
			// Backward compatibility: legacy rows may still contain plaintext.
			creds.APIKey2 = *apiKey2Encrypted
		} else {
			creds.APIKey2 = decrypted
		}
	}

	return creds, nil
}

func (s *Store) ResolveAPIKey(ctx context.Context, provider, fallback string) (string, error) {
	creds, err := s.GetProviderCredentials(ctx, provider)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(creds.APIKey) != "" {
		return creds.APIKey, nil
	}
	return strings.TrimSpace(fallback), nil
}

func (s *Store) UpsertProviderCredentials(ctx context.Context, provider, apiKey, apiKey2 string, endpointURL *string) error {
	provider = strings.TrimSpace(provider)
	if provider == "" {
		return fmt.Errorf("provider is required")
	}

	if _, err := s.db.Exec(ctx, `
		INSERT INTO system_api_keys (provider)
		VALUES ($1)
		ON CONFLICT (provider) DO NOTHING
	`, provider); err != nil {
		return err
	}

	if strings.TrimSpace(apiKey) != "" {
		encrypted, err := appcrypto.Encrypt(strings.TrimSpace(apiKey))
		if err != nil {
			return fmt.Errorf("%s api_key encrypt: %w", provider, err)
		}
		if _, err := s.db.Exec(ctx, `
			UPDATE system_api_keys
			SET api_key_encrypted = $1, updated_at = NOW()
			WHERE provider = $2
		`, encrypted, provider); err != nil {
			return err
		}
	}

	if strings.TrimSpace(apiKey2) != "" {
		encrypted, err := appcrypto.Encrypt(strings.TrimSpace(apiKey2))
		if err != nil {
			return fmt.Errorf("%s api_key2 encrypt: %w", provider, err)
		}
		if _, err := s.db.Exec(ctx, `
			UPDATE system_api_keys
			SET api_key_secondary_encrypted = $1, updated_at = NOW()
			WHERE provider = $2
		`, encrypted, provider); err != nil {
			return err
		}
	}

	if endpointURL != nil {
		if _, err := s.db.Exec(ctx, `
			UPDATE system_api_keys
			SET endpoint_url = $1, updated_at = NOW()
			WHERE provider = $2
		`, endpointURL, provider); err != nil {
			return err
		}
	}

	return nil
}
