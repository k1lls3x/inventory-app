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

	query :=`WITH period AS (
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
