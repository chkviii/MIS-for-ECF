-- 库存和礼品管理表

-- 库存表
CREATE TABLE IF NOT EXISTS inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    inventory_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(100),
    purchase_id INTEGER,
    location_id INTEGER,
    current_stock INTEGER DEFAULT 0,
    unit_cost DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'available',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (purchase_id) REFERENCES purchases(id),
    FOREIGN KEY (location_id) REFERENCES locations(id)
);

-- 礼品类型表
CREATE TABLE IF NOT EXISTS gift_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    unit_cost DECIMAL(8,2),
    requires_inventory BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 礼品记录表
CREATE TABLE IF NOT EXISTS gifts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    gift_id VARCHAR(20) UNIQUE NOT NULL,
    donation_id INTEGER,
    delivery_id INTEGER,
    gift_type_id INTEGER NOT NULL,
    quantity INTEGER DEFAULT 1,
    total_value DECIMAL(10,2),
    distribution_status VARCHAR(20) DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (donation_id) REFERENCES donations(id),
    FOREIGN KEY (delivery_id) REFERENCES deliveries(id),
    FOREIGN KEY (gift_type_id) REFERENCES gift_types(id)
);

-- 库存交易记录表
CREATE TABLE IF NOT EXISTS inventory_transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id VARCHAR(20) UNIQUE NOT NULL,
    inventory_id INTEGER NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    quantity_change INTEGER NOT NULL,
    transaction_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inventory_id) REFERENCES inventory(id)
);

-- 配送表
CREATE TABLE IF NOT EXISTS deliveries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    delivery_id VARCHAR(50) UNIQUE NOT NULL,
    inventory_id INTEGER NOT NULL,
    project_id INTEGER,
    location_id INTEGER,
    quantity INTEGER NOT NULL,
    recipient_name VARCHAR(200),
    delivery_date DATE,
    status VARCHAR(20) DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inventory_id) REFERENCES inventory(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (location_id) REFERENCES locations(id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_inventory_category ON inventory(category);
CREATE INDEX IF NOT EXISTS idx_inventory_status ON inventory(status);
CREATE INDEX IF NOT EXISTS idx_gifts_delivery ON gifts(delivery_id);
CREATE INDEX IF NOT EXISTS idx_deliveries_project ON deliveries(project_id);
