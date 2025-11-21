package models

import "time"

// Transaction 交易表
type Transaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TransactionID   string    `gorm:"size:50;unique;not null" json:"transaction_id"`
	Type            string    `gorm:"size:20;not null" json:"type"`
	Amount          float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Currency        string    `gorm:"size:3;default:USD" json:"currency"`
	FromEntity      string    `gorm:"size:200" json:"from_entity"`
	ToEntity        string    `gorm:"size:200" json:"to_entity"`
	Description     string    `json:"description"`
	ReferenceType   string    `gorm:"size:20" json:"reference_type"`
	ReferenceID     *uint     `json:"reference_id"`
	TransactionDate time.Time `json:"transaction_date"`
	Fingerprint     string    `gorm:"size:255" json:"fingerprint"`
	CreatedAt       time.Time `json:"created_at"`

	// 关联
	Purchases []Purchase `json:"purchases,omitempty"`
	Funds     []Fund     `json:"funds,omitempty"`
	Payrolls  []Payroll  `json:"payrolls,omitempty"`
}

// Donation 捐赠记录表
type Donation struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	DonationID             string    `gorm:"size:20;unique;not null" json:"donation_id"`
	DonorID                uint      `gorm:"not null" json:"donor_id"`
	Amount                 float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	DonationType           string    `gorm:"size:20;not null" json:"donation_type"`
	DonationMethod         string    `gorm:"size:20" json:"donation_method"`
	Category               string    `gorm:"size:20;not null" json:"category"`
	Subcategory            string    `gorm:"size:50" json:"subcategory"`
	Currency               string    `gorm:"size:3;default:USD" json:"currency"`
	ExchangeRate           float64   `gorm:"type:decimal(8,4);default:1.0000" json:"exchange_rate"`
	TaxDeductible          bool      `gorm:"default:true" json:"tax_deductible"`
	TaxDeductionAmount     float64   `gorm:"type:decimal(10,2)" json:"tax_deduction_amount"`
	ReceiptNumber          string    `gorm:"size:50;unique" json:"receipt_number"`
	ReceiptIssuedDate      *time.Time `json:"receipt_issued_date"`
	ProjectID              *uint     `json:"project_id"`
	CampaignID             string    `gorm:"size:50" json:"campaign_id"`
	TributeType            string    `gorm:"size:20" json:"tribute_type"`
	TributeName            string    `gorm:"size:200" json:"tribute_name"`
	TributeNotification    string    `json:"tribute_notification"`
	Notes                  string    `json:"notes"`
	DonationDate           time.Time `json:"donation_date"`
	ProcessedDate          *time.Time `json:"processed_date"`
	AcknowledgmentSentDate *time.Time `json:"acknowledgment_sent_date"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`

	// 关联
	Donor   Donor    `json:"donor,omitempty"`
	Project *Project `json:"project,omitempty"`
	Gifts   []Gift   `json:"gifts,omitempty"`
}

