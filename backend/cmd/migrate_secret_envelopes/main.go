package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/pkg/db"
	"github.com/learnweaver/backend/internal/pkg/secretcrypto"
	"github.com/learnweaver/backend/internal/pkg/secretmanager"
)

type secretTarget struct {
	Table      string
	Column     string
	ID         string
	Provider   string
	Ciphertext string
}

type migrationSummary struct {
	Scanned       int
	AlreadyV2     int
	Empty         int
	Candidates    int
	Migrated      int
	DecryptError  int
	EncryptError  int
	UpdateError   int
	V2VerifyError int
}

func main() {
	_ = godotenv.Load()

	var (
		apply        = flag.Bool("apply", false, "rewrite legacy v1 secrets to KMS v2 envelopes")
		dryRun       = flag.Bool("dry-run", true, "scan and decrypt-check without updating")
		includeUser  = flag.Bool("user-api-keys", true, "include user_api_keys.api_key_encrypted")
		includeAdmin = flag.Bool("system-api-keys", true, "include system_api_keys key columns")
		includeTOTP  = flag.Bool("totp-secrets", true, "include users.totp_secret")
		limit        = flag.Int("limit", 0, "maximum candidate rows to process; 0 means no limit")
		verifyV2     = flag.Bool("verify-v2", false, "decrypt-check existing v2 envelopes without updating")
	)
	flag.Parse()

	ctx := context.Background()
	pool, err := openPool(ctx)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer pool.Close()

	mode := "dry-run"
	if *apply && !*dryRun {
		mode = "apply"
	}

	summary, err := migrateSecrets(ctx, pool, migrationOptions{
		Apply:        *apply && !*dryRun,
		IncludeUser:  *includeUser,
		IncludeAdmin: *includeAdmin,
		IncludeTOTP:  *includeTOTP,
		Limit:        *limit,
		VerifyV2:     *verifyV2,
	})
	if err != nil {
		log.Fatalf("migrate secrets: %v", err)
	}

	log.Printf("mode=%s scanned=%d empty=%d already_v2=%d candidates=%d migrated=%d decrypt_errors=%d encrypt_errors=%d update_errors=%d v2_verify_errors=%d",
		mode,
		summary.Scanned,
		summary.Empty,
		summary.AlreadyV2,
		summary.Candidates,
		summary.Migrated,
		summary.DecryptError,
		summary.EncryptError,
		summary.UpdateError,
		summary.V2VerifyError,
	)

	if summary.DecryptError > 0 || summary.EncryptError > 0 || summary.UpdateError > 0 || summary.V2VerifyError > 0 {
		os.Exit(2)
	}
}

type migrationOptions struct {
	Apply        bool
	IncludeUser  bool
	IncludeAdmin bool
	IncludeTOTP  bool
	Limit        int
	VerifyV2     bool
}

func openPool(ctx context.Context) (*pgxpool.Pool, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	return db.NewPool(ctx, databaseURL)
}

func migrateSecrets(ctx context.Context, pool *pgxpool.Pool, opts migrationOptions) (migrationSummary, error) {
	targets, err := listTargets(ctx, pool, opts)
	if err != nil {
		return migrationSummary{}, err
	}

	summary := migrationSummary{Scanned: len(targets)}
	for _, target := range targets {
		state := classifySecret(target.Ciphertext)
		switch state {
		case "empty":
			summary.Empty++
			continue
		case "v2":
			summary.AlreadyV2++
			if opts.VerifyV2 {
				if _, err := secretmanager.Decrypt(target.Ciphertext); err != nil {
					summary.V2VerifyError++
					log.Printf("v2_verify_failed table=%s column=%s id=%s provider=%s", target.Table, target.Column, target.ID, target.Provider)
				}
			}
			continue
		}

		summary.Candidates++
		plaintext, err := secretmanager.Decrypt(target.Ciphertext)
		if err != nil {
			summary.DecryptError++
			log.Printf("decrypt_failed table=%s column=%s id=%s provider=%s", target.Table, target.Column, target.ID, target.Provider)
			continue
		}
		if !opts.Apply {
			continue
		}

		sealed, err := secretmanager.EncryptV2(plaintext)
		if err != nil {
			summary.EncryptError++
			log.Printf("encrypt_failed table=%s column=%s id=%s provider=%s", target.Table, target.Column, target.ID, target.Provider)
			continue
		}
		if err := updateTarget(ctx, pool, target, sealed); err != nil {
			summary.UpdateError++
			log.Printf("update_failed table=%s column=%s id=%s provider=%s", target.Table, target.Column, target.ID, target.Provider)
			continue
		}
		summary.Migrated++
	}
	return summary, nil
}

