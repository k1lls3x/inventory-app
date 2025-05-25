package repository
import (
	"inventory-app/backend/internal/db"
	"inventory-app/backend/internal/model"
	"log"
)

func GetAllMovementsThisMonth() ([]model.Movement, error) {
	query := `
		 SELECT
  i.inbound_id      AS movement_id,
  'inbound'         AS type,
  i.item_id         AS item_id,       -- вернуть item_id
  it.name           AS item_name,
  i.quantity,
  i.received_at     AS date,
  i.warehouse_id    AS warehouse_id,  -- вернуть warehouse_id
  w.name            AS warehouse_name,
  s.supplier_id     AS supplier_id,   -- вернуть supplier_id
  s.name            AS supplier_name,
  NULL              AS destination,
  NULL              AS shipped_at     -- чтобы Go не жаловался на пустой столбец
FROM inbound i
JOIN item      it ON i.item_id      = it.item_id
JOIN warehouse w  ON i.warehouse_id = w.warehouse_id
LEFT JOIN supplier s ON i.supplier_id = s.supplier_id
WHERE i.received_at >= date_trunc('month', CURRENT_DATE)

UNION ALL

SELECT
  o.outbound_id     AS movement_id,
  'outbound'        AS type,
  o.item_id         AS item_id,
  it.name           AS item_name,
  o.quantity,
  o.shipped_at      AS date,
  o.warehouse_id    AS warehouse_id,
  w.name            AS warehouse_name,
  NULL              AS supplier_id,
  NULL              AS supplier_name,
  o.destination,
  o.shipped_at      AS shipped_at
FROM outbound o
JOIN item      it ON o.item_id      = it.item_id
JOIN warehouse w  ON o.warehouse_id = w.warehouse_id
WHERE o.shipped_at >= date_trunc('month', CURRENT_DATE)

ORDER BY date DESC;
	`
	var movements []model.Movement
	err := db.DB.Select(&movements, query)
	if err != nil {
			log.Println("❌ Ошибка при получении движений:", err)
			return nil, err
	}
	return movements, nil
}
