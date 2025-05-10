package main

import (
				"fmt"
				"log"
				"path/filepath"
				"github.com/joho/godotenv"
				"inventory-app/backend/internal/db"
				"inventory-app/backend/internal/model"
)

func main() {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		log.Fatalf("❌ .env не загружен: %v", err)
	}

	db.Init()
	db.DB.Exec(`TRUNCATE inbound, stock, item, warehouse RESTART IDENTITY CASCADE`)
	// Вставляем 1 склад
	_, err = db.DB.Exec(`
		INSERT INTO warehouse (name, location) VALUES
		('Центральный склад', 'Москва')
	`)
	if err != nil {
		log.Fatalf("❌ Ошибка вставки склада: %v", err)
	}

	// Вставляем товары
	_, err = db.DB.Exec(`
		INSERT INTO item (sku, name, uom, reorder_level, reorder_qty, price, cost) VALUES
		('SKU001', 'Монитор 24"', 'шт', 5, 10, 12999.00, 11000.00),
		('SKU002', 'Клавиатура', 'шт', 10, 20, 3499.00, 2500.00),
		('SKU003', 'Мышь', 'шт', 15, 30, 1999.00, 1400.00),
		('SKU004', 'Флешка 32GB', 'шт', 10, 25, 599.00, 450.00),
		('SKU005', 'Коврик для мыши', 'шт', 20, 50, 299.00, 200.00)
	`)
	if err != nil {
		log.Fatalf("❌ Ошибка вставки товаров: %v", err)
	}

	// Вставка остатков
	_, err = db.DB.Exec(`
		INSERT INTO stock (item_id, quantity, warehouse_id) VALUES
		(1, 10, 1),
		(2, 25, 1),
		(3, 18, 1),
		(4, 5, 1),
		(5, 9, 1)
	`)
	if err != nil {
		log.Fatalf("❌ Ошибка вставки остатков: %v", err)
	}

	// Вставка поступлений
	_, err = db.DB.Exec(`
		INSERT INTO inbound (item_id, supplier_id, quantity, received_by, warehouse_id)
		VALUES
		(1, 1, 10, 2, 1),
		(2, 1, 25, 2, 1),
		(3, 2, 18, 2, 1)
	`)
	if err != nil {
		log.Fatalf("❌ Ошибка вставки поступлений: %v", err)
	}

	log.Println("✅ Все данные успешно вставлены")

	// Проверим вывод
	var items []model.Item
	err = db.DB.Select(&items, `SELECT * FROM item`)
	if err != nil {
		log.Fatalf("❌ Ошибка выборки: %v", err)
	}

	fmt.Println("📦 Список товаров:")
	for _, item := range items {
		fmt.Printf("- [%s] %s — %.2f\n", item.SKU, item.Name, item.Price)
	}
}
