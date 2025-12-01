package repo

import (
	"erp-backend/internal/models"

	"gorm.io/gorm"
)

// applyFilters applies query, number_range and date_range filters to the GORM tx.
// - query: map[string]interface{} -> for string values use LIKE (without adding wildcards), otherwise =
// - numberRange: map[string][]interface{} -> [min, max], apply >= min and <= max if present
// - dateRange: map[string][]string -> [start, end] in YYYY-MM-DD, apply >= start and <= end
func applyFilters(tx *gorm.DB, query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) *gorm.DB {
	if tx == nil {
		return tx
	}

	for key, value := range query {
		if value != "" && value != nil {
			if s, ok := value.(string); ok {
				// Use LIKE as requested but do not add wildcard characters
				tx = tx.Where(key+" LIKE ?", s)
			} else {
				tx = tx.Where(key+" = ?", value)
			}
		}
	}

	for key, rangeVals := range numberRange {
		if len(rangeVals) > 0 {
			if rangeVals[0] != nil && rangeVals[0] != "" {
				tx = tx.Where(key+" >= ?", rangeVals[0])
			}
		}
		if len(rangeVals) > 1 {
			if rangeVals[1] != nil && rangeVals[1] != "" {
				tx = tx.Where(key+" <= ?", rangeVals[1])
			}
		}
	}

	for key, rangeVals := range dateRange {
		if len(rangeVals) > 0 {
			if rangeVals[0] != "" {
				tx = tx.Where(key+" >= ?", rangeVals[0])
			}
		}
		if len(rangeVals) > 1 {
			if rangeVals[1] != "" {
				tx = tx.Where(key+" <= ?", rangeVals[1])
			}
		}
	}

	return tx
}

// UserRepository 用户仓库
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Search(query map[string]interface{}) ([]models.User, error) {
	tx := r.db.Model(&models.User{})

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var users []models.User
	err := tx.Find(&users).Error
	return users, err
}

func (r *UserRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.User, error) {
	tx := r.db.Model(&models.User{})

	tx = applyFilters(tx, query, numberRange, dateRange)

	var users []models.User
	err := tx.Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// ProjectRepository 项目仓储
type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) GetAll() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) Search(query map[string]interface{}) ([]models.Project, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var projects []models.Project
	err := tx.Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Project, error) {
	tx := r.db.Model(&models.Project{})

	tx = applyFilters(tx, query, numberRange, dateRange)

	var projects []models.Project
	err := tx.Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}

// DonorRepository 捐赠者仓储
type DonorRepository struct {
	db *gorm.DB
}

func NewDonorRepository(db *gorm.DB) *DonorRepository {
	return &DonorRepository{db: db}
}

func (r *DonorRepository) Create(donor *models.Donor) error {
	return r.db.Create(donor).Error
}

func (r *DonorRepository) GetAll() ([]models.Donor, error) {
	var donors []models.Donor
	err := r.db.Find(&donors).Error
	return donors, err
}

func (r *DonorRepository) Search(query map[string]interface{}) ([]models.Donor, error) {
	tx := r.db.Model(&models.Donor{})

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var donors []models.Donor
	err := tx.Find(&donors).Error
	return donors, err
}

func (r *DonorRepository) Update(donor *models.Donor) error {
	return r.db.Save(donor).Error
}

func (r *DonorRepository) Delete(id uint) error {
	return r.db.Delete(&models.Donor{}, id).Error
}

// DonationRepository 捐赠仓储
type DonationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) *DonationRepository {
	return &DonationRepository{db: db}
}

func (r *DonationRepository) Create(donation *models.Donation) error {
	return r.db.Create(donation).Error
}

func (r *DonationRepository) GetAll() ([]models.Donation, error) {
	var donations []models.Donation
	err := r.db.Find(&donations).Error
	return donations, err
}

func (r *DonationRepository) Search(query map[string]interface{}) ([]models.Donation, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var donations []models.Donation
	err := tx.Find(&donations).Error
	return donations, err
}

