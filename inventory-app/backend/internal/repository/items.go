package repository

import (
	"inventory-app/backend/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
  "database/sql"
)
type ItemRepository interface {
	GetItems(ctx context.Context) ([]model.Item, error)
	AddItem(ctx context.Context, item model.Item) error
	AddItemTx(ctx context.Context, tx *sqlx.Tx,item model.Item) error
	UpdateItem(ctx context.Context, item model.Item) error
	RemoveItem(ctx context.Context, sku string) error
	FindItems(ctx context.Context, filter model.ItemFilter) ([]model.Item, error)
	GetItemBriefList(ctx context.Context) ([]model.ItemBrief, error)
	ItemExistsTx(ctx context.Context, tx *sqlx.Tx, sku string) (bool, error)
}

type PgItemRepository struct {
	db *sqlx.DB
	log zerolog.Logger
}

func NewPgItemRepository(db *sqlx.DB, logger zerolog.Logger) *PgItemRepository{
	return &PgItemRepository{db: db, log : logger}
}

func (r *PgItemRepository) GetItems(ctx context.Context) ([]model.Item, error) {
  var items []model.Item
    query := `SELECT * FROM item ORDER BY name`
    err := r.db.SelectContext(ctx, &items, query)
    if err != nil {
        r.log.Error().Err(err).Msg("Ошибка получения товаров")
        return nil, err
    }
    r.log.Info().Int("count", len(items)).Msg("Товары успешно получены")
    return items, nil
}

func (r *PgItemRepository) GetItemBriefList(ctx context.Context) ([]model.ItemBrief, error) {
	var items []model.ItemBrief
	query := `SELECT item_id, name, sku FROM item ORDER BY name`
	err := r.db.SelectContext(ctx, &items, query)
	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка получения товаров")
		return nil, err
	}
	r.log.Info().Int("count", len(items)).Msg("Товары успешно получены")
	return items, nil
}

 func (r *PgItemRepository) UpdateItem(ctx context.Context,item model.Item) error  {
	query := `
		UPDATE item
			SET
				sku = :sku,
				name = :name,
				description = :description,
				uom = :uom,
				reorder_level = :reorder_level,
				reorder_qty = :reorder_qty,
				price = :price,
				cost = :cost
			WHERE item_id = :item_id
	`
	_, err := r.db.NamedExecContext(ctx,query,item)
		if err != nil {
			r.log.Error().Err(err).Str("sku", item.SKU).Msg("Ошибка при обновлении товара")
			return err
		}
		r.log.Info().Str("sku", item.SKU).Msg("Товар успешно обновлен")
    return nil
}

func (r *PgItemRepository) RemoveItem (ctx context.Context,sku string) error {
		res , err := r.db.ExecContext(ctx,`DELETE FROM item WHERE sku = $1`,sku)
		if err != nil {
			r.log.Error().Err(err).Str("sku", sku).Msg("Ошибка при удалении товара")
			return err
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			r.log.Warn().Str("sku", sku).Msg("Товар не найден для удаления")
			return sql.ErrNoRows
		}
		r.log.Info().Str("sku", sku).Msg("Товар успешно удален")
    return nil
}

func (r *PgItemRepository) AddItem(ctx context.Context,item model.Item) error {
		query := `
		INSERT INTO item (
				sku, name, description, uom,
				reorder_level, reorder_qty, price, cost, created_at
		) VALUES (
				:sku, :name, :description, :uom,
				:reorder_level, :reorder_qty, :price, :cost, NOW()
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, item)
	if err != nil {
		r.log.Error().Err(err).Str("sku", item.SKU).Msg("Ошибка при добавлении товара")
		return err
	}
	r.log.Info().Str("sku", item.SKU).Msg("Товар успешно добавлен")
	return nil
}

func (r *PgItemRepository) AddItemTx(ctx context.Context, tx *sqlx.Tx, item model.Item) error {
	query := `
			INSERT INTO item (
					sku, name, description, uom,
					reorder_level, reorder_qty, price, cost, created_at
			) VALUES (
					:sku, :name, :description, :uom,
					:reorder_level, :reorder_qty, :price, :cost, NOW()
			)
	`
	_, err := tx.NamedExecContext(ctx, query, item)
	if err != nil {
			r.log.Error().Err(err).Str("sku", item.SKU).Msg("Ошибка при добавлении товара (Tx)")
			return err
	}
	r.log.Info().Str("sku", item.SKU).Msg("Товар успешно добавлен (Tx)")
	return nil
}

func (r *PgItemRepository) FindItems(ctx context.Context, filter model.ItemFilter) ([]model.Item, error) {
	query := `
			SELECT * FROM item
			WHERE
					(:sku IS NULL OR sku = :sku) AND
					(:name IS NULL OR name ILIKE '%' || :name || '%') AND
					(:uom IS NULL OR uom = :uom) AND
					(:price_min IS NULL OR price >= :price_min) AND
					(:price_max IS NULL OR price <= :price_max)
			ORDER BY created_at DESC
	`
	rows, err := r.db.NamedQueryContext(ctx, query, filter)
	if err != nil {
			r.log.Error().Err(err).Msg("Ошибка при поиске товаров")
			return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
			var item model.Item
			if err := rows.StructScan(&item); err != nil {
					r.log.Error().Err(err).Msg("Ошибка при чтении результата поиска")
					return nil, err
			}
			items = append(items, item)
	}
	r.log.Info().Int("count", len(items)).Msg("Поиск товаров завершен")
	return items, nil
}

func (r *PgItemRepository) ItemExistsTx(ctx context.Context, tx *sqlx.Tx, sku string) (bool, error){
	var count int
	err := tx.GetContext(ctx, &count, "SELECT COUNT(1) FROM item WHERE sku = $1", sku)
	return count > 0, err
}
