package ncpkms

import "testing"

func TestCredentialFromEnvUsesUnifiedNCPByDefault(t *testing.T) {
	t.Setenv("NCP_ACCESS_KEY", " ncp-access ")
	t.Setenv("NCP_SECRET_KEY", " ncp-secret ")
	t.Setenv("NCP_OBJECT_STORAGE_ACCESS_KEY", "object-access")
	t.Setenv("NCP_OBJECT_STORAGE_SECRET_KEY", "object-secret")

	accessKey, secretKey, err := credentialFromEnv()
	if err != nil {
		t.Fatalf("credentialFromEnv() error = %v", err)
	}
	if accessKey != "ncp-access" || secretKey != "ncp-secret" {
		t.Fatalf("credentialFromEnv() = %q/%q, want unified NCP credentials", accessKey, secretKey)
	}
}

func TestCredentialFromEnvFallsBackToObjectStorage(t *testing.T) {
	t.Setenv("NCP_OBJECT_STORAGE_ACCESS_KEY", " object-access ")
	t.Setenv("NCP_OBJECT_STORAGE_SECRET_KEY", " object-secret ")

	accessKey, secretKey, err := credentialFromEnv()
	if err != nil {
		t.Fatalf("credentialFromEnv() error = %v", err)
	}
	if accessKey != "object-access" || secretKey != "object-secret" {
		t.Fatalf("credentialFromEnv() = %q/%q, want object storage compatibility credentials", accessKey, secretKey)
	}
}

func TestCredentialFromEnvDedicatedModeRequiresDedicatedKeys(t *testing.T) {
	t.Setenv("NCP_KMS_CREDENTIAL_MODE", "dedicated")
	t.Setenv("NCP_ACCESS_KEY", "ncp-access")
	t.Setenv("NCP_SECRET_KEY", "ncp-secret")

	if _, _, err := credentialFromEnv(); err == nil {
		t.Fatal("credentialFromEnv() error = nil, want missing dedicated credential error")
	}
}

func TestCredentialFromEnvDedicatedMode(t *testing.T) {
	t.Setenv("NCP_KMS_CREDENTIAL_MODE", "dedicated")
	t.Setenv("NCP_KMS_ACCESS_KEY", " kms-access ")
	t.Setenv("NCP_KMS_SECRET_KEY", " kms-secret ")
	t.Setenv("NCP_ACCESS_KEY", "ncp-access")
	t.Setenv("NCP_SECRET_KEY", "ncp-secret")

	accessKey, secretKey, err := credentialFromEnv()
	if err != nil {
		t.Fatalf("credentialFromEnv() error = %v", err)
	}
	if accessKey != "kms-access" || secretKey != "kms-secret" {
		t.Fatalf("credentialFromEnv() = %q/%q, want dedicated credentials", accessKey, secretKey)
	}
}
