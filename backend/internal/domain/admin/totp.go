package admin

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/learnweaver/backend/internal/pkg/secretmanager"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateTOTPSecret은 사용자의 이메일을 기반으로 새로운 TOTP 보안 비밀키를 생성합니다.
func GenerateTOTPSecret(userEmail string) (encryptedSecret, plaintextSecret, qrURL string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "LearnWeaver", //
		AccountName: userEmail,
	})
	if err != nil {
		return "", "", "", err
	}

	plaintextSecret = key.Secret()
	qrURL = key.URL()

	// 보안을 위해 Secret을 암호화하여 DB에 저장합니다.
	encryptedSecret, err = encryptTOTPSecret(plaintextSecret)
	if err != nil {
		return "", "", "", err
	}

	return encryptedSecret, plaintextSecret, qrURL, nil
}

// ValidateTOTP는 사용자가 입력한 6자리 코드를 암호화된 Secret과 대조하여 검증합니다.
func ValidateTOTP(encryptedSecret string, code string) bool {
	// 1. 암호화된 Secret 복호화
	plaintext, err := secretmanager.Decrypt(encryptedSecret)
	if err != nil {
		log.Print("[TOTP] Decrypt error")
		return false
	}

	// 2. 패딩 및 공백 제거 (가장 중요한 수정 사항)
	// AES 복호화 시 발생할 수 있는 Null 바이트(\x00)와 공백을 완전히 제거합니다.
	cleanSecret := strings.TrimSpace(plaintext)
	cleanSecret = strings.Trim(cleanSecret, "\x00")

	// 3. 사용자 입력값 전처리 (공백 제거)
	cleanCode := strings.TrimSpace(code)

	// 4. 상세 검증 (ValidateCustom 사용)
	// Skew: 1 설정을 통해 서버와 앱 간의 30초 정도의 미세한 시간 차이를 허용합니다.
	opts := totp.ValidateOpts{
		Period:    30,
		Skew:      1, // 앞뒤 30초씩 추가 허용 (총 90초 범위)
		Digits:    6,
		Algorithm: otp.AlgorithmSHA1,
	}

	result, err := totp.ValidateCustom(cleanCode, cleanSecret, time.Now().UTC(), opts)

	if err != nil {
		log.Print("[TOTP] Validation error")
		return false
	}

	return result
}

// IsTOTPEnabled는 해당 사용자가 TOTP를 활성화했는지 여부를 DB에서 확인합니다.
func IsTOTPEnabled(ctx context.Context, userID string, db *pgxpool.Pool) bool {
	var enabled bool
	err := db.QueryRow(ctx,
		"SELECT totp_enabled FROM users WHERE id = $1",
		userID,
	).Scan(&enabled)
	if err != nil {
		log.Printf("[TOTP] DB Query error for user %s: %v", userID, err)
		return false
	}
	return enabled
}

func encryptTOTPSecret(plaintextSecret string) (string, error) {
	return secretmanager.EncryptWithMode(plaintextSecret, os.Getenv("TOTP_SECRET_ENCRYPTION_MODE"))
}
