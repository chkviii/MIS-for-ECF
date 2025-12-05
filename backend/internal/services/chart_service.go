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
func (s *ChartService) DonationsByDonor(donorID uint, start, end *time.Time) ([]repo.LinePoint, error) {
	return s.repo.DonationsByDonor(donorID, start, end)
}

func (s *ChartService) DonorDonationsByProject(donorID uint, start, end *time.Time) ([]repo.PiePoint, error) {
	return s.repo.DonationsByProject(donorID, start, end)
}

// FundAllocationsByDate returns aggregated fund allocation points for a fund
func (s *ChartService) FundAllocationsByDate(start, end *time.Time) ([]repo.LinePoint, error) {
	return s.repo.FundAllocationsByDate(start, end)
}

func (s *ChartService) FundAllocationsByProject(start, end *time.Time) ([]repo.PiePoint, error) {
	return s.repo.FundAllocationsByProject(start, end)
}

// Expenses
func (s *ChartService) ExpensesByDate(start, end *time.Time) ([]repo.LinePoint, error) {
	return s.repo.ExpensesByDate(start, end)
}

func (s *ChartService) ExpensesByProject(start, end *time.Time) ([]repo.PiePoint, error) {
	return s.repo.ExpensesByProject(start, end)
}

// Donations
func (s *ChartService) DonationsByDate(start, end *time.Time) ([]repo.LinePoint, error) {
	return s.repo.DonationsByDate(start, end)
}

func (s *ChartService) DonationsByProject(start, end *time.Time) ([]repo.PiePoint, error) {
	return s.repo.DonationsByProject(0, start, end)
}

// VolunteerHoursByVolunteer returns aggregated volunteer hours points
func (s *ChartService) VolunteerHoursByVolunteer(volunteerID uint, start, end *time.Time) ([]repo.VolunteerHoursPoint, error) {
	return s.repo.VolunteerHoursByVolunteer(volunteerID, start, end)
}
