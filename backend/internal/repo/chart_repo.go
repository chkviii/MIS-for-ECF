package repo

import (
	"time"

	"erp-backend/internal/models"

	"gorm.io/gorm"
)

// ChartRepository provides DB aggregation methods used by charting/reporting
type ChartRepository struct {
	db *gorm.DB
}

func NewChartRepository(db *gorm.DB) *ChartRepository {
	return &ChartRepository{db: db}
}

// LinePoint represents aggregated data points per date
type LinePoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// PiePoint represents totals per project
type PiePoint struct {
	ProjectID   uint    `json:"project_id"`
	ProjectName string  `json:"project_name"`
	Value       float64 `json:"value"`
}

// DonationsByDonor returns donation totals per day (or per period) for a donor
func (r *ChartRepository) DonationsByDonor(donorID uint, start, end *time.Time) ([]LinePoint, error) {
	var rows []struct {
		Date  string  `gorm:"column:date"`
		Value float64 `gorm:"column:sum_amount"`
	}

	tx := r.db.Model(&models.Donation{}).Select("date(donation_date) as date, sum(amount) as sum_amount").Where("donor_id = ?", donorID)
	if start != nil {
		tx = tx.Where("date(donation_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		tx = tx.Where("date(donation_date) <= ?", end.Format("2006-01-02"))
	}
	tx = tx.Group("date(donation_date)").Order("date(donation_date)")

	if err := tx.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]LinePoint, 0, len(rows))
	for _, rrow := range rows {
		out = append(out, LinePoint{Date: rrow.Date, Value: rrow.Value})
	}
	return out, nil
}

func (r *ChartRepository) DonationsByProject(donorID uint, start, end *time.Time) ([]PiePoint, error) {
	var rows []struct {
		ProjectID   uint    `gorm:"column:project_id"`
		ProjectName string  `gorm:"column:project_name"`
		Value       float64 `gorm:"column:sum_amount"`
	}

	tx := r.db.Model(&models.Donation{}).
		Select("donations.project_id as project_id, projects.name as project_name, sum(donations.amount) as sum_amount").
		Joins("LEFT JOIN projects ON projects.id = donations.project_id")

	if donorID != 0 {
		tx = tx.Where("donor_id = ?", donorID)
	}
	if start != nil {
		tx = tx.Where("date(donation_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		tx = tx.Where("date(donation_date) <= ?", end.Format("2006-01-02"))
	}

	tx = tx.Group("donations.project_id, projects.name").Order("sum_amount DESC")

	if err := tx.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]PiePoint, 0, len(rows))
	for _, rrow := range rows {
		out = append(out, PiePoint{
			ProjectID:   rrow.ProjectID,
			ProjectName: rrow.ProjectName,
			Value:       rrow.Value,
		})
	}
	return out, nil
}

// LinePoint represents allocation totals per date for a fund

// FundAllocationsByDate aggregates fund allocations (FundProject or Fund movements) per date
func (r *ChartRepository) FundAllocationsByDate(start, end *time.Time) ([]LinePoint, error) {
	var rows []struct {
		Date  string  `gorm:"column:date"`
		Value float64 `gorm:"column:sum_amount"`
	}

	tx := r.db.Model(&models.FundProject{}).Select("date(allocation_date) as date, sum(allocated_amount) as sum_amount")
	if start != nil {
		tx = tx.Where("date(allocation_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		tx = tx.Where("date(allocation_date) <= ?", end.Format("2006-01-02"))
	}
	tx = tx.Group("date(allocation_date)").Order("date(allocation_date)")

	if err := tx.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]LinePoint, 0, len(rows))

	for _, rr := range rows {
		out = append(out, LinePoint{Date: rr.Date, Value: rr.Value})
	}

	return out, nil
}

// FundAllocationsByProject aggregates fund allocations per project
func (r *ChartRepository) FundAllocationsByProject(start, end *time.Time) ([]PiePoint, error) {
	var rows []struct {
		ProjectID   uint    `gorm:"column:project_id"`
		ProjectName string  `gorm:"column:project_name"`
		Value       float64 `gorm:"column:sum_amount"`
	}
	tx := r.db.Model(&models.FundProject{}).
		Select("fund_projects.project_id as project_id, projects.name as project_name, sum(fund_projects.allocated_amount) as sum_amount").
		Joins("LEFT JOIN projects ON projects.id = fund_projects.project_id")

	if start != nil {
		tx = tx.Where("date(allocation_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		tx = tx.Where("date(allocation_date) <= ?", end.Format("2006-01-02"))
	}

	tx = tx.Group("fund_projects.project_id").Order("sum_amount DESC")

	if err := tx.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]PiePoint, 0, len(rows))

	for _, rr := range rows {
		out = append(out, PiePoint{
			ProjectID:   rr.ProjectID,
			ProjectName: rr.ProjectName,
			Value:       rr.Value,
		})
	}

	return out, nil
}

