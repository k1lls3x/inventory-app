package dashboard

import (
	"inventory-app/backend/internal/db"
	"log"
	"inventory-app/backend/internal/model"
)

type DashboardData struct {
	TotalStock    float64 `json:"total_stock"`
	ItemCount     int     `json:"item_count"`
	MonthlyOrders int     `json:"monthly_orders"`
	NewItems      int     `json:"new_items"`
}
func LoadTurnoverByWarehouse() ([]model.ItemTurnoverByWarehouse, error) {
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
	err := db.DB.Select(&result, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}
// LoadTopItems возвращает топ-10 товаров по остатку
func LoadTopItems() ([]model.ItemWithStock, error) {
	query := `
	SELECT
		i.name,
		i.sku,
		'Main Warehouse' AS warehouse,  -- временно, если складов нет
		s.quantity
	FROM stock s
	JOIN item i ON i.item_id = s.item_id
	ORDER BY s.quantity DESC
	LIMIT 10;
	`

	var items []model.ItemWithStock
	err := db.DB.Select(&items, query)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func LoadDashboard() (*DashboardData, error) {
	var result DashboardData

	err := db.DB.Get(&result.TotalStock, `SELECT COALESCE(SUM(quantity), 0) FROM Stock`)
	if err != nil {
		log.Println("❌ Ошибка при получении total_stock:", err)
		return nil, err
	}

	err = db.DB.Get(&result.ItemCount, `SELECT COUNT(*) FROM Item`)
	if err != nil {
		log.Println("❌ Ошибка при получении item_count:", err)
		return nil, err
	}

	err = db.DB.Get(&result.MonthlyOrders, `
	SELECT COUNT(*) FROM inbound
	WHERE date_trunc('month', received_at) = date_trunc('month', CURRENT_DATE)
`)
	if err != nil {
		log.Println("❌ Ошибка при получении monthly_orders:", err)
		return nil, err
	}
	var newItems int
	err = db.DB.Get(&newItems, `
		SELECT COUNT(*) FROM item
		WHERE date_trunc('month', created_at) = date_trunc('month', CURRENT_DATE)
	`)
	if err != nil {
		log.Println("❌ Ошибка при получении количества новых товаров:", err)
		return nil, err
	}
	return &DashboardData{
		TotalStock:    result.TotalStock,
		ItemCount:     result.ItemCount,
		MonthlyOrders: result.MonthlyOrders,
		NewItems:      newItems,
	}, nil
}
