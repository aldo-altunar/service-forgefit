package ports

import "context"

type GoogleUserInfo struct {
	GoogleID  string
	Email     string
	Name      string
	AvatarURL string
}

type GoogleAuthService interface {
	VerifyIDToken(ctx context.Context, idToken string) (*GoogleUserInfo, error)
}
