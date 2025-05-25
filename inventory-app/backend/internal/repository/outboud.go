package repository

import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
)

// Получить все отгрузки с деталями (Outbounds)
func GetOutboundDetails() ([]model.OutboundDetails, error) {
	var items []model.OutboundDetails
	query := `
		SELECT
			o.outbound_id,
			to_char(o.shipped_at, 'YYYY-MM-DD') AS date,
			it.name AS name,
			it.sku AS sku,
			o.destination AS destination,
			o.quantity AS quantity,
			w.name AS warehouse
		FROM outbound o
		JOIN item it       ON o.item_id = it.item_id
		JOIN warehouse w   ON o.warehouse_id = w.warehouse_id
		ORDER BY o.shipped_at DESC;
	`
	err := db.DB.Select(&items, query)

	if err != nil {
		log.Println("❌ Произошла ошибка при получении отгрузок: ", err)
	}

	return items, err
}

// Добавить отгрузку (Outbound)
func AddOutbound(outb model.Outbound) error {
	query := `
		INSERT INTO outbound (item_id, quantity, shipped_at, destination, warehouse_id)
		VALUES ($1, $2, COALESCE($3, now()), $4, $5)
	`
	_, err := db.DB.Exec(query, outb.ItemID, outb.Quantity, outb.ShippedAt, outb.Destination, outb.WarehouseID)
	if err != nil {
		log.Println("❌ Произошла ошибка при добавлении отгрузки: ", err)
	}
	return err
}

// Получить отгрузки по дате
func GetOutboundDetailsByDate(date string) ([]model.OutboundDetails, error) {
	var items []model.OutboundDetails
	query := `
		SELECT
			o.outbound_id,
			to_char(o.shipped_at, 'DD.MM.YYYY') AS date,
			it.name AS name,
			it.sku AS sku,
			o.destination AS destination,
			o.quantity AS quantity,
			w.name AS warehouse
		FROM outbound o
		JOIN item it ON o.item_id = it.item_id
		JOIN warehouse w ON o.warehouse_id = w.warehouse_id
		WHERE DATE(o.shipped_at) = $1
		ORDER BY o.shipped_at DESC;
	`
	err := db.DB.Select(&items, query, date)
	if err != nil {
		log.Println("❌ Произошла ошибка при получении отгрузок по дате: ", err)
	}
	return items, err
}

// Удалить отгрузку
func DeleteOutbound(outboundId int) error {
	query := `
		DELETE FROM outbound WHERE outbound_id = $1;
	`
	_, err := db.DB.Exec(query, outboundId)
	if err != nil {
		log.Println("❌ Произошла ошибка при удалении отгрузки: ", err)
	}
	return err
}

// Редактировать отгрузку
func EditOutbound(outb model.Outbound) error {
	query := `
		UPDATE outbound
		SET item_id = $1,
			warehouse_id = $2,
			quantity = $3,
			shipped_at = $4,
			destination = $5
		WHERE outbound_id = $6;
	`
	_, err := db.DB.Exec(query, outb.ItemID, outb.WarehouseID, outb.Quantity, outb.ShippedAt, outb.Destination, outb.OutboundID)
	if err != nil {
		log.Println("❌ Произошла ошибка при изменении данных отгрузки: ", err)
	}
	return err
}
