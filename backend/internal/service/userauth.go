package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"time"

	"mypage-backend/internal/repo"
	"mypage-backend/internal/util"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo  *repo.UserRepository
	jwtSecret string
}

type PreLoginRequest struct {
	Username string `json:"username" validate:"required"`
}

type PreLoginResponse struct {
	SessID     string `json:"session_id"`
	SessExpire string `json:"session_expire"`
	SrvPubKey  string `json:"server_pubkey"` // server public key
	Challenge  string `json:"challenge"`     // encrypted session key
}

type LoginRequest struct {
	Solution string `json:"solution" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"email"`

	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	Success    bool   `json:"success"`
	SessExpire string `json:"session_expire"`
}

func NewUserService(userRepo *repo.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) Register(req RegisterRequest) (*AuthResponse, error) {
	// 检查用户是否已存在
	return nil, nil
}

// prelogin: check if user exists and create a session with challenge
func (s *UserService) PreLogin(req PreLoginRequest) (PreLoginResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Username)
	if err != nil {
		user, err = s.userRepo.GetByUsername(req.Username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return PreLoginResponse{}, errors.New("user not found")
			}
			return PreLoginResponse{}, err
		}
	}

	// create a new session for login action
	session := GetSessMgr().NewSession(util.UitoA(user.ID))

	// decode user's public key
	userPubKey, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return PreLoginResponse{}, err
	}

	// encrypt session key and send to client
	challenge, err := util.ECDHEncrypt(userPubKey, session.Key)
	if err != nil {
		return PreLoginResponse{}, err
	}

	return PreLoginResponse{
		SessID:     session.ID,
		SessExpire: time.Unix(session.ExpiresAt, 0).Format(time.RFC3339),
		SrvPubKey:  base64.StdEncoding.EncodeToString(util.SrvEncPubKey()),
		Challenge:  base64.StdEncoding.EncodeToString(challenge),
	}, nil
}

// login: verify the solution
func (s *UserService) Login(sessID string, req LoginRequest) (*AuthResponse, error) {
	sessKey, err := GetSessMgr().GetKey(sessID)
	if err != nil {
		return nil, err
	}

	answerBytes := util.Sha256(sessKey)

	// decrypt the solution using session key
	solutionBytes, err := util.AESDecrypt(sessKey, []byte(req.Solution))
	if err != nil {
		return nil, err
	}

	if solutionBytes == nil && bytes.Equal(solutionBytes, answerBytes) {
		t, err := GetSessMgr().Refresh(sessID, 60)
		if err != nil {
			return nil, err
		}
		return &AuthResponse{
			Success:    true,
			SessExpire: time.Unix(t, 0).Format(time.RFC3339),
		}, nil
	} else {
		t, err := GetSessMgr().GetExpTime(sessID)
		if err != nil {
			return nil, err
		}
		return &AuthResponse{
			Success:    false,
			SessExpire: time.Unix(t, 0).Format(time.RFC3339),
		}, nil
	}
}
