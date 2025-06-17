package service

import (
	"context"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type MovementService struct {
	repo repository.MovementRepository
	log  zerolog.Logger
	db   *sqlx.DB
}

func NewMovementService(repo repository.MovementRepository, db *sqlx.DB, log zerolog.Logger) *MovementService {
	return &MovementService{
		repo: repo,
		log:  log,
		db:   db,
	}
}

func (s *MovementService) GetAllMovementsThisMonth(ctx context.Context) ([]model.Movement, error) {
	movements, err := s.repo.GetAllMovementsThisMonth(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении движений за месяц")
		return nil, err
	}
	s.log.Info().Int("count", len(movements)).Msg("Движения за месяц успешно получены")
	return movements, nil
}
