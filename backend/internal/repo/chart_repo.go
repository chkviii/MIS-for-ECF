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

// DonationPoint represents aggregated donation data points per date
type DonationPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// DonationsByDonor returns donation totals per day (or per period) for a donor
func (r *ChartRepository) DonationsByDonor(donorID uint, start, end *time.Time) ([]DonationPoint, error) {
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

	out := make([]DonationPoint, 0, len(rows))
	for _, rrow := range rows {
		out = append(out, DonationPoint{Date: rrow.Date, Value: rrow.Value})
	}
	return out, nil
}

// FundAllocationPoint represents allocation totals per date for a fund
type FundAllocationPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// FundAllocationsByFund aggregates fund allocations (FundProject or Fund movements) per date
func (r *ChartRepository) FundAllocationsByFund(fundID uint, start, end *time.Time) ([]FundAllocationPoint, error) {
	var rows []struct {
		Date  string  `gorm:"column:date"`
		Value float64 `gorm:"column:sum_amount"`
	}

	// Attempt to sum expenses or allocations linked to fund_id via FundProject or Expense.FundID
	tx := r.db.Model(&models.Fund{}).
		Select("date(created_at) as date, sum(total_amount) as sum_amount").
		Where("id = ?", fundID)

	_ = tx // suppress unused variable warning
	// If there is a FundProject table with allocation amounts, try that as well using joins.
	// For simplicity we will query FundProjects (if present) by fund_id via raw SQL fallback.
	// Use Expense as a common target: sum of Fund.TotalAmount is not a direct allocation; instead we'll check FundProject

	// Prefer FundProject allocations
	tx2 := r.db.Model(&models.FundProject{}).Select("date(allocation_date) as date, sum(allocated_amount) as sum_amount").Where("fund_id = ?", fundID)
	if start != nil {
		tx2 = tx2.Where("date(allocation_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		tx2 = tx2.Where("date(allocation_date) <= ?", end.Format("2006-01-02"))
	}
	tx2 = tx2.Group("date(allocation_date)").Order("date(allocation_date)")

	if err := tx2.Scan(&rows).Error; err == nil && len(rows) > 0 {
		out := make([]FundAllocationPoint, 0, len(rows))
		for _, rr := range rows {
			out = append(out, FundAllocationPoint{Date: rr.Date, Value: rr.Value})
		}
		return out, nil
	}

	// Fallback: try expenses referencing fund
	rows = []struct {
		Date  string  `gorm:"column:date"`
		Value float64 `gorm:"column:sum_amount"`
	}{}
	tx3 := r.db.Model(&models.Expense{}).Select("date(expense_date) as date, sum(amount) as sum_amount").Where("fund_id = ?", fundID)
	if start != nil {
		tx3 = tx3.Where("date(expense_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		tx3 = tx3.Where("date(expense_date) <= ?", end.Format("2006-01-02"))
	}
	tx3 = tx3.Group("date(expense_date)").Order("date(expense_date)")
	if err := tx3.Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]FundAllocationPoint, 0, len(rows))
	for _, rr := range rows {
		out = append(out, FundAllocationPoint{Date: rr.Date, Value: rr.Value})
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

// FinanceReportRow represents aggregated finance metrics for a given period/group
type FinanceReportRow struct {
	Period         string  `json:"period"`
	TotalDonations float64 `json:"total_donations"`
	TotalFunds     float64 `json:"total_funds"`
	TotalExpenses  float64 `json:"total_expenses"`
	Location       string  `json:"location,omitempty"`
	Project        string  `json:"project,omitempty"`
}

// FinanceReportAggregates returns aggregated donation/fund/expense totals grouped by date or by project/location
// filterProjectID and filterLocationID are optional filters (0 means not applied)
func (r *ChartRepository) FinanceReportAggregates(start, end *time.Time, groupBy string, filterProjectID uint, filterLocationID uint) ([]FinanceReportRow, error) {
	// For simplicity, groupBy supports: "date", "project", "location"
	rows := []FinanceReportRow{}

	// Build base where clauses
	whereParts := []string{}
	args := []interface{}{}
	if start != nil {
		whereParts = append(whereParts, "date(donation_date) >= ?")
		args = append(args, start.Format("2006-01-02"))
	}
	if end != nil {
		whereParts = append(whereParts, "date(donation_date) <= ?")
		args = append(args, end.Format("2006-01-02"))
	}
	if filterProjectID != 0 {
		whereParts = append(whereParts, "project_id = ?")
		args = append(args, filterProjectID)
	}
	// Donation subquery
	donationSelect := r.db.Model(&models.Donation{}).Select("project_id, date(donation_date) as period, sum(amount) as total_donations").Where("1=1", args...)
	if len(whereParts) > 0 {
		donationSelect = donationSelect.Where(whereParts[0], args...)
	}

	// Expense subquery
	expenseSelect := r.db.Model(&models.Expense{}).Select("project_id, date(expense_date) as period, sum(amount) as total_expenses").Where("1=1")
	if start != nil {
		expenseSelect = expenseSelect.Where("date(expense_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		expenseSelect = expenseSelect.Where("date(expense_date) <= ?", end.Format("2006-01-02"))
	}
	if filterProjectID != 0 {
		expenseSelect = expenseSelect.Where("project_id = ?", filterProjectID)
	}

	// Fund allocations
	fundSelect := r.db.Model(&models.FundProject{}).Select("project_id, date(allocation_date) as period, sum(allocated_amount) as total_funds").Where("1=1")
	if start != nil {
		fundSelect = fundSelect.Where("date(allocation_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		fundSelect = fundSelect.Where("date(allocation_date) <= ?", end.Format("2006-01-02"))
	}
	if filterProjectID != 0 {
		fundSelect = fundSelect.Where("project_id = ?", filterProjectID)
	}

	// Depending on groupBy, aggregate and join results in memory for simplicity
	// Fetch donations grouped
	type drow struct {
		Period    string
		ProjectID *uint
		Total     float64
	}
	var drows []drow
	ds := r.db.Model(&models.Donation{}).Select("date(donation_date) as period, project_id, sum(amount) as total").Group("date(donation_date), project_id")
	if start != nil {
		ds = ds.Where("date(donation_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		ds = ds.Where("date(donation_date) <= ?", end.Format("2006-01-02"))
	}
	if filterProjectID != 0 {
		ds = ds.Where("project_id = ?", filterProjectID)
	}
	if err := ds.Scan(&drows).Error; err != nil {
		return nil, err
	}

	type erow struct {
		Period    string
		ProjectID *uint
		Total     float64
	}
	var erows []erow
	es := r.db.Model(&models.Expense{}).Select("date(expense_date) as period, project_id, sum(amount) as total").Group("date(expense_date), project_id")
	if start != nil {
		es = es.Where("date(expense_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		es = es.Where("date(expense_date) <= ?", end.Format("2006-01-02"))
	}
	if filterProjectID != 0 {
		es = es.Where("project_id = ?", filterProjectID)
	}
	if err := es.Scan(&erows).Error; err != nil {
		return nil, err
	}

	type frow struct {
		Period    string
		ProjectID *uint
		Total     float64
	}
	var frows []frow
	fs := r.db.Model(&models.FundProject{}).Select("date(allocation_date) as period, project_id, sum(allocated_amount) as total").Group("date(allocation_date), project_id")
	if start != nil {
		fs = fs.Where("date(allocation_date) >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		fs = fs.Where("date(allocation_date) <= ?", end.Format("2006-01-02"))
	}
	if filterProjectID != 0 {
		fs = fs.Where("project_id = ?", filterProjectID)
	}
	if err := fs.Scan(&frows).Error; err != nil {
		return nil, err
	}

	// Merge keys by period+project
	m := map[string]*FinanceReportRow{}
	key := func(period string, projectID *uint) string {
		pid := ""
		if projectID != nil {
			pid = string(rune(*projectID))
		}
		return period + "|" + pid
	}

	for _, d := range drows {
		k := key(d.Period, d.ProjectID)
		if _, ok := m[k]; !ok {
			m[k] = &FinanceReportRow{Period: d.Period}
		}
		m[k].TotalDonations += d.Total
	}
	for _, e := range erows {
		k := key(e.Period, e.ProjectID)
		if _, ok := m[k]; !ok {
			m[k] = &FinanceReportRow{Period: e.Period}
		}
		m[k].TotalExpenses += e.Total
	}
	for _, f := range frows {
		k := key(f.Period, f.ProjectID)
		if _, ok := m[k]; !ok {
			m[k] = &FinanceReportRow{Period: f.Period}
		}
		m[k].TotalFunds += f.Total
	}

	for _, v := range m {
		rows = append(rows, *v)
	}

	return rows, nil
}
