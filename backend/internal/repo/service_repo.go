package repo

import (
	"erp-backend/internal/models"
	"strings"
	"time"
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

// DonorGetProjects
func (r *DonorRepository) DonorGetProjects(donorID uint) ([]map[string]interface{}, error) {
	var rows []struct {
		ProjectID uint   `gorm:"column:id"`
		Name      string `gorm:"column:name"`
	}
	//Join with Project to get project details
	tx := r.db.Model(&models.Donation{}).
		Select("DISTINCT donations.project_id as id, projects.name").
		Joins("left join projects on donations.project_id = projects.id").
		Where("donor_id = ?", donorID).
		Find(&rows)
	if tx.Error != nil {
		return nil, tx.Error
	}

	out := make([]map[string]interface{}, 0, len(rows))

	for _, rrow := range rows {
		out = append(out, map[string]interface{}{
			"id":   rrow.ProjectID,
			"name": rrow.Name,
		})
	}

	return out, nil

}

func (r *DonorRepository) GetDonationsByDonorAndFilters(donorID uint, projectID uint, start, end time.Time) ([]map[string]interface{}, error) {
	var rows []struct {
		DonationID    string    `gorm:"column:donation_id"`
		DonationDate  time.Time `gorm:"column:donation_date"`
		Amount        float64   `gorm:"column:amount"`
		PaymentMethod string    `gorm:"column:payment_method"`
		ProjectName   string    `gorm:"column:project_name"`
	}

	tx := r.db.Model(&models.Donation{}).
		Select("donations.donation_id, donations.donation_date, donations.amount, donations.payment_method, projects.name as project_name").
		Joins("LEFT JOIN projects ON donations.project_id = projects.id").
		Where("donations.donor_id = ?", donorID)

	if projectID != 0 {
		tx = tx.Where("donations.project_id = ?", projectID)
	}
	if !start.IsZero() {
		tx = tx.Where("donations.donation_date >= ?", start)
	}
	if !end.IsZero() {
		tx = tx.Where("donations.donation_date <= ?", end)
	}

	if err := tx.Scan(&rows).Error; err != nil {
		return nil, err
	}

	donations := make([]map[string]interface{}, 0, len(rows))
	for _, rrow := range rows {
		donations = append(donations, map[string]interface{}{
			"donation_id":    rrow.DonationID,
			"donation_date":  rrow.DonationDate,
			"amount":         rrow.Amount,
			"payment_method": rrow.PaymentMethod,
			"project_name":   rrow.ProjectName,
		})
	}

	return donations, nil
}
