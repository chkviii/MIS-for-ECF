-- 财务相关表

-- 交易表
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id VARCHAR(50) UNIQUE NOT NULL,
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    from_entity VARCHAR(200),
    to_entity VARCHAR(200),
    description TEXT,
    reference_type VARCHAR(20),
    reference_id INTEGER,
    transaction_date DATE,
    fingerprint VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 捐赠记录表
CREATE TABLE IF NOT EXISTS donations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    donation_id VARCHAR(20) UNIQUE NOT NULL,
    donor_id INTEGER NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    donation_type VARCHAR(20) NOT NULL,
    donation_method VARCHAR(20),
    category VARCHAR(20) NOT NULL,
    subcategory VARCHAR(50),
    currency VARCHAR(3) DEFAULT 'USD',
    exchange_rate DECIMAL(8,4) DEFAULT 1.0000,
    tax_deductible BOOLEAN DEFAULT TRUE,
    tax_deduction_amount DECIMAL(10,2),
    receipt_number VARCHAR(50) UNIQUE,
    receipt_issued_date DATE,
    project_id INTEGER,
    campaign_id VARCHAR(50),
    tribute_type VARCHAR(20),
    tribute_name VARCHAR(200),
    tribute_notification TEXT,
    notes TEXT,
    donation_date DATE NOT NULL,
    processed_date DATE,
    acknowledgment_sent_date DATE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (donor_id) REFERENCES donors(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- 基金管理表
CREATE TABLE IF NOT EXISTS funds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fund_id VARCHAR(20) UNIQUE NOT NULL,
    donor_id INTEGER,
    name VARCHAR(200) NOT NULL,
    fund_type VARCHAR(20) NOT NULL,
    description TEXT,
    total_amount DECIMAL(12,2) NOT NULL,
    current_balance DECIMAL(12,2) DEFAULT 0,
    amount_left DECIMAL(12,2) NOT NULL,
    minimum_balance DECIMAL(12,2) DEFAULT 0,
    restrictions TEXT,
    established_date DATE,
    fund_date DATE,
    expiration_date DATE,
    currency VARCHAR(3) DEFAULT 'USD',
    purpose TEXT,
    transaction_id INTEGER,
    project_id INTEGER,
    fund_manager_id INTEGER,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (donor_id) REFERENCES donors(id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (fund_manager_id) REFERENCES employees(id)
);

-- 支出记录表
CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expense_id VARCHAR(20) UNIQUE NOT NULL,
    fund_id INTEGER NOT NULL,
    project_id INTEGER,
    employee_id INTEGER,
    vendor_name VARCHAR(200),
    description TEXT NOT NULL,
    expense_category VARCHAR(50),
    expense_type VARCHAR(50),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    payment_method VARCHAR(20),
    receipt_number VARCHAR(50),
    invoice_number VARCHAR(50),
    approval_status VARCHAR(20) DEFAULT 'pending',
    approved_by INTEGER,
    approval_date DATE,
    expense_date DATE NOT NULL,
    payment_date DATE,
    reimbursable BOOLEAN DEFAULT FALSE,
    reimbursed BOOLEAN DEFAULT FALSE,
    reimbursement_date DATE,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (fund_id) REFERENCES funds(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (approved_by) REFERENCES employees(id)
);

-- 采购表
CREATE TABLE IF NOT EXISTS purchases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_id VARCHAR(50) UNIQUE NOT NULL,
    total_spent DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    purchase_date DATE,
    transaction_id INTEGER,
    supplier_name VARCHAR(200),
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

-- 薪资表
CREATE TABLE IF NOT EXISTS payrolls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id INTEGER NOT NULL,
    employee_id INTEGER NOT NULL,
    pay_date DATE NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    pay_period_start DATE,
    pay_period_end DATE,
    deductions DECIMAL(10,2) DEFAULT 0.00,
    bonuses DECIMAL(10,2) DEFAULT 0.00,
    overtime_hours DECIMAL(5,2) DEFAULT 0.00,
    overtime_rate DECIMAL(8,2) DEFAULT 0.00,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_donations_donor ON donations(donor_id);
CREATE INDEX IF NOT EXISTS idx_donations_date ON donations(donation_date);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(transaction_date);
CREATE INDEX IF NOT EXISTS idx_expenses_project ON expenses(project_id);
CREATE INDEX IF NOT EXISTS idx_funds_status ON funds(status);