func (r *DonationRepository) Update(donation *models.Donation) error {
	return r.db.Save(donation).Error
}

func (r *DonationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Donation{}, id).Error
}

// VolunteerRepository 志愿者仓储
type VolunteerRepository struct {
	db *gorm.DB
}

func NewVolunteerRepository(db *gorm.DB) *VolunteerRepository {
	return &VolunteerRepository{db: db}
}

func (r *VolunteerRepository) Create(volunteer *models.Volunteer) error {
	return r.db.Create(volunteer).Error
}

func (r *VolunteerRepository) GetAll() ([]models.Volunteer, error) {
	var volunteers []models.Volunteer
	err := r.db.Find(&volunteers).Error
	return volunteers, err
}

func (r *VolunteerRepository) Search(query map[string]interface{}) ([]models.Volunteer, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var volunteers []models.Volunteer
	err := tx.Find(&volunteers).Error
	return volunteers, err
}

func (r *VolunteerRepository) Update(volunteer *models.Volunteer) error {
	return r.db.Save(volunteer).Error
}

func (r *VolunteerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Volunteer{}, id).Error
}

// EmployeeRepository 员工仓储
type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

func (r *EmployeeRepository) GetAll() ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) Search(query map[string]interface{}) ([]models.Employee, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var employees []models.Employee
	err := tx.Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) Update(employee *models.Employee) error {
	return r.db.Save(employee).Error
}

func (r *EmployeeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Employee{}, id).Error
}

// LocationRepository 地点仓储
type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Create(location *models.Location) error {
	return r.db.Create(location).Error
}

func (r *LocationRepository) GetAll() ([]models.Location, error) {
	var locations []models.Location
	err := r.db.Find(&locations).Error
	return locations, err
}

func (r *LocationRepository) Search(query map[string]interface{}) ([]models.Location, error) {
	tx := r.db.Model(&models.Location{})

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var locations []models.Location
	err := tx.Find(&locations).Error
	return locations, err
}

func (r *LocationRepository) Update(location *models.Location) error {
	return r.db.Save(location).Error
}

func (r *LocationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Location{}, id).Error
}

// FundRepository 基金仓储
type FundRepository struct {
	db *gorm.DB
}

func NewFundRepository(db *gorm.DB) *FundRepository {
	return &FundRepository{db: db}
}

func (r *FundRepository) Create(fund *models.Fund) error {
	return r.db.Create(fund).Error
}

func (r *FundRepository) GetAll() ([]models.Fund, error) {
	var funds []models.Fund
	err := r.db.Find(&funds).Error
	return funds, err
}

func (r *FundRepository) Search(query map[string]interface{}) ([]models.Fund, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var funds []models.Fund
	err := tx.Find(&funds).Error
	return funds, err
}

func (r *FundRepository) Update(fund *models.Fund) error {
	return r.db.Save(fund).Error
}

func (r *FundRepository) Delete(id uint) error {
	return r.db.Delete(&models.Fund{}, id).Error
}

// ExpenseRepository 支出仓储
type ExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Find(&expenses).Error
	return expenses, err
}

func (r *ExpenseRepository) Search(query map[string]interface{}) ([]models.Expense, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var expenses []models.Expense
	err := tx.Find(&expenses).Error
	return expenses, err
}

func (r *ExpenseRepository) Update(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

func (r *ExpenseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Expense{}, id).Error
}

// TransactionRepository 交易仓储
type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) GetAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) Search(query map[string]interface{}) ([]models.Transaction, error) {
	tx := r.db.Model(&models.Transaction{})

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var transactions []models.Transaction
	err := tx.Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *TransactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}

// PurchaseRepository 采购仓储
type PurchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) Create(purchase *models.Purchase) error {
	return r.db.Create(purchase).Error
}

func (r *PurchaseRepository) GetAll() ([]models.Purchase, error) {
	var purchases []models.Purchase
	err := r.db.Find(&purchases).Error
	return purchases, err
}

