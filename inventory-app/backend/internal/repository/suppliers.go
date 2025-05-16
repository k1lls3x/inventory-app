package repository

import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
)

func GetSuppliers() ([]model.Supplier, error){
		var result []model.Supplier
		query := `
				SELECT supplier_id, name, contact_info
		FROM supplier
		ORDER BY name;
		`
		err := db.DB.Select(&result, query)

		if err != nil {
			log.Println("❌ Ошибка при получении данных о поставщиках: ", err)
		}
		return result, err
}

