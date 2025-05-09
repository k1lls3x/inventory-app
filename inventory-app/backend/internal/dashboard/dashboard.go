package dashboard

import (
	"inventory-app/backend/internal/db"
	"log"
)

type DashboardData struct {
	TotalStock    float64 `json:"total_stock"`
	ItemCount     int     `json:"item_count"`
	MonthlyOrders int     `json:"monthly_orders"`
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
		SELECT COUNT(*) FROM PurchaseOrder
		WHERE date_trunc('month', order_date) = date_trunc('month', CURRENT_DATE)
	`)
	if err != nil {
		log.Println("❌ Ошибка при получении monthly_orders:", err)
		return nil, err
	}

	return &result, nil
}
