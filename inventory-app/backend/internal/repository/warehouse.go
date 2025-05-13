package repository
import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
)
func GetWarehouses() ([]model.Warehouse, error) {
	var warehouses []model.Warehouse
	query := `SELECT warehouse_id, name, location FROM Warehouse ORDER BY name`
	err := db.DB.Select(&warehouses, query)
	if err != nil {
		log.Println("❌ Ошибка при получении списка складов:", err)
	}
	return warehouses, nil
}
