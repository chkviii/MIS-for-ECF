package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"erp-backend/internal/models"
	"erp-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// ERPHandler 持有各服务的引用
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
	deliveryInventoryService    *services.DeliveryInventoryService
	scheduleService             *services.ScheduleService
}

// NewERPHandler 构造器
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
	deliveryInventoryService *services.DeliveryInventoryService,
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
		deliveryInventoryService:    deliveryInventoryService,
		scheduleService:             scheduleService,
	}
}

// generateID 生成文本 ID，格式：ABC + YYMMDD + HHMMSS (UTC)
func generateID(prefix string) string {
	p := strings.ToUpper(prefix)
	if len(p) > 3 {
		p = p[:3]
	} else if len(p) < 3 {
		p = fmt.Sprintf("%3s", p)
	}
	t := time.Now().UTC()
	date := t.Format("060102") // YYMMDD
	tm := t.Format("150405")   // HHMMSS
	return fmt.Sprintf("%s%s%s", p, date, tm)
}

// notImplementedFilter 返回 501
func notImplementedFilter(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Filter not implemented"})
}

// parseFilterParams reads query string params `query`, `number_range`, `date_range`
// and unmarshals them into appropriate Go maps.
func parseFilterParams(c *gin.Context) (map[string]interface{}, map[string][]interface{}, map[string][]string, error) {
	var q map[string]interface{}
	var nr map[string][]interface{}
	var dr map[string][]string

	if s := c.Query("query"); s != "" {
		if err := json.Unmarshal([]byte(s), &q); err != nil {
			for k, v := range q {
				q[k] = "%" + fmt.Sprint(v) + "%"
			}
			return nil, nil, nil, err
		}
	} else {
		q = map[string]interface{}{}
	}

	if s := c.Query("number_range"); s != "" {
		if err := json.Unmarshal([]byte(s), &nr); err != nil {
			return nil, nil, nil, err
		}
	} else {
		nr = map[string][]interface{}{}
	}

	if s := c.Query("date_range"); s != "" {
		if err := json.Unmarshal([]byte(s), &dr); err != nil {
			return nil, nil, nil, err
		}
	} else {
		dr = map[string][]string{}
	}

	// Validate and normalize date_range entries to YYYY-MM-DD
	const layout = "2006-01-02"
	for key, pair := range dr {
		if len(pair) == 0 {
			continue
		}
		// normalize each element if present, validate format
		for i := 0; i < len(pair) && i < 2; i++ {
			if pair[i] == "" {
				continue
			}
			t, err := time.Parse(layout, pair[i])
			if err != nil {
				return nil, nil, nil, fmt.Errorf("invalid date for %s: %s", key, pair[i])
			}
			pair[i] = t.Format(layout)
		}
		dr[key] = pair
	}

	return q, nr, dr, nil
}

// parseDateRange 简单解析 start_date 和 end_date（格式：YYYY-MM-DD），失败返回零值
func parseDateRange(c *gin.Context) (time.Time, time.Time) {
	var start, end time.Time
	const layout = "2006-01-02"
	if s := c.Query("start_date"); s != "" {
		if t, err := time.Parse(layout, s); err == nil {
			start = t
		}
	}
	if e := c.Query("end_date"); e != "" {
		if t, err := time.Parse(layout, e); err == nil {
			end = t
		}
	}
	return start, end
}

// --- For each entity: Create, GetAll, Update, Filter (501), Delete ---

