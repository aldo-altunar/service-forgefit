package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type userService struct {
	repo       ports.UserRepository
	hasher     ports.PasswordHasher
	tokenSvc   ports.TokenService
	googleAuth ports.GoogleAuthService
}

func NewUserService(
	repo ports.UserRepository,
	hasher ports.PasswordHasher,
	tokenSvc ports.TokenService,
	googleAuth ports.GoogleAuthService,
) ports.UserService {
	return &userService{
		repo:       repo,
		hasher:     hasher,
		tokenSvc:   tokenSvc,
		googleAuth: googleAuth,
	}
}

func (s *userService) Register(ctx context.Context, input ports.RegisterInput) (*ports.AuthResponse, error) {
	existingUser, err := s.repo.FindByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New().String(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: &hashedPassword,
		AuthProvider: domain.ProviderLocal,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.tokenSvc.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &ports.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *userService) Login(ctx context.Context, input ports.LoginInput) (*ports.AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	if user.PasswordHash == nil || !s.hasher.Compare(input.Password, *user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := s.tokenSvc.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &ports.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *userService) GoogleLogin(ctx context.Context, input ports.GoogleLoginInput) (*ports.AuthResponse, error) {
	googleUser, err := s.googleAuth.VerifyIDToken(ctx, input.IDToken)
	if err != nil {
		return nil, fmt.Errorf("google authentication failed: %w", err)
	}

	user, err := s.repo.FindByGoogleID(ctx, googleUser.GoogleID)
	if errors.Is(err, domain.ErrUserNotFound) || user == nil {
		// Also check if user exists by email
		userByEmail, emailErr := s.repo.FindByEmail(ctx, googleUser.Email)
		if emailErr == nil && userByEmail != nil {
			// Link Google ID to existing email account
			userByEmail.GoogleID = &googleUser.GoogleID
			if googleUser.AvatarURL != "" {
				userByEmail.AvatarURL = &googleUser.AvatarURL
			}
			if err := s.repo.Update(ctx, userByEmail); err != nil {
				return nil, fmt.Errorf("failed to link google account: %w", err)
			}
			user = userByEmail
		} else {
			// Register new user via Google
			newUser := &domain.User{
				ID:           uuid.New().String(),
				Name:         googleUser.Name,
				Email:        googleUser.Email,
				GoogleID:     &googleUser.GoogleID,
				AuthProvider: domain.ProviderGoogle,
				AvatarURL:    &googleUser.AvatarURL,
			}
			if err := s.repo.Create(ctx, newUser); err != nil {
				return nil, fmt.Errorf("failed to create google user: %w", err)
			}
			user = newUser
		}
	}

	token, err := s.tokenSvc.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &ports.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}
