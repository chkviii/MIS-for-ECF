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
	Phone        string    `gorm:"size:20" json:"phone"`
	Status       string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Donor 捐赠者表
type Donor struct {
	ID                       uint      `gorm:"primaryKey" json:"id"`
	DonorID                  string    `gorm:"size:20;unique;not null" json:"donor_id"`
	FirstName                string    `gorm:"size:100;not null" json:"first_name"`
	LastName                 string    `gorm:"size:100;not null" json:"last_name"`
	Email                    string    `gorm:"size:255;unique" json:"email"`
	Phone                    string    `gorm:"size:20" json:"phone"`
	Address                  string    `json:"address"`
	City                     string    `gorm:"size:100" json:"city"`
	State                    string    `gorm:"size:100" json:"state"`
	ZipCode                  string    `gorm:"size:20" json:"zip_code"`
	Country                  string    `gorm:"size:100" json:"country"`
	DonorType                string    `gorm:"size:20;default:individual" json:"donor_type"`
	Age                      int       `json:"age"`
	Gender                   string    `gorm:"size:10" json:"gender"`
	Occupation               string    `gorm:"size:100" json:"occupation"`
	CommunicationPreference  string    `gorm:"size:50;default:email" json:"communication_preference"`
	TaxID                    string    `gorm:"size:50" json:"tax_id"`
	EnrollmentDate           time.Time `json:"enrollment_date"`
	Status                   string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`

	// 关联
	Donations         []Donation         `json:"donations,omitempty"`
	Funds             []Fund             `json:"funds,omitempty"`
	Gifts             []Gift             `json:"gifts,omitempty"`
	DonationInventory []DonationInventory `json:"donation_inventory,omitempty"`
}

// Volunteer 志愿者表
type Volunteer struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	VolunteerID            string    `gorm:"size:20;unique;not null" json:"volunteer_id"`
	FirstName              string    `gorm:"size:100;not null" json:"first_name"`
	LastName               string    `gorm:"size:100;not null" json:"last_name"`
	Email                  string    `gorm:"size:255;unique" json:"email"`
	Phone                  string    `gorm:"size:20" json:"phone"`
	EmergencyContactName   string    `gorm:"size:200" json:"emergency_contact_name"`
	EmergencyContactPhone  string    `gorm:"size:20" json:"emergency_contact_phone"`
	Skills                 string    `json:"skills"`
	Certifications         string    `json:"certifications"`
	Availability           string    `json:"availability"`
	PreferredTasks         string    `json:"preferred_tasks"`
	LocationID             *uint     `json:"location_id"`
	BackgroundCheckStatus  string    `gorm:"size:20;default:pending" json:"background_check_status"`
	BackgroundCheckDate    *time.Time `json:"background_check_date"`
	TrainingCompletionDate *time.Time `json:"training_completion_date"`
	TotalVolunteerHours    float64   `gorm:"type:decimal(10,2);default:0" json:"total_volunteer_hours"`
	StartDate              *time.Time `json:"start_date"`
	Status                 string    `gorm:"size:20;default:active" json:"status"`
	Notes                  string    `json:"notes"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`

	// 关联
	Location  *Location           `json:"location,omitempty"`
	Projects  []Project           `gorm:"many2many:volunteer_projects" json:"projects,omitempty"`
}

// Employee 员工表
type Employee struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	EmployeeID         string    `gorm:"size:20;unique;not null" json:"employee_id"`
	FirstName          string    `gorm:"size:100;not null" json:"first_name"`
	LastName           string    `gorm:"size:100;not null" json:"last_name"`
	Email              string    `gorm:"size:255;unique" json:"email"`
	Phone              string    `gorm:"size:20" json:"phone"`
	Position           string    `gorm:"size:100" json:"position"`
	Department         string    `gorm:"size:100" json:"department"`
	Salary             float64   `gorm:"type:decimal(10,2)" json:"salary"`
	HourlyRate         float64   `gorm:"type:decimal(8,2)" json:"hourly_rate"`
	EmploymentType     string    `gorm:"size:20;default:full_time" json:"employment_type"`
	HireDate           *time.Time `json:"hire_date"`
	TerminationDate    *time.Time `json:"termination_date"`
	LocationID         *uint     `json:"location_id"`
	SupervisorID       *uint     `json:"supervisor_id"`
	PayrollID          string    `gorm:"size:50" json:"payroll_id"`
	TaxID              string    `gorm:"size:50" json:"tax_id"`
	Benefits           string    `json:"benefits"`
	PerformanceRating  float64   `gorm:"type:decimal(3,2)" json:"performance_rating"`
	WorkHoursPerWeek   int       `gorm:"default:40" json:"work_hours_per_week"`
	Status             string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// 关联
	Location           *Location          `json:"location,omitempty"`
	Supervisor         *Employee          `json:"supervisor,omitempty"`
	Projects           []Project          `gorm:"many2many:employee_projects" json:"projects,omitempty"`
	ManagedProjects    []Project          `gorm:"foreignKey:ProjectManagerID" json:"managed_projects,omitempty"`
	ManagedFunds       []Fund             `gorm:"foreignKey:FundManagerID" json:"managed_funds,omitempty"`
	Expenses           []Expense          `json:"expenses,omitempty"`
	Payrolls           []Payroll          `json:"payrolls,omitempty"`
}

// Location 地点表
type Location struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	LocationID       string    `gorm:"size:50;unique;not null" json:"location_id"`
	CountryCode      string    `gorm:"size:3" json:"country_code"`
	Name             string    `gorm:"size:200;not null" json:"name"`
	Type             string    `gorm:"size:50" json:"type"`
	ParentLocationID *uint     `json:"parent_location_id"`
	Address          string    `json:"address"`
	CreatedAt        time.Time `json:"created_at"`

	// 关联
	ParentLocation *Location   `json:"parent_location,omitempty"`
	ChildLocations []Location  `gorm:"foreignKey:ParentLocationID" json:"child_locations,omitempty"`
	Projects       []Project   `json:"projects,omitempty"`
	Volunteers     []Volunteer `json:"volunteers,omitempty"`
	Employees      []Employee  `json:"employees,omitempty"`
	Inventory      []Inventory `json:"inventory,omitempty"`
	Deliveries     []Delivery  `json:"deliveries,omitempty"`
}

// Project 项目表
type Project struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	ProjectID              string    `gorm:"size:50;unique;not null" json:"project_id"`
	Name                   string    `gorm:"size:200;not null" json:"name"`
	Description            string    `json:"description"`
	ProjectType            string    `gorm:"size:50" json:"project_type"`
	Budget                 float64   `gorm:"type:decimal(12,2)" json:"budget"`
	ActualCost             float64   `gorm:"type:decimal(12,2);default:0" json:"actual_cost"`
	FundingSource          string    `gorm:"size:100" json:"funding_source"`
	LocationID             *uint     `json:"location_id"`
	StartDate              *time.Time `json:"start_date"`
	EndDate                *time.Time `json:"end_date"`
	ExpectedCompletionDate *time.Time `json:"expected_completion_date"`
	Status                 string    `gorm:"size:20;default:planning" json:"status"`
	Priority               string    `gorm:"size:20;default:medium" json:"priority"`
	BeneficiariesCount     int       `gorm:"default:0" json:"beneficiaries_count"`
	SuccessMetrics         string    `json:"success_metrics"`
	ProjectManagerID       *uint     `json:"project_manager_id"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`

	// 关联
	Location           *Location           `json:"location,omitempty"`
	ProjectManager     *Employee           `json:"project_manager,omitempty"`
	Volunteers         []Volunteer         `gorm:"many2many:volunteer_projects" json:"volunteers,omitempty"`
	Employees          []Employee          `gorm:"many2many:employee_projects" json:"employees,omitempty"`
	Funds              []Fund              `gorm:"many2many:fund_projects" json:"funds,omitempty"`
	Donations          []Donation          `json:"donations,omitempty"`
	Expenses           []Expense           `json:"expenses,omitempty"`
	Deliveries         []Delivery          `json:"deliveries,omitempty"`
	DonationInventory  []DonationInventory `json:"donation_inventory,omitempty"`
}
