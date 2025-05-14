package repository

import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
)

func GetItems() ([]model.Item, error) {
	var items []model.Item
	query:= `SELECT * FROM Item ORDER BY name`
	err := db.DB.Select(&items,query)
		if err != nil {
			return nil,err
		}
	return items, err
}

func GetItemBriefList() ([]model.ItemBrief, error) {
	var items []model.ItemBrief
	query := `SELECT item_id, name, sku FROM item ORDER BY name`
	err := db.DB.Select(&items, query)
	if err != nil {
		return nil, err
	}
	return items, nil
}

 func UpdateItem(item model.Item) error  {
	query := `
		UPDATE Item
			SET
				sku = :sku,
				name = :name,
				description = :description,
				uom = :uom,
				reorder_level = :reorder_level,
				reorder_qty = :reorder_qty,
				price = :price,
				cost = :cost
			WHERE item_id = :item_id
	`
	_, err := db.DB.NamedExec(query,item)
		if err != nil {
			log.Println("❌ Произошла ошибка при обновлении товара: ", err)
		}
		return err
}

func RemoveItem (sku string) error {
		_ , err := db.DB.Exec(`DELETE FROM item WHERE sku = $1`,sku)
		if err != nil {
			log.Println("❌ Произошла ошибка при удалении товара: ", err)
		}
	return err
}

func AddItem(item model.Item) error {
	_, err := db.DB.NamedExec(`
		INSERT INTO item (
			sku, name, description, uom,
			reorder_level, reorder_qty, price, cost, created_at
		) VALUES (
			:sku, :name, :description, :uom,
			:reorder_level, :reorder_qty, :price, :cost, NOW()
		)
	`, item)

	if err != nil {
		log.Println("❌ Ошибка при добавлении товара:", err)
	}
	return err
}

func FindItems(filter model.ItemFilter) ([]model.Item, error) {
	query := `
		SELECT * FROM item
		WHERE
			(:sku IS NULL OR sku = :sku) AND
			(:name IS NULL OR name ILIKE '%' || :name || '%') AND
			(:uom IS NULL OR uom = :uom) AND
			(:price_min IS NULL OR price >= :price_min) AND
			(:price_max IS NULL OR price <= :price_max)
		ORDER BY created_at DESC
	`

	rows, err := db.DB.NamedQuery(query, filter)
	if err != nil {
		log.Println("❌ Ошибка при поиске товаров:", err)
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.StructScan(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
