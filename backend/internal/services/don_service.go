package services

import (
	"log"
	"time"

	"erp-backend/internal/repo"
)

// EmpService 员工服务
type DonService struct {
	donRepo  *repo.DonorRepository
	projRepo *repo.ProjectRepository
}

func NewDonService(donRepo *repo.DonorRepository, projRepo *repo.ProjectRepository, empProjRepo *repo.EmployeeProjectRepository) *DonService {
	return &DonService{
		donRepo:  donRepo,
		projRepo: projRepo,
	}
}

type project struct {
	ProjectID uint   `json:"id"`
	Name      string `json:"name"`
}

type DonProjectResponse struct {
	Projects []project `json:"projects"`
}

type DonDetailRequest struct {
	DonorID   uint
	ProjectID uint
	Start     time.Time
	End       time.Time
}

type DonDetail struct {
	ID      string    `json:"id"`
	Date    time.Time `json:"date"`
	Project string    `json:"project"`
	Amount  float64   `json:"amount"`
	Method  string    `json:"method"`
}

type DonDetailResponse struct {
	Details []DonDetail `json:"details"`
}

func (s *DonService) GetProjectsByDonor(donorID uint) (DonProjectResponse, error) {
	var projects DonProjectResponse
	donations, err := s.donRepo.DonorGetProjects(donorID)
	if err != nil {
		log.Printf("Error fetching donations: %v", err)
		return DonProjectResponse{}, err
	}

	for _, donation := range donations {
		proj := project{
			ProjectID: uint(donation["id"].(uint)),
			Name:      donation["name"].(string),
		}
		projects.Projects = append(projects.Projects, proj)
	}

	return projects, nil
}

func (s *DonService) GetDonationDetails(req DonDetailRequest) (DonDetailResponse, error) {
	var detail DonDetail
	var donationsResp DonDetailResponse
	donations, err := s.donRepo.GetDonationsByDonorAndFilters(req.DonorID, req.ProjectID, req.Start, req.End)
	if err != nil {
		log.Printf("Error fetching donation details: %v", err)
		return DonDetailResponse{}, err
	}

	for _, donation := range donations {
		detail = DonDetail{
			ID:      donation["donation_id"].(string),
			Date:    donation["donation_date"].(time.Time),
			Amount:  donation["amount"].(float64),
			Project: donation["project_name"].(string),
			Method:  donation["payment_method"].(string),
		}

		donationsResp.Details = append(donationsResp.Details, detail)
	}

	return donationsResp, nil
}
