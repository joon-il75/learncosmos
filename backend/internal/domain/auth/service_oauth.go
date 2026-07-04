package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/pkg/apperr"
)

func (s *Service) GenerateAuthURL(ctx context.Context, provider Provider, locale string) (string, error) {
	state := uuid.New().String()
	key := redisStatePrefix + state
	stateValue := normalizeOAuthLocale(locale)
	if err := s.redis.Set(ctx, key, stateValue, 10*time.Minute).Err(); err != nil {
		return "", err
	}

	pc := s.providerConfig(provider)

	// redirect_uri가 이미 인코딩된 경우 디코딩해서 plain 문자열로 정규화
	// params.Encode()가 한 번만 인코딩하도록 보장
	redirectURI, err := url.QueryUnescape(pc.RedirectURL)
	if err != nil {
		redirectURI = pc.RedirectURL
	}

	params := url.Values{}
	params.Set("client_id", pc.ClientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("state", state)

	switch provider {
	case ProviderGoogle:
		params.Set("scope", "openid email profile")
		params.Set("prompt", "select_account") // 로그아웃 후 자동 진입 방지
	case ProviderKakao:
		params.Set("scope", "profile_nickname profile_image")
		// prompt 파라미터 미적용 — select_account 미지원, login 강제 시 UX 불편
	case ProviderNaver:
		params.Set("scope", "name email profile_image")
		// 네이버는 prompt 파라미터 미지원 — auth_type=reprompt로 대체
		params.Set("auth_type", "reprompt")
	}

	return pc.AuthURL + "?" + params.Encode(), nil
}

func (s *Service) ValidateState(ctx context.Context, state string) (string, error) {
	key := redisStatePrefix + state
	val, err := s.redis.GetDel(ctx, key).Result()
	if err != nil || strings.TrimSpace(val) == "" {
		return "", apperr.ErrInvalidState
	}
	return normalizeOAuthLocale(val), nil
}

func normalizeOAuthLocale(value string) string {
	if strings.TrimSpace(strings.ToLower(value)) == "en" {
		return "en"
	}
	return "ko"
}

func (s *Service) HandleOAuthCallback(ctx context.Context, provider Provider, code, state string) (*LoginResponse, string, error) {
	pc := s.providerConfig(provider)

	accessToken, err := s.exchangeCode(pc, code)
	if err != nil {
		return nil, "", fmt.Errorf("exchangeCode: %w", err)
	}

	userInfo, err := s.fetchUserInfo(ctx, provider, pc, accessToken)
	if err != nil {
		return nil, "", fmt.Errorf("fetchUserInfo: %w", err)
	}

	user, err := s.repo.FindBySocialAccount(ctx, provider, userInfo.ProviderID)
	if err != nil && err != apperr.ErrNotFound {
		return nil, "", err
	}

	if user != nil {
		if user.Status == "withdrawn" {
			if err := s.repo.ReactivateUser(ctx, user.ID, userInfo.Email, userInfo.Nickname, userInfo.AvatarURL); err != nil {
				return nil, "", fmt.Errorf("ReactivateUser: %w", err)
			}
			user.Status = "active"
			user.WithdrawnAt = nil
			now := time.Now()
			user.ReactivatedAt = &now
		}
		// 재로그인: 닉네임·아바타 갱신, 이메일은 기존 값 있으면 보호
		if err := s.repo.UpdateOnLogin(ctx, user.ID, userInfo.Email, userInfo.Nickname, userInfo.AvatarURL); err != nil {
			return nil, "", fmt.Errorf("UpdateOnLogin: %w", err)
		}
		// 반환용 user 필드도 최신으로 갱신 (이메일은 보호 정책 반영)
		if user.Email == "" {
			user.Email = userInfo.Email
		}
		if user.Nickname == "" {
			user.Nickname = userInfo.Nickname
		}
		if !strings.HasPrefix(user.AvatarURL, "/api/v1/users/avatars/") {
			user.AvatarURL = userInfo.AvatarURL
		}
	}

	if user == nil {
		var displayID *string
		if provider == ProviderKakao && userInfo.Email == "" {
			id, err := s.generateKakaoDisplayID(ctx)
			if err != nil {
				fmt.Printf("display_id 생성 실패: %v\n", err)
			} else {
				displayID = &id
			}
		}
		user = &User{
			ID:        uuid.New().String(),
			Email:     userInfo.Email,
			Nickname:  userInfo.Nickname,
			AvatarURL: userInfo.AvatarURL,
			Role:      RoleLearner,
			DisplayID: displayID,
			CreatedAt: time.Now(),
		}
		if err := s.repo.Create(ctx, user); err != nil {
			return nil, "", err
		}
		sa := &SocialAccount{
			ID:         uuid.New().String(),
			UserID:     user.ID,
			Provider:   provider,
			ProviderID: userInfo.ProviderID,
			CreatedAt:  time.Now(),
		}
		if err := s.repo.CreateSocialAccount(ctx, sa); err != nil {
			return nil, "", err
		}
		welcomePoints, err := s.repo.GetPointSetting(ctx, "welcome_points")
		if err != nil || welcomePoints <= 0 {
			welcomePoints = initialFreePoints
		}
		if err := s.repo.GrantFreePoints(ctx, user.ID, welcomePoints); err != nil {
			return nil, "", err
		}
	}

	return s.IssueSessionForUser(ctx, user)
}

func (s *Service) providerConfig(provider Provider) ProviderConfig {
	switch provider {
	case ProviderGoogle:
		return s.oauth.Google
	case ProviderKakao:
		return s.oauth.Kakao
	case ProviderNaver:
		return s.oauth.Naver
	}
	return ProviderConfig{}
}

func (s *Service) exchangeCode(pc ProviderConfig, code string) (string, error) {
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", pc.ClientID)
	params.Set("client_secret", pc.ClientSecret)
	params.Set("redirect_uri", pc.RedirectURL)
	params.Set("code", code)

	resp, err := http.PostForm(pc.TokenURL, params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("token endpoint status=%d error=%s description=%s", resp.StatusCode, getString(result, "error"), getString(result, "error_description"))
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found in response status=%d error=%s description=%s", resp.StatusCode, getString(result, "error"), getString(result, "error_description"))
	}
	return token, nil
}

func (s *Service) fetchUserInfo(ctx context.Context, provider Provider, pc ProviderConfig, accessToken string) (*OAuthUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pc.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("userinfo endpoint status=%d error=%s description=%s", resp.StatusCode, getString(raw, "error"), getString(raw, "error_description"))
	}

	switch provider {
	case ProviderGoogle:
		return parseGoogleUserInfo(raw)
	case ProviderKakao:
		return parseKakaoUserInfo(raw)
	case ProviderNaver:
		return parseNaverUserInfo(raw)
	}
	return nil, fmt.Errorf("unknown provider: %s", provider)
}

