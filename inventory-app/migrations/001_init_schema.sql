-- Пользователи
CREATE TABLE IF NOT EXISTS "User" (
    user_id     SERIAL PRIMARY KEY,
    username    VARCHAR(50) UNIQUE NOT NULL,
    password_hash     TEXT NOT NULL,
    full_name   VARCHAR(100),
    role        VARCHAR(20) CHECK (role IN ('admin', 'manager', 'worker')) NOT NULL
);

-- Товары
CREATE TABLE IF NOT EXISTS item (
    item_id        SERIAL PRIMARY KEY,
    sku            VARCHAR(50) UNIQUE NOT NULL,
    name           VARCHAR(200) NOT NULL,
    description    TEXT,
    uom            VARCHAR(20) NOT NULL,
    reorder_level  INTEGER DEFAULT 0,
    reorder_qty    INTEGER DEFAULT 0,
    price          NUMERIC(12,2),
    cost           NUMERIC(12,2),
    created_at     TIMESTAMP DEFAULT now()
);

-- Поставщики
CREATE TABLE IF NOT EXISTS supplier (
    supplier_id    SERIAL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    contact_info   TEXT
);

-- Склады
CREATE TABLE IF NOT EXISTS warehouse (
    warehouse_id   SERIAL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    location       TEXT
);

-- Поступления
-- Поступления (inbound)
CREATE TABLE IF NOT EXISTS inbound (
    inbound_id     SERIAL PRIMARY KEY,
    item_id        INTEGER REFERENCES item(item_id),
    supplier_id    INTEGER REFERENCES supplier(supplier_id),
    quantity       INTEGER NOT NULL,
    received_at    TIMESTAMP DEFAULT now(),
    warehouse_id   INTEGER REFERENCES warehouse(warehouse_id)
);

-- Отгрузки (outbound)
CREATE TABLE IF NOT EXISTS outbound (
    outbound_id    SERIAL PRIMARY KEY,
    item_id        INTEGER REFERENCES item(item_id),
    quantity       INTEGER NOT NULL,
    shipped_at     TIMESTAMP DEFAULT now(),
    destination    TEXT,
    warehouse_id   INTEGER REFERENCES warehouse(warehouse_id)
);


-- Остатки
CREATE TABLE IF NOT EXISTS stock (
    stock_id       SERIAL PRIMARY KEY,
    item_id        INTEGER REFERENCES item(item_id),
    quantity       INTEGER NOT NULL DEFAULT 0,
    warehouse_id   INTEGER REFERENCES warehouse(warehouse_id),
    last_updated   TIMESTAMP DEFAULT now()
);

-- Гарантия уникальности
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'stock'
          AND constraint_type = 'UNIQUE'
          AND constraint_name = 'stock_item_id_key'
    ) THEN
        ALTER TABLE stock DROP CONSTRAINT stock_item_id_key;
    END IF;

    BEGIN
        ALTER TABLE stock ADD UNIQUE (item_id, warehouse_id);
    EXCEPTION WHEN duplicate_object THEN
        -- constraint уже существует — игнорируем
    END;
END$$;

