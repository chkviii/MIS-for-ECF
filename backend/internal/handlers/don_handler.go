package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"erp-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type DonHandler struct {
	donService *services.DonService
}

func NewDonHandler(donService *services.DonService) *DonHandler {
	return &DonHandler{
		donService: donService,
	}
}

func (h *DonHandler) GetProjectsByDonor(c *gin.Context) {
	donorID := c.GetUint("role_id")

	response, err := h.donService.GetProjectsByDonor(donorID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch donor projects: " + err.Error(),
		})
		return
	}

	jsonResponse, err := json.Marshal(response.Projects)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to marshal response: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"projects": json.RawMessage(jsonResponse),
	})
}

func (h *DonHandler) GetDonationDetails(c *gin.Context) {
	donorID := c.GetUint("role_id")
	projectIDStr := c.Query("project")
	startStr := c.Query("start")
	endStr := c.Query("end")

	request := services.DonDetailRequest{
		DonorID:   donorID,
		ProjectID: 0,
	}

	if projectIDStr != "" {
		projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid project_id parameter: " + err.Error(),
			})
			return
		}
		request.ProjectID = uint(projectID)
	}

	var layout = "2006-01-02"

	if startStr != "" {
		request.Start, _ = time.Parse(layout, startStr)
	}

	if endStr != "" {
		request.End, _ = time.Parse(layout, endStr)
	}

	log.Printf("Donation Details Request: %+v", request)

	response, err := h.donService.GetDonationDetails(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch donation details: " + err.Error(),
		})
		return
	}

	jsonResponse, err := json.Marshal(response.Details)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to marshal response: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"details": json.RawMessage(jsonResponse),
	})
}
