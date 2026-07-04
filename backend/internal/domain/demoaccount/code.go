package demoaccount

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

const loginCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func NormalizeLoginCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func HashLoginCode(code string, secret string) (string, error) {
	normalized := NormalizeLoginCode(code)
	secret = strings.TrimSpace(secret)
	if normalized == "" {
		return "", fmt.Errorf("login code is empty")
	}
	if secret == "" {
		return "", fmt.Errorf("demo login secret is empty")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(normalized))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GenerateLoginCode(accountNumber int) (string, error) {
	suffix, err := randomCodeSuffix(6)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("LW-DEMO-%02d-%s", accountNumber, suffix), nil
}

func GenerateLoginCodeForLabel(label string) (string, error) {
	accountNumber := trailingNumber(label)
	if accountNumber <= 0 {
		accountNumber = 1
	}
	return GenerateLoginCode(accountNumber)
}

func randomCodeSuffix(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid code suffix length")
	}

	var b strings.Builder
	b.Grow(length)
	max := big.NewInt(int64(len(loginCodeAlphabet)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b.WriteByte(loginCodeAlphabet[n.Int64()])
	}
	return b.String(), nil
}

func trailingNumber(value string) int {
	value = strings.TrimSpace(value)
	end := len(value)
	start := end
	for start > 0 {
		r, size := utf8LastRune(value[:start])
		if r == 0 || !unicode.IsDigit(r) {
			break
		}
		start -= size
	}
	if start == end {
		return 0
	}
	n, err := strconv.Atoi(value[start:end])
	if err != nil {
		return 0
	}
	return n
}

func utf8LastRune(value string) (rune, int) {
	for i := len(value); i > 0; i-- {
		if value[i-1] < 0x80 {
			return rune(value[i-1]), 1
		}
	}
	return 0, 0
}
