package repository

import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
	_"time"
)
func GetStockDetails() ([]model.ItemWithStock, error) {
	var result []model.ItemWithStock
	err := db.DB.Select(&result, `
		SELECT
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
		log.Println("❌ Ошибка при получении данных остатков:", err)
		return nil, err
	}
	return result, nil
}
// Добавить или обновить остаток
func AddStock(itemID, quantity, warehouseID int) error {
	_, err := db.DB.Exec(`
		INSERT INTO stock (item_id, quantity, warehouse_id, last_updated)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (item_id, warehouse_id)
		DO UPDATE SET quantity = EXCLUDED.quantity, last_updated = NOW();
	`, itemID, quantity, warehouseID)

	if err != nil {
		log.Println("❌ Ошибка при добавлении stock:", err)
	}
	return err
}

// Найти все остатки по складу
func FindStockByWarehouse(warehouseID int) ([]model.ItemWithStock, error) {
	var result []model.ItemWithStock
	err := db.DB.Select(&result, `
		SELECT
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
		log.Println("❌ Ошибка при получении остатков с деталями по складу:", err)
	}
	return result, err
}

// Обновить количество
func ChangeStock(itemID, warehouseID, newQuantity int) error {
	_, err := db.DB.Exec(`
		UPDATE stock
		SET quantity = $1, last_updated = NOW()
		WHERE item_id = $2 AND warehouse_id = $3
	`, newQuantity, itemID, warehouseID)

	if err != nil {
		log.Println("❌ Ошибка при изменении stock:", err)
	}
	return err
}

// Удалить запись stock
func RemoveStock(itemID, warehouseID int) error {
	_, err := db.DB.Exec(`
		DELETE FROM stock
		WHERE item_id = $1 AND warehouse_id = $2
	`, itemID, warehouseID)

	if err != nil {
		log.Println("❌ Ошибка при удалении stock:", err)
	}
	return err
}

// Получить все остатки
func GetStocks() ([]model.Stock, error) {
	var result []model.Stock
	err := db.DB.Select(&result, `SELECT * FROM stock`)
	if err != nil {
		log.Println("❌ Ошибка при получении всех stock:", err)
	}
	return result, err
}
