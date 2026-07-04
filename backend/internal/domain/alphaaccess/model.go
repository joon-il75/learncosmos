package alphaaccess

type StatusResponse struct {
	Granted      bool    `json:"granted"`
	GrantedAt    *string `json:"granted_at,omitempty"`
	RequiresCode bool    `json:"requires_code"`
	Role         string  `json:"role"`
}

type RedeemRequest struct {
	Code string `json:"code"`
}

type RedeemResponse struct {
	Granted   bool   `json:"granted"`
	GrantedAt string `json:"granted_at"`
	Message   string `json:"message"`
}

type CreateInviteCodeRequest struct {
	ExpiresAt  string `json:"expires_at"`
	MaxUses    int    `json:"max_uses"`
	SentToNote string `json:"sent_to_note"`
	AdminNote  string `json:"admin_note"`
}

type InviteCodeResponse struct {
	ID         string                `json:"id"`
	Code       string                `json:"code"`
	Status     string                `json:"status"`
	State      string                `json:"state"`
	MaxUses    int                   `json:"max_uses"`
	UsedCount  int                   `json:"used_count"`
	ExpiresAt  string                `json:"expires_at"`
	SentToNote string                `json:"sent_to_note"`
	AdminNote  string                `json:"admin_note"`
	CreatedBy  *string               `json:"created_by,omitempty"`
	CreatedAt  string                `json:"created_at"`
	UpdatedAt  string                `json:"updated_at"`
	Uses       []InviteCodeUseRecord `json:"uses"`
}

type InviteCodeUseRecord struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	Email      string  `json:"email"`
	Nickname   string  `json:"nickname"`
	DisplayID  *string `json:"display_id,omitempty"`
	Provider   string  `json:"provider"`
	ProviderID string  `json:"provider_id"`
	UsedAt     string  `json:"used_at"`
}