func (r *PurchaseRepository) Search(query map[string]interface{}) ([]models.Purchase, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var purchases []models.Purchase
	err := tx.Find(&purchases).Error
	return purchases, err
}

func (r *PurchaseRepository) Update(purchase *models.Purchase) error {
	return r.db.Save(purchase).Error
}

func (r *PurchaseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Purchase{}, id).Error
}

// PayrollRepository 薪资仓储
type PayrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) *PayrollRepository {
	return &PayrollRepository{db: db}
}

func (r *PayrollRepository) Create(payroll *models.Payroll) error {
	return r.db.Create(payroll).Error
}

func (r *PayrollRepository) GetAll() ([]models.Payroll, error) {
	var payrolls []models.Payroll
	err := r.db.Find(&payrolls).Error
	return payrolls, err
}

func (r *PayrollRepository) Search(query map[string]interface{}) ([]models.Payroll, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var payrolls []models.Payroll
	err := tx.Find(&payrolls).Error
	return payrolls, err
}

func (r *PayrollRepository) Update(payroll *models.Payroll) error {
	return r.db.Save(payroll).Error
}

func (r *PayrollRepository) Delete(id uint) error {
	return r.db.Delete(&models.Payroll{}, id).Error
}

// InventoryRepository 库存仓储
type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(inventory *models.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *InventoryRepository) GetAll() ([]models.Inventory, error) {
	var inventories []models.Inventory
	err := r.db.Find(&inventories).Error
	return inventories, err
}

func (r *InventoryRepository) Search(query map[string]interface{}) ([]models.Inventory, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var inventories []models.Inventory
	err := tx.Find(&inventories).Error
	return inventories, err
}

func (r *InventoryRepository) Update(inventory *models.Inventory) error {
	return r.db.Save(inventory).Error
}

func (r *InventoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Inventory{}, id).Error
}

// GiftTypeRepository 礼品类型仓储
type GiftTypeRepository struct {
	db *gorm.DB
}

func NewGiftTypeRepository(db *gorm.DB) *GiftTypeRepository {
	return &GiftTypeRepository{db: db}
}

func (r *GiftTypeRepository) Create(giftType *models.GiftType) error {
	return r.db.Create(giftType).Error
}

func (r *GiftTypeRepository) GetAll() ([]models.GiftType, error) {
	var giftTypes []models.GiftType
	err := r.db.Find(&giftTypes).Error
	return giftTypes, err
}

func (r *GiftTypeRepository) Search(query map[string]interface{}) ([]models.GiftType, error) {
	tx := r.db.Model(&models.GiftType{})

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var giftTypes []models.GiftType
	err := tx.Find(&giftTypes).Error
	return giftTypes, err
}

func (r *GiftTypeRepository) Update(giftType *models.GiftType) error {
	return r.db.Save(giftType).Error
}

func (r *GiftTypeRepository) Delete(id uint) error {
	return r.db.Delete(&models.GiftType{}, id).Error
}

// GiftRepository 礼品仓储
type GiftRepository struct {
	db *gorm.DB
}

func NewGiftRepository(db *gorm.DB) *GiftRepository {
	return &GiftRepository{db: db}
}

func (r *GiftRepository) Create(gift *models.Gift) error {
	return r.db.Create(gift).Error
}

func (r *GiftRepository) GetAll() ([]models.Gift, error) {
	var gifts []models.Gift
	err := r.db.Find(&gifts).Error
	return gifts, err
}

func (r *GiftRepository) Search(query map[string]interface{}) ([]models.Gift, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var gifts []models.Gift
	err := tx.Find(&gifts).Error
	return gifts, err
}

func (r *GiftRepository) Update(gift *models.Gift) error {
	return r.db.Save(gift).Error
}

func (r *GiftRepository) Delete(id uint) error {
	return r.db.Delete(&models.Gift{}, id).Error
}

// InventoryTransactionRepository 库存交易仓储
type InventoryTransactionRepository struct {
	db *gorm.DB
}

