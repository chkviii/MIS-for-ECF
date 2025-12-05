package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Username     string     `json:"username" gorm:"unique;not null"`
	PasswordHash string     `json:"password_hash" gorm:"column:password_hash;not null"`
	UserType     string     `json:"user_type" gorm:"column:user_type;not null"`
	Status       string     `json:"status" gorm:"default:active"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Location 地点模型
type Location struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	LocationID  string    `json:"location_id" gorm:"uniqueIndex;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Type        string    `json:"type"`
	Address     string    `json:"address"`
	CountryCode string    `json:"country_code" gorm:"size:3"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Location) TableName() string {
	return "locations"
}

// Project 项目模型
type Project struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	ProjectID   string     `json:"project_id" gorm:"uniqueIndex;not null"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	ProjectType string     `json:"project_type"`
	Budget      float64    `json:"budget"`
	ActualCost  float64    `json:"actual_cost" gorm:"default:0"`
	LocationID  *uint      `json:"location_id"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Status      string     `json:"status" gorm:"default:planning"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	Location *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

// Donor 捐赠者模型
type Donor struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         *uint     `json:"user_id" gorm:"uniqueIndex"`
	DonorID        string    `json:"donor_id" gorm:"uniqueIndex;not null"`
	FirstName      string    `json:"first_name" gorm:"not null"`
	LastName       string    `json:"last_name" gorm:"not null"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	DonorType      string    `json:"donor_type" gorm:"default:individual"`
	TotalDonated   float64   `json:"total_donated" gorm:"default:0"`
	EnrollmentDate time.Time `json:"enrollment_date" gorm:"default:CURRENT_DATE"`
	Status         string    `json:"status" gorm:"default:active"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (Donor) TableName() string {
	return "donors"
}

// Volunteer 志愿者模型
type Volunteer struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           *uint     `json:"user_id" gorm:"uniqueIndex"`
	VolunteerID      string    `json:"volunteer_id" gorm:"uniqueIndex;not null"`
	FirstName        string    `json:"first_name" gorm:"not null"`
	LastName         string    `json:"last_name" gorm:"not null"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	LocationID       *uint     `json:"location_id"`
	Skills           string    `json:"skills"`
	Availability     string    `json:"availability"`
	HoursContributed float64   `json:"hours_contributed" gorm:"default:0"`
	Status           string    `json:"status" gorm:"default:active"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Location *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// TableName 指定表名
func (Volunteer) TableName() string {
	return "volunteers"
}

// Employee 员工模型
type Employee struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     *uint     `json:"user_id" gorm:"uniqueIndex"`
	EmployeeID string    `json:"employee_id" gorm:"uniqueIndex;not null"`
	FirstName  string    `json:"first_name" gorm:"not null"`
	LastName   string    `json:"last_name" gorm:"not null"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Position   string    `json:"position"`
	Department string    `json:"department"`
	Salary     float64   `json:"salary"`
	HireDate   time.Time `json:"hire_date" gorm:"default:CURRENT_DATE"`
	LocationID *uint     `json:"location_id"`
	Status     string    `json:"status" gorm:"default:active"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Location *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// TableName 指定表名
func (Employee) TableName() string {
	return "employees"
}
