package service

import (
	"context"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type WarehouseService struct {
	repo repository.WarehouseRepository
	log  zerolog.Logger
	db   *sqlx.DB
}

func NewWarehouseService(repo repository.WarehouseRepository, db *sqlx.DB, log zerolog.Logger) *WarehouseService {
	return &WarehouseService{
		repo: repo,
		log:  log,
		db:   db,
	}
}

func (s *WarehouseService) GetWarehouses(ctx context.Context) ([]model.Warehouse, error) {
	warehouses, err := s.repo.GetWarehouses(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении списка складов")
		return nil, err
	}
	s.log.Info().Int("count", len(warehouses)).Msg("Список складов успешно получен")
	return warehouses, nil
}

func (s *WarehouseService) AddWarehouse(ctx context.Context, warehouse model.Warehouse) error {
	err := s.repo.AddWarehouse(ctx, warehouse)
	if err != nil {
		s.log.Error().Str("name", warehouse.Name).Err(err).Msg("Ошибка при добавлении склада")
		return err
	}
	s.log.Info().Str("name", warehouse.Name).Msg("Склад успешно добавлен")
	return nil
}

func (s *WarehouseService) EditWarehouse(ctx context.Context, warehouse model.Warehouse) error {
	err := s.repo.EditWarehouse(ctx, warehouse)
	if err != nil {
		s.log.Error().Int("warehouse_id", warehouse.WarehouseID).Err(err).Msg("Ошибка при редактировании склада")
		return err
	}
	s.log.Info().Int("warehouse_id", warehouse.WarehouseID).Msg("Склад успешно отредактирован")
	return nil
}
