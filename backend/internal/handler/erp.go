package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"mypage-backend/internal/models"
	"mypage-backend/internal/services"
)

type ERPHandler struct {
	projectService              *services.ProjectService
	donorService                *services.DonorService
	donationService             *services.DonationService
	volunteerService            *services.VolunteerService
	employeeService             *services.EmployeeService
	locationService             *services.LocationService
	fundService                 *services.FundService
	expenseService              *services.ExpenseService
	transactionService          *services.TransactionService
	purchaseService             *services.PurchaseService
	payrollService              *services.PayrollService
	inventoryService            *services.InventoryService
	giftTypeService             *services.GiftTypeService
	giftService                 *services.GiftService
	inventoryTransactionService *services.InventoryTransactionService
	deliveryService             *services.DeliveryService
	volunteerProjectService     *services.VolunteerProjectService
	employeeProjectService      *services.EmployeeProjectService
	fundProjectService          *services.FundProjectService
	donationInventoryService    *services.DonationInventoryService
	scheduleService             *services.ScheduleService
}

func NewERPHandler(
	projectService *services.ProjectService,
	donorService *services.DonorService,
	donationService *services.DonationService,
	volunteerService *services.VolunteerService,
	employeeService *services.EmployeeService,
	locationService *services.LocationService,
	fundService *services.FundService,
	expenseService *services.ExpenseService,
	transactionService *services.TransactionService,
	purchaseService *services.PurchaseService,
	payrollService *services.PayrollService,
	inventoryService *services.InventoryService,
	giftTypeService *services.GiftTypeService,
	giftService *services.GiftService,
	inventoryTransactionService *services.InventoryTransactionService,
	deliveryService *services.DeliveryService,
	volunteerProjectService *services.VolunteerProjectService,
	employeeProjectService *services.EmployeeProjectService,
	fundProjectService *services.FundProjectService,
	donationInventoryService *services.DonationInventoryService,
	scheduleService *services.ScheduleService,
) *ERPHandler {
	return &ERPHandler{
		projectService:              projectService,
		donorService:                donorService,
		donationService:             donationService,
		volunteerService:            volunteerService,
		employeeService:             employeeService,
		locationService:             locationService,
		fundService:                 fundService,
		expenseService:              expenseService,
		transactionService:          transactionService,
		purchaseService:             purchaseService,
		payrollService:              payrollService,
		inventoryService:            inventoryService,
		giftTypeService:             giftTypeService,
		giftService:                 giftService,
		inventoryTransactionService: inventoryTransactionService,
		deliveryService:             deliveryService,
		volunteerProjectService:     volunteerProjectService,
		employeeProjectService:      employeeProjectService,
		fundProjectService:          fundProjectService,
		donationInventoryService:    donationInventoryService,
		scheduleService:             scheduleService,
	}
}

// Helper function: Parse date range
func parseDateRange(c *gin.Context) (*time.Time, *time.Time) {
	var startDate, endDate *time.Time

	if startStr := c.Query("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = &t
		}
	}

	if endStr := c.Query("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = &t
		}
	}

	return startDate, endDate
}

// ==================== Project Management ====================

func (h *ERPHandler) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.projectService.CreateProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": project, "message": "Project created successfully"})
}

func (h *ERPHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.projectService.GetProject(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})
}

func (h *ERPHandler) GetAllProjects(c *gin.Context) {
	status := c.Query("status")
	var projects []models.Project
	var err error

	if status != "" {
		projects, err = h.projectService.GetProjectsByStatus(status)
	} else {
		projects, err = h.projectService.GetAllProjects()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": projects, "count": len(projects)})
}

func (h *ERPHandler) SearchProjects(c *gin.Context) {
	query := make(map[string]interface{})

	if name := c.Query("name"); name != "" {
		query["name"] = name
	}
	if status := c.Query("status"); status != "" {
		query["status"] = status
	}
	if projectID := c.Query("project_id"); projectID != "" {
		query["project_id"] = projectID
	}

	startDate, endDate := parseDateRange(c)

	projects, err := h.projectService.SearchProjects(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": projects, "count": len(projects)})
}

func (h *ERPHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ID = uint(id)
	if err := h.projectService.UpdateProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project, "message": "Project updated successfully"})
}

