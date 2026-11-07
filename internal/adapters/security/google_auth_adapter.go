package security

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/idtoken"

	"ForgeFit_Servidor/internal/core/ports"
)

type googleAuthAdapter struct {
	clientID string
}

func NewGoogleAuthAdapter(clientID string) ports.GoogleAuthService {
	return &googleAuthAdapter{clientID: clientID}
}

func (g *googleAuthAdapter) VerifyIDToken(ctx context.Context, tokenString string) (*ports.GoogleUserInfo, error) {
	var payload *idtoken.Payload
	var err error

	if g.clientID != "" {
		payload, err = idtoken.Validate(ctx, tokenString, g.clientID)
	} else {
		// If clientID is not set in config, validate without audience check for local dev testing
		payload, err = idtoken.Validate(ctx, tokenString, "")
	}

	if err != nil {
		return nil, fmt.Errorf("invalid google id_token: %w", err)
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("google id_token payload missing email")
	}

	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	sub, ok := payload.Subject, payload.Subject != ""
	if !ok {
		sub = payload.Claims["sub"].(string)
	}

	return &ports.GoogleUserInfo{
		GoogleID:  sub,
		Email:     email,
		Name:      name,
		AvatarURL: picture,
	}, nil
}
