package services

import (
	"erp-backend/internal/models"
	"erp-backend/internal/repo"
)

// ==================== Service Definitions ====================

// ProjectService 项目服务
type ProjectService struct {
	repo *repo.ProjectRepository
}

func NewProjectService(projectRepo *repo.ProjectRepository) *ProjectService {
	return &ProjectService{repo: projectRepo}
}

// DonorService 捐赠者服务
type DonorService struct {
	repo *repo.DonorRepository
}

func NewDonorService(donorRepo *repo.DonorRepository) *DonorService {
	return &DonorService{repo: donorRepo}
}

// DonationService 捐赠服务
type DonationService struct {
	repo *repo.DonationRepository
}

func NewDonationService(donationRepo *repo.DonationRepository) *DonationService {
	return &DonationService{repo: donationRepo}
}

// VolunteerService 志愿者服务
type VolunteerService struct {
	repo *repo.VolunteerRepository
}

func NewVolunteerService(volunteerRepo *repo.VolunteerRepository) *VolunteerService {
	return &VolunteerService{repo: volunteerRepo}
}

// EmployeeService 员工服务
type EmployeeService struct {
	repo *repo.EmployeeRepository
}

func NewEmployeeService(employeeRepo *repo.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: employeeRepo}
}

// LocationService 地点服务
type LocationService struct {
	repo *repo.LocationRepository
}

func NewLocationService(locationRepo *repo.LocationRepository) *LocationService {
	return &LocationService{repo: locationRepo}
}

// FundService 基金服务
type FundService struct {
	repo *repo.FundRepository
}

func NewFundService(fundRepo *repo.FundRepository) *FundService {
	return &FundService{repo: fundRepo}
}

// ExpenseService 支出服务
type ExpenseService struct {
	repo *repo.ExpenseRepository
}

func NewExpenseService(expenseRepo *repo.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: expenseRepo}
}

// TransactionService 交易服务
type TransactionService struct {
	repo *repo.TransactionRepository
}

func NewTransactionService(transactionRepo *repo.TransactionRepository) *TransactionService {
	return &TransactionService{repo: transactionRepo}
}

// PurchaseService 采购服务
type PurchaseService struct {
	repo *repo.PurchaseRepository
}

func NewPurchaseService(purchaseRepo *repo.PurchaseRepository) *PurchaseService {
	return &PurchaseService{repo: purchaseRepo}
}

// PayrollService 薪资服务
type PayrollService struct {
	repo *repo.PayrollRepository
}

func NewPayrollService(payrollRepo *repo.PayrollRepository) *PayrollService {
	return &PayrollService{repo: payrollRepo}
}

// InventoryService 库存服务
type InventoryService struct {
	repo *repo.InventoryRepository
}

func NewInventoryService(inventoryRepo *repo.InventoryRepository) *InventoryService {
	return &InventoryService{repo: inventoryRepo}
}

// GiftTypeService 礼品类型服务
type GiftTypeService struct {
	repo *repo.GiftTypeRepository
}

func NewGiftTypeService(giftTypeRepo *repo.GiftTypeRepository) *GiftTypeService {
	return &GiftTypeService{repo: giftTypeRepo}
}

// GiftService 礼品服务
type GiftService struct {
	repo *repo.GiftRepository
}

func NewGiftService(giftRepo *repo.GiftRepository) *GiftService {
	return &GiftService{repo: giftRepo}
}

// InventoryTransactionService 库存交易服务
type InventoryTransactionService struct {
	repo *repo.InventoryTransactionRepository
}

func NewInventoryTransactionService(inventoryTransactionRepo *repo.InventoryTransactionRepository) *InventoryTransactionService {
	return &InventoryTransactionService{repo: inventoryTransactionRepo}
}

// DeliveryService 配送服务
type DeliveryService struct {
	repo *repo.DeliveryRepository
}

func NewDeliveryService(deliveryRepo *repo.DeliveryRepository) *DeliveryService {
	return &DeliveryService{repo: deliveryRepo}
}

// VolunteerProjectService 志愿者-项目服务
type VolunteerProjectService struct {
	repo *repo.VolunteerProjectRepository
}

func NewVolunteerProjectService(volunteerProjectRepo *repo.VolunteerProjectRepository) *VolunteerProjectService {
	return &VolunteerProjectService{repo: volunteerProjectRepo}
}

// EmployeeProjectService 员工-项目服务
type EmployeeProjectService struct {
	repo *repo.EmployeeProjectRepository
}

func NewEmployeeProjectService(employeeProjectRepo *repo.EmployeeProjectRepository) *EmployeeProjectService {
	return &EmployeeProjectService{repo: employeeProjectRepo}
}