func (h *ERPHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := h.projectService.DeleteProject(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// ==================== Donor Management ====================

func (h *ERPHandler) CreateDonor(c *gin.Context) {
	var donor models.Donor
	if err := c.ShouldBindJSON(&donor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.donorService.CreateDonor(&donor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": donor, "message": "Donor created successfully"})
}

func (h *ERPHandler) GetDonor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
		return
	}

	donor, err := h.donorService.GetDonor(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Donor not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": donor})
}

func (h *ERPHandler) GetAllDonors(c *gin.Context) {
	donors, err := h.donorService.GetAllDonors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": donors, "count": len(donors)})
}

func (h *ERPHandler) SearchDonors(c *gin.Context) {
	query := make(map[string]interface{})

	if firstName := c.Query("first_name"); firstName != "" {
		query["first_name"] = firstName
	}
	if lastName := c.Query("last_name"); lastName != "" {
		query["last_name"] = lastName
	}
	if email := c.Query("email"); email != "" {
		query["email"] = email
	}
	if status := c.Query("status"); status != "" {
		query["status"] = status
	}
	if donorID := c.Query("donor_id"); donorID != "" {
		query["donor_id"] = donorID
	}

	startDate, endDate := parseDateRange(c)

	donors, err := h.donorService.SearchDonors(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": donors, "count": len(donors)})
}

func (h *ERPHandler) UpdateDonor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
		return
	}

	var donor models.Donor
	if err := c.ShouldBindJSON(&donor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	donor.ID = uint(id)
	if err := h.donorService.UpdateDonor(&donor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": donor, "message": "Donor updated successfully"})
}

func (h *ERPHandler) DeleteDonor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
		return
	}

	if err := h.donorService.DeleteDonor(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donor deleted successfully"})
}

// ==================== Donation Management ====================

func (h *ERPHandler) CreateDonation(c *gin.Context) {
	var donation models.Donation
	if err := c.ShouldBindJSON(&donation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.donationService.CreateDonation(&donation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": donation, "message": "Donation created successfully"})
}

func (h *ERPHandler) GetAllDonations(c *gin.Context) {
	donorIDStr := c.Query("donor_id")

	var donations []models.Donation
	var err error

	if donorIDStr != "" {
		donorID, parseErr := strconv.ParseUint(donorIDStr, 10, 32)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donor ID"})
			return
		}
		donations, err = h.donationService.GetDonationsByDonor(uint(donorID))
	} else {
		donations, err = h.donationService.GetAllDonations()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": donations, "count": len(donations)})
}

func (h *ERPHandler) SearchDonations(c *gin.Context) {
	query := make(map[string]interface{})

	if donationID := c.Query("donation_id"); donationID != "" {
		query["donation_id"] = donationID
	}
	if donationType := c.Query("donation_type"); donationType != "" {
		query["donation_type"] = donationType
	}
	if paymentMethod := c.Query("payment_method"); paymentMethod != "" {
		query["payment_method"] = paymentMethod
	}

	startDate, endDate := parseDateRange(c)

	donations, err := h.donationService.SearchDonations(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": donations, "count": len(donations)})
}

func (h *ERPHandler) UpdateDonation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	var donation models.Donation
	if err := c.ShouldBindJSON(&donation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	donation.ID = uint(id)
	if err := h.donationService.UpdateDonation(&donation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": donation, "message": "Donation updated successfully"})
}

func (h *ERPHandler) DeleteDonation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID"})
		return
	}

	if err := h.donationService.DeleteDonation(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donation deleted successfully"})
}

// ==================== Volunteer Management ====================

func (h *ERPHandler) CreateVolunteer(c *gin.Context) {
	var volunteer models.Volunteer
	if err := c.ShouldBindJSON(&volunteer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.volunteerService.CreateVolunteer(&volunteer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": volunteer, "message": "Volunteer created successfully"})
}

func (h *ERPHandler) GetAllVolunteers(c *gin.Context) {
	volunteers, err := h.volunteerService.GetAllVolunteers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": volunteers, "count": len(volunteers)})
}

func (h *ERPHandler) SearchVolunteers(c *gin.Context) {
	query := make(map[string]interface{})

	if firstName := c.Query("first_name"); firstName != "" {
		query["first_name"] = firstName
	}
	if lastName := c.Query("last_name"); lastName != "" {
		query["last_name"] = lastName
	}
	if status := c.Query("status"); status != "" {
		query["status"] = status
	}
	if volunteerID := c.Query("volunteer_id"); volunteerID != "" {
		query["volunteer_id"] = volunteerID
	}

	volunteers, err := h.volunteerService.SearchVolunteers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": volunteers, "count": len(volunteers)})
}

func (h *ERPHandler) UpdateVolunteer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	var volunteer models.Volunteer
	if err := c.ShouldBindJSON(&volunteer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volunteer.ID = uint(id)
	if err := h.volunteerService.UpdateVolunteer(&volunteer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": volunteer, "message": "Volunteer updated successfully"})
}

func (h *ERPHandler) DeleteVolunteer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid volunteer ID"})
		return
	}

	if err := h.volunteerService.DeleteVolunteer(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Volunteer deleted successfully"})
}

// ==================== Employee Management ====================

func (h *ERPHandler) CreateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.employeeService.CreateEmployee(&employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": employee, "message": "Employee created successfully"})
}

func (h *ERPHandler) GetAllEmployees(c *gin.Context) {
	employees, err := h.employeeService.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employees, "count": len(employees)})
}

func (h *ERPHandler) SearchEmployees(c *gin.Context) {
	query := make(map[string]interface{})

	if firstName := c.Query("first_name"); firstName != "" {
		query["first_name"] = firstName
	}
	if lastName := c.Query("last_name"); lastName != "" {
		query["last_name"] = lastName
	}
	if position := c.Query("position"); position != "" {
		query["position"] = position
	}
	if status := c.Query("status"); status != "" {
		query["status"] = status
	}
	if employeeID := c.Query("employee_id"); employeeID != "" {
		query["employee_id"] = employeeID
	}

	employees, err := h.employeeService.SearchEmployees(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employees, "count": len(employees)})
}

func (h *ERPHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee.ID = uint(id)
	if err := h.employeeService.UpdateEmployee(&employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employee, "message": "Employee updated successfully"})
}

func (h *ERPHandler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if err := h.employeeService.DeleteEmployee(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

// ==================== Location Management ====================

func (h *ERPHandler) CreateLocation(c *gin.Context) {
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.locationService.CreateLocation(&location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": location, "message": "Location created successfully"})
}

func (h *ERPHandler) GetAllLocations(c *gin.Context) {
	locations, err := h.locationService.GetAllLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": locations, "count": len(locations)})
}

func (h *ERPHandler) SearchLocations(c *gin.Context) {
	query := make(map[string]interface{})

	if name := c.Query("name"); name != "" {
		query["name"] = name
	}
	if address := c.Query("address"); address != "" {
		query["address"] = address
	}
	if locationType := c.Query("location_type"); locationType != "" {
		query["location_type"] = locationType
	}
	if locationID := c.Query("location_id"); locationID != "" {
		query["location_id"] = locationID
	}

	locations, err := h.locationService.SearchLocations(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": locations, "count": len(locations)})
}

func (h *ERPHandler) UpdateLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location.ID = uint(id)
	if err := h.locationService.UpdateLocation(&location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": location, "message": "Location updated successfully"})
}

func (h *ERPHandler) DeleteLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	if err := h.locationService.DeleteLocation(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location deleted successfully"})
}

// ==================== Fund Management ====================
func (h *ERPHandler) CreateFund(c *gin.Context) {
	var fund models.Fund
	if err := c.ShouldBindJSON(&fund); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.fundService.CreateFund(&fund); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": fund, "message": "Fund created successfully"})
}

func (h *ERPHandler) GetAllFunds(c *gin.Context) {
	funds, err := h.fundService.GetAllFunds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": funds, "count": len(funds)})
}

func (h *ERPHandler) SearchFunds(c *gin.Context) {
	query := make(map[string]interface{})

	if name := c.Query("name"); name != "" {
		query["name"] = name
	}
	if fundID := c.Query("fund_id"); fundID != "" {
		query["fund_id"] = fundID
	}
	if status := c.Query("status"); status != "" {
		query["status"] = status
	}

	startDate, endDate := parseDateRange(c)

	funds, err := h.fundService.SearchFunds(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": funds, "count": len(funds)})
}

func (h *ERPHandler) GetFund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fund ID"})
		return
	}

	fund, err := h.fundService.GetFund(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fund not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fund})
}

func (h *ERPHandler) UpdateFund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fund ID"})
		return
	}

	var fund models.Fund
	if err := c.ShouldBindJSON(&fund); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fund.ID = uint(id)
	if err := h.fundService.UpdateFund(&fund); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fund, "message": "Fund updated successfully"})
}

func (h *ERPHandler) DeleteFund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fund ID"})
		return
	}

	if err := h.fundService.DeleteFund(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Fund deleted successfully"})
}

// ==================== Expense Management ====================
func (h *ERPHandler) CreateExpense(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.expenseService.CreateExpense(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": expense, "message": "Expense created successfully"})
}

func (h *ERPHandler) GetAllExpenses(c *gin.Context) {
	expenses, err := h.expenseService.GetAllExpenses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expenses, "count": len(expenses)})
}

