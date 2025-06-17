package service


import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"
)

type UserService struct {
repo repository.UserRepository
authService repository.Auth
log zerolog.Logger
db *sqlx.DB
}

func NewUserService(repo repository.UserRepository,logger zerolog.Logger) *UserService{
return &UserService{
	repo: repo,
	log: logger,
}
}

func (s *UserService) RemoveUser(ctx context.Context, userId int) error {
	err := s.repo.RemoveUser(ctx, userId)
	if err != nil {
			s.log.Error().Err(err).Int("user_id", userId).Msg("Ошибка удаления пользователя")
	}
	return err
}

func (s *UserService) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
			s.log.Error().Err(err).Msg("Ошибка получения списка пользователей")
			return nil, err
	}
	return users, nil
}

