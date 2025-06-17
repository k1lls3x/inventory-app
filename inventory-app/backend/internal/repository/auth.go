package repository

import (
	"context"
	"inventory-app/backend/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Auth interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	ChangePassword(ctx context.Context, login, oldPassword, newPassword string) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	Register(ctx context.Context, user *model.User, rawPassword string) error
	Authorize(ctx context.Context, login, password string) (*model.User, bool)
}

type AuthRepository struct {
	db  *sqlx.DB
	log zerolog.Logger
}

func NewAuthRepository(db *sqlx.DB, logger zerolog.Logger) *AuthRepository {
	return &AuthRepository{db: db, log: logger}
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	query := `SELECT user_id, username, password_hash, role, full_name FROM "User" WHERE username = $1`
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		r.log.Error().Err(err).Str("username", username).Msg("Ошибка получения пользователя по имени")
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) RegisterUser(ctx context.Context, user *model.User, hashedPassword string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO "User"(username, password_hash, role, full_name) VALUES ($1, $2, $3, $4)`,
		user.Username, hashedPassword, user.Role, user.FullName,
	)
	if err != nil {
		r.log.Error().Err(err).Str("username", user.Username).Msg("Ошибка регистрации пользователя")
	}
	return err
}

func (r *AuthRepository) ChangePassword(ctx context.Context, username, newPasswordHash string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE "User" SET password_hash = $1 WHERE username = $2`, newPasswordHash, username)
	if err != nil {
		r.log.Error().Err(err).Str("username", username).Msg("Ошибка смены пароля")
	}
	return err
}
