package repository

import (
	"context"
	"inventory-app/backend/internal/model"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type StockRepository interface {
	GetStockDetails(ctx context.Context) ([]model.ItemWithStock, error)
	GetWeeklyStockTrend(ctx context.Context) ([]model.DailyStock, error)
	AddStock(ctx context.Context, itemID, quantity, warehouseID int) error
	FindStockByWarehouse(ctx context.Context, warehouseID int) ([]model.ItemWithStock, error)
	ChangeStock(ctx context.Context, itemID, warehouseID, newQuantity int) error
	RemoveStock(ctx context.Context, stockID int) error
	GetStocks(ctx context.Context) ([]model.Stock, error)
}

type PgStockRepository struct {
	db  *sqlx.DB
	log zerolog.Logger
}

func NewStockRepository(db *sqlx.DB, log zerolog.Logger) *PgStockRepository {
	return &PgStockRepository{db: db, log: log}
}

func (r *PgStockRepository) GetStockDetails(ctx context.Context) ([]model.ItemWithStock, error) {
	var result []model.ItemWithStock
	err := r.db.SelectContext(ctx, &result, `
	SELECT
		s.stock_id,
		i.item_id,
		w.warehouse_id,
		i.name,
		i.sku,
		w.name AS warehouse,
		s.quantity
	FROM stock s
	JOIN item i ON i.item_id = s.item_id
	JOIN warehouse w ON w.warehouse_id = s.warehouse_id
	ORDER BY i.name
	`)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении данных остатков")
		return nil, err
	}
	r.log.Info().Int("count", len(result)).Msg("Остатки успешно получены")
	return result, nil
}

func (r *PgStockRepository) GetWeeklyStockTrend(ctx context.Context) ([]model.DailyStock, error) {
	query := `WITH period AS (
    SELECT
        (CURRENT_DATE - INTERVAL '6 days')::date AS start_date,
        CURRENT_DATE::date AS end_date
),
init_stock AS (
    SELECT
        COALESCE(SUM(quantity), 0) AS qty
    FROM inbound
    WHERE received_at < (SELECT start_date FROM period)
),
init_out AS (
    SELECT
        COALESCE(SUM(quantity), 0) AS qty
    FROM outbound
    WHERE shipped_at < (SELECT start_date FROM period)
),
days AS (
    SELECT generate_series(
        (SELECT start_date FROM period),
        (SELECT end_date FROM period),
        INTERVAL '1 day'
    )::date AS d
),
daily_in AS (
    SELECT received_at::date AS d, SUM(quantity) AS qty
    FROM inbound
    WHERE received_at >= (SELECT start_date FROM period)
    GROUP BY received_at::date
),
daily_out AS (
    SELECT shipped_at::date AS d, SUM(quantity) AS qty
    FROM outbound
    WHERE shipped_at >= (SELECT start_date FROM period)
    GROUP BY shipped_at::date
),
accum AS (
    SELECT
        days.d,
        COALESCE(di.qty, 0) AS in_qty,
        COALESCE(dout.qty, 0) AS out_qty
    FROM days
    LEFT JOIN daily_in di ON di.d = days.d
    LEFT JOIN daily_out dout ON dout.d = days.d
),
running_total AS (
    SELECT
        d,
        SUM(in_qty) OVER (ORDER BY d) - SUM(out_qty) OVER (ORDER BY d) AS net_change
    FROM accum
)
SELECT
    a.d AS date,
    (SELECT qty FROM init_stock) - (SELECT qty FROM init_out) + r.net_change AS total
FROM days a
JOIN running_total r ON a.d = r.d
ORDER BY a.d;
`
	var trend []model.DailyStock
	err := r.db.SelectContext(ctx, &trend, query)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении weekly stock trend")
		return nil, err
	}
	r.log.Info().Int("count", len(trend)).Msg("Weekly stock trend успешно получен")
	return trend, nil
}

func (r *PgStockRepository) AddStock(ctx context.Context, itemID, quantity, warehouseID int) error {
	// Проверка: существует ли товар
	var exists bool
	err := r.db.QueryRowxContext(ctx, `SELECT EXISTS(SELECT 1 FROM item WHERE item_id = $1)`, itemID).Scan(&exists)
	if err != nil {
		r.log.Error().Int("item_id", itemID).Err(err).Msg("Ошибка при проверке наличия item")
		return err
	}
	if !exists {
		r.log.Warn().Int("item_id", itemID).Msg("Товар не существует")
		return fmt.Errorf("Товар с ID %d не существует", itemID)
	}

	// Вставка/обновление остатка
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO stock (item_id, quantity, warehouse_id, last_updated)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (item_id, warehouse_id)
		DO UPDATE SET quantity = EXCLUDED.quantity, last_updated = NOW();
	`, itemID, quantity, warehouseID)

	if err != nil {
		r.log.Error().
			Int("item_id", itemID).
			Int("warehouse_id", warehouseID).
			Err(err).
			Msg("Ошибка при добавлении/обновлении остатка")
		return err
	}
	r.log.Info().
		Int("item_id", itemID).
		Int("warehouse_id", warehouseID).
		Msg("Остаток успешно добавлен/обновлён")
	return nil
}

func (r *PgStockRepository) FindStockByWarehouse(ctx context.Context, warehouseID int) ([]model.ItemWithStock, error) {
	var result []model.ItemWithStock
	err := r.db.SelectContext(ctx, &result, `
	SELECT
		s.stock_id,
		i.item_id,
		w.warehouse_id,
		i.name,
		i.sku,
		w.name AS warehouse,
		s.quantity
	FROM stock s
	JOIN item i ON s.item_id = i.item_id
	JOIN warehouse w ON s.warehouse_id = w.warehouse_id
	WHERE s.warehouse_id = $1
	`, warehouseID)
	if err != nil {
		r.log.Error().Int("warehouse_id", warehouseID).Err(err).Msg("Ошибка при получении остатков по складу")
		return nil, err
	}
	r.log.Info().Int("warehouse_id", warehouseID).Int("count", len(result)).Msg("Остатки по складу успешно получены")
	return result, err
}

func (r *PgStockRepository) ChangeStock(ctx context.Context, itemID, warehouseID, newQuantity int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE stock
		SET quantity = $1, last_updated = NOW()
		WHERE item_id = $2 AND warehouse_id = $3
	`, newQuantity, itemID, warehouseID)

	if err != nil {
		r.log.Error().
			Int("item_id", itemID).
			Int("warehouse_id", warehouseID).
			Err(err).
			Msg("Ошибка при изменении stock")
		return err
	}
	r.log.Info().
		Int("item_id", itemID).
		Int("warehouse_id", warehouseID).
		Int("new_quantity", newQuantity).
		Msg("Остаток успешно изменён")
	return nil
}

func (r *PgStockRepository) RemoveStock(ctx context.Context, stockID int) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM stock
		WHERE stock_id = $1
	`, stockID)
	if err != nil {
		r.log.Error().Int("stock_id", stockID).Err(err).Msg("Ошибка при удалении stock")
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		r.log.Warn().Int("stock_id", stockID).Msg("Запись stock не найдена для удаления")
	} else {
		r.log.Info().Int("stock_id", stockID).Msg("Запись stock успешно удалена")
	}
	return nil
}

func (r *PgStockRepository) GetStocks(ctx context.Context) ([]model.Stock, error) {
	var result []model.Stock
	err := r.db.SelectContext(ctx, &result, `SELECT * FROM stock`)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении всех stock")
		return nil, err
	}
	r.log.Info().Int("count", len(result)).Msg("Все остатки успешно получены")
	return result, err
}
