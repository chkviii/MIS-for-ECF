package handlers

import (
	"net/http"
	"strconv"
	"time"

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

// GET /api/v1/charts/donations-by-donor?donor_id=1&start=2025-01-01&end=2025-12-31
func (h *ChartHandler) DonationsByDonor(c *gin.Context) {
	donorIDStr := c.Query("donor_id")
	if donorIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "donor_id is required"})
		return
	}
	did, err := strconv.ParseUint(donorIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid donor_id"})
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

	pts, err := h.chartService.DonationsByDonor(uint(did), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pts})
}

// GET /api/v1/charts/fund-allocations?fund_id=1&start=...&end=...
func (h *ChartHandler) FundAllocations(c *gin.Context) {
	fundIDStr := c.Query("fund_id")
	if fundIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fund_id is required"})
		return
	}
	fid, err := strconv.ParseUint(fundIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fund_id"})
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

	pts, err := h.chartService.FundAllocationsByFund(uint(fid), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pts})
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

// GET /api/v1/charts/finance-report?start=...&end=...&group_by=date|project|location&project_id=&location_id=
func (h *ChartHandler) FinanceReport(c *gin.Context) {
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
	groupBy := c.DefaultQuery("group_by", "date")
	projID := uint(0)
	if v := c.Query("project_id"); v != "" {
		pv, err := strconv.ParseUint(v, 10, 32)
		if err == nil {
			projID = uint(pv)
		}
	}
	locID := uint(0)
	if v := c.Query("location_id"); v != "" {
		lv, err := strconv.ParseUint(v, 10, 32)
		if err == nil {
			locID = uint(lv)
		}
	}

	rows, err := h.chartService.FinanceReportAggregates(start, end, groupBy, projID, locID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows})
}