// DonationsByDate 聚合所有捐赠按日期的总额（不按 donor 过滤）
func (r *ChartRepository) DonationsByDate(start, end *time.Time) ([]LinePoint, error) {
	var rows []struct {
		Date  string  `gorm:"column:date"`
		Value float64 `gorm:"column:sum_amount"`
	}

	tx := r.db.Model(&models.Donation{}).
		Select("date(donation_date) as date, sum(amount) as sum_amount")

	if start != nil {
		tx = tx.Where("date(donation_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		tx = tx.Where("date(donation_date) <= ?", end.Format("2006-01-02"))
	}

	tx = tx.Group("date(donation_date)").Order("date(donation_date)")

	if err := tx.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]LinePoint, 0, len(rows))
	for _, rrow := range rows {
		out = append(out, LinePoint{Date: rrow.Date, Value: rrow.Value})
	}
	return out, nil
}

// ExpensesByDate 聚合所有费用按日期的总额
func (r *ChartRepository) ExpensesByDate(start, end *time.Time) ([]LinePoint, error) {
	var rows []struct {
		Date  string  `gorm:"column:date"`
		Value float64 `gorm:"column:sum_amount"`
	}

	tx := r.db.Model(&models.Expense{}).
		Select("date(expense_date) as date, sum(amount) as sum_amount")

	if start != nil {
		tx = tx.Where("date(expense_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		tx = tx.Where("date(expense_date) <= ?", end.Format("2006-01-02"))
	}

	tx = tx.Group("date(expense_date)").Order("date(expense_date)")

	if err := tx.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]LinePoint, 0, len(rows))
	for _, rrow := range rows {
		out = append(out, LinePoint{Date: rrow.Date, Value: rrow.Value})
	}
	return out, nil
}

// ExpensesByProject 按项目聚合费用总额，返回 project_id/project_name/value（三个字段均非指针）
func (r *ChartRepository) ExpensesByProject(start, end *time.Time) ([]PiePoint, error) {
	var rows []struct {
		ProjectID   uint    `gorm:"column:project_id"`
		ProjectName string  `gorm:"column:project_name"`
		Value       float64 `gorm:"column:sum_amount"`
	}

	tx := r.db.Model(&models.Expense{}).
		Select("expenses.project_id as project_id, projects.name as project_name, sum(expenses.amount) as sum_amount").
		Joins("LEFT JOIN projects ON projects.id = expenses.project_id")

	if start != nil {
		tx = tx.Where("date(expense_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		tx = tx.Where("date(expense_date) <= ?", end.Format("2006-01-02"))
	}

	tx = tx.Group("expenses.project_id, projects.name").Order("sum_amount DESC")

	if err := tx.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]PiePoint, 0, len(rows))
	for _, rrow := range rows {
		out = append(out, PiePoint{
			ProjectID:   rrow.ProjectID,
			ProjectName: rrow.ProjectName,
			Value:       rrow.Value,
		})
	}
	return out, nil
}

// VolunteerHoursPoint represents aggregated volunteer hours per date
type VolunteerHoursPoint struct {
	Date  string  `json:"date"`
	Hours float64 `json:"hours"`
}

// VolunteerHoursByVolunteer aggregates hours_worked from schedules per day for a volunteer
func (r *ChartRepository) VolunteerHoursByVolunteer(volunteerID uint, start, end *time.Time) ([]VolunteerHoursPoint, error) {
	var rows []struct {
		Date  string  `gorm:"column:date"`
		Hours float64 `gorm:"column:sum_hours"`
	}

	tx := r.db.Model(&models.Schedule{}).Select("date(shift_date) as date, sum(hours_worked) as sum_hours").Where("person_type = ? AND person_id = ?", "volunteer", volunteerID)
	if start != nil {
		tx = tx.Where("date(shift_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		tx = tx.Where("date(shift_date) <= ?", end.Format("2006-01-02"))
	}
	tx = tx.Group("date(shift_date)").Order("date(shift_date)")

	if err := tx.Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]VolunteerHoursPoint, 0, len(rows))
	for _, rr := range rows {
		out = append(out, VolunteerHoursPoint{Date: rr.Date, Hours: rr.Hours})
	}
	return out, nil
}
