-- 关联表

-- 志愿者-项目关联表
CREATE TABLE IF NOT EXISTS volunteer_projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    volunteer_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,
    role VARCHAR(100),
    contract_start DATE,
    contract_end DATE,
    work_unit VARCHAR(50),
    total_amount DECIMAL(10,2),
    contract_date DATE,
    contract_detail TEXT,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (volunteer_id) REFERENCES volunteers(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    UNIQUE(volunteer_id, project_id)
);

-- 员工-项目关联表
CREATE TABLE IF NOT EXISTS employee_projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,
    title VARCHAR(100),
    start_date DATE,
    end_date DATE,
    work_unit VARCHAR(50),
    allocated_amount DECIMAL(10,2),
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    UNIQUE(employee_id, project_id)
);

-- 资金-项目关联表
CREATE TABLE IF NOT EXISTS fund_projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,
    fund_id INTEGER NOT NULL,
    allocated_amount DECIMAL(12,2) NOT NULL,
    allocation_date DATE DEFAULT (DATE('now')),
    purpose TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (fund_id) REFERENCES funds(id)
);

-- 捐赠-库存关联表（实物捐赠）
CREATE TABLE IF NOT EXISTS donation_inventory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    donor_id INTEGER NOT NULL,
    inventory_id INTEGER NOT NULL,
    donation_date DATE,
    project_id INTEGER,
    quantity INTEGER DEFAULT 1,
    estimated_value DECIMAL(10,2),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (donor_id) REFERENCES donors(id),
    FOREIGN KEY (inventory_id) REFERENCES inventory(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- 调度表
CREATE TABLE IF NOT EXISTS schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    schedule_id VARCHAR(20) UNIQUE NOT NULL,
    person_id INTEGER NOT NULL,
    person_type VARCHAR(20) NOT NULL,
    project_id INTEGER,
    shift_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    hours_worked DECIMAL(5,2),
    status VARCHAR(20) DEFAULT 'scheduled',
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_volunteer_projects_volunteer ON volunteer_projects(volunteer_id);
CREATE INDEX IF NOT EXISTS idx_volunteer_projects_project ON volunteer_projects(project_id);
CREATE INDEX IF NOT EXISTS idx_employee_projects_employee ON employee_projects(employee_id);
CREATE INDEX IF NOT EXISTS idx_employee_projects_project ON employee_projects(project_id);
CREATE INDEX IF NOT EXISTS idx_schedules_person ON schedules(person_id, person_type);
CREATE INDEX IF NOT EXISTS idx_schedules_date ON schedules(shift_date);