func (h *ERPHandler) SearchExpenses(c *gin.Context) {
	query := make(map[string]interface{})

	if expenseID := c.Query("expense_id"); expenseID != "" {
		query["expense_id"] = expenseID
	}
	if category := c.Query("category"); category != "" {
		query["category"] = category
	}
	if approvalStatus := c.Query("approval_status"); approvalStatus != "" {
		query["approval_status"] = approvalStatus
	}

	startDate, endDate := parseDateRange(c)

	expenses, err := h.expenseService.SearchExpenses(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expenses, "count": len(expenses)})
}

func (h *ERPHandler) UpdateExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.ID = uint(id)
	if err := h.expenseService.UpdateExpense(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expense, "message": "Expense updated successfully"})
}

func (h *ERPHandler) DeleteExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	if err := h.expenseService.DeleteExpense(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

// ==================== Transaction Management ====================
func (h *ERPHandler) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.transactionService.CreateTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": transaction, "message": "Transaction created successfully"})
}

func (h *ERPHandler) GetAllTransactions(c *gin.Context) {
	transactions, err := h.transactionService.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions, "count": len(transactions)})
}

func (h *ERPHandler) SearchTransactions(c *gin.Context) {
	query := make(map[string]interface{})

	if transactionID := c.Query("transaction_id"); transactionID != "" {
		query["transaction_id"] = transactionID
	}
	if transactionType := c.Query("transaction_type"); transactionType != "" {
		query["transaction_type"] = transactionType
	}

	startDate, endDate := parseDateRange(c)

	transactions, err := h.transactionService.SearchTransactions(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions, "count": len(transactions)})
}

