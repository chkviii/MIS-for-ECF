-- 财务管理表 - 交易、捐赠、基金、支出、采购、薪资

-- 交易表
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id VARCHAR(50) UNIQUE NOT NULL,
    transaction_record TEXT,
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    from_currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    to_currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    from_entity VARCHAR(200),
    to_entity VARCHAR(200),
    transaction_date DATE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 捐赠记录表
CREATE TABLE IF NOT EXISTS donations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    donation_id VARCHAR(20) UNIQUE NOT NULL,
    donor_id INTEGER NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    donation_type VARCHAR(20) NOT NULL,
    category VARCHAR(20) NOT NULL,
    project_id INTEGER,
    donation_date DATE NOT NULL,
    transaction_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (donor_id) REFERENCES donors(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

-- 基金管理表
CREATE TABLE IF NOT EXISTS funds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fund_id VARCHAR(20) UNIQUE NOT NULL,
    donor_id INTEGER,
    project_id INTEGER,
    transaction_id INTEGER,
    name VARCHAR(200) NOT NULL,
    fund_type VARCHAR(20) NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    current_balance DECIMAL(12,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (donor_id) REFERENCES donors(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

-- 支出记录表
CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expense_id VARCHAR(20) UNIQUE NOT NULL,
    fund_id INTEGER NOT NULL,
    project_id INTEGER,
    employee_id INTEGER,
    transaction_id INTEGER,
    description TEXT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    expense_date DATE NOT NULL,
    approval_status VARCHAR(20) DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (fund_id) REFERENCES funds(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

-- 采购表
CREATE TABLE IF NOT EXISTS purchases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_id VARCHAR(50) UNIQUE NOT NULL,
    transaction_id INTEGER,
    total_spent DECIMAL(12,2) NOT NULL,
    supplier_name VARCHAR(200),
    purchase_date DATE,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

-- 薪资表
CREATE TABLE IF NOT EXISTS payrolls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id INTEGER NOT NULL,
    employee_id INTEGER NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    pay_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id)
);

-- 创建索引（优化查询性能）
CREATE INDEX IF NOT EXISTS idx_donations_donor ON donations(donor_id);
CREATE INDEX IF NOT EXISTS idx_donations_project ON donations(project_id);
CREATE INDEX IF NOT EXISTS idx_donations_date ON donations(donation_date);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(transaction_date);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_expenses_fund ON expenses(fund_id);
CREATE INDEX IF NOT EXISTS idx_expenses_project ON expenses(project_id);
CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(expense_date);
CREATE INDEX IF NOT EXISTS idx_funds_status ON funds(status);
CREATE INDEX IF NOT EXISTS idx_payrolls_employee ON payrolls(employee_id);
CREATE INDEX IF NOT EXISTS idx_payrolls_date ON payrolls(pay_date);
