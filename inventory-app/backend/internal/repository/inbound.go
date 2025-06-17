package repository

import (
	"inventory-app/backend/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type InboundRepository interface {
 GetInboundDetails(ctx context.Context)([]model.InboundDetails,error)
 AddInbound(ctx context.Context, inb model.Inbound) error
 AddInboundTx(ctx context.Context, tx *sqlx.Tx, inb model.Inbound) error
 GetInboundDetailsByDate(ctx context.Context,date string)([]model.InboundDetails,error)
 DeleteInbound(ctx context.Context,inboundId int) error
 EditInbound(ctx context.Context,inb model.Inbound) error
}


type PgInboundRepository struct{
	db *sqlx.DB
	log zerolog.Logger
}

func NewInboundRepository(db *sqlx.DB, logger zerolog.Logger) *PgInboundRepository{
	return &PgInboundRepository{
		db : db,
		log : logger,
	}
}

func (r *PgInboundRepository) GetInboundDetails(ctx context.Context)([]model.InboundDetails,error){
	var items []model.InboundDetails
	query:= `
			SELECT
			i.inbound_id,
			to_char(i.received_at, 'YYYY-MM-DD') AS date,
			it.name AS name,
			it.sku AS sku,
			s.name AS supplier,
			i.quantity AS quantity,
			w.name AS warehouse
	FROM inbound i
	JOIN item it       ON i.item_id = it.item_id
	JOIN supplier s    ON i.supplier_id = s.supplier_id
	JOIN warehouse w   ON i.warehouse_id = w.warehouse_id
	ORDER BY i.received_at DESC;
	`
	err:= r.db.SelectContext(ctx,&items,query)

	if err != nil {
		r.log.Error().Err(err).Msg("Ошибка при получении поставок")
	}

	return items,err
}

func (r *PgInboundRepository) AddInbound(ctx context.Context,inb model.Inbound) error {
		query:=`
			INSERT INTO inbound (item_id, supplier_id, quantity, received_at, warehouse_id)
							VALUES ($1, $2, $3, COALESCE($4, now()), $5)
		`
		_, err := r.db.ExecContext(ctx,query, inb.ItemID, inb.SupplierID, inb.Quantity, inb.ReceivedAt, inb.WarehouseID)
		if err != nil {
			r.log.Error().
			Int("item_id", inb.ItemID).
			Int("supplier_id", inb.SupplierID).
			Int("quantity", inb.Quantity).
			Str("received_at", inb.ReceivedAt.Format("2006-01-02")).
			Int("warehouse_id", inb.WarehouseID).
			Err(err).
			Msg("Ошибка при добавлении поставки (AddInbound)")
			return err
		}
		r.log.Info().
		Int("item_id", inb.ItemID).
		Int("supplier_id", inb.SupplierID).
		Int("quantity", inb.Quantity).
		Str("received_at", inb.ReceivedAt.Format("2006-01-02")).
		Int("warehouse_id", inb.WarehouseID).
		Msg("Поставка успешно добавлена")
		return nil
}

func (r *PgInboundRepository) GetInboundDetailsByDate(ctx context.Context,date string)([]model.InboundDetails,error){
	var items []model.InboundDetails
	query := `
						SELECT
					   	i.inbound_id,
							to_char(i.received_at, 'DD.MM.YYYY') AS date, -- вот так!
							it.name AS name,
							it.sku AS sku,
							s.name AS supplier,
							i.quantity AS quantity,
							w.name AS warehouse
					FROM inbound i
					JOIN item it ON i.item_id = it.item_id
					JOIN supplier s ON i.supplier_id = s.supplier_id
					JOIN warehouse w ON i.warehouse_id = w.warehouse_id
					WHERE DATE(i.received_at) = $1
					ORDER BY i.received_at DESC;
	`
	err := r.db.SelectContext(ctx,&items, query,date)
	if err != nil {
		r.log.Error().Err(err).Str("date", date).Msg("Ошибка при получении поставок по дате")
	}
	return items, err
}

func (r *PgInboundRepository) AddInboundTx(ctx context.Context, tx *sqlx.Tx, inb model.Inbound) error {
	query := `
		INSERT INTO inbound (item_id, supplier_id, quantity, received_at, warehouse_id)
		VALUES ($1, $2, $3, COALESCE($4, now()), $5)
	`
	_, err := tx.ExecContext(ctx, query, inb.ItemID, inb.SupplierID, inb.Quantity, inb.ReceivedAt, inb.WarehouseID)
	if err != nil {
		r.log.Error().
			Int("item_id", inb.ItemID).
			Int("supplier_id", inb.SupplierID).
			Int("quantity", inb.Quantity).
			Str("received_at", inb.ReceivedAt.Format("2006-01-02")).
			Int("warehouse_id", inb.WarehouseID).
			Err(err).
			Msg("Ошибка при добавлении поставки")
		return err
	}
	r.log.Info().
		Int("item_id", inb.ItemID).
		Int("supplier_id", inb.SupplierID).
		Int("quantity", inb.Quantity).
		Str("received_at", inb.ReceivedAt.Format("2006-01-02")).
		Int("warehouse_id", inb.WarehouseID).
		Msg("Поставка успешно добавлена")
	return nil
}

func (r *PgInboundRepository) DeleteInbound(ctx context.Context,inboundId int) error {
	query := `DELETE FROM inbound WHERE inbound_id = $1;`
	res, err := r.db.ExecContext(ctx, query, inboundId)
	if err != nil {
		r.log.Error().
			Int("inbound_id", inboundId).
			Err(err).
			Msg("Ошибка при удалении поставки")
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		r.log.Warn().
			Int("inbound_id", inboundId).
			Msg("Поставка не найдена для удаления")
	} else {
		r.log.Info().
			Int("inbound_id", inboundId).
			Msg("Поставка успешно удалена")
	}
	return nil
}

func (r *PgInboundRepository) EditInbound(ctx context.Context,inb model.Inbound) error {
	query := `
	UPDATE inbound
	SET item_id = $1,
		supplier_id = $2,
		warehouse_id = $3,
		quantity = $4,
		received_at = $5
	WHERE inbound_id = $6;
`
res, err := r.db.ExecContext(ctx, query, inb.ItemID, inb.SupplierID, inb.WarehouseID, inb.Quantity, inb.ReceivedAt, inb.InboundID)
if err != nil {
	r.log.Error().
		Int("inbound_id", inb.InboundID).
		Int("item_id", inb.ItemID).
		Err(err).
		Msg("Ошибка при изменении данных поставки")
	return err
}
count, _ := res.RowsAffected()
if count == 0 {
	r.log.Warn().
		Int("inbound_id", inb.InboundID).
		Msg("Поставка не найдена для обновления")
} else {
	r.log.Info().
		Int("inbound_id", inb.InboundID).
		Msg("Поставка успешно обновлена")
}
return nil
}
