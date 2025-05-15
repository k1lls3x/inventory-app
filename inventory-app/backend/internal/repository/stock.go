package repository

import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
	_"time"
	"fmt"
)

func GetStockDetails() ([]model.ItemWithStock, error) {
	var result []model.ItemWithStock
	err := db.DB.Select(&result, `
	SELECT
		s.stock_id,         -- ← обязательно!
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
		log.Println("❌ Ошибка при получении данных остатков:", err)
		return nil, err
	}
	return result, nil
}

func GetWeeklyStockTrend() ([]model.DailyStock, error) {

	query := `
		SELECT
			d::date AS date,
			COALESCE(SUM(s.quantity), 0) AS total
		FROM generate_series(
			CURRENT_DATE - interval '6 days',
			CURRENT_DATE,
			interval '1 day'
		) d
		LEFT JOIN stock s ON date_trunc('day', s.last_updated) <= d
		GROUP BY d
		ORDER BY d;
	`

	var trend []model.DailyStock
	err := db.DB.Select(&trend, query)
	if err != nil {
		log.Println("❌ Ошибка при получении weekly stock trend:", err)
		return nil, err
	}

	return trend, nil
}

// Добавить или обновить остаток
func AddStock(itemID, quantity, warehouseID int) error {
	// Проверка: существует ли товар
	var exists bool
	err := db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM item WHERE item_id = $1)`, itemID).Scan(&exists)
	if err != nil {
		log.Println("❌ Ошибка при проверке item:", err)
		return err
	}
	if !exists {
		return fmt.Errorf("Товар с ID %d не существует", itemID)
	}

	// Вставка/обновление остатка
	_, err = db.DB.Exec(`
		INSERT INTO stock (item_id, quantity, warehouse_id, last_updated)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (item_id, warehouse_id)
		DO UPDATE SET quantity = EXCLUDED.quantity, last_updated = NOW();
	`, itemID, quantity, warehouseID)

	if err != nil {
		log.Println("❌ Ошибка при добавлении остатка:", err)
	}
	return err
}

// Найти все остатки по складу
func FindStockByWarehouse(warehouseID int) ([]model.ItemWithStock, error) {
	var result []model.ItemWithStock
	err := db.DB.Select(&result, `
	SELECT
		s.stock_id,         -- ← обязательно!
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
		log.Println("❌ Ошибка при получении остатков с деталями по складу:", err)
	}
	return result, err
}


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
func RemoveStock(stockID int) error {
	_, err := db.DB.Exec(`
		DELETE FROM stock
		WHERE stock_id = $1
	`, stockID)

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