func (h *ERPHandler) UpdateTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction.ID = uint(id)
	if err := h.transactionService.UpdateTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction, "message": "Transaction updated successfully"})
}

func (h *ERPHandler) DeleteTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	if err := h.transactionService.DeleteTransaction(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

// ==================== Purchase Management ====================
func (h *ERPHandler) CreatePurchase(c *gin.Context) {
	var purchase models.Purchase
	if err := c.ShouldBindJSON(&purchase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.purchaseService.CreatePurchase(&purchase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": purchase, "message": "Purchase created successfully"})
}

func (h *ERPHandler) GetAllPurchases(c *gin.Context) {
	purchases, err := h.purchaseService.GetAllPurchases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": purchases, "count": len(purchases)})
}

func (h *ERPHandler) SearchPurchases(c *gin.Context) {
	query := make(map[string]interface{})

	if purchaseID := c.Query("purchase_id"); purchaseID != "" {
		query["purchase_id"] = purchaseID
	}
	if supplier := c.Query("supplier"); supplier != "" {
		query["supplier"] = supplier
	}

	startDate, endDate := parseDateRange(c)

	purchases, err := h.purchaseService.SearchPurchases(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": purchases, "count": len(purchases)})
}

func (h *ERPHandler) UpdatePurchase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase ID"})
		return
	}

	var purchase models.Purchase
	if err := c.ShouldBindJSON(&purchase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchase.ID = uint(id)
	if err := h.purchaseService.UpdatePurchase(&purchase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": purchase, "message": "Purchase updated successfully"})
}

func (h *ERPHandler) DeletePurchase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase ID"})
		return
	}

	if err := h.purchaseService.DeletePurchase(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase deleted successfully"})
}

// ==================== Payroll Management ====================
func (h *ERPHandler) CreatePayroll(c *gin.Context) {
	var payroll models.Payroll
	if err := c.ShouldBindJSON(&payroll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.payrollService.CreatePayroll(&payroll); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": payroll, "message": "Payroll created successfully"})
}

func (h *ERPHandler) GetAllPayrolls(c *gin.Context) {
	payrolls, err := h.payrollService.GetAllPayrolls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payrolls, "count": len(payrolls)})
}

func (h *ERPHandler) SearchPayrolls(c *gin.Context) {
	query := make(map[string]interface{})

	if paymentMethod := c.Query("payment_method"); paymentMethod != "" {
		query["payment_method"] = paymentMethod
	}

	startDate, endDate := parseDateRange(c)

	payrolls, err := h.payrollService.SearchPayrolls(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payrolls, "count": len(payrolls)})
}

func (h *ERPHandler) UpdatePayroll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	var payroll models.Payroll
	if err := c.ShouldBindJSON(&payroll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payroll.ID = uint(id)
	if err := h.payrollService.UpdatePayroll(&payroll); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payroll, "message": "Payroll updated successfully"})
}

