package service

import (
	"context"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type SupplierService struct {
	repo repository.SupplierRepository
	log  zerolog.Logger
	db   *sqlx.DB
}

func NewSupplierService(repo repository.SupplierRepository, db *sqlx.DB, log zerolog.Logger) *SupplierService {
	return &SupplierService{
		repo: repo,
		log:  log,
		db:   db,
	}
}

func (s *SupplierService) GetSuppliers(ctx context.Context) ([]model.Supplier, error) {
	result, err := s.repo.GetSuppliers(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении списка поставщиков")
		return nil, err
	}
	s.log.Info().Int("count", len(result)).Msg("Список поставщиков успешно получен")
	return result, nil
}

func (s *SupplierService) AddSupplier(ctx context.Context, sp model.Supplier) error {
	err := s.repo.AddSupplier(ctx, sp)
	if err != nil {
		s.log.Error().Str("name", sp.Name).Err(err).Msg("Ошибка при добавлении поставщика")
		return err
	}
	s.log.Info().Str("name", sp.Name).Msg("Поставщик успешно добавлен")
	return nil
}

func (s *SupplierService) EditSupplier(ctx context.Context, sp model.Supplier) error {
	err := s.repo.EditSupplier(ctx, sp)
	if err != nil {
		s.log.Error().Int("supplier_id", sp.SupplierID).Err(err).Msg("Ошибка при изменении поставщика")
		return err
	}
	s.log.Info().Int("supplier_id", sp.SupplierID).Msg("Поставщик успешно изменён")
	return nil
}

func (s *SupplierService) RemoveSupplier(ctx context.Context, supplierId int) error {
	err := s.repo.RemoveSupplier(ctx, supplierId)
	if err != nil {
		s.log.Error().Int("supplier_id", supplierId).Err(err).Msg("Ошибка при удалении поставщика")
		return err
	}
	s.log.Info().Int("supplier_id", supplierId).Msg("Поставщик успешно удалён")
	return nil
}
