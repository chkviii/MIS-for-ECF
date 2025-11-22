package services

import (
	"errors"
	"fmt"
	"time"

	"mypage-backend/internal/models"
	"mypage-backend/internal/repo"
)

type ProjectService struct {
	projectRepo  *repo.ProjectRepository
	locationRepo *repo.LocationRepository
}

func NewProjectService(projectRepo *repo.ProjectRepository, locationRepo *repo.LocationRepository) *ProjectService {
	return &ProjectService{
		projectRepo:  projectRepo,
		locationRepo: locationRepo,
	}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(project *models.Project) error {
	if project.Name == "" {
		return errors.New("project name cannot be empty")
	}

	if project.Budget < 0 {
		return errors.New("project budget cannot be negative")
	}

	// Validate location
	if project.LocationID != nil {
		_, err := s.locationRepo.GetByID(*project.LocationID)
		if (err != nil) {
			return errors.New("invalid location ID")
		}
	}

	// Generate project ID
	project.ProjectID = fmt.Sprintf("PRJ-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	project.Status = "planning"

	return s.projectRepo.Create(project)
}

// GetProject retrieves a project
func (s *ProjectService) GetProject(id uint) (*models.Project, error) {
	return s.projectRepo.GetByID(id)
}

// GetAllProjects retrieves all projects
func (s *ProjectService) GetAllProjects() ([]models.Project, error) {
	return s.projectRepo.GetAll()
}

// GetProjectsByStatus retrieves projects by status
func (s *ProjectService) GetProjectsByStatus(status string) ([]models.Project, error) {
	return s.projectRepo.GetByStatus(status)
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(project *models.Project) error {
	// Validate project exists
	_, err := s.projectRepo.GetByID(project.ID)
	if err != nil {
		return errors.New("project not found")
	}

	return s.projectRepo.Update(project)
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(id uint) error {
	return s.projectRepo.Delete(id)
}

// AssignVolunteer assigns a volunteer to a project
func (s *ProjectService) AssignVolunteer(volunteerID, projectID uint, assignment *models.VolunteerProject) error {
	assignment.VolunteerID = volunteerID
	assignment.ProjectID = projectID
	assignment.Status = "active"
	return s.projectRepo.AssignVolunteer(assignment)
}

// AssignEmployee assigns an employee to a project
func (s *ProjectService) AssignEmployee(employeeID, projectID uint, assignment *models.EmployeeProject) error {
	assignment.EmployeeID = employeeID
	assignment.ProjectID = projectID
	return s.projectRepo.AssignEmployee(assignment)
}

// DonorService Donor service
type DonorService struct {
	donorRepo *repo.DonorRepository
}

func NewDonorService(donorRepo *repo.DonorRepository) *DonorService {
	return &DonorService{donorRepo: donorRepo}
}

func (s *DonorService) CreateDonor(donor *models.Donor) error {
	if donor.FirstName == "" || donor.LastName == "" {
		return errors.New("donor name cannot be empty")
	}

	// Generate donor ID
	donor.DonorID = fmt.Sprintf("DNR-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	donor.Status = "active"

	return s.donorRepo.Create(donor)
}

func (s *DonorService) GetDonor(id uint) (*models.Donor, error) {
	return s.donorRepo.GetByID(id)
}

func (s *DonorService) GetAllDonors() ([]models.Donor, error) {
	return s.donorRepo.GetAll()
}

func (s *DonorService) UpdateDonor(donor *models.Donor) error {
	return s.donorRepo.Update(donor)
}

func (s *DonorService) DeleteDonor(id uint) error {
	return s.donorRepo.Delete(id)
}

func (s *DonorService) GetDonorWithDonations(id uint) (*models.Donor, error) {
	return s.donorRepo.GetWithDonations(id)
}

func (s *DonorService) SearchDonors(query map[string]interface{}, startDate, endDate *time.Time) ([]models.Donor, error) {
	return s.donorRepo.Search(query, startDate, endDate)
}

// DonationService Donation service
type DonationService struct {
	donationRepo *repo.DonationRepository
	donorRepo    *repo.DonorRepository
}

func NewDonationService(donationRepo *repo.DonationRepository, donorRepo *repo.DonorRepository) *DonationService {
	return &DonationService{
		donationRepo: donationRepo,
		donorRepo:    donorRepo,
	}
}

func (s *DonationService) CreateDonation(donation *models.Donation) error {
	if donation.Amount <= 0 {
		return errors.New("donation amount must be greater than 0")
	}

	// Validate donor exists
	_, err := s.donorRepo.GetByID(donation.DonorID)
	if err != nil {
		return errors.New("donor not found")
	}

	// Generate donation ID
	donation.DonationID = fmt.Sprintf("DON-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	
	// If no donation date provided, use current time
	if donation.DonationDate.IsZero() {
		donation.DonationDate = time.Now()
	}

	return s.donationRepo.Create(donation)
}

func (s *DonationService) GetDonation(id uint) (*models.Donation, error) {
	return s.donationRepo.GetByID(id)
}

func (s *DonationService) GetAllDonations() ([]models.Donation, error) {
	return s.donationRepo.GetAll()
}

func (s *DonationService) GetDonationsByDonor(donorID uint) ([]models.Donation, error) {
	return s.donationRepo.GetByDonorID(donorID)
}

func (s *DonationService) SearchDonations(query map[string]interface{}, startDate, endDate *time.Time) ([]models.Donation, error) {
	return s.donationRepo.Search(query, startDate, endDate)
}

func (s *DonationService) UpdateDonation(donation *models.Donation) error {
	return s.donationRepo.Update(donation)
}

func (s *DonationService) DeleteDonation(id uint) error {
	return s.donationRepo.Delete(id)
}

// VolunteerService Volunteer service
type VolunteerService struct {
	volunteerRepo *repo.VolunteerRepository
}

func NewVolunteerService(volunteerRepo *repo.VolunteerRepository) *VolunteerService {
	return &VolunteerService{volunteerRepo: volunteerRepo}
}

func (s *VolunteerService) CreateVolunteer(volunteer *models.Volunteer) error {
	if volunteer.FirstName == "" || volunteer.LastName == "" {
		return errors.New("volunteer name cannot be empty")
	}

	// Generate volunteer ID
	volunteer.VolunteerID = fmt.Sprintf("VOL-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	volunteer.Status = "active"

	return s.volunteerRepo.Create(volunteer)
}

func (s *VolunteerService) GetVolunteer(id uint) (*models.Volunteer, error) {
	return s.volunteerRepo.GetByID(id)
}

func (s *VolunteerService) GetAllVolunteers() ([]models.Volunteer, error) {
	return s.volunteerRepo.GetAll()
}

func (s *VolunteerService) UpdateVolunteer(volunteer *models.Volunteer) error {
	return s.volunteerRepo.Update(volunteer)
}

func (s *VolunteerService) DeleteVolunteer(id uint) error {
	return s.volunteerRepo.Delete(id)
}

func (s *VolunteerService) SearchVolunteers(query map[string]interface{}) ([]models.Volunteer, error) {
	return s.volunteerRepo.Search(query)
}

// EmployeeService Employee service
type EmployeeService struct {
	employeeRepo *repo.EmployeeRepository
}

func NewEmployeeService(employeeRepo *repo.EmployeeRepository) *EmployeeService {
	return &EmployeeService{employeeRepo: employeeRepo}
}

func (s *EmployeeService) CreateEmployee(employee *models.Employee) error {
	if employee.FirstName == "" || employee.LastName == "" {
		return errors.New("employee name cannot be empty")
	}

	// Generate employee ID
	employee.EmployeeID = fmt.Sprintf("EMP-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	employee.Status = "active"

	return s.employeeRepo.Create(employee)
}

func (s *EmployeeService) GetEmployee(id uint) (*models.Employee, error) {
	return s.employeeRepo.GetByID(id)
}

func (s *EmployeeService) GetAllEmployees() ([]models.Employee, error) {
	return s.employeeRepo.GetAll()
}

func (s *EmployeeService) UpdateEmployee(employee *models.Employee) error {
	return s.employeeRepo.Update(employee)
}

func (s *EmployeeService) DeleteEmployee(id uint) error {
	return s.employeeRepo.Delete(id)
}

func (s *EmployeeService) SearchEmployees(query map[string]interface{}) ([]models.Employee, error) {
	return s.employeeRepo.Search(query)
}

// LocationService Location service
type LocationService struct {
	locationRepo *repo.LocationRepository
}

func NewLocationService(locationRepo *repo.LocationRepository) *LocationService {
	return &LocationService{locationRepo: locationRepo}
}

func (s *LocationService) CreateLocation(location *models.Location) error {
	if location.Name == "" {
		return errors.New("location name cannot be empty")
	}

	// Generate location ID
	location.LocationID = fmt.Sprintf("LOC-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	return s.locationRepo.Create(location)
}

func (s *LocationService) GetLocation(id uint) (*models.Location, error) {
	return s.locationRepo.GetByID(id)
}

func (s *LocationService) GetAllLocations() ([]models.Location, error) {
	return s.locationRepo.GetAll()
}

func (s *LocationService) UpdateLocation(location *models.Location) error {
	return s.locationRepo.Update(location)
}

func (s *LocationService) DeleteLocation(id uint) error {
	return s.locationRepo.Delete(id)
}

func (s *LocationService) SearchLocations(query map[string]interface{}) ([]models.Location, error) {
	return s.locationRepo.Search(query)
}

func (s *ProjectService) SearchProjects(query map[string]interface{}, startDate, endDate *time.Time) ([]models.Project, error) {
	return s.projectRepo.Search(query, startDate, endDate)
}

// FundService Fund service
type FundService struct {
	fundRepo *repo.FundRepository
}

func NewFundService(fundRepo *repo.FundRepository) *FundService {
	return &FundService{fundRepo: fundRepo}
}

func (s *FundService) CreateFund(fund *models.Fund) error {
	if fund.Name == "" {
		return errors.New("fund name cannot be empty")
	}
	fund.FundID = fmt.Sprintf("FND-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	fund.Status = "active"
	return s.fundRepo.Create(fund)
}

func (s *FundService) GetFund(id uint) (*models.Fund, error) {
	return s.fundRepo.GetByID(id)
}

func (s *FundService) GetAllFunds() ([]models.Fund, error) {
	return s.fundRepo.GetAll()
}

func (s *FundService) SearchFunds(query map[string]interface{}, startDate, endDate *time.Time) ([]models.Fund, error) {
	return s.fundRepo.Search(query, startDate, endDate)
}

func (s *FundService) UpdateFund(fund *models.Fund) error {
	return s.fundRepo.Update(fund)
}

func (s *FundService) DeleteFund(id uint) error {
	return s.fundRepo.Delete(id)
}

// ExpenseService Expense service
type ExpenseService struct {
	expenseRepo *repo.ExpenseRepository
}

func NewExpenseService(expenseRepo *repo.ExpenseRepository) *ExpenseService {
	return &ExpenseService{expenseRepo: expenseRepo}
}

func (s *ExpenseService) CreateExpense(expense *models.Expense) error {
	if expense.Amount <= 0 {
		return errors.New("expense amount must be greater than 0")
	}
	expense.ExpenseID = fmt.Sprintf("EXP-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	expense.ApprovalStatus = "pending"
	return s.expenseRepo.Create(expense)
}

func (s *ExpenseService) GetExpense(id uint) (*models.Expense, error) {
	return s.expenseRepo.GetByID(id)
}

func (s *ExpenseService) GetAllExpenses() ([]models.Expense, error) {
	return s.expenseRepo.GetAll()
}

func (s *ExpenseService) SearchExpenses(query map[string]interface{}, startDate, endDate *time.Time) ([]models.Expense, error) {
	return s.expenseRepo.Search(query, startDate, endDate)
}

func (s *ExpenseService) UpdateExpense(expense *models.Expense) error {
	return s.expenseRepo.Update(expense)
}

func (s *ExpenseService) DeleteExpense(id uint) error {
	return s.expenseRepo.Delete(id)
}

// TransactionService Transaction service
type TransactionService struct {
	transactionRepo *repo.TransactionRepository
}

func NewTransactionService(transactionRepo *repo.TransactionRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo}
}

func (s *TransactionService) CreateTransaction(transaction *models.Transaction) error {
	if transaction.Amount <= 0 {
		return errors.New("transaction amount must be greater than 0")
	}
	transaction.TransactionID = fmt.Sprintf("TXN-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	return s.transactionRepo.Create(transaction)
}

func (s *TransactionService) GetTransaction(id uint) (*models.Transaction, error) {
	return s.transactionRepo.GetByID(id)
}

func (s *TransactionService) GetAllTransactions() ([]models.Transaction, error) {
	return s.transactionRepo.GetAll()
}

func (s *TransactionService) SearchTransactions(query map[string]interface{}, startDate, endDate *time.Time) ([]models.Transaction, error) {
	return s.transactionRepo.Search(query, startDate, endDate)
}

func (s *TransactionService) UpdateTransaction(transaction *models.Transaction) error {
	return s.transactionRepo.Update(transaction)
}

func (s *TransactionService) DeleteTransaction(id uint) error {
	return s.transactionRepo.Delete(id)
}

// PurchaseService Purchase service
type PurchaseService struct {
	purchaseRepo *repo.PurchaseRepository
}

func NewPurchaseService(purchaseRepo *repo.PurchaseRepository) *PurchaseService {
	return &PurchaseService{purchaseRepo: purchaseRepo}
}

func (s *PurchaseService) CreatePurchase(purchase *models.Purchase) error {
	if purchase.TotalSpent <= 0 {
		return errors.New("purchase amount must be greater than 0")
	}
	purchase.PurchaseID = fmt.Sprintf("PUR-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	return s.purchaseRepo.Create(purchase)
}

func (s *PurchaseService) GetPurchase(id uint) (*models.Purchase, error) {
	return s.purchaseRepo.GetByID(id)
}

func (s *PurchaseService) GetAllPurchases() ([]models.Purchase, error) {
	return s.purchaseRepo.GetAll()
}

func (s *PurchaseService) SearchPurchases(query map[string]interface{}, startDate, endDate *time.Time) ([]models.Purchase, error) {
	return s.purchaseRepo.Search(query, startDate, endDate)
}

func (s *PurchaseService) UpdatePurchase(purchase *models.Purchase) error {
	return s.purchaseRepo.Update(purchase)
}

func (s *PurchaseService) DeletePurchase(id uint) error {
	return s.purchaseRepo.Delete(id)
}

// PayrollService Payroll service
type PayrollService struct {
	payrollRepo *repo.PayrollRepository
}

func NewPayrollService(payrollRepo *repo.PayrollRepository) *PayrollService {
	return &PayrollService{payrollRepo: payrollRepo}
}

func (s *PayrollService) CreatePayroll(payroll *models.Payroll) error {
	if payroll.Amount <= 0 {
		return errors.New("payroll amount must be greater than 0")
	}
	return s.payrollRepo.Create(payroll)
}

func (s *PayrollService) GetPayroll(id uint) (*models.Payroll, error) {
	return s.payrollRepo.GetByID(id)
}

func (s *PayrollService) GetAllPayrolls() ([]models.Payroll, error) {
	return s.payrollRepo.GetAll()
}

func (s *PayrollService) SearchPayrolls(query map[string]interface{}, startDate, endDate *time.Time) ([]models.Payroll, error) {
	return s.payrollRepo.Search(query, startDate, endDate)
}

func (s *PayrollService) UpdatePayroll(payroll *models.Payroll) error {
	return s.payrollRepo.Update(payroll)
}

func (s *PayrollService) DeletePayroll(id uint) error {
	return s.payrollRepo.Delete(id)
}

// InventoryService Inventory service
type InventoryService struct {
	inventoryRepo *repo.InventoryRepository
}

func NewInventoryService(inventoryRepo *repo.InventoryRepository) *InventoryService {
	return &InventoryService{inventoryRepo: inventoryRepo}
}

func (s *InventoryService) CreateInventory(inventory *models.Inventory) error {
	if inventory.Name == "" {
		return errors.New("inventory name cannot be empty")
	}
	inventory.InventoryID = fmt.Sprintf("INV-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	inventory.Status = "available"
	return s.inventoryRepo.Create(inventory)
}

func (s *InventoryService) GetInventory(id uint) (*models.Inventory, error) {
	return s.inventoryRepo.GetByID(id)
}

func (s *InventoryService) GetAllInventories() ([]models.Inventory, error) {
	return s.inventoryRepo.GetAll()
}

func (s *InventoryService) SearchInventories(query map[string]interface{}) ([]models.Inventory, error) {
	return s.inventoryRepo.Search(query)
}

func (s *InventoryService) UpdateInventory(inventory *models.Inventory) error {
	return s.inventoryRepo.Update(inventory)
}

func (s *InventoryService) DeleteInventory(id uint) error {
	return s.inventoryRepo.Delete(id)
}

// GiftTypeService Gift type service
type GiftTypeService struct {
	giftTypeRepo *repo.GiftTypeRepository
}

func NewGiftTypeService(giftTypeRepo *repo.GiftTypeRepository) *GiftTypeService {
	return &GiftTypeService{giftTypeRepo: giftTypeRepo}
}

func (s *GiftTypeService) CreateGiftType(giftType *models.GiftType) error {
	if giftType.Name == "" {
		return errors.New("gift type name cannot be empty")
	}
	return s.giftTypeRepo.Create(giftType)
}

func (s *GiftTypeService) GetGiftType(id uint) (*models.GiftType, error) {
	return s.giftTypeRepo.GetByID(id)
}

func (s *GiftTypeService) GetAllGiftTypes() ([]models.GiftType, error) {
	return s.giftTypeRepo.GetAll()
}

func (s *GiftTypeService) UpdateGiftType(giftType *models.GiftType) error {
	return s.giftTypeRepo.Update(giftType)
}

func (s *GiftTypeService) DeleteGiftType(id uint) error {
	return s.giftTypeRepo.Delete(id)
}

// GiftService Gift service
type GiftService struct {
	giftRepo *repo.GiftRepository
}

func NewGiftService(giftRepo *repo.GiftRepository) *GiftService {
	return &GiftService{giftRepo: giftRepo}
}

func (s *GiftService) CreateGift(gift *models.Gift) error {
	gift.GiftID = fmt.Sprintf("GFT-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	gift.DistributionStatus = "pending"
	return s.giftRepo.Create(gift)
}

func (s *GiftService) GetGift(id uint) (*models.Gift, error) {
	return s.giftRepo.GetByID(id)
}

func (s *GiftService) GetAllGifts() ([]models.Gift, error) {
	return s.giftRepo.GetAll()
}

func (s *GiftService) SearchGifts(query map[string]interface{}, startDate, endDate *time.Time) ([]models.Gift, error) {
	return s.giftRepo.Search(query, startDate, endDate)
}

func (s *GiftService) UpdateGift(gift *models.Gift) error {
	return s.giftRepo.Update(gift)
}

func (s *GiftService) DeleteGift(id uint) error {
	return s.giftRepo.Delete(id)
}

// InventoryTransactionService Inventory transaction service
type InventoryTransactionService struct {
	inventoryTransactionRepo *repo.InventoryTransactionRepository
}

func NewInventoryTransactionService(inventoryTransactionRepo *repo.InventoryTransactionRepository) *InventoryTransactionService {
	return &InventoryTransactionService{inventoryTransactionRepo: inventoryTransactionRepo}
}

func (s *InventoryTransactionService) CreateInventoryTransaction(transaction *models.InventoryTransaction) error {
	transaction.TransactionID = fmt.Sprintf("ITX-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	return s.inventoryTransactionRepo.Create(transaction)
}

func (s *InventoryTransactionService) GetInventoryTransaction(id uint) (*models.InventoryTransaction, error) {
	return s.inventoryTransactionRepo.GetByID(id)
}

func (s *InventoryTransactionService) GetAllInventoryTransactions() ([]models.InventoryTransaction, error) {
	return s.inventoryTransactionRepo.GetAll()
}

func (s *InventoryTransactionService) SearchInventoryTransactions(query map[string]interface{}, startDate, endDate *time.Time) ([]models.InventoryTransaction, error) {
	return s.inventoryTransactionRepo.Search(query, startDate, endDate)
}

func (s *InventoryTransactionService) UpdateInventoryTransaction(transaction *models.InventoryTransaction) error {
	return s.inventoryTransactionRepo.Update(transaction)
}

func (s *InventoryTransactionService) DeleteInventoryTransaction(id uint) error {
	return s.inventoryTransactionRepo.Delete(id)
}

// DeliveryService Delivery service
type DeliveryService struct {
	deliveryRepo *repo.DeliveryRepository
}

func NewDeliveryService(deliveryRepo *repo.DeliveryRepository) *DeliveryService {
	return &DeliveryService{deliveryRepo: deliveryRepo}
}

func (s *DeliveryService) CreateDelivery(delivery *models.Delivery) error {
	delivery.DeliveryID = fmt.Sprintf("DLV-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	delivery.Status = "pending"
	return s.deliveryRepo.Create(delivery)
}

func (s *DeliveryService) GetDelivery(id uint) (*models.Delivery, error) {
	return s.deliveryRepo.GetByID(id)
}

func (s *DeliveryService) GetAllDeliveries() ([]models.Delivery, error) {
	return s.deliveryRepo.GetAll()
}

func (s *DeliveryService) SearchDeliveries(query map[string]interface{}, startDate, endDate *time.Time) ([]models.Delivery, error) {
	return s.deliveryRepo.Search(query, startDate, endDate)
}

func (s *DeliveryService) UpdateDelivery(delivery *models.Delivery) error {
	return s.deliveryRepo.Update(delivery)
}

func (s *DeliveryService) DeleteDelivery(id uint) error {
	return s.deliveryRepo.Delete(id)
}

// VolunteerProjectService Volunteer-Project association service
type VolunteerProjectService struct {
	repo *repo.VolunteerProjectRepository
}

func NewVolunteerProjectService(repo *repo.VolunteerProjectRepository) *VolunteerProjectService {
	return &VolunteerProjectService{repo: repo}
}

func (s *VolunteerProjectService) Create(vp *models.VolunteerProject) error {
	if vp.VolunteerID == 0 || vp.ProjectID == 0 {
		return errors.New("volunteer ID and project ID cannot be empty")
	}
	vp.Status = "active"
	return s.repo.Create(vp)
}

func (s *VolunteerProjectService) GetByID(id uint) (*models.VolunteerProject, error) {
	return s.repo.GetByID(id)
}

func (s *VolunteerProjectService) GetAll() ([]models.VolunteerProject, error) {
	return s.repo.GetAll()
}

func (s *VolunteerProjectService) GetByVolunteerID(volunteerID uint) ([]models.VolunteerProject, error) {
	return s.repo.GetByVolunteerID(volunteerID)
}

func (s *VolunteerProjectService) GetByProjectID(projectID uint) ([]models.VolunteerProject, error) {
	return s.repo.GetByProjectID(projectID)
}

func (s *VolunteerProjectService) Update(vp *models.VolunteerProject) error {
	return s.repo.Update(vp)
}

func (s *VolunteerProjectService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// EmployeeProjectService Employee-Project association service
type EmployeeProjectService struct {
	repo *repo.EmployeeProjectRepository
}

func NewEmployeeProjectService(repo *repo.EmployeeProjectRepository) *EmployeeProjectService {
	return &EmployeeProjectService{repo: repo}
}

func (s *EmployeeProjectService) Create(ep *models.EmployeeProject) error {
	if ep.EmployeeID == 0 || ep.ProjectID == 0 {
		return errors.New("employee ID and project ID cannot be empty")
	}
	return s.repo.Create(ep)
}

func (s *EmployeeProjectService) GetByID(id uint) (*models.EmployeeProject, error) {
	return s.repo.GetByID(id)
}

func (s *EmployeeProjectService) GetAll() ([]models.EmployeeProject, error) {
	return s.repo.GetAll()
}

func (s *EmployeeProjectService) GetByEmployeeID(employeeID uint) ([]models.EmployeeProject, error) {
	return s.repo.GetByEmployeeID(employeeID)
}

func (s *EmployeeProjectService) GetByProjectID(projectID uint) ([]models.EmployeeProject, error) {
	return s.repo.GetByProjectID(projectID)
}

func (s *EmployeeProjectService) Update(ep *models.EmployeeProject) error {
	return s.repo.Update(ep)
}

func (s *EmployeeProjectService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// FundProjectService Fund-Project association service
type FundProjectService struct {
	repo *repo.FundProjectRepository
}

func NewFundProjectService(repo *repo.FundProjectRepository) *FundProjectService {
	return &FundProjectService{repo: repo}
}

func (s *FundProjectService) Create(fp *models.FundProject) error {
	if fp.TransactionID == 0 || fp.ProjectID == 0 || fp.FundID == 0 {
		return errors.New("transaction ID, project ID and fund ID cannot be empty")
	}
	if fp.AllocatedAmount <= 0 {
		return errors.New("allocated amount must be greater than 0")
	}
	return s.repo.Create(fp)
}

func (s *FundProjectService) GetByID(id uint) (*models.FundProject, error) {
	return s.repo.GetByID(id)
}

func (s *FundProjectService) GetAll() ([]models.FundProject, error) {
	return s.repo.GetAll()
}

func (s *FundProjectService) GetByProjectID(projectID uint) ([]models.FundProject, error) {
	return s.repo.GetByProjectID(projectID)
}

func (s *FundProjectService) GetByFundID(fundID uint) ([]models.FundProject, error) {
	return s.repo.GetByFundID(fundID)
}

func (s *FundProjectService) Update(fp *models.FundProject) error {
	return s.repo.Update(fp)
}

func (s *FundProjectService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// DonationInventoryService Donation-Inventory association service (in-kind donation)
type DonationInventoryService struct {
	repo *repo.DonationInventoryRepository
}

func NewDonationInventoryService(repo *repo.DonationInventoryRepository) *DonationInventoryService {
	return &DonationInventoryService{repo: repo}
}

func (s *DonationInventoryService) Create(di *models.DonationInventory) error {
	if di.DonorID == 0 || di.InventoryID == 0 {
		return errors.New("donor ID and inventory ID cannot be empty")
	}
	if di.Quantity <= 0 {
		di.Quantity = 1
	}
	return s.repo.Create(di)
}

func (s *DonationInventoryService) GetByID(id uint) (*models.DonationInventory, error) {
	return s.repo.GetByID(id)
}

func (s *DonationInventoryService) GetAll() ([]models.DonationInventory, error) {
	return s.repo.GetAll()
}

func (s *DonationInventoryService) GetByDonorID(donorID uint) ([]models.DonationInventory, error) {
	return s.repo.GetByDonorID(donorID)
}

func (s *DonationInventoryService) GetByInventoryID(inventoryID uint) ([]models.DonationInventory, error) {
	return s.repo.GetByInventoryID(inventoryID)
}

func (s *DonationInventoryService) Update(di *models.DonationInventory) error {
	return s.repo.Update(di)
}

func (s *DonationInventoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// ScheduleService Schedule service
type ScheduleService struct {
	repo *repo.ScheduleRepository
}

func NewScheduleService(repo *repo.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) Create(schedule *models.Schedule) error {
	if schedule.PersonID == 0 || schedule.PersonType == "" {
		return errors.New("person ID and person type cannot be empty")
	}
	if schedule.PersonType != "volunteer" && schedule.PersonType != "employee" {
		return errors.New("person type must be volunteer or employee")
	}
	
	// Generate schedule ID
	schedule.ScheduleID = fmt.Sprintf("SCH-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	schedule.Status = "scheduled"
	
	// Calculate work hours
	if schedule.HoursWorked == 0 && schedule.StartTime != "" && schedule.EndTime != "" {
		// Time calculation logic can be added here
		schedule.HoursWorked = 8.0 // Default value
	}
	
	return s.repo.Create(schedule)
}

func (s *ScheduleService) GetByID(id uint) (*models.Schedule, error) {
	return s.repo.GetByID(id)
}

func (s *ScheduleService) GetAll() ([]models.Schedule, error) {
	return s.repo.GetAll()
}

func (s *ScheduleService) GetByPersonID(personID uint, personType string) ([]models.Schedule, error) {
	return s.repo.GetByPersonID(personID, personType)
}

func (s *ScheduleService) GetByProjectID(projectID uint) ([]models.Schedule, error) {
	return s.repo.GetByProjectID(projectID)
}

func (s *ScheduleService) GetByDateRange(startDate, endDate time.Time) ([]models.Schedule, error) {
	return s.repo.GetByDateRange(startDate, endDate)
}

func (s *ScheduleService) Update(schedule *models.Schedule) error {
	return s.repo.Update(schedule)
}

func (s *ScheduleService) Delete(id uint) error {
	return s.repo.Delete(id)
}
