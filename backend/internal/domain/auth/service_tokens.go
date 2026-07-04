package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/learnweaver/backend/internal/pkg/apperr"
)

func (s *Service) IssueSessionForUser(ctx context.Context, user *User) (*LoginResponse, string, error) {
	pair, err := s.issueTokenPair(user)
	if err != nil {
		return nil, "", err
	}

	if err := s.storeRefreshToken(ctx, user.ID, pair.RefreshToken); err != nil {
		return nil, "", err
	}

	consentStatus, err := s.repo.GetConsentStatus(ctx, user.ID, normalizeOAuthLocale(user.UILocale))
	if err != nil {
		return nil, "", err
	}

	resp := &LoginResponse{
		AccessToken: pair.AccessToken,
		User: UserInfo{
			ID:                     user.ID,
			Email:                  user.Email,
			Nickname:               user.Nickname,
			AvatarURL:              user.AvatarURL,
			Role:                   user.Role,
			TermsAgreed:            consentStatus.TermsAgreed,
			PrivacyAgreed:          consentStatus.PrivacyAgreed,
			RequiredConsentPending: consentStatus.RequiredConsentPending,
			UILocale:               user.UILocale,
			LearningLanguage:       user.LearningLanguage,
			LanguageSetupRequired:  user.LanguageSetupCompletedAt == nil,
			IsNewSocialSignup:      consentStatus.RequiredConsentPending,
		},
	}
	return resp, pair.RefreshToken, nil
}

func (s *Service) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return "", apperr.ErrUnauthorized
	}

	stored, err := s.redis.Get(ctx, redisRefreshPrefix+claims.UserID).Result()
	if err != nil || stored != refreshToken {
		return "", apperr.ErrUnauthorized
	}

	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return "", err
	}

	return s.issueAccessToken(user)
}

func (s *Service) Logout(ctx context.Context, userID string) error {
	return s.redis.Del(ctx, redisRefreshPrefix+userID).Err()
}

func (s *Service) ParseToken(tokenStr string) (*TokenClaims, error) {
	return s.parseToken(tokenStr)
}

func (s *Service) issueTokenPair(user *User) (*TokenPair, error) {
	access, err := s.issueAccessToken(user)
	if err != nil {
		return nil, err
	}
	refresh, err := s.issueRefreshToken(user)
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *Service) issueAccessToken(user *User) (string, error) {
	claims := &TokenClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func (s *Service) issueRefreshToken(user *User) (string, error) {
	claims := &TokenClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func (s *Service) parseToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func (s *Service) storeRefreshToken(ctx context.Context, userID, token string) error {
	return s.redis.Set(ctx, redisRefreshPrefix+userID, token, refreshTokenTTL).Err()
}

func (s *Service) IssueSuperAdminToken(user User) (string, error) {
	claims := &TokenClaims{
		UserID: user.ID,
		Role:   RoleSuperAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func (s *Service) IssueAdminTOTPToken(userID string) (string, error) {
	claims := &TokenClaims{
		UserID: userID,
		Role:   RoleAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}
