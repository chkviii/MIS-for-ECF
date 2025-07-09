package handler

import (
	"mypage-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req service.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的请求格式",
		})
	}

	// 基本验证
	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "用户名和密码不能为空",
		})
	}

	if len(req.Username) < 3 || len(req.Username) > 20 {
		return c.Status(400).JSON(fiber.Map{
			"error": "用户名长度必须在3-20个字符之间",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{
			"error": "密码长度至少6个字符",
		})
	}

	result, err := h.userService.Register(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "注册成功",
		"data":    result,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req service.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的请求格式",
		})
	}

	// 基本验证
	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "用户名和密码不能为空",
		})
	}

	result, err := h.userService.Login(req)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "登录成功",
		"data":    result,
	})
}

func (h *AuthHandler) VerifyToken(c *fiber.Ctx) error {
	// 从中间件获取用户信息
	userID := c.Locals("user_id").(uint)
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)

	return c.JSON(fiber.Map{
		"message": "Token验证成功",
		"user": fiber.Map{
			"id":       userID,
			"username": username,
			"role":     role,
		},
	})
}