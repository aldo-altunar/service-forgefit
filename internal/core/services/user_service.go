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
	repo     ports.UserRepository
	hasher   ports.PasswordHasher
	tokenSvc ports.TokenService
}

func NewUserService(
	repo ports.UserRepository,
	hasher ports.PasswordHasher,
	tokenSvc ports.TokenService,
) ports.UserService {
	return &userService{
		repo:     repo,
		hasher:   hasher,
		tokenSvc: tokenSvc,
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

func (s *userService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (s *userService) UpdateAccountSettings(ctx context.Context, userID string, input ports.UpdateAccountSettingsInput) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	existing, err := s.repo.FindByEmail(ctx, input.Email)
	if err == nil && existing.ID != userID {
		return nil, domain.ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	user.Name = input.Name
	user.Email = input.Email
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteAccount(ctx context.Context, userID string) error {
	return s.repo.DeleteAccount(ctx, userID)
}
