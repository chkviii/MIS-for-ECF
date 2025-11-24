package services

import (
	"errors"
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

	// 创建用户
	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		UserType:     req.UserType,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	// 根据用户类型创建对应的详细信息
	switch req.UserType {
	case "employee":
		employee := &models.Employee{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
			Status:    "pending", // set to pending for employee
			CreatedAt: time.Now(),
		}
		if err := s.employeeRepo.Create(employee); err != nil {
			return nil, errors.New("failed to create employee information: " + err.Error())
		}

	case "volunteer":
		volunteer := &models.Volunteer{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
			Status:    "active",
			CreatedAt: time.Now(),
		}
		if err := s.volunteerRepo.Create(volunteer); err != nil {
			return nil, errors.New("failed to create volunteer information: " + err.Error())
		}

	case "donor":
		donor := &models.Donor{
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
	token, err := utils.GenerateToken(int64(user.ID), user.Username, user.UserType)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token:    token,
		UserType: user.UserType,
		UserID:   user.ID,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// 查找用户
	users, err := s.userRepo.Search(map[string]interface{}{"username": req.Username})
	if err != nil && len(users) != 1 {
		return nil, errors.New("wrong username or password")
	}

	user := users[0]

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("wrong username or password")
	}

	// 检查用户是否激活
	if user.Status != "active" {
		return nil, errors.New("user has been disabled")
	}

	// 生成 JWT Token
	token, err := utils.GenerateToken(int64(user.ID), user.Username, user.UserType)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token:    token,
		UserType: user.UserType,
		UserID:   user.ID,
	}, nil
}
