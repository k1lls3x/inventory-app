package repository

import (
	"context"
	"inventory-app/backend/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type WarehouseRepository interface {
	GetWarehouses(ctx context.Context) ([]model.Warehouse, error)
	AddWarehouse(ctx context.Context, warehouse model.Warehouse) error
	EditWarehouse(ctx context.Context, warehouse model.Warehouse) error
}

type PgWarehouseRepository struct {
	db  *sqlx.DB
	log zerolog.Logger
}

func NewWarehouseRepository(db *sqlx.DB, log zerolog.Logger) *PgWarehouseRepository {
	return &PgWarehouseRepository{db: db, log: log}
}

func (r *PgWarehouseRepository) GetWarehouses(ctx context.Context) ([]model.Warehouse, error) {
	var warehouses []model.Warehouse
	query := `SELECT warehouse_id, name, location FROM warehouse ORDER BY name`
	err := r.db.SelectContext(ctx, &warehouses, query)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении списка складов")
		return nil, err
	}
	r.log.Info().Int("count", len(warehouses)).Msg("Список складов успешно получен")
	return warehouses, nil
}

func (r *PgWarehouseRepository) AddWarehouse(ctx context.Context, warehouse model.Warehouse) error {
	query := `
	INSERT INTO warehouse (name, location)
	VALUES (:name, :location)
	`
	_, err := r.db.NamedExecContext(ctx, query, warehouse)
	if err != nil {
		r.log.Error().Str("name", warehouse.Name).Err(err).Msg("Ошибка при добавлении склада")
		return err
	}
	r.log.Info().Str("name", warehouse.Name).Msg("Склад успешно добавлен")
	return nil
}

func (r *PgWarehouseRepository) EditWarehouse(ctx context.Context, warehouse model.Warehouse) error {
	query := `
		UPDATE warehouse
		SET name = :name,
			location = :location
		WHERE warehouse_id = :warehouse_id
	`
	_, err := r.db.NamedExecContext(ctx, query, warehouse)
	if err != nil {
		r.log.Error().Int("warehouse_id", warehouse.WarehouseID).Err(err).Msg("Ошибка при редактировании склада")
		return err
	}
	r.log.Info().Int("warehouse_id", warehouse.WarehouseID).Msg("Склад успешно отредактирован")
	return nil
}
