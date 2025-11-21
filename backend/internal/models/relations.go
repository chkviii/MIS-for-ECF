package models

import "time"

// VolunteerProject 志愿者-项目关联表
type VolunteerProject struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	VolunteerID    uint       `gorm:"not null" json:"volunteer_id"`
	ProjectID      uint       `gorm:"not null" json:"project_id"`
	Role           string     `gorm:"size:100" json:"role"`
	ContractStart  *time.Time `json:"contract_start"`
	ContractEnd    *time.Time `json:"contract_end"`
	WorkUnit       string     `gorm:"size:50" json:"work_unit"`
	TotalAmount    float64    `gorm:"type:decimal(10,2)" json:"total_amount"`
	ContractDate   *time.Time `json:"contract_date"`
	ContractDetail string     `json:"contract_detail"`
	Status         string     `gorm:"size:20;default:active" json:"status"`
	CreatedAt      time.Time  `json:"created_at"`

	// 关联
	Volunteer Volunteer `json:"volunteer,omitempty"`
	Project   Project   `json:"project,omitempty"`
}

// EmployeeProject 员工-项目关联表
type EmployeeProject struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	EmployeeID      uint       `gorm:"not null" json:"employee_id"`
	ProjectID       uint       `gorm:"not null" json:"project_id"`
	Title           string     `gorm:"size:100" json:"title"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	WorkUnit        string     `gorm:"size:50" json:"work_unit"`
	AllocatedAmount float64    `gorm:"type:decimal(10,2)" json:"allocated_amount"`
	LastUpdated     time.Time  `json:"last_updated"`
	CreatedAt       time.Time  `json:"created_at"`

	// 关联
	Employee Employee `json:"employee,omitempty"`
	Project  Project  `json:"project,omitempty"`
}

// FundProject 资金-项目关联表
type FundProject struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TransactionID   uint      `gorm:"not null" json:"transaction_id"`
	ProjectID       uint      `gorm:"not null" json:"project_id"`
	FundID          uint      `gorm:"not null" json:"fund_id"`
	AllocatedAmount float64   `gorm:"type:decimal(12,2);not null" json:"allocated_amount"`
	AllocationDate  time.Time `json:"allocation_date"`
	Purpose         string    `json:"purpose"`
	CreatedAt       time.Time `json:"created_at"`

	// 关联
	Transaction Transaction `json:"transaction,omitempty"`
	Project     Project     `json:"project,omitempty"`
	Fund        Fund        `json:"fund,omitempty"`
}

// DonationInventory 捐赠-库存关联表（实物捐赠）
type DonationInventory struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	DonorID        uint       `gorm:"not null" json:"donor_id"`
	InventoryID    uint       `gorm:"not null" json:"inventory_id"`
	DonationDate   *time.Time `json:"donation_date"`
	ProjectID      *uint      `json:"project_id"`
	Quantity       int        `gorm:"default:1" json:"quantity"`
	EstimatedValue float64    `gorm:"type:decimal(10,2)" json:"estimated_value"`
	CreatedAt      time.Time  `json:"created_at"`

	// 关联
	Donor     Donor      `json:"donor,omitempty"`
	Inventory Inventory  `json:"inventory,omitempty"`
	Project   *Project   `json:"project,omitempty"`
}

// Schedule 调度表
type Schedule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ScheduleID  string    `gorm:"size:20;unique;not null" json:"schedule_id"`
	PersonID    uint      `gorm:"not null" json:"person_id"`
	PersonType  string    `gorm:"size:20;not null" json:"person_type"`
	ProjectID   *uint     `json:"project_id"`
	ShiftDate   time.Time `json:"shift_date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	HoursWorked float64   `gorm:"type:decimal(5,2)" json:"hours_worked"`
	Status      string    `gorm:"size:20;default:scheduled" json:"status"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联
	Project *Project `json:"project,omitempty"`
}
