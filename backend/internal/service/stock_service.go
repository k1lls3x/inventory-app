package service

import (
	"context"
	"inventory-app/backend/internal/model"
	"inventory-app/backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type StockService struct {
	repo repository.StockRepository
	log  zerolog.Logger
	db   *sqlx.DB
}

func NewStockService(repo repository.StockRepository, db *sqlx.DB, log zerolog.Logger) *StockService {
	return &StockService{
		repo: repo,
		log:  log,
		db:   db,
	}
}

func (s *StockService) GetStockDetails(ctx context.Context) ([]model.ItemWithStock, error) {
	result, err := s.repo.GetStockDetails(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении данных остатков")
		return nil, err
	}
	s.log.Info().Int("count", len(result)).Msg("Остатки успешно получены")
	return result, nil
}

func (s *StockService) GetWeeklyStockTrend(ctx context.Context) ([]model.DailyStock, error) {
	trend, err := s.repo.GetWeeklyStockTrend(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении weekly stock trend")
		return nil, err
	}
	s.log.Info().Int("count", len(trend)).Msg("Weekly stock trend успешно получен")
	return trend, nil
}

func (s *StockService) AddStock(ctx context.Context, itemID, quantity, warehouseID int) error {
	err := s.repo.AddStock(ctx, itemID, quantity, warehouseID)
	if err != nil {
		s.log.Error().
			Int("item_id", itemID).
			Int("warehouse_id", warehouseID).
			Err(err).Msg("Ошибка при добавлении/обновлении остатка")
		return err
	}
	s.log.Info().
		Int("item_id", itemID).
		Int("warehouse_id", warehouseID).
		Msg("Остаток успешно добавлен/обновлён")
	return nil
}

func (s *StockService) FindStockByWarehouse(ctx context.Context, warehouseID int) ([]model.ItemWithStock, error) {
	result, err := s.repo.FindStockByWarehouse(ctx, warehouseID)
	if err != nil {
		s.log.Error().Int("warehouse_id", warehouseID).Err(err).Msg("Ошибка при получении остатков по складу")
		return nil, err
	}
	s.log.Info().Int("warehouse_id", warehouseID).Int("count", len(result)).Msg("Остатки по складу успешно получены")
	return result, nil
}

func (s *StockService) ChangeStock(ctx context.Context, itemID, warehouseID, newQuantity int) error {
	err := s.repo.ChangeStock(ctx, itemID, warehouseID, newQuantity)
	if err != nil {
		s.log.Error().
			Int("item_id", itemID).
			Int("warehouse_id", warehouseID).
			Int("new_quantity", newQuantity).
			Err(err).Msg("Ошибка при изменении stock")
		return err
	}
	s.log.Info().
		Int("item_id", itemID).
		Int("warehouse_id", warehouseID).
		Int("new_quantity", newQuantity).
		Msg("Остаток успешно изменён")
	return nil
}

func (s *StockService) RemoveStock(ctx context.Context, stockID int) error {
	err := s.repo.RemoveStock(ctx, stockID)
	if err != nil {
		s.log.Error().Int("stock_id", stockID).Err(err).Msg("Ошибка при удалении stock")
		return err
	}
	s.log.Info().Int("stock_id", stockID).Msg("Запись stock успешно удалена")
	return nil
}

func (s *StockService) GetStocks(ctx context.Context) ([]model.Stock, error) {
	result, err := s.repo.GetStocks(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Ошибка при получении всех stock")
		return nil, err
	}
	s.log.Info().Int("count", len(result)).Msg("Все остатки успешно получены")
	return result, nil
}
