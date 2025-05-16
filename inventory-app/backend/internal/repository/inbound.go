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
			i.inbound_id,
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

func AddInbound (inb model.Inbound) error {
		query:=`
		INSERT INTO inbound (item_id, supplier_id, quantity, received_at, received_by, warehouse_id)
					VALUES ($1, $2, $3, COALESCE($4, now()), $5, $6)
		`
		_, err := db.DB.Exec(query, inb.ItemID, inb.SupplierID, inb.Quantity, inb.ReceivedAt,inb.ReceivedBy, inb.WarehouseID)
		if err != nil {
			log.Println("❌ Произошла ошибка при добавлении поставок: ", err)
		}
		return err
}

func GetInboundDetailsByDate(date string)([]model.InboundDetails,error){
	var items []model.InboundDetails
	query := `
						SELECT
					   	i.inbound_id,
							to_char(i.received_at, 'DD.MM.YYYY') AS date, -- вот так!
							it.name AS name,
							it.sku AS sku,
							s.name AS supplier,
							i.quantity AS quantity,
							w.name AS warehouse
					FROM inbound i
					JOIN item it ON i.item_id = it.item_id
					JOIN supplier s ON i.supplier_id = s.supplier_id
					JOIN warehouse w ON i.warehouse_id = w.warehouse_id
					WHERE DATE(i.received_at) = $1
					ORDER BY i.received_at DESC;
	`
	err := db.DB.Select(&items, query,date)
	if err != nil {
		log.Println("❌ Произошла ошибка при получении поставок", err)
	}
	return items, err
}

func DeleteInbound(inboundId int) error {
		query := `
		DELETE FROM inbound WHERE inbound_id = $1;
		`
		_, err := db.DB.Exec(query, inboundId)
		if err != nil {
			log.Println("❌ Произошла ошибка при удалении поставки", err)
		}
		return err
}

func EditInbound(inb model.Inbound) error {
	query := `
			UPDATE inbound
			SET item_id = $1,
					supplier_id = $2,
					warehouse_id = $3,
					quantity = $4,
					received_at = $5
			WHERE inbound_id = $6;
	`
	_, err := db.DB.Exec(query, inb.ItemID, inb.SupplierID, inb.WarehouseID, inb.Quantity, inb.ReceivedAt, inb.InboundID)
	if err != nil {
		log.Println("❌ Произошла ошибка при изменении данных поставки", err)
	}
	return err
}