// Fund 基金管理表
type Fund struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	FundID          string    `gorm:"size:20;unique;not null" json:"fund_id"`
	DonorID         *uint     `json:"donor_id"`
	Name            string    `gorm:"size:200;not null" json:"name"`
	FundType        string    `gorm:"size:20;not null" json:"fund_type"`
	Description     string    `json:"description"`
	TotalAmount     float64   `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	CurrentBalance  float64   `gorm:"type:decimal(12,2);default:0" json:"current_balance"`
	AmountLeft      float64   `gorm:"type:decimal(12,2);not null" json:"amount_left"`
	MinimumBalance  float64   `gorm:"type:decimal(12,2);default:0" json:"minimum_balance"`
	Restrictions    string    `json:"restrictions"`
	EstablishedDate *time.Time `json:"established_date"`
	FundDate        *time.Time `json:"fund_date"`
	ExpirationDate  *time.Time `json:"expiration_date"`
	Currency        string    `gorm:"size:3;default:USD" json:"currency"`
	Purpose         string    `json:"purpose"`
	TransactionID   *uint     `json:"transaction_id"`
	ProjectID       *uint     `json:"project_id"`
	FundManagerID   *uint     `json:"fund_manager_id"`
	Status          string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// 关联
	Donor       *Donor       `json:"donor,omitempty"`
	Transaction *Transaction `json:"transaction,omitempty"`
	Project     *Project     `json:"project,omitempty"`
	FundManager *Employee    `json:"fund_manager,omitempty"`
	Projects    []Project    `gorm:"many2many:fund_projects" json:"projects,omitempty"`
	Expenses    []Expense    `json:"expenses,omitempty"`
}

// Expense 支出记录表
type Expense struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	ExpenseID        string    `gorm:"size:20;unique;not null" json:"expense_id"`
	FundID           uint      `gorm:"not null" json:"fund_id"`
	ProjectID        *uint     `json:"project_id"`
	EmployeeID       *uint     `json:"employee_id"`
	VendorName       string    `gorm:"size:200" json:"vendor_name"`
	Description      string    `gorm:"not null" json:"description"`
	ExpenseCategory  string    `gorm:"size:50" json:"expense_category"`
	ExpenseType      string    `gorm:"size:50" json:"expense_type"`
	Amount           float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency         string    `gorm:"size:3;default:USD" json:"currency"`
	PaymentMethod    string    `gorm:"size:20" json:"payment_method"`
	ReceiptNumber    string    `gorm:"size:50" json:"receipt_number"`
	InvoiceNumber    string    `gorm:"size:50" json:"invoice_number"`
	ApprovalStatus   string    `gorm:"size:20;default:pending" json:"approval_status"`
	ApprovedBy       *uint     `json:"approved_by"`
	ApprovalDate     *time.Time `json:"approval_date"`
	ExpenseDate      time.Time `json:"expense_date"`
	PaymentDate      *time.Time `json:"payment_date"`
	Reimbursable     bool      `gorm:"default:false" json:"reimbursable"`
	Reimbursed       bool      `gorm:"default:false" json:"reimbursed"`
	ReimbursementDate *time.Time `json:"reimbursement_date"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// 关联
	Fund     Fund      `json:"fund,omitempty"`
	Project  *Project  `json:"project,omitempty"`
	Employee *Employee `json:"employee,omitempty"`
	Approver *Employee `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

// Purchase 采购表
type Purchase struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	PurchaseID    string    `gorm:"size:50;unique;not null" json:"purchase_id"`
	TotalSpent    float64   `gorm:"type:decimal(12,2);not null" json:"total_spent"`
	Currency      string    `gorm:"size:3;default:USD" json:"currency"`
	PurchaseDate  *time.Time `json:"purchase_date"`
	TransactionID *uint     `json:"transaction_id"`
	SupplierName  string    `gorm:"size:200" json:"supplier_name"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`

	// 关联
	Transaction *Transaction `json:"transaction,omitempty"`
	Inventory   []Inventory  `json:"inventory,omitempty"`
}

// Payroll 薪资表
type Payroll struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TransactionID   uint      `gorm:"not null" json:"transaction_id"`
	EmployeeID      uint      `gorm:"not null" json:"employee_id"`
	PayDate         time.Time `json:"pay_date"`
	Amount          float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	PayPeriodStart  *time.Time `json:"pay_period_start"`
	PayPeriodEnd    *time.Time `json:"pay_period_end"`
	Deductions      float64   `gorm:"type:decimal(10,2);default:0.00" json:"deductions"`
	Bonuses         float64   `gorm:"type:decimal(10,2);default:0.00" json:"bonuses"`
	OvertimeHours   float64   `gorm:"type:decimal(5,2);default:0.00" json:"overtime_hours"`
	OvertimeRate    float64   `gorm:"type:decimal(8,2);default:0.00" json:"overtime_rate"`
	CreatedAt       time.Time `json:"created_at"`

	// 关联
	Transaction Transaction `json:"transaction,omitempty"`
	Employee    Employee    `json:"employee,omitempty"`
}

// Inventory 库存表
type Inventory struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	InventoryID        string    `gorm:"size:50;unique;not null" json:"inventory_id"`
	Name               string    `gorm:"size:200;not null" json:"name"`
	Category           string    `gorm:"size:100" json:"category"`
	Subcategory        string    `gorm:"size:50" json:"subcategory"`
	SKU                string    `gorm:"size:50;unique" json:"sku"`
	PurchaseAmount     float64   `gorm:"type:decimal(10,2)" json:"purchase_amount"`
	RemainAmount       int       `gorm:"default:0" json:"remain_amount"`
	CurrentStock       int       `gorm:"default:0" json:"current_stock"`
	MinimumStockLevel  int       `gorm:"default:0" json:"minimum_stock_level"`
	MaximumStockLevel  int       `json:"maximum_stock_level"`
	PurchaseID         *uint     `json:"purchase_id"`
	LocationID         *uint     `json:"location_id"`
	UnitCost           float64   `gorm:"type:decimal(10,2)" json:"unit_cost"`
	TotalValue         float64   `gorm:"type:decimal(10,2)" json:"total_value"`
	SupplierName       string    `gorm:"size:200" json:"supplier_name"`
	SupplierContact    string    `json:"supplier_contact"`
	StorageLocation    string    `gorm:"size:100" json:"storage_location"`
	Depreciation       float64   `gorm:"type:decimal(5,2);default:0.00" json:"depreciation"`
	ExpirationDate     *time.Time `json:"expiration_date"`
	LastInventoryDate  *time.Time `json:"last_inventory_date"`
	Status             string    `gorm:"size:20;default:available" json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// 关联
	Purchase              *Purchase             `json:"purchase,omitempty"`
	Location              *Location             `json:"location,omitempty"`
	Deliveries            []Delivery            `json:"deliveries,omitempty"`
	DonationInventory     []DonationInventory   `json:"donation_inventory,omitempty"`
	InventoryTransactions []InventoryTransaction `json:"inventory_transactions,omitempty"`
}

