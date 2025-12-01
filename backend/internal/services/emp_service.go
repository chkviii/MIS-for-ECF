package services

import (
	"log"
	"strings"

	"encoding/json"
	"erp-backend/internal/models"
	"erp-backend/internal/repo"
)

// EmpService 员工服务
type EmpService struct {
	empRepo     *repo.EmployeeRepository
	projRepo    *repo.ProjectRepository
	empProjRepo *repo.EmployeeProjectRepository
}

type EmpInternalProjectRequest struct {
	EmployeeID uint `json:"role_id" binding:"required"`
}

type EmpInternalProject struct {
	ProjectID uint   `json:"project_id"`
	Name      string `json:"name"`
	Link      string `json:"link"`
}

type EmpInternalProjectsResponse struct {
	Projects []EmpInternalProject `json:"projects"`
}

type EmpProjectRequest struct {
	EmployeeID uint `json:"role_id" binding:"required"`
}

type EmpProjectsResponse struct {
	Projects []models.EmployeeProject `json:"projects"`
}

func NewEmpService(empRepo *repo.EmployeeRepository, projRepo *repo.ProjectRepository, empProjRepo *repo.EmployeeProjectRepository) *EmpService {
	return &EmpService{
		empRepo:     empRepo,
		projRepo:    projRepo,
		empProjRepo: empProjRepo,
	}
}

// GetInternalProjects 获取员工参与的内部项目
func (s *EmpService) GetInternalProjects(req *EmpInternalProjectRequest) ([]models.Project, error) {
	log.Printf("Fetching internal projects for employee ID: %d", req.EmployeeID)
	empProjects, err := s.empProjRepo.Search(map[string]interface{}{"employee_id": req.EmployeeID})
	if err != nil {
		log.Printf("Error fetching employee projects: %v", err)
		return nil, err
	}
	if len(empProjects) == 0 {
		log.Printf("No projects found for employee ID: %d", req.EmployeeID)
		return nil, nil
	}
	var projects []models.Project
	var EmpInternalProjects EmpInternalProject
	var response EmpInternalProjectsResponse

	for _, ep := range empProjects {
		project, err := s.projRepo.Search(map[string]interface{}{"id": ep.ProjectID})
		if err != nil {
			log.Printf("Error fetching project ID %d: %v", ep.ProjectID, err)
			return nil, err
		}

		if len(project) == 1 {
			//check if project is internal
			projectType := project[0].ProjectType
			if strings.Contains(projectType, "\"internal\"") {
				var data map[string]interface{}
				err := json.Unmarshal([]byte(project[0].Description), &data)
				if err != nil {
					log.Printf("Error marshalling project description for project ID %d: %v", ep.ProjectID, err)
					return nil, err
				}
				if link, ok := data["link"].(string); !ok {
					link = ""
				} else {
					EmpInternalProjects.Link = link
				}

				EmpInternalProjects.ProjectID = project[0].ID
				EmpInternalProjects.Name = project[0].Name

				response.Projects = append(response.Projects, EmpInternalProjects)

			}
		} else {
			log.Printf("Project ID %d not found or multiple entries exist", ep.ProjectID)
		}

	}

	log.Printf("Found %d internal projects for employee ID: %d", len(response.Projects), req.EmployeeID)
	return projects, nil
}

// GetProjects 获取员工参与的所有项目 (不包含Internal Projects)
func (s *EmpService) GetProjects(req *EmpProjectRequest) ([]models.EmployeeProject, error) {
	log.Printf("Fetching all projects for employee ID: %d", req.EmployeeID)
	empProjects, err := s.empProjRepo.EmpGetProjects(req.EmployeeID)
	if err != nil {
		log.Printf("Error fetching employee projects: %v", err)
		return nil, err
	}
	log.Printf("Found %d projects for employee ID: %d", len(empProjects), req.EmployeeID)
	return empProjects, nil
}
