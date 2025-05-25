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

func AddWarehouse(warehouse model.Warehouse) error {
	query:= `
	INSERT INTO warehouse (name, location)
	VALUES (:name, :location)
	`
	_, err := db.DB.NamedExec(query, warehouse)

	if err != nil {
		log.Println("❌ Ошибка при добавлении склада:", err)
	}
	return err
}

func EditWarehouse(warehouse model.Warehouse) error {
	query := `
		UPDATE warehouse
		SET name = :name,
		    location = :location
		WHERE warehouse_id = :warehouse_id
	`
	_, err := db.DB.NamedExec(query, warehouse)
	if err != nil {
		log.Println("❌ Ошибка при редактировании склада:", err)
	}
	return err
}

