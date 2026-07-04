package user

import (
	"time"

	"github.com/learnweaver/backend/internal/domain/auth"
)

type User struct {
	ID            string     `db:"id"`
	Email         string     `db:"email"`
	Nickname      string     `db:"nickname"`
	Role          auth.Role  `db:"role"`
	PremiumAccess bool       `db:"premium_access"`
	AvatarURL     string     `db:"avatar_url"`
	LastLoginAt   *time.Time `db:"last_login_at"`
	CreatedAt     time.Time  `db:"created_at"`
}

type SocialAccount struct {
	ID         string        `db:"id"`
	UserID     string        `db:"user_id"`
	Provider   auth.Provider `db:"provider"`
	ProviderID string        `db:"provider_id"`
	CreatedAt  time.Time     `db:"created_at"`
}

type AIPointWallet struct {
	ID          string     `db:"id"`
	UserID      string     `db:"user_id"`
	FreeBalance int        `db:"free_balance"`
	PaidBalance int        `db:"paid_balance"`
	FreeResetAt *time.Time `db:"free_reset_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
