package models

import (
	"time"
)

// User 用户表
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;unique;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:20;not null" json:"role"`
	Email        string    `gorm:"size:255;unique" json:"email"`
	FirstName    string    `gorm:"size:100" json:"first_name"`
	LastName     string    `gorm:"size:100" json:"last_name"`
	Status       string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Donor 捐赠者表
type Donor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DonorID   string    `gorm:"size:20;unique;not null" json:"donor_id"`
	FirstName string    `gorm:"size:100;not null" json:"first_name"`
	LastName  string    `gorm:"size:100;not null" json:"last_name"`
	Email     string    `gorm:"size:255;unique" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	DonorType string    `gorm:"size:20;default:individual" json:"donor_type"`
	Status    string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Donations []Donation `json:"donations,omitempty"`
	Funds     []Fund     `json:"funds,omitempty"`
	Gifts     []Gift     `json:"gifts,omitempty"`
}

// Volunteer 志愿者表
type Volunteer struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VolunteerID string    `gorm:"size:20;unique;not null" json:"volunteer_id"`
	FirstName   string    `gorm:"size:100;not null" json:"first_name"`
	LastName    string    `gorm:"size:100;not null" json:"last_name"`
	Email       string    `gorm:"size:255;unique" json:"email"`
	Phone       string    `gorm:"size:20" json:"phone"`
	Skills      string    `json:"skills"`
	LocationID  *uint     `json:"location_id"`
	Status      string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Location *Location `json:"location,omitempty"`
}

// Employee 员工表
type Employee struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID string    `gorm:"size:20;unique;not null" json:"employee_id"`
	FirstName  string    `gorm:"size:100;not null" json:"first_name"`
	LastName   string    `gorm:"size:100;not null" json:"last_name"`
	Email      string    `gorm:"size:255;unique" json:"email"`
	Phone      string    `gorm:"size:20" json:"phone"`
	Position   string    `gorm:"size:100" json:"position"`
	Department string    `gorm:"size:100" json:"department"`
	Salary     float64   `gorm:"type:decimal(10,2)" json:"salary"`
	LocationID *uint     `json:"location_id"`
	Status     string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联
	Location *Location `json:"location,omitempty"`
	Expenses []Expense `json:"expenses,omitempty"`
	Payrolls []Payroll `json:"payrolls,omitempty"`
}

// Location 地点表
type Location struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LocationID  string    `gorm:"size:50;unique;not null" json:"location_id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Type        string    `gorm:"size:50" json:"type"`
	Address     string    `json:"address"`
	CountryCode string    `gorm:"size:3" json:"country_code"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联
	Projects   []Project   `json:"projects,omitempty"`
	Volunteers []Volunteer `json:"volunteers,omitempty"`
	Employees  []Employee  `json:"employees,omitempty"`
	Inventory  []Inventory `json:"inventory,omitempty"`
	Deliveries []Delivery  `json:"deliveries,omitempty"`
}

// Project 项目表
type Project struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ProjectID   string     `gorm:"size:50;unique;not null" json:"project_id"`
	Name        string     `gorm:"size:200;not null" json:"name"`
	Description string     `json:"description"`
	ProjectType string     `gorm:"size:50" json:"project_type"`
	Budget      float64    `gorm:"type:decimal(12,2)" json:"budget"`
	ActualCost  float64    `gorm:"type:decimal(12,2);default:0" json:"actual_cost"`
	LocationID  *uint      `json:"location_id"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Status      string     `gorm:"size:20;default:planning" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联
	Location   *Location  `json:"location,omitempty"`
	Funds      []Fund     `json:"funds,omitempty"`
	Donations  []Donation `json:"donations,omitempty"`
	Expenses   []Expense  `json:"expenses,omitempty"`
	Deliveries []Delivery `json:"deliveries,omitempty"`
}
