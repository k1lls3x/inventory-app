CREATE TABLE "User" (
  user_id        SERIAL PRIMARY KEY,
  username       VARCHAR(50) UNIQUE NOT NULL,
  full_name      VARCHAR(100) NOT NULL,
  role           VARCHAR(30) NOT NULL
);

CREATE TABLE Supplier (
  supplier_id    SERIAL PRIMARY KEY,
  name           VARCHAR(100) NOT NULL,
  contact_info   TEXT
);

-- Основные справочники
CREATE TABLE Item (
  item_id        SERIAL PRIMARY KEY,
  sku            VARCHAR(50) UNIQUE NOT NULL,
  name           VARCHAR(200) NOT NULL,
  description    TEXT,
  uom            VARCHAR(20) NOT NULL,
  reorder_level  INTEGER DEFAULT 0,
  reorder_qty    INTEGER DEFAULT 0,
  price          NUMERIC(12,2),
  cost           NUMERIC(12,2)
);

CREATE TABLE Warehouse (
  warehouse_id       SERIAL PRIMARY KEY,
  name               VARCHAR(100) NOT NULL,
  location           TEXT,
  manager_user_id    INTEGER,
  CONSTRAINT fk_wh_manager FOREIGN KEY(manager_user_id) REFERENCES "User"(user_id)
);

-- Остатки и ячейки (опционально)
CREATE TABLE BinLocation (
  bin_id         SERIAL PRIMARY KEY,
  warehouse_id   INTEGER NOT NULL,
  code           VARCHAR(50) NOT NULL,
  description    TEXT,
  CONSTRAINT fk_bin_wh FOREIGN KEY(warehouse_id) REFERENCES Warehouse(warehouse_id)
);

CREATE TABLE Stock (
  stock_id       SERIAL PRIMARY KEY,
  item_id        INTEGER NOT NULL,
  warehouse_id   INTEGER NOT NULL,
  quantity       NUMERIC(12,2) NOT NULL DEFAULT 0,
  last_updated   TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_stock_item FOREIGN KEY(item_id) REFERENCES Item(item_id),
  CONSTRAINT fk_stock_wh   FOREIGN KEY(warehouse_id) REFERENCES Warehouse(warehouse_id),
  UNIQUE(item_id, warehouse_id)
);

-- Движения товара
CREATE TABLE StockTransaction (
  tx_id            SERIAL PRIMARY KEY,
  item_id          INTEGER NOT NULL,
  warehouse_id     INTEGER NOT NULL,
  tx_type          VARCHAR(20) NOT NULL,  -- 'RECEIPT','ISSUE','TRANSFER'
  quantity         NUMERIC(12,2) NOT NULL,
  tx_date          TIMESTAMP NOT NULL DEFAULT NOW(),
  related_tx_id    INTEGER,               -- для перемещений, связанных транзакций
  user_id          INTEGER,
  comment          TEXT,
  CONSTRAINT fk_tx_item  FOREIGN KEY(item_id) REFERENCES Item(item_id),
  CONSTRAINT fk_tx_wh    FOREIGN KEY(warehouse_id) REFERENCES Warehouse(warehouse_id),
  CONSTRAINT fk_tx_user  FOREIGN KEY(user_id) REFERENCES "User"(user_id),
  CONSTRAINT fk_tx_self  FOREIGN KEY(related_tx_id) REFERENCES StockTransaction(tx_id)
);

-- Закупочные заказы
CREATE TABLE PurchaseOrder (
  po_id           SERIAL PRIMARY KEY,
  supplier_id     INTEGER NOT NULL,
  order_date      DATE NOT NULL DEFAULT CURRENT_DATE,
  expected_date   DATE,
  status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
  CONSTRAINT fk_po_supplier FOREIGN KEY(supplier_id) REFERENCES Supplier(supplier_id)
);

CREATE TABLE PurchaseOrderLine (
  poline_id       SERIAL PRIMARY KEY,
  po_id           INTEGER NOT NULL,
  item_id         INTEGER NOT NULL,
  quantity        NUMERIC(12,2) NOT NULL,
  unit_price      NUMERIC(12,2),
  CONSTRAINT fk_pol_po   FOREIGN KEY(po_id) REFERENCES PurchaseOrder(po_id),
  CONSTRAINT fk_pol_item FOREIGN KEY(item_id) REFERENCES Item(item_id)
);
