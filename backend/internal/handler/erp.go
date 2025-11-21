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

// 辅助函数：解析日期范围
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

// ==================== 项目管理 ====================

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

	c.JSON(http.StatusCreated, gin.H{"data": project, "message": "项目创建成功"})
}

func (h *ERPHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	project, err := h.projectService.GetProject(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
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

func (h *ERPHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": project, "message": "项目更新成功"})
}

func (h *ERPHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	if err := h.projectService.DeleteProject(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "项目删除成功"})
}

// ==================== 捐赠者管理 ====================

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

	c.JSON(http.StatusCreated, gin.H{"data": donor, "message": "捐赠者创建成功"})
}

func (h *ERPHandler) GetDonor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的捐赠者ID"})
		return
	}

	donor, err := h.donorService.GetDonor(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "捐赠者不存在"})
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

func (h *ERPHandler) UpdateDonor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的捐赠者ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": donor, "message": "捐赠者更新成功"})
}

func (h *ERPHandler) DeleteDonor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的捐赠者ID"})
		return
	}

	if err := h.donorService.DeleteDonor(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "捐赠者删除成功"})
}

// ==================== 捐赠管理 ====================

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

	c.JSON(http.StatusCreated, gin.H{"data": donation, "message": "捐赠记录创建成功"})
}

func (h *ERPHandler) GetAllDonations(c *gin.Context) {
	donorIDStr := c.Query("donor_id")

	var donations []models.Donation
	var err error

	if donorIDStr != "" {
		donorID, parseErr := strconv.ParseUint(donorIDStr, 10, 32)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的捐赠者ID"})
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

func (h *ERPHandler) UpdateDonation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的捐赠ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": donation, "message": "捐赠记录更新成功"})
}

func (h *ERPHandler) DeleteDonation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的捐赠ID"})
		return
	}

	if err := h.donationService.DeleteDonation(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "捐赠记录删除成功"})
}

// ==================== 志愿者管理 ====================

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

	c.JSON(http.StatusCreated, gin.H{"data": volunteer, "message": "志愿者创建成功"})
}

func (h *ERPHandler) GetAllVolunteers(c *gin.Context) {
	volunteers, err := h.volunteerService.GetAllVolunteers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": volunteers, "count": len(volunteers)})
}

func (h *ERPHandler) UpdateVolunteer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的志愿者ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": volunteer, "message": "志愿者更新成功"})
}

func (h *ERPHandler) DeleteVolunteer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的志愿者ID"})
		return
	}

	if err := h.volunteerService.DeleteVolunteer(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "志愿者删除成功"})
}

// ==================== 员工管理 ====================

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

	c.JSON(http.StatusCreated, gin.H{"data": employee, "message": "员工创建成功"})
}

func (h *ERPHandler) GetAllEmployees(c *gin.Context) {
	employees, err := h.employeeService.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employees, "count": len(employees)})
}

func (h *ERPHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的员工ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": employee, "message": "员工更新成功"})
}

func (h *ERPHandler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的员工ID"})
		return
	}

	if err := h.employeeService.DeleteEmployee(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "员工删除成功"})
}

// ==================== 地点管理 ====================

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

	c.JSON(http.StatusCreated, gin.H{"data": location, "message": "地点创建成功"})
}

func (h *ERPHandler) GetAllLocations(c *gin.Context) {
	locations, err := h.locationService.GetAllLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": locations, "count": len(locations)})
}

func (h *ERPHandler) UpdateLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的地点ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": location, "message": "地点更新成功"})
}

func (h *ERPHandler) DeleteLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的地点ID"})
		return
	}

	if err := h.locationService.DeleteLocation(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "地点删除成功"})
}

// ==================== 基金管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": fund, "message": "基金创建成功"})
}

func (h *ERPHandler) GetAllFunds(c *gin.Context) {
	funds, err := h.fundService.GetAllFunds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": funds, "count": len(funds)})
}

func (h *ERPHandler) GetFund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的基金ID"})
		return
	}

	fund, err := h.fundService.GetFund(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "基金不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": fund})
}

func (h *ERPHandler) UpdateFund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的基金ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": fund, "message": "基金更新成功"})
}

func (h *ERPHandler) DeleteFund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的基金ID"})
		return
	}

	if err := h.fundService.DeleteFund(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "基金删除成功"})
}

// ==================== 支出管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": expense, "message": "支出记录创建成功"})
}

func (h *ERPHandler) GetAllExpenses(c *gin.Context) {
	expenses, err := h.expenseService.GetAllExpenses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expenses, "count": len(expenses)})
}

func (h *ERPHandler) UpdateExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的支出ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": expense, "message": "支出记录更新成功"})
}

func (h *ERPHandler) DeleteExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的支出ID"})
		return
	}

	if err := h.expenseService.DeleteExpense(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "支出记录删除成功"})
}

// ==================== 交易管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": transaction, "message": "交易记录创建成功"})
}

func (h *ERPHandler) GetAllTransactions(c *gin.Context) {
	transactions, err := h.transactionService.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions, "count": len(transactions)})
}

func (h *ERPHandler) UpdateTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的交易ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": transaction, "message": "交易记录更新成功"})
}

func (h *ERPHandler) DeleteTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的交易ID"})
		return
	}

	if err := h.transactionService.DeleteTransaction(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "交易记录删除成功"})
}

