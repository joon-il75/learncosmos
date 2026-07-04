package admin

import (
	"strings"

	"github.com/learnweaver/backend/internal/pkg/secretmanager"
	"testing"
	"time"

	"github.com/learnweaver/backend/internal/pkg/crypto"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func TestValidateTOTPAcceptsLegacyEncryptedSecret(t *testing.T) {
	t.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	secret := "JBSWY3DPEHPK3PXP"
	encrypted, err := crypto.Encrypt(secret)
	if err != nil {
		t.Fatalf("encrypt test secret: %v", err)
	}

	code, err := totp.GenerateCodeCustom(secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    6,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		t.Fatalf("generate totp code: %v", err)
	}

	if !ValidateTOTP(encrypted, code) {
		t.Fatal("ValidateTOTP() = false, want true for legacy encrypted secret")
	}
}

func TestEncryptTOTPSecretDefaultsToLegacyEncryption(t *testing.T) {
	t.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	t.Setenv("TOTP_SECRET_ENCRYPTION_MODE", "")

	secret := "JBSWY3DPEHPK3PXP"
	encrypted, err := encryptTOTPSecret(secret)
	if err != nil {
		t.Fatalf("encryptTOTPSecret() error = %v", err)
	}
	if strings.HasPrefix(encrypted, "lwsec:v2:") {
		t.Fatal("encryptTOTPSecret() used v2 without explicit TOTP_SECRET_ENCRYPTION_MODE")
	}
	decrypted, err := secretmanager.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("decrypt encrypted TOTP secret: %v", err)
	}
	if decrypted != secret {
		t.Fatalf("decrypted TOTP secret = %q, want original", decrypted)
	}
}
