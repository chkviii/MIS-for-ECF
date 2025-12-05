package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"erp-backend/internal/models"
	"erp-backend/internal/repo"
	"erp-backend/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务
type AuthService struct {
	userRepo      *repo.UserRepository
	employeeRepo  *repo.EmployeeRepository
	volunteerRepo *repo.VolunteerRepository
	donorRepo     *repo.DonorRepository
}

// NewAuthService 创建认证服务实例
// 依赖多个 Repository 来操作不同的表
func NewAuthService(
	userRepo *repo.UserRepository,
	employeeRepo *repo.EmployeeRepository,
	volunteerRepo *repo.VolunteerRepository,
	donorRepo *repo.DonorRepository,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		employeeRepo:  employeeRepo,
		volunteerRepo: volunteerRepo,
		donorRepo:     donorRepo,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	UserType  string `json:"user_type" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token    string `json:"token"`
	UserType string `json:"user_type"`
	UserID   uint   `json:"user_id"`
	RoleID   uint   `json:"role_id"`
}

// generateID 生成唯一ID
func generateID(prefix string) string {
	p := strings.ToUpper(prefix)
	if len(p) > 3 {
		p = p[:3]
	} else if len(p) < 3 {
		p = fmt.Sprintf("%3s", p)
	}
	t := time.Now().UTC()
	date := t.Format("060102") // YYMMDD
	tm := t.Format("150405")   // HHMMSS
	return fmt.Sprintf("%s%s%s", p, date, tm)
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// 检查用户名是否已存在
	existingUsers, err := s.userRepo.Search(map[string]interface{}{"username": req.Username})
	if err == nil && len(existingUsers) > 0 {
		return nil, errors.New("username already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	var status string

	if req.UserType == "donor" || req.UserType == "volunteer" {
		status = "active"
	} else {
		status = "pending"
	}

	// 创建用户
	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		UserType:     req.UserType,
		Status:       status,
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	// 根据用户类型创建对应的详细信息
	switch req.UserType {
	case "employee":
		employee := &models.Employee{
			UserID:     &user.ID,
			EmployeeID: generateID("EMP"),
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Email:      req.Email,
			Phone:      req.Phone,
			Status:     "pending", // set to pending for employee
			CreatedAt:  time.Now(),
		}
		if err := s.employeeRepo.Create(employee); err != nil {
			return nil, errors.New("failed to create employee information: " + err.Error())
		}

	case "volunteer":
		volunteer := &models.Volunteer{
			UserID:      &user.ID,
			VolunteerID: generateID("VOL"),
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       req.Email,
			Phone:       req.Phone,
			Status:      "active",
			CreatedAt:   time.Now(),
		}
		if err := s.volunteerRepo.Create(volunteer); err != nil {
			return nil, errors.New("failed to create volunteer information: " + err.Error())
		}

	case "donor":
		donor := &models.Donor{
			UserID:    &user.ID,
			DonorID:   generateID("DON"),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
			Status:    "active",
			CreatedAt: time.Now(),
		}
		if err := s.donorRepo.Create(donor); err != nil {
			return nil, errors.New("failed to create donor information: " + err.Error())
		}

	default:
		return nil, errors.New("invalid user type")
	}

	// 生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.UserType, 0)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token:    token,
		UserType: user.UserType,
		UserID:   user.ID,
		RoleID:   0,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// 查找用户
	users, err := s.userRepo.Search(map[string]interface{}{"username": req.Username})
	if err != nil || len(users) != 1 {
		return nil, errors.New("wrong username or password")
	}

	user := users[0]

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("wrong username or password")
	}

	// 检查用户是否激活
	if user.Status != "active" {
		return nil, errors.New("user account not active")
	}

	var role_id uint

	switch user.UserType {
	case "employee":
		employees, err := s.employeeRepo.Search(map[string]interface{}{"user_id": user.ID})
		if err != nil || len(employees) != 1 {
			return nil, errors.New("employee profile not found")
		} else {
			role_id = employees[0].ID
		}
	case "volunteer":
		volunteers, err := s.volunteerRepo.Search(map[string]interface{}{"user_id": user.ID})
		if err != nil || len(volunteers) != 1 {
			return nil, errors.New("volunteer profile not found")
		} else {
			role_id = volunteers[0].ID
		}
	case "donor":
		donors, err := s.donorRepo.Search(map[string]interface{}{"user_id": user.ID})
		if err != nil || len(donors) != 1 {
			return nil, errors.New("donor profile not found")
		} else {
			role_id = donors[0].ID
		}
	default:
		return nil, errors.New("invalid user type")
	}

	// 生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.UserType, role_id)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// update last login time
	now := time.Now().UTC()
	user.LastLogin = &now
	s.userRepo.Update(&user)

	return &AuthResponse{
		Token:    token,
		UserType: user.UserType,
		UserID:   user.ID,
		RoleID:   role_id,
	}, nil
}