func NewInventoryTransactionRepository(db *gorm.DB) *InventoryTransactionRepository {
	return &InventoryTransactionRepository{db: db}
}

func (r *InventoryTransactionRepository) Create(transaction *models.InventoryTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *InventoryTransactionRepository) GetAll() ([]models.InventoryTransaction, error) {
	var transactions []models.InventoryTransaction
	err := r.db.Find(&transactions).Error
	return transactions, err
}

func (r *InventoryTransactionRepository) Search(query map[string]interface{}) ([]models.InventoryTransaction, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var transactions []models.InventoryTransaction
	err := tx.Find(&transactions).Error
	return transactions, err
}

func (r *InventoryTransactionRepository) Update(transaction *models.InventoryTransaction) error {
	return r.db.Save(transaction).Error
}

func (r *InventoryTransactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.InventoryTransaction{}, id).Error
}

// DeliveryRepository 配送仓储
type DeliveryRepository struct {
	db *gorm.DB
}

func NewDeliveryRepository(db *gorm.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

func (r *DeliveryRepository) Create(delivery *models.Delivery) error {
	return r.db.Create(delivery).Error
}

func (r *DeliveryRepository) GetAll() ([]models.Delivery, error) {
	var deliveries []models.Delivery
	err := r.db.Find(&deliveries).Error
	return deliveries, err
}

func (r *DeliveryRepository) Search(query map[string]interface{}) ([]models.Delivery, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var deliveries []models.Delivery
	err := tx.Find(&deliveries).Error
	return deliveries, err
}

func (r *DeliveryRepository) Update(delivery *models.Delivery) error {
	return r.db.Save(delivery).Error
}

func (r *DeliveryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Delivery{}, id).Error
}

// VolunteerProjectRepository 志愿者-项目关联仓储
type VolunteerProjectRepository struct {
	db *gorm.DB
}

func NewVolunteerProjectRepository(db *gorm.DB) *VolunteerProjectRepository {
	return &VolunteerProjectRepository{db: db}
}

func (r *VolunteerProjectRepository) Create(vp *models.VolunteerProject) error {
	return r.db.Create(vp).Error
}

func (r *VolunteerProjectRepository) GetAll() ([]models.VolunteerProject, error) {
	var vps []models.VolunteerProject
	err := r.db.Find(&vps).Error
	return vps, err
}

func (r *VolunteerProjectRepository) Search(query map[string]interface{}) ([]models.VolunteerProject, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var volunteerProjects []models.VolunteerProject
	err := tx.Find(&volunteerProjects).Error
	return volunteerProjects, err
}

func (r *VolunteerProjectRepository) Update(vp *models.VolunteerProject) error {
	return r.db.Save(vp).Error
}

func (r *VolunteerProjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.VolunteerProject{}, id).Error
}

// EmployeeProjectRepository 员工-项目关联仓储
type EmployeeProjectRepository struct {
	db *gorm.DB
}

func NewEmployeeProjectRepository(db *gorm.DB) *EmployeeProjectRepository {
	return &EmployeeProjectRepository{db: db}
}

func (r *EmployeeProjectRepository) Create(ep *models.EmployeeProject) error {
	return r.db.Create(ep).Error
}

func (r *EmployeeProjectRepository) GetAll() ([]models.EmployeeProject, error) {
	var eps []models.EmployeeProject
	err := r.db.Find(&eps).Error
	return eps, err
}

func (r *EmployeeProjectRepository) Search(query map[string]interface{}) ([]models.EmployeeProject, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var employeeProjects []models.EmployeeProject
	err := tx.Find(&employeeProjects).Error
	return employeeProjects, err
}

func (r *EmployeeProjectRepository) Update(ep *models.EmployeeProject) error {
	return r.db.Save(ep).Error
}

func (r *EmployeeProjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.EmployeeProject{}, id).Error
}

// FundProjectRepository 基金-项目关联仓储
type FundProjectRepository struct {
	db *gorm.DB
}