// FundProjectService 基金-项目服务
type FundProjectService struct {
	repo *repo.FundProjectRepository
}

func NewFundProjectService(fundProjectRepo *repo.FundProjectRepository) *FundProjectService {
	return &FundProjectService{repo: fundProjectRepo}
}

// DonationInventoryService 捐赠-库存服务
type DonationInventoryService struct {
	repo *repo.DonationInventoryRepository
}

func NewDonationInventoryService(donationInventoryRepo *repo.DonationInventoryRepository) *DonationInventoryService {
	return &DonationInventoryService{repo: donationInventoryRepo}
}

// DeliveryInventoryService 配送-库存服务
type DeliveryInventoryService struct {
	repo *repo.DeliveryInventoryRepository
}

func NewDeliveryInventoryService(deliveryInventoryRepo *repo.DeliveryInventoryRepository) *DeliveryInventoryService {
	return &DeliveryInventoryService{repo: deliveryInventoryRepo}
}

// ScheduleService 日程服务
type ScheduleService struct {
	repo *repo.ScheduleRepository
}

func NewScheduleService(scheduleRepo *repo.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: scheduleRepo}
}

// ==================== Project Service Methods ====================

func (s *ProjectService) Create(project *models.Project) error {
	return s.repo.Create(project)
}

func (s *ProjectService) GetAll() ([]models.Project, error) {
	return s.repo.GetAll()
}

func (s *ProjectService) Search(query map[string]interface{}) ([]models.Project, error) {
	return s.repo.Search(query)
}

func (s *ProjectService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Project, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *ProjectService) Update(project *models.Project) error {
	return s.repo.Update(project)
}

