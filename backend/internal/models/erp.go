package models

// import (
// 	"time"
// )

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	UserType  string `json:"user_type" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token   string      `json:"token"`
	User    *User       `json:"user"`
	Profile interface{} `json:"profile"`
}