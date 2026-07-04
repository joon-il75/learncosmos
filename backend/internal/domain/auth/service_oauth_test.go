package auth

import "testing"

func TestParseKakaoUserInfoAllowsMissingEmail(t *testing.T) {
	raw := map[string]interface{}{
		"id": float64(12345),
		"kakao_account": map[string]interface{}{
			"profile": map[string]interface{}{
				"nickname":          "Kakao Learner",
				"profile_image_url": "https://example.com/avatar.png",
			},
		},
	}

	info, err := parseKakaoUserInfo(raw)
	if err != nil {
		t.Fatalf("parseKakaoUserInfo returned error: %v", err)
	}
	if info.Provider != ProviderKakao {
		t.Fatalf("provider = %q, want %q", info.Provider, ProviderKakao)
	}
	if info.ProviderID == "" {
		t.Fatal("provider id should be present")
	}
	if info.Email != "" {
		t.Fatalf("email = %q, want empty string for missing Kakao email", info.Email)
	}
	if info.Nickname != "Kakao Learner" {
		t.Fatalf("nickname = %q", info.Nickname)
	}
}
