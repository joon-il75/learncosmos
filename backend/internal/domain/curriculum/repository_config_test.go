package curriculum

import (
	"errors"
	"testing"
)

func TestBuildUserRuntimeAIConfigDoesNotDecryptDisabledKey(t *testing.T) {
	encrypted := "encrypted-api-key"
	called := false

	config, err := buildUserRuntimeAIConfig("openai", &encrypted, nil, false, func(string) (string, error) {
		called = true
		return "should-not-be-used", nil
	})
	if err != nil {
		t.Fatalf("buildUserRuntimeAIConfig returned error: %v", err)
	}
	if called {
		t.Fatalf("disabled BYOK key should not be decrypted")
	}
	if config.Mode != "managed_credit" {
		t.Fatalf("Mode = %q, want managed_credit", config.Mode)
	}
	if config.APIKey != "" {
		t.Fatalf("APIKey = %q, want empty", config.APIKey)
	}
}

func TestBuildUserRuntimeAIConfigDecryptsEnabledKey(t *testing.T) {
	encrypted := "encrypted-api-key"

	config, err := buildUserRuntimeAIConfig(" openai ", &encrypted, nil, true, func(value string) (string, error) {
		if value != encrypted {
			t.Fatalf("decrypt value = %q, want %q", value, encrypted)
		}
		return " decrypted-api-key ", nil
	})
	if err != nil {
		t.Fatalf("buildUserRuntimeAIConfig returned error: %v", err)
	}
	if config.Mode != "byok" {
		t.Fatalf("Mode = %q, want byok", config.Mode)
	}
	if config.Provider != "openai" {
		t.Fatalf("Provider = %q, want openai", config.Provider)
	}
	if config.APIKey != "decrypted-api-key" {
		t.Fatalf("APIKey = %q, want decrypted-api-key", config.APIKey)
	}
}

func TestBuildUserRuntimeAIConfigReturnsDecryptError(t *testing.T) {
	encrypted := "encrypted-api-key"
	decryptErr := errors.New("decrypt failed")

	_, err := buildUserRuntimeAIConfig("openai", &encrypted, nil, true, func(string) (string, error) {
		return "", decryptErr
	})
	if !errors.Is(err, decryptErr) {
		t.Fatalf("error = %v, want decryptErr", err)
	}
}
