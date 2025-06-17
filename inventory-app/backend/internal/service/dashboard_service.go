package service

import (
	"context"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type DashboardService struct {
	repo repository.DashboardRepository
	log  zerolog.Logger
	db   *sqlx.DB
}

func NewDashboardService(repo repository.DashboardRepository, db *sqlx.DB, log zerolog.Logger) *DashboardService {
	return &DashboardService{
		repo: repo,
		log:  log,
		db:   db,
	}
}

func (s *DashboardService) LoadTurnoverByWarehouse(ctx context.Context) ([]model.ItemTurnoverByWarehouse, error) {
	result, err := s.repo.LoadTurnoverByWarehouse(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении оборота по складам")
		return nil, err
	}
	s.log.Info().Int("count", len(result)).Msg("Оборот по складам успешно получен")
	return result, nil
}

func (s *DashboardService) LoadTopItems(ctx context.Context) ([]model.ItemWithStock, error) {
	items, err := s.repo.LoadTopItems(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении топ-товаров по остаткам")
		return nil, err
	}
	s.log.Info().Int("count", len(items)).Msg("Топ товаров по остаткам успешно получен")
	return items, nil
}

func (s *DashboardService) LoadDashboard(ctx context.Context) (*model.DashboardData, error) {
	data, err := s.repo.LoadDashboard(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении данных дашборда")
		return nil, err
	}
	s.log.Info().
		Float64("total_stock", data.TotalStock).
		Int("item_count", data.ItemCount).
		Int("monthly_orders", data.MonthlyOrders).
		Int("new_items", data.NewItems).
		Msg("Данные дашборда успешно получены")
	return data, nil
}
