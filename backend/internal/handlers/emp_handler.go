package handlers

import (
	"net/http"

	"erp-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// EmpHandler 员工处理器
type EmpHandler struct {
	empService *services.EmpService
}

func NewEmpHandler(empService *services.EmpService) *EmpHandler {
	return &EmpHandler{
		empService: empService,
	}
}

// Get Internal Projects

func (h *EmpHandler) GetInternalProjects(c *gin.Context) {
	var req services.EmpInternalProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Request parameter error: " + err.Error(),
		})
		return
	}

	response, err := h.empService.GetInternalProjects(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch internal projects: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// Get Projects

func (h *EmpHandler) GetProjects(c *gin.Context) {
	var req services.EmpProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Request parameter error: " + err.Error(),
		})
		return
	}
	response, err := h.empService.GetProjects(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch projects: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