func (h *ERPHandler) DeletePayroll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return
	}

	if err := h.payrollService.DeletePayroll(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payroll deleted successfully"})
}

// ==================== Inventory Management ====================
func (h *ERPHandler) CreateInventory(c *gin.Context) {
	var inventory models.Inventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.inventoryService.CreateInventory(&inventory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": inventory, "message": "Inventory created successfully"})
}

func (h *ERPHandler) GetAllInventories(c *gin.Context) {
	inventories, err := h.inventoryService.GetAllInventories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": inventories, "count": len(inventories)})
}

func (h *ERPHandler) SearchInventories(c *gin.Context) {
	query := make(map[string]interface{})

	if name := c.Query("name"); name != "" {
		query["name"] = name
	}
	if inventoryID := c.Query("inventory_id"); inventoryID != "" {
		query["inventory_id"] = inventoryID
	}
	if status := c.Query("status"); status != "" {
		query["status"] = status
	}

	inventories, err := h.inventoryService.SearchInventories(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": inventories, "count": len(inventories)})
}

func (h *ERPHandler) UpdateInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
		return
	}

	var inventory models.Inventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inventory.ID = uint(id)
	if err := h.inventoryService.UpdateInventory(&inventory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": inventory, "message": "Inventory updated successfully"})
}

func (h *ERPHandler) DeleteInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
		return
	}

	if err := h.inventoryService.DeleteInventory(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory deleted successfully"})
}

// ==================== Gift Type Management ====================
func (h *ERPHandler) CreateGiftType(c *gin.Context) {
	var giftType models.GiftType
	if err := c.ShouldBindJSON(&giftType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.giftTypeService.CreateGiftType(&giftType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": giftType, "message": "Gift type created successfully"})
}

func (h *ERPHandler) GetAllGiftTypes(c *gin.Context) {
	giftTypes, err := h.giftTypeService.GetAllGiftTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": giftTypes, "count": len(giftTypes)})
}

func (h *ERPHandler) UpdateGiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gift type ID"})
		return
	}

	var giftType models.GiftType
	if err := c.ShouldBindJSON(&giftType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	giftType.ID = uint(id)
	if err := h.giftTypeService.UpdateGiftType(&giftType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": giftType, "message": "Gift type updated successfully"})
}

func (h *ERPHandler) DeleteGiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gift type ID"})
		return
	}

	if err := h.giftTypeService.DeleteGiftType(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gift type deleted successfully"})
}

// ==================== Gift Management ====================
func (h *ERPHandler) CreateGift(c *gin.Context) {
	var gift models.Gift
	if err := c.ShouldBindJSON(&gift); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.giftService.CreateGift(&gift); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": gift, "message": "Gift created successfully"})
}

func (h *ERPHandler) GetAllGifts(c *gin.Context) {
	gifts, err := h.giftService.GetAllGifts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gifts, "count": len(gifts)})
}

