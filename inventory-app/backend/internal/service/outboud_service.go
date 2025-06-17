package service

import (
	"context"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type OutboundService struct {
	repo repository.OutboundRepository
	log  zerolog.Logger
	db   *sqlx.DB
}

func NewOutboundService(repo repository.OutboundRepository, db *sqlx.DB, log zerolog.Logger) *OutboundService {
	return &OutboundService{
		repo: repo,
		log:  log,
		db:   db,
	}
}

func (s *OutboundService) GetOutboundDetails(ctx context.Context) ([]model.OutboundDetails, error) {
	items, err := s.repo.GetOutboundDetails(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении отгрузок")
		return nil, err
	}
	s.log.Info().Int("count", len(items)).Msg("Список отгрузок успешно получен")
	return items, nil
}

func (s *OutboundService) AddOutbound(ctx context.Context, outb model.Outbound) error {
	err := s.repo.AddOutbound(ctx, outb)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при добавлении отгрузки")
		return err
	}
	s.log.Info().
		Int("item_id", outb.ItemID).
		Int("warehouse_id", outb.WarehouseID).
		Msg("Отгрузка успешно добавлена")
	return nil
}

func (s *OutboundService) GetOutboundDetailsByDate(ctx context.Context, date string) ([]model.OutboundDetails, error) {
	items, err := s.repo.GetOutboundDetailsByDate(ctx, date)
	if err != nil {
		s.log.Error().Str("date", date).Err(err).Msg("Ошибка при получении отгрузок по дате")
		return nil, err
	}
	s.log.Info().Str("date", date).Int("count", len(items)).Msg("Список отгрузок по дате успешно получен")
	return items, nil
}

func (s *OutboundService) DeleteOutbound(ctx context.Context, outboundId int) error {
	err := s.repo.DeleteOutbound(ctx, outboundId)
	if err != nil {
		s.log.Error().Int("outbound_id", outboundId).Err(err).Msg("Ошибка при удалении отгрузки")
		return err
	}
	s.log.Info().Int("outbound_id", outboundId).Msg("Отгрузка успешно удалена")
	return nil
}

func (s *OutboundService) EditOutbound(ctx context.Context, outb model.Outbound) error {
	err := s.repo.EditOutbound(ctx, outb)
	if err != nil {
		s.log.Error().Int("outbound_id", outb.OutboundID).Err(err).Msg("Ошибка при изменении данных отгрузки")
		return err
	}
	s.log.Info().Int("outbound_id", outb.OutboundID).Msg("Данные отгрузки успешно изменены")
	return nil
}
