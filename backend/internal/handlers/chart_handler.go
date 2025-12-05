package handlers

import (
	"net/http"
	"strconv"
	"time"

	"erp-backend/internal/repo"
	"erp-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// ChartHandler handles chart-related API endpoints
type ChartHandler struct {
	chartService *services.ChartService
}

func NewChartHandler(cs *services.ChartService) *ChartHandler {
	return &ChartHandler{chartService: cs}
}

func parseDatePtr(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func serializeLinePts(pts []repo.LinePoint) (out []map[string]interface{}, total float64, avg float64) {
	points := make([]map[string]interface{}, len(pts))

	count := float64(len(pts))

	for i, pt := range pts {
		points[i] = map[string]interface{}{"t": pt.Date, "v": pt.Value}
		total += pt.Value
		avg += total / count
	}

	series := map[string]interface{}{
		"name":   "Values",
		"points": points,
	}

	return []map[string]interface{}{series}, total, avg
}

func serializePiePts(pts []repo.PiePoint) (pie []map[string]interface{}, total float64, avg float64) {
	out := make([]map[string]interface{}, len(pts))
	count := float64(len(pts))

	for i, pt := range pts {
		out[i] = map[string]interface{}{"label": pt.ProjectName + " (ID: " + strconv.Itoa(int(pt.ProjectID)) + ")", "value": pt.Value}
		total += pt.Value
		avg += total / count
	}
	return out, total, avg
}

// GET /api/v1/don/charts/donations-by-donor?start=2025-01-01&end=2025-12-31
func (h *ChartHandler) DonorDonations(c *gin.Context) {
	donorID, _ := c.Get("role_id")

	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	pts, err := h.chartService.DonationsByDonor(donorID.(uint), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	series, total, avg := serializeLinePts(pts)

	c.JSON(http.StatusOK, gin.H{
		"title":  "Donations by Date",
		"series": series,
		"total":  total,
		"avg":    avg,
	})
}

// Get /api/v1/don/charts/donations-by-project?start=...&end=...
func (h *ChartHandler) DonorDonationsByProject(c *gin.Context) {
	donorID, _ := c.Get("role_id")

	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}

	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	pts, err := h.chartService.DonorDonationsByProject(donorID.(uint), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pie, total, avg := serializePiePts(pts)

	c.JSON(http.StatusOK, gin.H{
		"title": "Donations by Project",
		"pie":   pie,
		"total": total,
		"avg":   avg,
	})
}

// GET /api/v1/fin/charts/line/fund?start=...&end=...
func (h *ChartHandler) FundAllocations(c *gin.Context) {

	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	pts, err := h.chartService.FundAllocationsByDate(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	series, total, avg := serializeLinePts(pts)

	c.JSON(http.StatusOK, gin.H{
		"title":  "Fund Allocations by Date",
		"series": series,
		"total":  total,
		"avg":    avg,
	})
}

// GET /api/v1/fin/charts/pie/fund?start=...&end=...
func (h *ChartHandler) FundAllocationsByProject(c *gin.Context) {
	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	pts, err := h.chartService.FundAllocationsByProject(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pie, total, avg := serializePiePts(pts)

	c.JSON(http.StatusOK, gin.H{
		"title": "Fund Allocations by Project",
		"pie":   pie,
		"total": total,
		"avg":   avg,
	})
}

// GET /api/v1/fin/charts/line/expenses?start=...&end=...
func (h *ChartHandler) Expenses(c *gin.Context) {
	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}
	pts, err := h.chartService.ExpensesByDate(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	series, total, avg := serializeLinePts(pts)

	c.JSON(http.StatusOK, gin.H{
		"title":  "Expenses by Date",
		"series": series,
		"total":  total,
		"avg":    avg,
	})
}

// GET /api/v1/fin/charts/pie/expenses?start=...&end=...
func (h *ChartHandler) ExpensesByProject(c *gin.Context) {
	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}
	pts, err := h.chartService.ExpensesByProject(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pie, total, avg := serializePiePts(pts)

	c.JSON(http.StatusOK, gin.H{
		"title": "Expenses by Project",
		"pie":   pie,
		"total": total,
		"avg":   avg,
	})
}

// GET /api/v1/charts/donations?start=...&end=...
func (h *ChartHandler) Donations(c *gin.Context) {
	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}
	pts, err := h.chartService.DonationsByDate(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	series, total, avg := serializeLinePts(pts)

	c.JSON(http.StatusOK, gin.H{
		"title":  "Donations by Date",
		"series": series,
		"total":  total,
		"avg":    avg,
	})
}

// GET /api/v1/fin/charts/donations-by-project?start=...&end=...
func (h *ChartHandler) DonationsByProject(c *gin.Context) {

	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}

	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	pts, err := h.chartService.DonationsByProject(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pie, total, avg := serializePiePts(pts)

	c.JSON(http.StatusOK, gin.H{
		"title": "Donations by Project",
		"pie":   pie,
		"total": total,
		"avg":   avg,
	})
}

// GET /api/v1/charts/volunteer-hours?volunteer_id=1&start=...&end=...
func (h *ChartHandler) VolunteerHours(c *gin.Context) {
	vidStr := c.Query("volunteer_id")
	if vidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "volunteer_id is required"})
		return
	}
	vid, err := strconv.ParseUint(vidStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid volunteer_id"})
		return
	}
	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	pts, err := h.chartService.VolunteerHoursByVolunteer(uint(vid), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pts})
}

// GET /api/v1/don/charts/line/donations-by-donor?start=...&end=...
func (h *ChartHandler) DonorDonationsLineChart(c *gin.Context) {
	donorID, _ := c.Get("role_id")
	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}
	pts, err := h.chartService.DonationsByDonor(donorID.(uint), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	series, total, avg := serializeLinePts(pts)
	c.JSON(http.StatusOK, gin.H{
		"title":  "Donations by Date",
		"series": series,
		"total":  total,
		"avg":    avg,
	})
}

// Get /api/v1/don/charts/pie/donations-by-project?start=...&end=...
func (h *ChartHandler) DonorDonationsByProjectPieChart(c *gin.Context) {
	donorID, _ := c.Get("role_id")

	start, err := parseDatePtr(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}
	end, err := parseDatePtr(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}
	pts, err := h.chartService.DonorDonationsByProject(donorID.(uint), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pie, total, avg := serializePiePts(pts)
	c.JSON(http.StatusOK, gin.H{
		"title": "Donations by Project",
		"pie":   pie,
		"total": total,
		"avg":   avg,
	})
}
