package repo

import (
	"erp-backend/internal/models"
	"strings"
)

// EmpGetProjects
func (r *EmployeeProjectRepository) EmpGetProjects(employeeID uint) ([]models.EmployeeProject, error) {
	var empProjects []models.EmployeeProject
	tx := r.db.Preload("Project").Where("employee_id = ?", employeeID).Find(&empProjects)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// Define the tag to filter out internal projects
	tag := "%\"internal\"%"
	var notInternalProjects []models.EmployeeProject
	//filter out internal projects
	for i := 0; i < len(empProjects); i++ {
		if empProjects[i].Project != nil && !strings.Contains(empProjects[i].Project.ProjectType, tag) {
			notInternalProjects = append(notInternalProjects, empProjects[i])
		}
	}

	return notInternalProjects, nil
}
