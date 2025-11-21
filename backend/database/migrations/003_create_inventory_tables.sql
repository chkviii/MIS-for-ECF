-- 库存和礼品管理表

-- 库存表
CREATE TABLE IF NOT EXISTS inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    inventory_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(100),
    subcategory VARCHAR(50),
    sku VARCHAR(50) UNIQUE,
    purchase_amount DECIMAL(10,2),
    remain_amount INTEGER DEFAULT 0,
    current_stock INTEGER DEFAULT 0,
    minimum_stock_level INTEGER DEFAULT 0,
    maximum_stock_level INTEGER,
    purchase_id INTEGER,
    location_id INTEGER,
    unit_cost DECIMAL(10,2),
    total_value DECIMAL(10,2),
    supplier_name VARCHAR(200),
    supplier_contact TEXT,
    storage_location VARCHAR(100),
    depreciation DECIMAL(5,2) DEFAULT 0.00,
    expiration_date DATE,
    last_inventory_date DATE,
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
    description TEXT,
    category VARCHAR(50),
    unit_cost DECIMAL(8,2),
    tax_deductible_value DECIMAL(8,2),
    requires_inventory BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 礼品记录表
CREATE TABLE IF NOT EXISTS gifts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    gift_id VARCHAR(20) UNIQUE NOT NULL,
    donor_id INTEGER,
    donation_id INTEGER,
    gift_type_id INTEGER NOT NULL,
    quantity INTEGER DEFAULT 1,
    unit_value DECIMAL(8,2),
    total_value DECIMAL(10,2),
    is_free BOOLEAN DEFAULT FALSE,
    distribution_method VARCHAR(20),
    distribution_status VARCHAR(20) DEFAULT 'pending',
    recipient_name VARCHAR(200),
    recipient_address TEXT,
    tracking_number VARCHAR(100),
    distributed_date DATE,
    distributed_by INTEGER,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (donor_id) REFERENCES donors(id),
    FOREIGN KEY (donation_id) REFERENCES donations(id),
    FOREIGN KEY (gift_type_id) REFERENCES gift_types(id),
    FOREIGN KEY (distributed_by) REFERENCES employees(id)
);

-- 库存交易记录表
CREATE TABLE IF NOT EXISTS inventory_transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id VARCHAR(20) UNIQUE NOT NULL,
    inventory_id INTEGER NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    quantity_change INTEGER NOT NULL,
    unit_cost DECIMAL(8,2),
    total_cost DECIMAL(10,2),
    reference_type VARCHAR(20),
    reference_id INTEGER,
    notes TEXT,
    processed_by INTEGER,
    transaction_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inventory_id) REFERENCES inventory(id),
    FOREIGN KEY (processed_by) REFERENCES employees(id)
);

-- 配送表
CREATE TABLE IF NOT EXISTS deliveries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    delivery_id VARCHAR(50) UNIQUE NOT NULL,
    inventory_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    project_id INTEGER,
    delivery_date DATE,
    status VARCHAR(20) DEFAULT 'pending',
    location_id INTEGER,
    recipient_name VARCHAR(200),
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inventory_id) REFERENCES inventory(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (location_id) REFERENCES locations(id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_inventory_category ON inventory(category);
CREATE INDEX IF NOT EXISTS idx_inventory_status ON inventory(status);
CREATE INDEX IF NOT EXISTS idx_gifts_donor ON gifts(donor_id);
CREATE INDEX IF NOT EXISTS idx_deliveries_project ON deliveries(project_id);