func (h *ERPHandler) SearchGifts(c *gin.Context) {
	query := make(map[string]interface{})

	if giftID := c.Query("gift_id"); giftID != "" {
		query["gift_id"] = giftID
	}
	if distributionStatus := c.Query("distribution_status"); distributionStatus != "" {
		query["distribution_status"] = distributionStatus
	}

	startDate, endDate := parseDateRange(c)

	gifts, err := h.giftService.SearchGifts(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gifts, "count": len(gifts)})
}

func (h *ERPHandler) UpdateGift(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gift ID"})
		return
	}

	var gift models.Gift
	if err := c.ShouldBindJSON(&gift); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gift.ID = uint(id)
	if err := h.giftService.UpdateGift(&gift); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gift, "message": "Gift updated successfully"})
}

func (h *ERPHandler) DeleteGift(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gift ID"})
		return
	}

	if err := h.giftService.DeleteGift(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gift deleted successfully"})
}

// ==================== Inventory Transaction Management ====================
func (h *ERPHandler) CreateInventoryTransaction(c *gin.Context) {
	var transaction models.InventoryTransaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.inventoryTransactionService.CreateInventoryTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": transaction, "message": "Inventory transaction created successfully"})
}

func (h *ERPHandler) GetAllInventoryTransactions(c *gin.Context) {
	transactions, err := h.inventoryTransactionService.GetAllInventoryTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions, "count": len(transactions)})
}

func (h *ERPHandler) SearchInventoryTransactions(c *gin.Context) {
	query := make(map[string]interface{})

	if transactionID := c.Query("transaction_id"); transactionID != "" {
		query["transaction_id"] = transactionID
	}
	if transactionType := c.Query("transaction_type"); transactionType != "" {
		query["transaction_type"] = transactionType
	}

	startDate, endDate := parseDateRange(c)

	transactions, err := h.inventoryTransactionService.SearchInventoryTransactions(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions, "count": len(transactions)})
}

func (h *ERPHandler) UpdateInventoryTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory transaction ID"})
		return
	}

	var transaction models.InventoryTransaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction.ID = uint(id)
	if err := h.inventoryTransactionService.UpdateInventoryTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction, "message": "Inventory transaction updated successfully"})
}

func (h *ERPHandler) DeleteInventoryTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory transaction ID"})
		return
	}

	if err := h.inventoryTransactionService.DeleteInventoryTransaction(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory transaction deleted successfully"})
}

// ==================== Delivery Management ====================
func (h *ERPHandler) CreateDelivery(c *gin.Context) {
	var delivery models.Delivery
	if err := c.ShouldBindJSON(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.deliveryService.CreateDelivery(&delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": delivery, "message": "Delivery created successfully"})
}

func (h *ERPHandler) GetAllDeliveries(c *gin.Context) {
	deliveries, err := h.deliveryService.GetAllDeliveries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deliveries, "count": len(deliveries)})
}

func (h *ERPHandler) SearchDeliveries(c *gin.Context) {
	query := make(map[string]interface{})

	if deliveryID := c.Query("delivery_id"); deliveryID != "" {
		query["delivery_id"] = deliveryID
	}
	if status := c.Query("status"); status != "" {
		query["status"] = status
	}

	startDate, endDate := parseDateRange(c)

	deliveries, err := h.deliveryService.SearchDeliveries(query, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deliveries, "count": len(deliveries)})
}

func (h *ERPHandler) UpdateDelivery(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery ID"})
		return
	}

	var delivery models.Delivery
	if err := c.ShouldBindJSON(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery.ID = uint(id)
	if err := h.deliveryService.UpdateDelivery(&delivery); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": delivery, "message": "Delivery updated successfully"})
}

func (h *ERPHandler) DeleteDelivery(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery ID"})
		return
	}

	if err := h.deliveryService.DeleteDelivery(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery deleted successfully"})
}

