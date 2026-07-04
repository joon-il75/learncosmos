package secretmanager

import (
	"errors"
	"os"
	"strings"

	appcrypto "github.com/learnweaver/backend/internal/pkg/crypto"
	"github.com/learnweaver/backend/internal/pkg/ncpkms"
	"github.com/learnweaver/backend/internal/pkg/secretcrypto"
)

const encryptionModeKMSV2 = "kms_v2"

type kmsFactory func() (secretcrypto.KMSClient, error)
type encryptV1Func func(string) (string, error)

// Encrypt stores new BYOK secrets with the configured runtime mode. The safe
// default is legacy v1 encryption; set BYOK_SECRET_ENCRYPTION_MODE=kms_v2 to opt
// in to NCP KMS v2 envelopes for newly saved BYOK values.
func Encrypt(plaintext string) (string, error) {
	return EncryptWithMode(plaintext, encryptionModeFromEnv())
}

// EncryptWithMode stores a secret using an explicit mode. Any mode other than
// kms_v2 falls back to legacy v1 encryption, so callers can opt in per domain.
func EncryptWithMode(plaintext string, mode string) (string, error) {
	return encryptWithMode(plaintext, mode, appcrypto.Encrypt, func() (secretcrypto.KMSClient, error) {
		return ncpkms.NewFromEnv()
	})
}

// Decrypt supports legacy app-level encrypted values and v2 KMS envelopes.
func Decrypt(encoded string) (string, error) {
	if !secretcrypto.IsV2(encoded) {
		return appcrypto.Decrypt(encoded)
	}
	client, err := ncpkms.NewFromEnv()
	if err != nil {
		return "", err
	}
	return secretcrypto.DecryptAuto(encoded, appcrypto.Decrypt, client)
}

// EncryptV2 creates a KMS-backed v2 envelope. It is intentionally exposed for
// explicit migration steps, not used unless a caller opts in.
func EncryptV2(plaintext string) (string, error) {
	client, err := ncpkms.NewFromEnv()
	if err != nil {
		return "", err
	}
	return secretcrypto.EncryptV2(plaintext, client)
}

func encryptWithMode(plaintext string, mode string, encryptV1 encryptV1Func, newKMS kmsFactory) (string, error) {
	if strings.TrimSpace(strings.ToLower(mode)) != encryptionModeKMSV2 {
		if encryptV1 == nil {
			return "", errors.New("v1 encrypt function is required")
		}
		return encryptV1(plaintext)
	}
	if newKMS == nil {
		return "", errors.New("kms factory is required")
	}
	client, err := newKMS()
	if err != nil {
		return "", err
	}
	return secretcrypto.EncryptV2(plaintext, client)
}

func encryptionModeFromEnv() string {
	mode := strings.TrimSpace(os.Getenv("BYOK_SECRET_ENCRYPTION_MODE"))
	if mode == "" {
		mode = strings.TrimSpace(os.Getenv("SECRET_ENCRYPTION_MODE"))
	}
	return strings.ToLower(mode)
}