// ==================== 采购管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": purchase, "message": "采购记录创建成功"})
}

func (h *ERPHandler) GetAllPurchases(c *gin.Context) {
	purchases, err := h.purchaseService.GetAllPurchases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": purchases, "count": len(purchases)})
}

func (h *ERPHandler) UpdatePurchase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的采购ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": purchase, "message": "采购记录更新成功"})
}

func (h *ERPHandler) DeletePurchase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的采购ID"})
		return
	}

	if err := h.purchaseService.DeletePurchase(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "采购记录删除成功"})
}

// ==================== 薪资管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": payroll, "message": "薪资记录创建成功"})
}

func (h *ERPHandler) GetAllPayrolls(c *gin.Context) {
	payrolls, err := h.payrollService.GetAllPayrolls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payrolls, "count": len(payrolls)})
}

func (h *ERPHandler) UpdatePayroll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的薪资ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": payroll, "message": "薪资记录更新成功"})
}

func (h *ERPHandler) DeletePayroll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的薪资ID"})
		return
	}

	if err := h.payrollService.DeletePayroll(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "薪资记录删除成功"})
}

// ==================== 库存管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": inventory, "message": "库存创建成功"})
}

func (h *ERPHandler) GetAllInventories(c *gin.Context) {
	inventories, err := h.inventoryService.GetAllInventories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": inventories, "count": len(inventories)})
}

func (h *ERPHandler) UpdateInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的库存ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": inventory, "message": "库存更新成功"})
}

func (h *ERPHandler) DeleteInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的库存ID"})
		return
	}

	if err := h.inventoryService.DeleteInventory(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "库存删除成功"})
}

// ==================== 礼品类型管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": giftType, "message": "礼品类型创建成功"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的礼品类型ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": giftType, "message": "礼品类型更新成功"})
}

func (h *ERPHandler) DeleteGiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的礼品类型ID"})
		return
	}

	if err := h.giftTypeService.DeleteGiftType(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "礼品类型删除成功"})
}

// ==================== 礼品管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": gift, "message": "礼品记录创建成功"})
}

func (h *ERPHandler) GetAllGifts(c *gin.Context) {
	gifts, err := h.giftService.GetAllGifts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gifts, "count": len(gifts)})
}

func (h *ERPHandler) UpdateGift(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的礼品ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": gift, "message": "礼品记录更新成功"})
}

func (h *ERPHandler) DeleteGift(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的礼品ID"})
		return
	}

	if err := h.giftService.DeleteGift(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "礼品记录删除成功"})
}

// ==================== 库存交易管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": transaction, "message": "库存交易记录创建成功"})
}

func (h *ERPHandler) GetAllInventoryTransactions(c *gin.Context) {
	transactions, err := h.inventoryTransactionService.GetAllInventoryTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions, "count": len(transactions)})
}

func (h *ERPHandler) UpdateInventoryTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的库存交易ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": transaction, "message": "库存交易记录更新成功"})
}

func (h *ERPHandler) DeleteInventoryTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的库存交易ID"})
		return
	}

	if err := h.inventoryTransactionService.DeleteInventoryTransaction(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "库存交易记录删除成功"})
}

// ==================== 配送管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": delivery, "message": "配送记录创建成功"})
}

func (h *ERPHandler) GetAllDeliveries(c *gin.Context) {
	deliveries, err := h.deliveryService.GetAllDeliveries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deliveries, "count": len(deliveries)})
}

func (h *ERPHandler) UpdateDelivery(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配送ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": delivery, "message": "配送记录更新成功"})
}

func (h *ERPHandler) DeleteDelivery(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配送ID"})
		return
	}

	if err := h.deliveryService.DeleteDelivery(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "配送记录删除成功"})
}

// ==================== 志愿者-项目关联管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": vp, "message": "志愿者项目分配创建成功"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": vp, "message": "志愿者项目分配更新成功"})
}

func (h *ERPHandler) DeleteVolunteerProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.volunteerProjectService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "志愿者项目分配删除成功"})
}

// ==================== 员工-项目关联管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": ep, "message": "员工项目分配创建成功"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": ep, "message": "员工项目分配更新成功"})
}

func (h *ERPHandler) DeleteEmployeeProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.employeeProjectService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "员工项目分配删除成功"})
}

// ==================== 基金-项目关联管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": fp, "message": "基金项目分配创建成功"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": fp, "message": "基金项目分配更新成功"})
}

func (h *ERPHandler) DeleteFundProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.fundProjectService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "基金项目分配删除成功"})
}

// ==================== 捐赠-库存关联管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": di, "message": "实物捐赠记录创建成功"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": di, "message": "实物捐赠记录更新成功"})
}

func (h *ERPHandler) DeleteDonationInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.donationInventoryService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "实物捐赠记录删除成功"})
}

// ==================== 调度管理 ====================
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

	c.JSON(http.StatusCreated, gin.H{"data": schedule, "message": "调度记录创建成功"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的调度ID"})
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

	c.JSON(http.StatusOK, gin.H{"data": schedule, "message": "调度记录更新成功"})
}

func (h *ERPHandler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的调度ID"})
		return
	}

	if err := h.scheduleService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "调度记录删除成功"})
}