// GiftType 礼品类型表
type GiftType struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Name                string    `gorm:"size:100;not null" json:"name"`
	Description         string    `json:"description"`
	Category            string    `gorm:"size:50" json:"category"`
	UnitCost            float64   `gorm:"type:decimal(8,2)" json:"unit_cost"`
	TaxDeductibleValue  float64   `gorm:"type:decimal(8,2)" json:"tax_deductible_value"`
	RequiresInventory   bool      `gorm:"default:true" json:"requires_inventory"`
	CreatedAt           time.Time `json:"created_at"`

	// 关联
	Gifts []Gift `json:"gifts,omitempty"`
}

// Gift 礼品记录表
type Gift struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	GiftID              string    `gorm:"size:20;unique;not null" json:"gift_id"`
	DonorID             *uint     `json:"donor_id"`
	DonationID          *uint     `json:"donation_id"`
	GiftTypeID          uint      `gorm:"not null" json:"gift_type_id"`
	Quantity            int       `gorm:"default:1" json:"quantity"`
	UnitValue           float64   `gorm:"type:decimal(8,2)" json:"unit_value"`
	TotalValue          float64   `gorm:"type:decimal(10,2)" json:"total_value"`
	IsFree              bool      `gorm:"default:false" json:"is_free"`
	DistributionMethod  string    `gorm:"size:20" json:"distribution_method"`
	DistributionStatus  string    `gorm:"size:20;default:pending" json:"distribution_status"`
	RecipientName       string    `gorm:"size:200" json:"recipient_name"`
	RecipientAddress    string    `json:"recipient_address"`
	TrackingNumber      string    `gorm:"size:100" json:"tracking_number"`
	DistributedDate     *time.Time `json:"distributed_date"`
	DistributedBy       *uint     `json:"distributed_by"`
	Notes               string    `json:"notes"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// 关联
	Donor       *Donor    `json:"donor,omitempty"`
	Donation    *Donation `json:"donation,omitempty"`
	GiftType    GiftType  `json:"gift_type,omitempty"`
	Distributor *Employee `gorm:"foreignKey:DistributedBy" json:"distributor,omitempty"`
}

// InventoryTransaction 库存交易记录表
type InventoryTransaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TransactionID   string    `gorm:"size:20;unique;not null" json:"transaction_id"`
	InventoryID     uint      `gorm:"not null" json:"inventory_id"`
	TransactionType string    `gorm:"size:20;not null" json:"transaction_type"`
	QuantityChange  int       `gorm:"not null" json:"quantity_change"`
	UnitCost        float64   `gorm:"type:decimal(8,2)" json:"unit_cost"`
	TotalCost       float64   `gorm:"type:decimal(10,2)" json:"total_cost"`
	ReferenceType   string    `gorm:"size:20" json:"reference_type"`
	ReferenceID     *uint     `json:"reference_id"`
	Notes           string    `json:"notes"`
	ProcessedBy     *uint     `json:"processed_by"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`

	// 关联
	Inventory Inventory `json:"inventory,omitempty"`
	Processor *Employee `gorm:"foreignKey:ProcessedBy" json:"processor,omitempty"`
}

// Delivery 配送表
type Delivery struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DeliveryID    string    `gorm:"size:50;unique;not null" json:"delivery_id"`
	InventoryID   uint      `gorm:"not null" json:"inventory_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	ProjectID     *uint     `json:"project_id"`
	DeliveryDate  *time.Time `json:"delivery_date"`
	Status        string    `gorm:"size:20;default:pending" json:"status"`
	LocationID    *uint     `json:"location_id"`
	RecipientName string    `gorm:"size:200" json:"recipient_name"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`

	// 关联
	Inventory *Inventory `json:"inventory,omitempty"`
	Project   *Project   `json:"project,omitempty"`
	Location  *Location  `json:"location,omitempty"`
}
