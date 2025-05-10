package repository

import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
)

func AddStock( itemID, quantity, warehouseID int) error {
	_ ,err:= db.DB.Exec(`
											INSERT INTO Stock (item_id, quantity, warehouse_id)
											VALUES ($1, $2, $3)
											ON CONFLICT (item_id, warehouse_id) DO UPDATE SET quantity = EXCLUDED.quantity;
										`, itemID, quantity, warehouseID)
	return err
}

func FindStockByWarehouse(warehouseID int) ([]model.Stock, error) {
	var result []model.Stock
	err := db.DB.Select(&result, `SELECT * FROM Stock WHERE warehouse_id = $1`, warehouseID)
	return result, err
}

func ChangeStock(itemID, warehouseID, newQuantity int) error {
		_, err := db.DB.Exec(`
													UPDATE Stock
													SET quantity = $1
													WHERE item_id = $2 AND warehouse_id = $3
												`, newQuantity, itemID, warehouseID)
	if err != nil {
	log.Println("❌ Произошла ошибка при изменение данных: ", err)
	}
	return err
}

func RemoveStock(itemID, warehouseID int) error {
	_, err := db.DB.Exec(`
			DELETE FROM Stock
			WHERE item_id = $1 AND warehouse_id = $2
	`, itemID, warehouseID)
	return err
}

func GetStocks() ([]model.Stock,error){
		var result []model.Stock
		query:= `
			SELECT * FROM Stock
		`
		err := db.DB.Select(&result, query)
		if err != nil {
			return nil, err
		}
		return result,err
}

