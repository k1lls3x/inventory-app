package repository

import (
	"context"
	"inventory-app/backend/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type DashboardRepository interface {
	 LoadTurnoverByWarehouse(ctx context.Context) ([]model.ItemTurnoverByWarehouse, error)
   LoadTopItems(ctx context.Context) ([]model.ItemWithStock, error)
	 LoadDashboard(ctx context.Context) (*model.DashboardData, error)
}

type PgDashboardRepository struct {
	db *sqlx.DB
	log  zerolog.Logger
}

func NewDashboardRepository(db *sqlx.DB, log zerolog.Logger) *PgDashboardRepository {
	return &PgDashboardRepository{
		db:  db,
		log: log,
	}
}

func (r *PgDashboardRepository) LoadTurnoverByWarehouse(ctx context.Context) ([]model.ItemTurnoverByWarehouse, error) {
	query := `
						SELECT
							w.name AS warehouse,
							i.name AS name,
							i.sku AS sku,
							COALESCE(SUM(inb.quantity), 0) AS received,
							COALESCE(SUM(outb.quantity), 0) AS shipped
						FROM warehouse w
						JOIN stock s ON s.warehouse_id = w.warehouse_id
						JOIN item i ON i.item_id = s.item_id
						LEFT JOIN inbound inb ON inb.item_id = i.item_id AND inb.warehouse_id = w.warehouse_id
							AND date_trunc('month', inb.received_at) = date_trunc('month', current_date)
						LEFT JOIN outbound outb ON outb.item_id = i.item_id AND outb.warehouse_id = w.warehouse_id
							AND date_trunc('month', outb.shipped_at) = date_trunc('month', current_date)
						GROUP BY w.name, i.name, i.sku
						ORDER BY w.name, i.name;
						`

	var result []model.ItemTurnoverByWarehouse
	err := r.db.SelectContext(ctx,&result, query)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении оборота по складам (LoadTurnoverByWarehouse)")
		return nil, err
	}
	r.log.Info().Int("count", len(result)).Msg("Оборот по складам успешно получен")
	return result, nil
}
// LoadTopItems возвращает топ-10 товаров по остатку
func (r *PgDashboardRepository) LoadTopItems(ctx context.Context) ([]model.ItemWithStock, error) {
		query := `
		SELECT
			i.name,
			i.sku,
			w.name AS warehouse,
			s.quantity
		FROM stock s
		JOIN item i ON i.item_id = s.item_id
		JOIN warehouse w ON w.warehouse_id = s.warehouse_id
		ORDER BY s.quantity DESC
		LIMIT 10;
		`

	var items []model.ItemWithStock
	err := r.db.SelectContext(ctx, &items, query)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении топ-товаров по остаткам (LoadTopItems)")
		return nil, err
	}
	r.log.Info().Int("count", len(items)).Msg("Топ товаров по остаткам успешно получен")
	return items, nil
}

func (r *PgDashboardRepository) LoadDashboard(ctx context.Context) (*model.DashboardData, error) {
	var result model.DashboardData

	err := r.db.GetContext(ctx,&result.TotalStock, `SELECT COALESCE(SUM(quantity), 0) FROM Stock`)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении total_stock (LoadDashboard)")
		return nil, err
	}

	err =r.db.GetContext(ctx,&result.ItemCount, `SELECT COUNT(*) FROM Item`)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении item_count (LoadDashboard)")
		return nil, err
	}
	err = r.db.GetContext(ctx,&result.MonthlyOrders, `
	SELECT COUNT(*) FROM inbound
	WHERE date_trunc('month', received_at) = date_trunc('month', CURRENT_DATE)
	`)

	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении monthly_orders (LoadDashboard)")
		return nil, err
	}

	var newItems int
	err = r.db.GetContext(ctx,&newItems, `
		SELECT COUNT(*) FROM item
		WHERE date_trunc('month', created_at) = date_trunc('month', CURRENT_DATE)
	`)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении количества новых товаров (LoadDashboard)")
		return nil, err
	}
		r.log.Info().
		Float64("total_stock", result.TotalStock).
		Int("item_count", result.ItemCount).
		Int("monthly_orders", result.MonthlyOrders).
		Int("new_items", newItems).
		Msg("Данные дашборда успешно получены")

	return &model.DashboardData{
		TotalStock:    result.TotalStock,
		ItemCount:     result.ItemCount,
		MonthlyOrders: result.MonthlyOrders,
		NewItems:      newItems,
	}, nil
}
