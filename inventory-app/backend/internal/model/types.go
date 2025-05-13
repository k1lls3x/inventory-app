package model
import "time"

type User struct {
	UserID   int    `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`
	FullName string `db:"full_name" json:"full_name"`
	Role     string `db:"role" json:"role"`
}

type Item struct {
	ItemID       int     `db:"item_id" json:"item_id"`
	SKU          string  `db:"sku" json:"sku"`
	Name         string  `db:"name" json:"name"`
	Description  *string `db:"description" json:"description"`
	UOM          string  `db:"uom" json:"uom"`
	ReorderLevel int     `db:"reorder_level" json:"reorder_level"`
	ReorderQty   int     `db:"reorder_qty" json:"reorder_qty"`
	Price        float64 `db:"price" json:"price"`
	Cost         float64 `db:"cost" json:"cost"`
	CreatedAt    time.Time `db:"created_at"`
}

type Supplier struct {
	SupplierID  int     `db:"supplier_id" json:"supplier_id"`
	Name        string  `db:"name" json:"name"`
	ContactInfo *string `db:"contact_info" json:"contact_info"`
}

type Inbound struct {
	InboundID    int       `db:"inbound_id" json:"inbound_id"`
	ItemID       int       `db:"item_id" json:"item_id"`
	SupplierID   int       `db:"supplier_id" json:"supplier_id"`
	Quantity     int       `db:"quantity" json:"quantity"`
	ReceivedAt   time.Time `db:"received_at" json:"received_at"`
	ReceivedBy   int       `db:"received_by" json:"received_by"`
	WarehouseID  int       `db:"warehouse_id" json:"warehouse_id"`
}

type Outbound struct {
	OutboundID   int       `db:"outbound_id" json:"outbound_id"`
	ItemID       int       `db:"item_id" json:"item_id"`
	Quantity     int       `db:"quantity" json:"quantity"`
	ShippedAt    time.Time `db:"shipped_at" json:"shipped_at"`
	ShippedBy    int       `db:"shipped_by" json:"shipped_by"`
	Destination  *string   `db:"destination" json:"destination"`
	WarehouseID  int       `db:"warehouse_id" json:"warehouse_id"`
}

type Stock struct {
	StockID     int `db:"stock_id" json:"stock_id"`
	ItemID      int `db:"item_id" json:"item_id"`
	Quantity    int `db:"quantity" json:"quantity"`
	WarehouseID int `db:"warehouse_id" json:"warehouse_id"`
	LastUpdated time.Time `db:"last_updated" json:"last_updated"`
}

type Warehouse struct {
	WarehouseID int     `db:"warehouse_id" json:"warehouse_id"`
	Name        string  `db:"name" json:"name"`
	Location    *string `db:"location" json:"location"`
}



//---------------------------------------------------------------\\
type ItemTurnoverByWarehouse struct {
	Warehouse string `db:"warehouse" json:"warehouse"`
	Name      string `db:"name" json:"name"`
	SKU       string `db:"sku" json:"sku"`
	Received  int    `db:"received" json:"received"`
	Shipped   int    `db:"shipped" json:"shipped"`
}

type ItemWithStock struct {
	Name      string `json:"name"`
	SKU       string `json:"sku"`
	Warehouse string `json:"warehouse"`
	Quantity  int    `json:"quantity"`
}

type ItemFilter struct {
	SKU       *string
	Name      *string
	UOM       *string
	PriceMin  *float64
	PriceMax  *float64
}