func parseGoogleUserInfo(raw map[string]interface{}) (*OAuthUserInfo, error) {
	return &OAuthUserInfo{
		Provider:   ProviderGoogle,
		ProviderID: getString(raw, "sub"),
		Email:      getString(raw, "email"),
		Nickname:   getString(raw, "name"),
		AvatarURL:  getString(raw, "picture"),
	}, nil
}

func parseKakaoUserInfo(raw map[string]interface{}) (*OAuthUserInfo, error) {
	id := fmt.Sprintf("%v", raw["id"])
	var email, nickname, avatar string
	if account, ok := raw["kakao_account"].(map[string]interface{}); ok {
		email = getString(account, "email")
		if profile, ok := account["profile"].(map[string]interface{}); ok {
			nickname = getString(profile, "nickname")
			avatar = getString(profile, "profile_image_url")
		}
	}
	return &OAuthUserInfo{
		Provider:   ProviderKakao,
		ProviderID: id,
		Email:      email,
		Nickname:   nickname,
		AvatarURL:  avatar,
	}, nil
}

func parseNaverUserInfo(raw map[string]interface{}) (*OAuthUserInfo, error) {
	resp, ok := raw["response"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid naver response")
	}
	return &OAuthUserInfo{
		Provider:   ProviderNaver,
		ProviderID: getString(resp, "id"),
		Email:      getString(resp, "email"),
		Nickname:   getString(resp, "name"),
		AvatarURL:  getString(resp, "profile_image"),
	}, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func (s *Service) generateKakaoDisplayID(ctx context.Context) (string, error) {
	for i := 0; i < 10; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
		if err != nil {
			return "", err
		}
		candidate := fmt.Sprintf("Kakao#%06d", n.Int64())
		exists, err := s.repo.DisplayIDExists(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("display_id 생성 실패: 10회 시도 초과")
}
