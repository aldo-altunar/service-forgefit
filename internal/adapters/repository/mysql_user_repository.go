package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type mysqlUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) ports.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, google_id, auth_provider, avatar_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.GoogleID,
		user.AuthProvider,
		user.AvatarURL,
	)

	if err != nil {
		return fmt.Errorf("error inserting user into database: %w", err)
	}

	return nil
}

func (r *mysqlUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, google_id, auth_provider, avatar_url, created_at, updated_at
		FROM users WHERE email = ? LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, email)
	user, err := r.scanUserRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("error querying user by email: %w", err)
	}

	return user, nil
}

func (r *mysqlUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, google_id, auth_provider, avatar_url, created_at, updated_at
		FROM users WHERE id = ? LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	user, err := r.scanUserRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("error querying user by id: %w", err)
	}

	return user, nil
}

func (r *mysqlUserRepository) FindByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, google_id, auth_provider, avatar_url, created_at, updated_at
		FROM users WHERE google_id = ? LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, googleID)
	user, err := r.scanUserRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("error querying user by google_id: %w", err)
	}

	return user, nil
}

func (r *mysqlUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET name = ?, email = ?, password_hash = ?, google_id = ?, auth_provider = ?, avatar_url = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.GoogleID,
		user.AuthProvider,
		user.AvatarURL,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating user in database: %w", err)
	}

	return nil
}

func (r *mysqlUserRepository) scanUserRow(row *sql.Row) (*domain.User, error) {
	var user domain.User
	var passwordHash sql.NullString
	var googleID sql.NullString
	var avatarURL sql.NullString

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&passwordHash,
		&googleID,
		&user.AuthProvider,
		&avatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if passwordHash.Valid {
		user.PasswordHash = &passwordHash.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}

	return &user, nil
}
