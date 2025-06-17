package repository

import (
	"context"
	"inventory-app/backend/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type UserRepository interface {
	RemoveUser(ctx context.Context, userId int) error
	GetUsers(ctx context.Context) ([]model.User, error)
	ChangeUserData(ctx context.Context, u model.UserUpdate, hashedPassword *string) error
}

type PgUserRepository struct {
	db *sqlx.DB
	log zerolog.Logger
}

func NewUserRepository(db *sqlx.DB, log zerolog.Logger) *PgUserRepository {
	return &PgUserRepository{db: db, log: log}
}

func (r *PgUserRepository) RemoveUser(ctx context.Context, userId int) error {
	query := `DELETE FROM "User" WHERE user_id = $1`
	res, err := r.db.ExecContext(ctx, query, userId)
	if err != nil {
			r.log.Error().Err(err).Int("user_id", userId).Msg("Ошибка при удалении пользователя")
			return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
			r.log.Warn().Int("user_id", userId).Msg("Пользователь не найден для удаления")
	} else {
			r.log.Info().Int("user_id", userId).Msg("Пользователь успешно удалён")
	}
	return nil
}

func (r *PgUserRepository) GetUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	query := `SELECT user_id, username, full_name, role FROM "User"`
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
			r.log.Error().Err(err).Msg("Ошибка при получении пользователей")
			return nil, err
	}
	r.log.Info().Int("count", len(users)).Msg("Пользователи успешно получены")
	return users, nil
}

func (r *PgUserRepository) ChangeUserData(ctx context.Context, u model.UserUpdate, hashedPassword *string) error {
	if hashedPassword != nil {
			_, err := r.db.ExecContext(ctx, `
					UPDATE "User" SET username=$1, full_name=$2, role=$3, password_hash=$4 WHERE user_id=$5`,
					u.Username, u.FullName, u.Role, *hashedPassword, u.UserID,
			)
			if err != nil {
					r.log.Error().Err(err).Int("user_id", u.UserID).Msg("Ошибка обновления пользователя с новым паролем")
					return err
			}
			r.log.Info().Int("user_id", u.UserID).Msg("Данные пользователя и пароль успешно обновлены")
			return nil
	} else {
			_, err := r.db.ExecContext(ctx, `
					UPDATE "User" SET username=$1, full_name=$2, role=$3 WHERE user_id=$4`,
					u.Username, u.FullName, u.Role, u.UserID,
			)
			if err != nil {
					r.log.Error().Err(err).Int("user_id", u.UserID).Msg("Ошибка обновления пользователя")
					return err
			}
			r.log.Info().Int("user_id", u.UserID).Msg("Данные пользователя успешно обновлены")
			return nil
	}
}