func NewFundProjectRepository(db *gorm.DB) *FundProjectRepository {
	return &FundProjectRepository{db: db}
}

func (r *FundProjectRepository) Create(fp *models.FundProject) error {
	return r.db.Create(fp).Error
}

func (r *FundProjectRepository) GetAll() ([]models.FundProject, error) {
	var fps []models.FundProject
	err := r.db.Find(&fps).Error
	return fps, err
}

func (r *FundProjectRepository) Search(query map[string]interface{}) ([]models.FundProject, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var fundProjects []models.FundProject
	err := tx.Find(&fundProjects).Error
	return fundProjects, err
}

func (r *FundProjectRepository) Update(fp *models.FundProject) error {
	return r.db.Save(fp).Error
}

func (r *FundProjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.FundProject{}, id).Error
}

// DonationInventoryRepository 捐赠-库存关联仓储
type DonationInventoryRepository struct {
	db *gorm.DB
}

func NewDonationInventoryRepository(db *gorm.DB) *DonationInventoryRepository {
	return &DonationInventoryRepository{db: db}
}

func (r *DonationInventoryRepository) Create(di *models.DonationInventory) error {
	return r.db.Create(di).Error
}

func (r *DonationInventoryRepository) GetAll() ([]models.DonationInventory, error) {
	var dis []models.DonationInventory
	err := r.db.Find(&dis).Error
	return dis, err
}

func (r *DonationInventoryRepository) Search(query map[string]interface{}) ([]models.DonationInventory, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var donationInventories []models.DonationInventory
	err := tx.Find(&donationInventories).Error
	return donationInventories, err
}

func (r *DonationInventoryRepository) Update(di *models.DonationInventory) error {
	return r.db.Save(di).Error
}

func (r *DonationInventoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.DonationInventory{}, id).Error
}

// DeliveryInventoryRepository 配送-库存关联仓储
type DeliveryInventoryRepository struct {
	db *gorm.DB
}

func NewDeliveryInventoryRepository(db *gorm.DB) *DeliveryInventoryRepository {
	return &DeliveryInventoryRepository{db: db}
}

func (r *DeliveryInventoryRepository) Create(di *models.DeliveryInventory) error {
	return r.db.Create(di).Error
}
func (r *DeliveryInventoryRepository) GetAll() ([]models.DeliveryInventory, error) {
	var dis []models.DeliveryInventory
	err := r.db.Find(&dis).Error
	return dis, err
}

func (r *DeliveryInventoryRepository) Search(query map[string]interface{}) ([]models.DeliveryInventory, error) {
	tx := r.db
	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}
	var deliveryInventories []models.DeliveryInventory
	err := tx.Find(&deliveryInventories).Error
	return deliveryInventories, err
}

func (r *DeliveryInventoryRepository) Update(di *models.DeliveryInventory) error {
	return r.db.Save(di).Error
}

func (r *DeliveryInventoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.DeliveryInventory{}, id).Error
}

// ScheduleRepository 调度仓储
type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(schedule *models.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *ScheduleRepository) GetAll() ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) Search(query map[string]interface{}) ([]models.Schedule, error) {
	tx := r.db

	for key, value := range query {
		if value != "" && value != nil {
			tx = tx.Where(key+" = ?", value)
		}
	}

	var schedules []models.Schedule
	err := tx.Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) Update(schedule *models.Schedule) error {
	return r.db.Save(schedule).Error
}

func (r *ScheduleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Schedule{}, id).Error
}

// ------- Generic Filter methods for repositories (use applyFilters) -------