// Project
func (h *ERPHandler) CreateProject(c *gin.Context) {
	var m models.Project
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.ProjectID == "" {
		m.ProjectID = generateID("PRO")
	}
	if err := h.projectService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllProjects(c *gin.Context) {
	list, err := h.projectService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Project
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.projectService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterProjects(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.projectService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.projectService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Donor
func (h *ERPHandler) CreateDonor(c *gin.Context) {
	var m models.Donor
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.DonorID == "" {
		m.DonorID = generateID("DNR")
	}
	if err := h.donorService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllDonors(c *gin.Context) {
	list, err := h.donorService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateDonor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Donor
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.donorService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterDonors(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.donorService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteDonor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.donorService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Donation
func (h *ERPHandler) CreateDonation(c *gin.Context) {
	var m models.Donation
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.DonationID == "" {
		m.DonationID = generateID("DON")
	}
	if err := h.donationService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllDonations(c *gin.Context) {
	list, err := h.donationService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateDonation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Donation
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.donationService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterDonations(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.donationService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteDonation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.donationService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Volunteer
func (h *ERPHandler) CreateVolunteer(c *gin.Context) {
	var m models.Volunteer
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.VolunteerID == "" {
		m.VolunteerID = generateID("VOL")
	}
	if err := h.volunteerService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllVolunteers(c *gin.Context) {
	list, err := h.volunteerService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateVolunteer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Volunteer
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.volunteerService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterVolunteers(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.volunteerService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteVolunteer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.volunteerService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Employee
func (h *ERPHandler) CreateEmployee(c *gin.Context) {
	var m models.Employee
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.EmployeeID == "" {
		m.EmployeeID = generateID("EMP")
	}
	if err := h.employeeService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllEmployees(c *gin.Context) {
	list, err := h.employeeService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Employee
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.employeeService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterEmployees(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.employeeService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.employeeService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Location
func (h *ERPHandler) CreateLocation(c *gin.Context) {
	var m models.Location
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.LocationID == "" {
		m.LocationID = generateID("LOC")
	}
	if err := h.locationService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllLocations(c *gin.Context) {
	list, err := h.locationService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Location
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.locationService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterLocations(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.locationService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.locationService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Fund
func (h *ERPHandler) CreateFund(c *gin.Context) {
	var m models.Fund
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.FundID == "" {
		m.FundID = generateID("FND")
	}
	if err := h.fundService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllFunds(c *gin.Context) {
	list, err := h.fundService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateFund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Fund
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.fundService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterFunds(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.fundService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteFund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.fundService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Expense
func (h *ERPHandler) CreateExpense(c *gin.Context) {
	var m models.Expense
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.ExpenseID == "" {
		m.ExpenseID = generateID("EXP")
	}
	if err := h.expenseService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllExpenses(c *gin.Context) {
	list, err := h.expenseService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Expense
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.expenseService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterExpenses(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.expenseService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.expenseService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Transaction
func (h *ERPHandler) CreateTransaction(c *gin.Context) {
	var m models.Transaction
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.TransactionID == "" {
		m.TransactionID = generateID("TRX")
	}
	if err := h.transactionService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllTransactions(c *gin.Context) {
	list, err := h.transactionService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Transaction
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.transactionService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterTransactions(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.transactionService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.transactionService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Purchase
func (h *ERPHandler) CreatePurchase(c *gin.Context) {
	var m models.Purchase
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.PurchaseID == "" {
		m.PurchaseID = generateID("PUR")
	}
	if err := h.purchaseService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllPurchases(c *gin.Context) {
	list, err := h.purchaseService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdatePurchase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Purchase
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.purchaseService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterPurchases(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.purchaseService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeletePurchase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.purchaseService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Payroll (no textual ID)
func (h *ERPHandler) CreatePayroll(c *gin.Context) {
	var m models.Payroll
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.payrollService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllPayrolls(c *gin.Context) {
	list, err := h.payrollService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdatePayroll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Payroll
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.payrollService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterPayrolls(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.payrollService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeletePayroll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.payrollService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Inventory
func (h *ERPHandler) CreateInventory(c *gin.Context) {
	var m models.Inventory
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.InventoryID == "" {
		m.InventoryID = generateID("INV")
	}
	if err := h.inventoryService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllInventories(c *gin.Context) {
	list, err := h.inventoryService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Inventory
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.inventoryService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterInventories(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.inventoryService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.inventoryService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GiftType
func (h *ERPHandler) CreateGiftType(c *gin.Context) {
	var m models.GiftType
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.giftTypeService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllGiftTypes(c *gin.Context) {
	list, err := h.giftTypeService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateGiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.GiftType
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.giftTypeService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterGiftTypes(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.giftTypeService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteGiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.giftTypeService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Gift
func (h *ERPHandler) CreateGift(c *gin.Context) {
	var m models.Gift
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.GiftID == "" {
		m.GiftID = generateID("GFT")
	}
	if err := h.giftService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllGifts(c *gin.Context) {
	list, err := h.giftService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateGift(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Gift
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.giftService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterGifts(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.giftService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteGift(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.giftService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// InventoryTransaction
func (h *ERPHandler) CreateInventoryTransaction(c *gin.Context) {
	var m models.InventoryTransaction
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.inventoryTransactionService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllInventoryTransactions(c *gin.Context) {
	list, err := h.inventoryTransactionService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateInventoryTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.InventoryTransaction
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.inventoryTransactionService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterInventoryTransactions(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.inventoryTransactionService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteInventoryTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.inventoryTransactionService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Delivery
func (h *ERPHandler) CreateDelivery(c *gin.Context) {
	var m models.Delivery
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.DeliveryID == "" {
		m.DeliveryID = generateID("DLY")
	}
	if err := h.deliveryService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllDeliveries(c *gin.Context) {
	list, err := h.deliveryService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateDelivery(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Delivery
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.deliveryService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterDeliveries(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.deliveryService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteDelivery(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.deliveryService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// VolunteerProject, EmployeeProject, FundProject, DonationInventory, Schedule
// For conciseness these follow the same pattern and call service methods named Create/GetAll/Update/Delete

func (h *ERPHandler) CreateVolunteerProject(c *gin.Context) {
	var m models.VolunteerProject
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.volunteerProjectService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllVolunteerProjects(c *gin.Context) {
	list, err := h.volunteerProjectService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateVolunteerProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.VolunteerProject
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.volunteerProjectService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterVolunteerProjects(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.volunteerProjectService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
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
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *ERPHandler) CreateEmployeeProject(c *gin.Context) {
	var m models.EmployeeProject
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.employeeProjectService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllEmployeeProjects(c *gin.Context) {
	list, err := h.employeeProjectService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateEmployeeProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.EmployeeProject
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.employeeProjectService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterEmployeeProjects(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.employeeProjectService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
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
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *ERPHandler) CreateFundProject(c *gin.Context) {
	var m models.FundProject
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.fundProjectService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllFundProjects(c *gin.Context) {
	list, err := h.fundProjectService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateFundProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.FundProject
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.fundProjectService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterFundProjects(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.fundProjectService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
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
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *ERPHandler) CreateDonationInventory(c *gin.Context) {
	var m models.DonationInventory
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.donationInventoryService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllDonationInventories(c *gin.Context) {
	list, err := h.donationInventoryService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateDonationInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.DonationInventory
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.donationInventoryService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterDonationInventories(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.donationInventoryService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
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
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *ERPHandler) CreateDeliveryInventory(c *gin.Context) {
	var m models.DeliveryInventory
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.deliveryInventoryService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllDeliveryInventories(c *gin.Context) {
	list, err := h.deliveryInventoryService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateDeliveryInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.DeliveryInventory
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.deliveryInventoryService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterDeliveryInventories(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.deliveryInventoryService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteDeliveryInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.deliveryInventoryService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *ERPHandler) CreateSchedule(c *gin.Context) {
	var m models.Schedule
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.ScheduleID == "" {
		m.ScheduleID = generateID("SCH")
	}
	if err := h.scheduleService.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *ERPHandler) GetAllSchedules(c *gin.Context) {
	list, err := h.scheduleService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) UpdateSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var m models.Schedule
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m.ID = uint(id)
	if err := h.scheduleService.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *ERPHandler) FilterSchedules(c *gin.Context) {
	query, numberRange, dateRange, err := parseFilterParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter parameters: " + err.Error()})
		return
	}
	list, err := h.scheduleService.Filter(query, numberRange, dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "count": len(list)})
}

func (h *ERPHandler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.scheduleService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
