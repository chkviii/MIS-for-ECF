-- 核心实体表

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    email VARCHAR(255) UNIQUE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 捐赠者表
CREATE TABLE IF NOT EXISTS donors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    donor_id VARCHAR(20) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    zip_code VARCHAR(20),
    country VARCHAR(100),
    donor_type VARCHAR(20) DEFAULT 'individual',
    age INTEGER,
    gender VARCHAR(10),
    occupation VARCHAR(100),
    communication_preference VARCHAR(50) DEFAULT 'email',
    tax_id VARCHAR(50),
    enrollment_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 志愿者表
CREATE TABLE IF NOT EXISTS volunteers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    volunteer_id VARCHAR(20) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
    emergency_contact_name VARCHAR(200),
    emergency_contact_phone VARCHAR(20),
    skills TEXT,
    certifications TEXT,
    availability TEXT,
    preferred_tasks TEXT,
    location_id INTEGER,
    background_check_status VARCHAR(20) DEFAULT 'pending',
    background_check_date DATE,
    training_completion_date DATE,
    total_volunteer_hours DECIMAL(10,2) DEFAULT 0,
    start_date DATE,
    status VARCHAR(20) DEFAULT 'active',
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (location_id) REFERENCES locations(id)
);

-- 员工表
CREATE TABLE IF NOT EXISTS employees (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id VARCHAR(20) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
    position VARCHAR(100),
    department VARCHAR(100),
    salary DECIMAL(10,2),
    hourly_rate DECIMAL(8,2),
    employment_type VARCHAR(20) DEFAULT 'full_time',
    hire_date DATE,
    termination_date DATE,
    location_id INTEGER,
    supervisor_id INTEGER,
    payroll_id VARCHAR(50),
    tax_id VARCHAR(50),
    benefits TEXT,
    performance_rating DECIMAL(3,2),
    work_hours_per_week INTEGER DEFAULT 40,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (location_id) REFERENCES locations(id),
    FOREIGN KEY (supervisor_id) REFERENCES employees(id)
);

-- 地点表
CREATE TABLE IF NOT EXISTS locations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    location_id VARCHAR(50) UNIQUE NOT NULL,
    country_code VARCHAR(3),
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50),
    parent_location_id INTEGER,
    address TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_location_id) REFERENCES locations(id)
);

-- 项目表
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    project_type VARCHAR(50),
    budget DECIMAL(12,2),
    actual_cost DECIMAL(12,2) DEFAULT 0,
    funding_source VARCHAR(100),
    location_id INTEGER,
    start_date DATE,
    end_date DATE,
    expected_completion_date DATE,
    status VARCHAR(20) DEFAULT 'planning',
    priority VARCHAR(20) DEFAULT 'medium',
    beneficiaries_count INTEGER DEFAULT 0,
    success_metrics TEXT,
    project_manager_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (location_id) REFERENCES locations(id),
    FOREIGN KEY (project_manager_id) REFERENCES employees(id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_donors_email ON donors(email);
CREATE INDEX IF NOT EXISTS idx_volunteers_email ON volunteers(email);
CREATE INDEX IF NOT EXISTS idx_employees_email ON employees(email);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
