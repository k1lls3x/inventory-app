package repository
import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
)
func GetInboundDetails()([]model.InboundDetails,error){
	var items []model.InboundDetails
	query:= `
		SELECT
    to_char(i.received_at, 'YYYY-MM-DD') AS date,
    it.name AS name,
    it.sku AS sku,
    s.name AS supplier,
    i.quantity AS quantity,
    w.name AS warehouse
FROM inbound i
JOIN item it       ON i.item_id = it.item_id
JOIN supplier s    ON i.supplier_id = s.supplier_id
JOIN warehouse w   ON i.warehouse_id = w.warehouse_id
ORDER BY i.received_at DESC;
	`
	err:= db.DB.Select(&items,query)

	if err != nil {
		log.Println("❌ Произошла ошибка при получении поставок: ", err)
	}

	return items,err
}