func (r *DonorRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Donor, error) {
	tx := r.db.Model(&models.Donor{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var donors []models.Donor
	err := tx.Find(&donors).Error
	return donors, err
}

func (r *DonationRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Donation, error) {
	tx := r.db.Model(&models.Donation{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var donations []models.Donation
	err := tx.Find(&donations).Error
	return donations, err
}

func (r *VolunteerRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Volunteer, error) {
	tx := r.db.Model(&models.Volunteer{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var volunteers []models.Volunteer
	err := tx.Find(&volunteers).Error
	return volunteers, err
}

func (r *EmployeeRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Employee, error) {
	tx := r.db.Model(&models.Employee{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var employees []models.Employee
	err := tx.Find(&employees).Error
	return employees, err
}

func (r *LocationRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Location, error) {
	tx := r.db.Model(&models.Location{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var locations []models.Location
	err := tx.Find(&locations).Error
	return locations, err
}

func (r *FundRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Fund, error) {
	tx := r.db.Model(&models.Fund{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var funds []models.Fund
	err := tx.Find(&funds).Error
	return funds, err
}

func (r *ExpenseRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Expense, error) {
	tx := r.db.Model(&models.Expense{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var expenses []models.Expense
	err := tx.Find(&expenses).Error
	return expenses, err
}

func (r *TransactionRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Transaction, error) {
	tx := r.db.Model(&models.Transaction{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var transactions []models.Transaction
	err := tx.Find(&transactions).Error
	return transactions, err
}

func (r *PurchaseRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Purchase, error) {
	tx := r.db.Model(&models.Purchase{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var purchases []models.Purchase
	err := tx.Find(&purchases).Error
	return purchases, err
}

func (r *PayrollRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Payroll, error) {
	tx := r.db.Model(&models.Payroll{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var payrolls []models.Payroll
	err := tx.Find(&payrolls).Error
	return payrolls, err
}

func (r *InventoryRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Inventory, error) {
	tx := r.db.Model(&models.Inventory{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var inventories []models.Inventory
	err := tx.Find(&inventories).Error
	return inventories, err
}

func (r *GiftTypeRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.GiftType, error) {
	tx := r.db.Model(&models.GiftType{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var giftTypes []models.GiftType
	err := tx.Find(&giftTypes).Error
	return giftTypes, err
}

func (r *GiftRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Gift, error) {
	tx := r.db.Model(&models.Gift{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var gifts []models.Gift
	err := tx.Find(&gifts).Error
	return gifts, err
}

func (r *InventoryTransactionRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.InventoryTransaction, error) {
	tx := r.db.Model(&models.InventoryTransaction{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var transactions []models.InventoryTransaction
	err := tx.Find(&transactions).Error
	return transactions, err
}

func (r *DeliveryRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Delivery, error) {
	tx := r.db.Model(&models.Delivery{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var deliveries []models.Delivery
	err := tx.Find(&deliveries).Error
	return deliveries, err
}

func (r *VolunteerProjectRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.VolunteerProject, error) {
	tx := r.db.Model(&models.VolunteerProject{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var vps []models.VolunteerProject
	err := tx.Find(&vps).Error
	return vps, err
}

func (r *EmployeeProjectRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.EmployeeProject, error) {
	tx := r.db.Model(&models.EmployeeProject{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var eps []models.EmployeeProject
	err := tx.Find(&eps).Error
	return eps, err
}

func (r *FundProjectRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.FundProject, error) {
	tx := r.db.Model(&models.FundProject{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var fps []models.FundProject
	err := tx.Find(&fps).Error
	return fps, err
}

func (r *DonationInventoryRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.DonationInventory, error) {
	tx := r.db.Model(&models.DonationInventory{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var dis []models.DonationInventory
	err := tx.Find(&dis).Error
	return dis, err
}

func (r *DeliveryInventoryRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.DeliveryInventory, error) {
	tx := r.db.Model(&models.DeliveryInventory{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var dis []models.DeliveryInventory
	err := tx.Find(&dis).Error
	return dis, err
}

func (r *ScheduleRepository) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Schedule, error) {
	tx := r.db.Model(&models.Schedule{})
	tx = applyFilters(tx, query, numberRange, dateRange)
	var schedules []models.Schedule
	err := tx.Find(&schedules).Error
	return schedules, err
}
