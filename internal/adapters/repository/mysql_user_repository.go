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
		INSERT INTO users (id, name, email, password_hash, auth_provider, avatar_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
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
		SELECT id, name, email, password_hash, auth_provider, avatar_url, created_at, updated_at
		FROM users WHERE email = ? AND is_deleted = FALSE LIMIT 1
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
		SELECT id, name, email, password_hash, auth_provider, avatar_url, created_at, updated_at
		FROM users WHERE id = ? AND is_deleted = FALSE LIMIT 1
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

func (r *mysqlUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET name = ?, email = ?, password_hash = ?, auth_provider = ?, avatar_url = ?, updated_at = NOW()
		WHERE id = ? AND is_deleted = FALSE
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.AuthProvider,
		user.AvatarURL,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating user in database: %w", err)
	}

	return nil
}

func (r *mysqlUserRepository) IsAccountActive(ctx context.Context, userID string) (bool, error) {
	var active bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND is_deleted = FALSE)`, userID,
	).Scan(&active)
	return active, err
}

func (r *mysqlUserRepository) DeleteAccount(ctx context.Context, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, query := range []string{
		`DELETE FROM user_profiles WHERE user_id = ?`,
		`DELETE FROM daily_activities WHERE user_id = ?`,
		`DELETE FROM routines WHERE user_id = ?`,
		`DELETE FROM runs WHERE user_id = ?`,
		`DELETE FROM daily_nutrition WHERE user_id = ?`,
		`DELETE FROM hydration_logs WHERE user_id = ?`,
		`DELETE FROM meal_logs WHERE user_id = ?`,
		`DELETE FROM notification_preferences WHERE user_id = ?`,
	} {
		if _, err := tx.ExecContext(ctx, query, userID); err != nil {
			return fmt.Errorf("error deleting account data: %w", err)
		}
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE users
		SET name = 'Deleted user', email = CONCAT('deleted+', id, '@deleted.invalid'),
		    password_hash = SHA2(UUID(), 256), auth_provider = 'local', avatar_url = NULL,
		    is_deleted = TRUE, updated_at = NOW()
		WHERE id = ? AND is_deleted = FALSE
	`, userID)
	if err != nil {
		return fmt.Errorf("error anonymizing account: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return tx.Commit()
}

func (r *mysqlUserRepository) scanUserRow(row *sql.Row) (*domain.User, error) {
	var user domain.User
	var passwordHash sql.NullString
	var avatarURL sql.NullString

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&passwordHash,
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
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}

	return &user, nil
}
