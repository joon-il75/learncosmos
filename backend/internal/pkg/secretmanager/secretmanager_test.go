package secretmanager

import (
	"errors"
	"testing"

	"github.com/learnweaver/backend/internal/pkg/secretcrypto"
)

type fakeKMS struct{}

func (fakeKMS) EncryptSecret(plaintext string) (string, string, error) {
	return "cipher:" + plaintext, "test-key", nil
}

func (fakeKMS) DecryptSecret(ciphertext string, keyTag string) (string, error) {
	return ciphertext + ":" + keyTag, nil
}

func TestEncryptWithModeDefaultsToV1(t *testing.T) {
	got, err := encryptWithMode("secret", "", func(value string) (string, error) {
		return "v1:" + value, nil
	}, func() (secretcrypto.KMSClient, error) {
		t.Fatal("kms factory should not be called for default mode")
		return nil, nil
	})
	if err != nil {
		t.Fatalf("encryptWithMode() error = %v", err)
	}
	if got != "v1:secret" {
		t.Fatalf("encryptWithMode() = %q, want v1 encryption", got)
	}
}

func TestEncryptWithModeUsesKMSV2WhenExplicit(t *testing.T) {
	got, err := encryptWithMode("secret", "kms_v2", func(value string) (string, error) {
		t.Fatal("v1 encrypt should not be called for kms_v2 mode")
		return "", nil
	}, func() (secretcrypto.KMSClient, error) {
		return fakeKMS{}, nil
	})
	if err != nil {
		t.Fatalf("encryptWithMode() error = %v", err)
	}
	if !secretcrypto.IsV2(got) {
		t.Fatalf("encryptWithMode() = %q, want v2 envelope", got)
	}
}

func TestEncryptWithModeReturnsKMSFactoryError(t *testing.T) {
	wantErr := errors.New("kms unavailable")
	_, err := encryptWithMode("secret", "kms_v2", nil, func() (secretcrypto.KMSClient, error) {
		return nil, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("encryptWithMode() error = %v, want %v", err, wantErr)
	}
}
