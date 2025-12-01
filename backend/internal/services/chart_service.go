package services

import (
	"time"

	"erp-backend/internal/repo"
)

// ChartService exposes aggregation methods for charting
type ChartService struct {
	repo *repo.ChartRepository
}

// NewChartService creates a new ChartService
func NewChartService(r *repo.ChartRepository) *ChartService {
	return &ChartService{repo: r}
}

// DonationsByDonor returns aggregated donation points for a donor
func (s *ChartService) DonationsByDonor(donorID uint, start, end *time.Time) ([]repo.DonationPoint, error) {
	return s.repo.DonationsByDonor(donorID, start, end)
}

// FundAllocationsByFund returns aggregated fund allocation points for a fund
func (s *ChartService) FundAllocationsByFund(fundID uint, start, end *time.Time) ([]repo.FundAllocationPoint, error) {
	return s.repo.FundAllocationsByFund(fundID, start, end)
}

// VolunteerHoursByVolunteer returns aggregated volunteer hours points
func (s *ChartService) VolunteerHoursByVolunteer(volunteerID uint, start, end *time.Time) ([]repo.VolunteerHoursPoint, error) {
	return s.repo.VolunteerHoursByVolunteer(volunteerID, start, end)
}

// FinanceReportAggregates returns finance aggregates grouped as requested
func (s *ChartService) FinanceReportAggregates(start, end *time.Time, groupBy string, filterProjectID uint, filterLocationID uint) ([]repo.FinanceReportRow, error) {
	return s.repo.FinanceReportAggregates(start, end, groupBy, filterProjectID, filterLocationID)
}