func listTargets(ctx context.Context, pool *pgxpool.Pool, opts migrationOptions) ([]secretTarget, error) {
	var targets []secretTarget
	remaining := opts.Limit
	appendTargets := func(rows pgx.Rows, err error) error {
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var target secretTarget
			if err := rows.Scan(&target.Table, &target.Column, &target.ID, &target.Provider, &target.Ciphertext); err != nil {
				return err
			}
			targets = append(targets, target)
			if opts.Limit > 0 {
				remaining--
			}
		}
		return rows.Err()
	}
	queryLimit := func() string {
		if opts.Limit > 0 && remaining > 0 {
			return fmt.Sprintf(" LIMIT %d", remaining)
		}
		return ""
	}
	if opts.IncludeUser && (opts.Limit == 0 || remaining > 0) {
		err := appendTargets(pool.Query(ctx, `
			SELECT 'user_api_keys', 'api_key_encrypted', id::text, provider, COALESCE(api_key_encrypted, '')
			FROM user_api_keys
			WHERE api_key_encrypted IS NOT NULL AND btrim(api_key_encrypted) <> ''
			ORDER BY updated_at ASC, id ASC`+queryLimit()))
		if err != nil {
			return nil, err
		}
	}
	if opts.IncludeAdmin && (opts.Limit == 0 || remaining > 0) {
		err := appendTargets(pool.Query(ctx, `
			SELECT 'system_api_keys', 'api_key_encrypted', id::text, provider, COALESCE(api_key_encrypted, '')
			FROM system_api_keys
			WHERE api_key_encrypted IS NOT NULL AND btrim(api_key_encrypted) <> ''
			ORDER BY updated_at ASC, id ASC`+queryLimit()))
		if err != nil {
			return nil, err
		}
	}
	if opts.IncludeAdmin && (opts.Limit == 0 || remaining > 0) {
		err := appendTargets(pool.Query(ctx, `
			SELECT 'system_api_keys', 'api_key_secondary_encrypted', id::text, provider, COALESCE(api_key_secondary_encrypted, '')
			FROM system_api_keys
			WHERE api_key_secondary_encrypted IS NOT NULL AND btrim(api_key_secondary_encrypted) <> ''
			ORDER BY updated_at ASC, id ASC`+queryLimit()))
		if err != nil {
			return nil, err
		}
	}

	if opts.IncludeTOTP && (opts.Limit == 0 || remaining > 0) {
		err := appendTargets(pool.Query(ctx, `
			SELECT 'users', 'totp_secret', id::text, COALESCE(email, ''), COALESCE(totp_secret, '')
			FROM users
			WHERE totp_secret IS NOT NULL AND btrim(totp_secret) <> ''
			ORDER BY updated_at ASC, id ASC`+queryLimit()))
		if err != nil {
			return nil, err
		}
	}
	return targets, nil
}

func updateTarget(ctx context.Context, pool *pgxpool.Pool, target secretTarget, sealed string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	switch target.Table + "." + target.Column {
	case "user_api_keys.api_key_encrypted":
		cmdTag, err := pool.Exec(ctx, `
			UPDATE user_api_keys
			SET api_key_encrypted = $1, updated_at = NOW()
			WHERE id = $2 AND api_key_encrypted = $3
		`, sealed, target.ID, target.Ciphertext)
		if err != nil {
			return err
		}
		if cmdTag.RowsAffected() != 1 {
			return fmt.Errorf("unexpected rows affected: %d", cmdTag.RowsAffected())
		}
		return nil
	case "system_api_keys.api_key_encrypted":
		cmdTag, err := pool.Exec(ctx, `
			UPDATE system_api_keys
			SET api_key_encrypted = $1, updated_at = NOW()
			WHERE id = $2 AND api_key_encrypted = $3
		`, sealed, target.ID, target.Ciphertext)
		if err != nil {
			return err
		}
		if cmdTag.RowsAffected() != 1 {
			return fmt.Errorf("unexpected rows affected: %d", cmdTag.RowsAffected())
		}
		return nil

	case "users.totp_secret":
		cmdTag, err := pool.Exec(ctx, `
			UPDATE users
			SET totp_secret = $1, updated_at = NOW()
			WHERE id = $2 AND totp_secret = $3
		`, sealed, target.ID, target.Ciphertext)
		if err != nil {
			return err
		}
		if cmdTag.RowsAffected() != 1 {
			return fmt.Errorf("unexpected rows affected: %d", cmdTag.RowsAffected())
		}
		return nil
	case "system_api_keys.api_key_secondary_encrypted":
		cmdTag, err := pool.Exec(ctx, `
			UPDATE system_api_keys
			SET api_key_secondary_encrypted = $1, updated_at = NOW()
			WHERE id = $2 AND api_key_secondary_encrypted = $3
		`, sealed, target.ID, target.Ciphertext)
		if err != nil {
			return err
		}
		if cmdTag.RowsAffected() != 1 {
			return fmt.Errorf("unexpected rows affected: %d", cmdTag.RowsAffected())
		}
		return nil
	default:
		return fmt.Errorf("unsupported target: %s.%s", target.Table, target.Column)
	}
}

func classifySecret(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "empty"
	}
	if secretcrypto.IsV2(trimmed) {
		return "v2"
	}
	return "v1"
}
