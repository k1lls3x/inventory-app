package service

import (
	"context"
	"fmt"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.AuthRepository
	log  zerolog.Logger
}

func NewAuthService(repo *repository.AuthRepository, logger zerolog.Logger) *AuthService {
	return &AuthService{
		repo: repo,
		log:  logger,
	}
}

// Хэширование пароля
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверка пароля
func (s *AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Получить пользователя по имени
func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		s.log.Error().Err(err).Str("username", username).Msg("Ошибка поиска пользователя")
		return nil, err
	}
	return user, nil
}

// Регистрация пользователя
func (s *AuthService) Register(ctx context.Context, user *model.User, rawPassword string) error {
	hash, err := s.HashPassword(rawPassword)
	if err != nil {
		s.log.Error().Err(err).Str("username", user.Username).Msg("Ошибка хэширования пароля")
		return err
	}
	err = s.repo.RegisterUser(ctx, user, hash)
	if err != nil {
		s.log.Error().Err(err).Str("username", user.Username).Msg("Ошибка регистрации пользователя")
	}
	return err
}

// Смена пароля
func (s *AuthService) ChangePassword(ctx context.Context, login, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserByUsername(ctx, login)
	if err != nil {
		s.log.Error().Err(err).Str("username", login).Msg("Пользователь не найден для смены пароля")
		return err
	}
	if !s.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return fmt.Errorf("старый пароль неверный")
	}
	newPasswordHash, err := s.HashPassword(newPassword)
	if err != nil {
		s.log.Error().Err(err).Str("username", login).Msg("Ошибка хэширования нового пароля")
		return err
	}
	err = s.repo.ChangePassword(ctx, login, newPasswordHash)
	if err != nil {
		s.log.Error().Err(err).Str("username", login).Msg("Ошибка смены пароля")
	}
	return err
}

// Авторизация
func (s *AuthService) Authorize(ctx context.Context, login, password string) (*model.User, bool) {
	user, err := s.repo.GetUserByUsername(ctx, login)
	if err != nil {
		s.log.Error().Err(err).Str("username", login).Msg("Ошибка поиска пользователя для авторизации")
		return nil, false
	}
	if s.CheckPasswordHash(password, user.PasswordHash) {
		return user, true
	}
	return nil, false
}
func (s *UserService) ChangeUserData(ctx context.Context, u model.UserUpdate) error {
	var hashedPassword *string
	if u.Password != "" {
			hash, err := s.authService.HashPassword(u.Password)
			if err != nil {
					s.log.Error().Err(err).Int("user_id", u.UserID).Msg("Ошибка хэширования пароля")
					return err
			}
			hashedPassword = &hash
	}
	err := s.repo.ChangeUserData(ctx, u, hashedPassword)
	if err != nil {
			s.log.Error().Err(err).Int("user_id", u.UserID).Msg("Ошибка обновления данных пользователя")
	}
	return err
}
