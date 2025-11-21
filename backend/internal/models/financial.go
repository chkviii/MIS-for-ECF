package models

import "time"

// Transaction 交易表
type Transaction struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	TransactionID   string     `gorm:"size:50;unique;not null" json:"transaction_id"`
	Type            string     `gorm:"size:20;not null" json:"type"`
	Amount          float64    `gorm:"type:decimal(12,2);not null" json:"amount"`
	FromEntity      string     `gorm:"size:200" json:"from_entity"`
	ToEntity        string     `gorm:"size:200" json:"to_entity"`
	TransactionDate *time.Time `json:"transaction_date"`
	CreatedAt       time.Time  `json:"created_at"`

	// 关联
	Purchases []Purchase `json:"purchases,omitempty"`
	Payrolls  []Payroll  `json:"payrolls,omitempty"`
}

// Donation 捐赠记录表
type Donation struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	DonationID   string     `gorm:"size:20;unique;not null" json:"donation_id"`
	DonorID      uint       `gorm:"not null" json:"donor_id"`
	Amount       float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	DonationType string     `gorm:"size:20;not null" json:"donation_type"`
	Category     string     `gorm:"size:20;not null" json:"category"`
	ProjectID    *uint      `json:"project_id"`
	DonationDate time.Time  `json:"donation_date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// 关联
	Donor   Donor    `json:"donor,omitempty"`
	Project *Project `json:"project,omitempty"`
	Gifts   []Gift   `json:"gifts,omitempty"`
}

// Fund 基金管理表
type Fund struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	FundID         string    `gorm:"size:20;unique;not null" json:"fund_id"`
	DonorID        *uint     `json:"donor_id"`
	ProjectID      *uint     `json:"project_id"`
	Name           string    `gorm:"size:200;not null" json:"name"`
	FundType       string    `gorm:"size:20;not null" json:"fund_type"`
	TotalAmount    float64   `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	CurrentBalance float64   `gorm:"type:decimal(12,2);default:0" json:"current_balance"`
	Status         string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联
	Donor    *Donor    `json:"donor,omitempty"`
	Project  *Project  `json:"project,omitempty"`
	Expenses []Expense `json:"expenses,omitempty"`
}

// Expense 支出记录表
type Expense struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ExpenseID      string     `gorm:"size:20;unique;not null" json:"expense_id"`
	FundID         uint       `gorm:"not null" json:"fund_id"`
	ProjectID      *uint      `json:"project_id"`
	EmployeeID     *uint      `json:"employee_id"`
	Description    string     `gorm:"not null" json:"description"`
	Amount         float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	ExpenseDate    time.Time  `json:"expense_date"`
	ApprovalStatus string     `gorm:"size:20;default:pending" json:"approval_status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// 关联
	Fund     Fund      `json:"fund,omitempty"`
	Project  *Project  `json:"project,omitempty"`
	Employee *Employee `json:"employee,omitempty"`
}

// Purchase 采购表
type Purchase struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	PurchaseID    string     `gorm:"size:50;unique;not null" json:"purchase_id"`
	TransactionID *uint      `json:"transaction_id"`
	TotalSpent    float64    `gorm:"type:decimal(12,2);not null" json:"total_spent"`
	SupplierName  string     `gorm:"size:200" json:"supplier_name"`
	PurchaseDate  *time.Time `json:"purchase_date"`
	Description   string     `json:"description"`
	CreatedAt     time.Time  `json:"created_at"`

	// 关联
	Transaction *Transaction `json:"transaction,omitempty"`
	Inventory   []Inventory  `json:"inventory,omitempty"`
}

// Payroll 薪资表
type Payroll struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID uint      `gorm:"not null" json:"transaction_id"`
	EmployeeID    uint      `gorm:"not null" json:"employee_id"`
	Amount        float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	PayDate       time.Time `json:"pay_date"`
	Deductions    float64   `gorm:"type:decimal(10,2);default:0.00" json:"deductions"`
	Bonuses       float64   `gorm:"type:decimal(10,2);default:0.00" json:"bonuses"`
	CreatedAt     time.Time `json:"created_at"`

	// 关联
	Transaction Transaction `json:"transaction,omitempty"`
	Employee    Employee    `json:"employee,omitempty"`
}

// Inventory 库存表
type Inventory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	InventoryID  string    `gorm:"size:50;unique;not null" json:"inventory_id"`
	Name         string    `gorm:"size:200;not null" json:"name"`
	Category     string    `gorm:"size:100" json:"category"`
	PurchaseID   *uint     `json:"purchase_id"`
	LocationID   *uint     `json:"location_id"`
	CurrentStock int       `gorm:"default:0" json:"current_stock"`
	UnitCost     float64   `gorm:"type:decimal(10,2)" json:"unit_cost"`
	Status       string    `gorm:"size:20;default:available" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	Purchase              *Purchase              `json:"purchase,omitempty"`
	Location              *Location              `json:"location,omitempty"`
	Deliveries            []Delivery             `json:"deliveries,omitempty"`
	InventoryTransactions []InventoryTransaction `json:"inventory_transactions,omitempty"`
}

// GiftType 礼品类型表
type GiftType struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"size:100;not null" json:"name"`
	Category          string    `gorm:"size:50" json:"category"`
	UnitCost          float64   `gorm:"type:decimal(8,2)" json:"unit_cost"`
	RequiresInventory bool      `gorm:"default:true" json:"requires_inventory"`
	CreatedAt         time.Time `json:"created_at"`

	// 关联
	Gifts []Gift `json:"gifts,omitempty"`
}

// Gift 礼品记录表
type Gift struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	GiftID             string    `gorm:"size:20;unique;not null" json:"gift_id"`
	DonorID            *uint     `json:"donor_id"`
	DonationID         *uint     `json:"donation_id"`
	GiftTypeID         uint      `gorm:"not null" json:"gift_type_id"`
	Quantity           int       `gorm:"default:1" json:"quantity"`
	TotalValue         float64   `gorm:"type:decimal(10,2)" json:"total_value"`
	DistributionStatus string    `gorm:"size:20;default:pending" json:"distribution_status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// 关联
	Donor    *Donor    `json:"donor,omitempty"`
	Donation *Donation `json:"donation,omitempty"`
	GiftType GiftType  `json:"gift_type,omitempty"`
}

// InventoryTransaction 库存交易记录表
type InventoryTransaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TransactionID   string    `gorm:"size:20;unique;not null" json:"transaction_id"`
	InventoryID     uint      `gorm:"not null" json:"inventory_id"`
	TransactionType string    `gorm:"size:20;not null" json:"transaction_type"`
	QuantityChange  int       `gorm:"not null" json:"quantity_change"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`

	// 关联
	Inventory Inventory `json:"inventory,omitempty"`
}

// Delivery 配送表
type Delivery struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	DeliveryID    string     `gorm:"size:50;unique;not null" json:"delivery_id"`
	InventoryID   uint       `gorm:"not null" json:"inventory_id"`
	ProjectID     *uint      `json:"project_id"`
	LocationID    *uint      `json:"location_id"`
	Quantity      int        `gorm:"not null" json:"quantity"`
	RecipientName string     `gorm:"size:200" json:"recipient_name"`
	DeliveryDate  *time.Time `json:"delivery_date"`
	Status        string     `gorm:"size:20;default:pending" json:"status"`
	CreatedAt     time.Time  `json:"created_at"`

	// 关联
	Inventory *Inventory `json:"inventory,omitempty"`
	Project   *Project   `json:"project,omitempty"`
	Location  *Location  `json:"location,omitempty"`
}
