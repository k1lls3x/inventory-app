-- Описание: база данных для управления складом

-- Пользователи
CREATE TABLE IF NOT EXISTS "User" (
    user_id     SERIAL PRIMARY KEY,
    username    VARCHAR(50) UNIQUE NOT NULL,
    password    TEXT NOT NULL,
    full_name   VARCHAR(100),
    role        VARCHAR(20) CHECK (role IN ('admin', 'manager', 'worker')) NOT NULL
);

-- Товары
CREATE TABLE IF NOT EXISTS Item (
    item_id         SERIAL PRIMARY KEY,
    sku             VARCHAR(50) UNIQUE NOT NULL,
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    uom             VARCHAR(20) NOT NULL,
    reorder_level   INTEGER DEFAULT 0,
    reorder_qty     INTEGER DEFAULT 0,
    price           NUMERIC(12,2),
    cost            NUMERIC(12,2)
);

-- Поставщики
CREATE TABLE IF NOT EXISTS Supplier (
    supplier_id     SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    contact_info    TEXT
);

-- Поступления
CREATE TABLE IF NOT EXISTS Inbound (
    inbound_id      SERIAL PRIMARY KEY,
    item_id         INTEGER REFERENCES Item(item_id),
    supplier_id     INTEGER REFERENCES Supplier(supplier_id),
    quantity        INTEGER NOT NULL,
    received_at     TIMESTAMP DEFAULT now(),
    received_by     INTEGER REFERENCES "User"(user_id)
);

-- Отгрузки
CREATE TABLE IF NOT EXISTS Outbound (
    outbound_id     SERIAL PRIMARY KEY,
    item_id         INTEGER REFERENCES Item(item_id),
    quantity        INTEGER NOT NULL,
    shipped_at      TIMESTAMP DEFAULT now(),
    shipped_by      INTEGER REFERENCES "User"(user_id),
    destination     TEXT
);

-- Остатки
CREATE TABLE IF NOT EXISTS Stock (
    stock_id        SERIAL PRIMARY KEY,
    item_id         INTEGER REFERENCES Item(item_id),
    quantity        INTEGER NOT NULL DEFAULT 0
);

-- Склады
CREATE TABLE IF NOT EXISTS Warehouse (
    warehouse_id   SERIAL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    location       TEXT
);

-- Добавление связи со складами (если не добавлено ранее)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='inbound' AND column_name='warehouse_id'
    ) THEN
        ALTER TABLE Inbound ADD COLUMN warehouse_id INTEGER REFERENCES Warehouse(warehouse_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='outbound' AND column_name='warehouse_id'
    ) THEN
        ALTER TABLE Outbound ADD COLUMN warehouse_id INTEGER REFERENCES Warehouse(warehouse_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='stock' AND column_name='warehouse_id'
    ) THEN
        ALTER TABLE Stock ADD COLUMN warehouse_id INTEGER REFERENCES Warehouse(warehouse_id);
    END IF;
END$$;

-- Удаление старого ограничения, если есть, и установка нового
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'stock' AND constraint_type = 'UNIQUE' AND constraint_name = 'stock_item_id_key'
    ) THEN
        ALTER TABLE Stock DROP CONSTRAINT stock_item_id_key;
    END IF;

    BEGIN
        ALTER TABLE Stock ADD UNIQUE (item_id, warehouse_id);
    EXCEPTION WHEN duplicate_object THEN
        -- уже добавлено — игнорируем
    END;
END$$;