// ==================== Volunteer-Project Association Management ====================
func (h *ERPHandler) CreateVolunteerProject(c *gin.Context) {
	var vp models.VolunteerProject
	if err := c.ShouldBindJSON(&vp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.volunteerProjectService.Create(&vp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": vp, "message": "Volunteer project assignment created successfully"})
}

func (h *ERPHandler) GetAllVolunteerProjects(c *gin.Context) {
	vps, err := h.volunteerProjectService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": vps, "count": len(vps)})
}

func (h *ERPHandler) UpdateVolunteerProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var vp models.VolunteerProject
	if err := c.ShouldBindJSON(&vp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vp.ID = uint(id)
	if err := h.volunteerProjectService.Update(&vp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": vp, "message": "Volunteer project assignment updated successfully"})
}

func (h *ERPHandler) DeleteVolunteerProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.volunteerProjectService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Volunteer project assignment deleted successfully"})
}

// ==================== Employee-Project Association Management ====================
func (h *ERPHandler) CreateEmployeeProject(c *gin.Context) {
	var ep models.EmployeeProject
	if err := c.ShouldBindJSON(&ep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.employeeProjectService.Create(&ep); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": ep, "message": "Employee project assignment created successfully"})
}

func (h *ERPHandler) GetAllEmployeeProjects(c *gin.Context) {
	eps, err := h.employeeProjectService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": eps, "count": len(eps)})
}

func (h *ERPHandler) UpdateEmployeeProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var ep models.EmployeeProject
	if err := c.ShouldBindJSON(&ep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ep.ID = uint(id)
	if err := h.employeeProjectService.Update(&ep); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ep, "message": "Employee project assignment updated successfully"})
}

func (h *ERPHandler) DeleteEmployeeProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.employeeProjectService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee project assignment deleted successfully"})
}

// ==================== Fund-Project Association Management ====================
func (h *ERPHandler) CreateFundProject(c *gin.Context) {
	var fp models.FundProject
	if err := c.ShouldBindJSON(&fp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.fundProjectService.Create(&fp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": fp, "message": "Fund project allocation created successfully"})
}

func (h *ERPHandler) GetAllFundProjects(c *gin.Context) {
	fps, err := h.fundProjectService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fps, "count": len(fps)})
}

func (h *ERPHandler) UpdateFundProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var fp models.FundProject
	if err := c.ShouldBindJSON(&fp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fp.ID = uint(id)
	if err := h.fundProjectService.Update(&fp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fp, "message": "Fund project allocation updated successfully"})
}

func (h *ERPHandler) DeleteFundProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.fundProjectService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Fund project allocation deleted successfully"})
}

// ==================== Donation-Inventory Association Management ====================
func (h *ERPHandler) CreateDonationInventory(c *gin.Context) {
	var di models.DonationInventory
	if err := c.ShouldBindJSON(&di); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.donationInventoryService.Create(&di); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": di, "message": "In-kind donation created successfully"})
}

func (h *ERPHandler) GetAllDonationInventories(c *gin.Context) {
	dis, err := h.donationInventoryService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dis, "count": len(dis)})
}

func (h *ERPHandler) UpdateDonationInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var di models.DonationInventory
	if err := c.ShouldBindJSON(&di); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	di.ID = uint(id)
	if err := h.donationInventoryService.Update(&di); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": di, "message": "In-kind donation updated successfully"})
}

func (h *ERPHandler) DeleteDonationInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.donationInventoryService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "In-kind donation deleted successfully"})
}

// ==================== Schedule Management ====================
func (h *ERPHandler) CreateSchedule(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.scheduleService.Create(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": schedule, "message": "Schedule created successfully"})
}

func (h *ERPHandler) GetAllSchedules(c *gin.Context) {
	schedules, err := h.scheduleService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": schedules, "count": len(schedules)})
}

func (h *ERPHandler) UpdateSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule.ID = uint(id)
	if err := h.scheduleService.Update(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": schedule, "message": "Schedule updated successfully"})
}

func (h *ERPHandler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	if err := h.scheduleService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
}
