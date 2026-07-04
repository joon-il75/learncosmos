package secretcrypto

import (
	"errors"
	"strings"
	"testing"
)

type fakeKMS struct {
	keyTag          string
	ciphertext      string
	plaintext       string
	encryptCalls    int
	decryptCalls    int
	lastDecryptKey  string
	lastDecryptText string
	encryptErr      error
	decryptErr      error
}

func (f *fakeKMS) EncryptSecret(plaintext string) (string, string, error) {
	f.encryptCalls++
	f.plaintext = plaintext
	if f.encryptErr != nil {
		return "", "", f.encryptErr
	}
	return f.ciphertext, f.keyTag, nil
}

func (f *fakeKMS) DecryptSecret(ciphertext string, keyTag string) (string, error) {
	f.decryptCalls++
	f.lastDecryptText = ciphertext
	f.lastDecryptKey = keyTag
	if f.decryptErr != nil {
		return "", f.decryptErr
	}
	return f.plaintext, nil
}

func TestEncryptV2CreatesEnvelopeWithoutPlaintext(t *testing.T) {
	kms := &fakeKMS{keyTag: "kms-key-tag", ciphertext: "kms-ciphertext"}
	encoded, err := EncryptV2("secret-api-key", kms)
	if err != nil {
		t.Fatalf("EncryptV2() error = %v", err)
	}
	if !IsV2(encoded) {
		t.Fatalf("EncryptV2() did not return v2 envelope: %q", encoded)
	}
	if strings.Contains(encoded, "secret-api-key") {
		t.Fatalf("v2 envelope leaked plaintext")
	}
	envelope, err := DecodeEnvelope(encoded)
	if err != nil {
		t.Fatalf("DecodeEnvelope() error = %v", err)
	}
	if envelope.Version != 2 || envelope.Provider != ProviderNCPKMS || envelope.KeyTag != "kms-key-tag" || envelope.Ciphertext != "kms-ciphertext" {
		t.Fatalf("unexpected envelope: %+v", envelope)
	}
	if kms.encryptCalls != 1 {
		t.Fatalf("encrypt calls = %d, want 1", kms.encryptCalls)
	}
}

func TestDecryptAutoUsesV1ForLegacyCiphertext(t *testing.T) {
	kms := &fakeKMS{}
	legacyCiphertext := "legacy-v1-ciphertext"
	got, err := DecryptAuto(legacyCiphertext, func(encoded string) (string, error) {
		if encoded != legacyCiphertext {
			t.Fatalf("v1 decrypt got %q, want %q", encoded, legacyCiphertext)
		}
		return "legacy-plaintext", nil
	}, kms)
	if err != nil {
		t.Fatalf("DecryptAuto() error = %v", err)
	}
	if got != "legacy-plaintext" {
		t.Fatalf("DecryptAuto() = %q, want legacy plaintext", got)
	}
	if kms.decryptCalls != 0 {
		t.Fatalf("kms decrypt calls = %d, want 0", kms.decryptCalls)
	}
}

func TestDecryptAutoUsesKMSForV2Envelope(t *testing.T) {
	kms := &fakeKMS{keyTag: "kms-key-tag", ciphertext: "kms-ciphertext", plaintext: "v2-plaintext"}
	encoded, err := EncryptV2("v2-plaintext", kms)
	if err != nil {
		t.Fatalf("EncryptV2() error = %v", err)
	}
	kms.encryptCalls = 0
	got, err := DecryptAuto(encoded, func(string) (string, error) {
		t.Fatal("v1 decrypt should not be called for v2 envelope")
		return "", nil
	}, kms)
	if err != nil {
		t.Fatalf("DecryptAuto() error = %v", err)
	}
	if got != "v2-plaintext" {
		t.Fatalf("DecryptAuto() = %q, want v2 plaintext", got)
	}
	if kms.decryptCalls != 1 || kms.lastDecryptKey != "kms-key-tag" || kms.lastDecryptText != "kms-ciphertext" {
		t.Fatalf("unexpected kms decrypt call: calls=%d key=%q text=%q", kms.decryptCalls, kms.lastDecryptKey, kms.lastDecryptText)
	}
}

func TestDecryptAutoRejectsInvalidV2WithoutV1Fallback(t *testing.T) {
	_, err := DecryptAuto("lwsec:v2:not-valid-base64", func(string) (string, error) {
		t.Fatal("v1 decrypt should not be called for invalid v2 envelope")
		return "", nil
	}, &fakeKMS{})
	if err == nil {
		t.Fatal("DecryptAuto() error = nil, want invalid v2 error")
	}
}

func TestEncryptV2ReturnsKMSError(t *testing.T) {
	want := errors.New("kms unavailable")
	_, err := EncryptV2("secret", &fakeKMS{encryptErr: want})
	if err == nil || !strings.Contains(err.Error(), "kms encrypt secret") {
		t.Fatalf("EncryptV2() error = %v, want wrapped kms error", err)
	}
}
