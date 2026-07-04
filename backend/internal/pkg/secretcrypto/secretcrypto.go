package secretcrypto

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const envelopePrefix = "lwsec:v2:"

const ProviderNCPKMS = "ncp-kms"

type V1DecryptFunc func(string) (string, error)

type KMSClient interface {
	EncryptSecret(plaintext string) (ciphertext string, keyTag string, err error)
	DecryptSecret(ciphertext string, keyTag string) (plaintext string, err error)
}

type Envelope struct {
	Version    int    `json:"version"`
	Provider   string `json:"provider"`
	KeyTag     string `json:"key_tag"`
	Ciphertext string `json:"ciphertext"`
}

func EncryptV2(plaintext string, kms KMSClient) (string, error) {
	if kms == nil {
		return "", errors.New("kms client is required")
	}
	ciphertext, keyTag, err := kms.EncryptSecret(plaintext)
	if err != nil {
		return "", fmt.Errorf("kms encrypt secret: %w", err)
	}
	envelope := Envelope{
		Version:    2,
		Provider:   ProviderNCPKMS,
		KeyTag:     strings.TrimSpace(keyTag),
		Ciphertext: strings.TrimSpace(ciphertext),
	}
	if envelope.KeyTag == "" || envelope.Ciphertext == "" {
		return "", errors.New("kms envelope is incomplete")
	}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return "", fmt.Errorf("marshal envelope: %w", err)
	}
	return envelopePrefix + base64.RawURLEncoding.EncodeToString(payload), nil
}

func DecryptAuto(encoded string, decryptV1 V1DecryptFunc, kms KMSClient) (string, error) {
	if strings.HasPrefix(encoded, envelopePrefix) {
		return decryptV2(encoded, kms)
	}
	if decryptV1 == nil {
		return "", errors.New("v1 decrypt function is required")
	}
	return decryptV1(encoded)
}

func IsV2(encoded string) bool {
	return strings.HasPrefix(encoded, envelopePrefix)
}

func DecodeEnvelope(encoded string) (*Envelope, error) {
	if !strings.HasPrefix(encoded, envelopePrefix) {
		return nil, errors.New("secret is not a v2 envelope")
	}
	payload, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(encoded, envelopePrefix))
	if err != nil {
		return nil, errors.New("invalid v2 envelope encoding")
	}
	var envelope Envelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return nil, errors.New("invalid v2 envelope payload")
	}
	if envelope.Version != 2 || envelope.Provider != ProviderNCPKMS || strings.TrimSpace(envelope.KeyTag) == "" || strings.TrimSpace(envelope.Ciphertext) == "" {
		return nil, errors.New("unsupported v2 envelope")
	}
	return &envelope, nil
}

func decryptV2(encoded string, kms KMSClient) (string, error) {
	if kms == nil {
		return "", errors.New("kms client is required")
	}
	envelope, err := DecodeEnvelope(encoded)
	if err != nil {
		return "", err
	}
	plaintext, err := kms.DecryptSecret(envelope.Ciphertext, envelope.KeyTag)
	if err != nil {
		return "", fmt.Errorf("kms decrypt secret: %w", err)
	}
	return plaintext, nil
}