func (s *ProjectService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Donor Service Methods ====================

func (s *DonorService) Create(donor *models.Donor) error {
	return s.repo.Create(donor)
}

func (s *DonorService) GetAll() ([]models.Donor, error) {
	return s.repo.GetAll()
}

func (s *DonorService) Search(query map[string]interface{}) ([]models.Donor, error) {
	return s.repo.Search(query)
}

func (s *DonorService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Donor, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *DonorService) Update(donor *models.Donor) error {
	return s.repo.Update(donor)
}

func (s *DonorService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Donation Service Methods ====================

func (s *DonationService) Create(donation *models.Donation) error {
	return s.repo.Create(donation)
}

func (s *DonationService) GetAll() ([]models.Donation, error) {
	return s.repo.GetAll()
}

func (s *DonationService) Search(query map[string]interface{}) ([]models.Donation, error) {
	return s.repo.Search(query)
}

func (s *DonationService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Donation, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *DonationService) Update(donation *models.Donation) error {
	return s.repo.Update(donation)
}

func (s *DonationService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Volunteer Service Methods ====================

func (s *VolunteerService) Create(volunteer *models.Volunteer) error {
	return s.repo.Create(volunteer)
}

func (s *VolunteerService) GetAll() ([]models.Volunteer, error) {
	return s.repo.GetAll()
}

func (s *VolunteerService) Search(query map[string]interface{}) ([]models.Volunteer, error) {
	return s.repo.Search(query)
}

func (s *VolunteerService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Volunteer, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *VolunteerService) Update(volunteer *models.Volunteer) error {
	return s.repo.Update(volunteer)
}

func (s *VolunteerService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Employee Service Methods ====================

func (s *EmployeeService) Create(employee *models.Employee) error {
	return s.repo.Create(employee)
}

func (s *EmployeeService) GetAll() ([]models.Employee, error) {
	return s.repo.GetAll()
}

func (s *EmployeeService) Search(query map[string]interface{}) ([]models.Employee, error) {
	return s.repo.Search(query)
}

func (s *EmployeeService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Employee, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *EmployeeService) Update(employee *models.Employee) error {
	return s.repo.Update(employee)
}

func (s *EmployeeService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Location Service Methods ====================

func (s *LocationService) Create(location *models.Location) error {
	return s.repo.Create(location)
}

func (s *LocationService) GetAll() ([]models.Location, error) {
	return s.repo.GetAll()
}

func (s *LocationService) Search(query map[string]interface{}) ([]models.Location, error) {
	return s.repo.Search(query)
}

func (s *LocationService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Location, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *LocationService) Update(location *models.Location) error {
	return s.repo.Update(location)
}

func (s *LocationService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Fund Service Methods ====================

func (s *FundService) Create(fund *models.Fund) error {
	return s.repo.Create(fund)
}

func (s *FundService) GetAll() ([]models.Fund, error) {
	return s.repo.GetAll()
}

func (s *FundService) Search(query map[string]interface{}) ([]models.Fund, error) {
	return s.repo.Search(query)
}

func (s *FundService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Fund, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *FundService) Update(fund *models.Fund) error {
	return s.repo.Update(fund)
}

func (s *FundService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Expense Service Methods ====================

func (s *ExpenseService) Create(expense *models.Expense) error {
	return s.repo.Create(expense)
}

func (s *ExpenseService) GetAll() ([]models.Expense, error) {
	return s.repo.GetAll()
}

func (s *ExpenseService) Search(query map[string]interface{}) ([]models.Expense, error) {
	return s.repo.Search(query)
}

func (s *ExpenseService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Expense, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *ExpenseService) Update(expense *models.Expense) error {
	return s.repo.Update(expense)
}

func (s *ExpenseService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Transaction Service Methods ====================

func (s *TransactionService) Create(transaction *models.Transaction) error {
	return s.repo.Create(transaction)
}

func (s *TransactionService) GetAll() ([]models.Transaction, error) {
	return s.repo.GetAll()
}

func (s *TransactionService) Search(query map[string]interface{}) ([]models.Transaction, error) {
	return s.repo.Search(query)
}

func (s *TransactionService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Transaction, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *TransactionService) Update(transaction *models.Transaction) error {
	return s.repo.Update(transaction)
}

func (s *TransactionService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Purchase Service Methods ====================

func (s *PurchaseService) Create(purchase *models.Purchase) error {
	return s.repo.Create(purchase)
}

func (s *PurchaseService) GetAll() ([]models.Purchase, error) {
	return s.repo.GetAll()
}

func (s *PurchaseService) Search(query map[string]interface{}) ([]models.Purchase, error) {
	return s.repo.Search(query)
}

func (s *PurchaseService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Purchase, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *PurchaseService) Update(purchase *models.Purchase) error {
	return s.repo.Update(purchase)
}

func (s *PurchaseService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Payroll Service Methods ====================

func (s *PayrollService) Create(payroll *models.Payroll) error {
	return s.repo.Create(payroll)
}

func (s *PayrollService) GetAll() ([]models.Payroll, error) {
	return s.repo.GetAll()
}

func (s *PayrollService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Payroll, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *PayrollService) Search(query map[string]interface{}) ([]models.Payroll, error) {
	return s.repo.Search(query)
}

func (s *PayrollService) Update(payroll *models.Payroll) error {
	return s.repo.Update(payroll)
}

func (s *PayrollService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Inventory Service Methods ====================

func (s *InventoryService) Create(inventory *models.Inventory) error {
	return s.repo.Create(inventory)
}

func (s *InventoryService) GetAll() ([]models.Inventory, error) {
	return s.repo.GetAll()
}

func (s *InventoryService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Inventory, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *InventoryService) Search(query map[string]interface{}) ([]models.Inventory, error) {
	return s.repo.Search(query)
}

func (s *InventoryService) Update(inventory *models.Inventory) error {
	return s.repo.Update(inventory)
}

func (s *InventoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== GiftType Service Methods ====================

func (s *GiftTypeService) Create(giftType *models.GiftType) error {
	return s.repo.Create(giftType)
}

func (s *GiftTypeService) GetAll() ([]models.GiftType, error) {
	return s.repo.GetAll()
}

func (s *GiftTypeService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.GiftType, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *GiftTypeService) Search(query map[string]interface{}) ([]models.GiftType, error) {
	return s.repo.Search(query)
}

func (s *GiftTypeService) Update(giftType *models.GiftType) error {
	return s.repo.Update(giftType)
}

func (s *GiftTypeService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Gift Service Methods ====================

func (s *GiftService) Create(gift *models.Gift) error {
	return s.repo.Create(gift)
}

func (s *GiftService) GetAll() ([]models.Gift, error) {
	return s.repo.GetAll()
}

func (s *GiftService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Gift, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *GiftService) Search(query map[string]interface{}) ([]models.Gift, error) {
	return s.repo.Search(query)
}

func (s *GiftService) Update(gift *models.Gift) error {
	return s.repo.Update(gift)
}

func (s *GiftService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== InventoryTransaction Service Methods ====================

func (s *InventoryTransactionService) Create(transaction *models.InventoryTransaction) error {
	return s.repo.Create(transaction)
}

func (s *InventoryTransactionService) GetAll() ([]models.InventoryTransaction, error) {
	return s.repo.GetAll()
}

func (s *InventoryTransactionService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.InventoryTransaction, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *InventoryTransactionService) Search(query map[string]interface{}) ([]models.InventoryTransaction, error) {
	return s.repo.Search(query)
}

func (s *InventoryTransactionService) Update(transaction *models.InventoryTransaction) error {
	return s.repo.Update(transaction)
}

func (s *InventoryTransactionService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Delivery Service Methods ====================

func (s *DeliveryService) Create(delivery *models.Delivery) error {
	return s.repo.Create(delivery)
}

func (s *DeliveryService) GetAll() ([]models.Delivery, error) {
	return s.repo.GetAll()
}

func (s *DeliveryService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Delivery, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *DeliveryService) Search(query map[string]interface{}) ([]models.Delivery, error) {
	return s.repo.Search(query)
}

func (s *DeliveryService) Update(delivery *models.Delivery) error {
	return s.repo.Update(delivery)
}

func (s *DeliveryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== VolunteerProject Service Methods ====================

func (s *VolunteerProjectService) Create(vp *models.VolunteerProject) error {
	return s.repo.Create(vp)
}

func (s *VolunteerProjectService) GetAll() ([]models.VolunteerProject, error) {
	return s.repo.GetAll()
}

func (s *VolunteerProjectService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.VolunteerProject, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *VolunteerProjectService) Search(query map[string]interface{}) ([]models.VolunteerProject, error) {
	return s.repo.Search(query)
}

func (s *VolunteerProjectService) Update(vp *models.VolunteerProject) error {
	return s.repo.Update(vp)
}

func (s *VolunteerProjectService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== EmployeeProject Service Methods ====================

func (s *EmployeeProjectService) Create(ep *models.EmployeeProject) error {
	return s.repo.Create(ep)
}

func (s *EmployeeProjectService) GetAll() ([]models.EmployeeProject, error) {
	return s.repo.GetAll()
}

func (s *EmployeeProjectService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.EmployeeProject, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *EmployeeProjectService) Search(query map[string]interface{}) ([]models.EmployeeProject, error) {
	return s.repo.Search(query)
}

func (s *EmployeeProjectService) Update(ep *models.EmployeeProject) error {
	return s.repo.Update(ep)
}

func (s *EmployeeProjectService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== FundProject Service Methods ====================

func (s *FundProjectService) Create(fp *models.FundProject) error {
	return s.repo.Create(fp)
}

func (s *FundProjectService) GetAll() ([]models.FundProject, error) {
	return s.repo.GetAll()
}

func (s *FundProjectService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.FundProject, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *FundProjectService) Search(query map[string]interface{}) ([]models.FundProject, error) {
	return s.repo.Search(query)
}

func (s *FundProjectService) Update(fp *models.FundProject) error {
	return s.repo.Update(fp)
}

func (s *FundProjectService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== DonationInventory Service Methods ====================

func (s *DonationInventoryService) Create(di *models.DonationInventory) error {
	return s.repo.Create(di)
}

func (s *DonationInventoryService) GetAll() ([]models.DonationInventory, error) {
	return s.repo.GetAll()
}

func (s *DonationInventoryService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.DonationInventory, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *DonationInventoryService) Search(query map[string]interface{}) ([]models.DonationInventory, error) {
	return s.repo.Search(query)
}

func (s *DonationInventoryService) Update(di *models.DonationInventory) error {
	return s.repo.Update(di)
}

func (s *DonationInventoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== DeliveryInventory Service Methods ====================

func (s *DeliveryInventoryService) Create(di *models.DeliveryInventory) error {
	return s.repo.Create(di)
}

func (s *DeliveryInventoryService) GetAll() ([]models.DeliveryInventory, error) {
	return s.repo.GetAll()
}

func (s *DeliveryInventoryService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.DeliveryInventory, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *DeliveryInventoryService) Search(query map[string]interface{}) ([]models.DeliveryInventory, error) {
	return s.repo.Search(query)
}

func (s *DeliveryInventoryService) Update(di *models.DeliveryInventory) error {
	return s.repo.Update(di)
}

func (s *DeliveryInventoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ==================== Schedule Service Methods ====================

func (s *ScheduleService) Create(schedule *models.Schedule) error {
	return s.repo.Create(schedule)
}

func (s *ScheduleService) GetAll() ([]models.Schedule, error) {
	return s.repo.GetAll()
}

func (s *ScheduleService) Filter(query map[string]interface{}, numberRange map[string][]interface{}, dateRange map[string][]string) ([]models.Schedule, error) {
	return s.repo.Filter(query, numberRange, dateRange)
}

func (s *ScheduleService) Search(query map[string]interface{}) ([]models.Schedule, error) {
	return s.repo.Search(query)
}

func (s *ScheduleService) Update(schedule *models.Schedule) error {
	return s.repo.Update(schedule)
}

func (s *ScheduleService) Delete(id uint) error {
	return s.repo.Delete(id)
}
