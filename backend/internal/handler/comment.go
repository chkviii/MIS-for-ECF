package handler

import (
	"mypage-backend/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) GetComments(c *fiber.Ctx) error {
	articleID := c.Query("article_id", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result, err := h.commentService.GetComments(articleID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "获取评论失败",
		})
	}

	return c.JSON(fiber.Map{
		"message": "获取评论成功",
		"data":    result,
	})
}

func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req service.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的请求格式",
		})
	}

	// 基本验证
	if req.ArticleID == "" || req.Content == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "文章ID和评论内容不能为空",
		})
	}

	if len(req.Content) > 500 {
		return c.Status(400).JSON(fiber.Map{
			"error": "评论内容不能超过500个字符",
		})
	}

	comment, err := h.commentService.CreateComment(userID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "创建评论失败",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "评论创建成功",
		"data":    comment,
	})
}

func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	userRole := c.Locals("role").(string)

	commentID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的评论ID",
		})
	}

	if err := h.commentService.DeleteComment(uint(commentID), userID, userRole); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "删除评论失败",
		})
	}

	return c.JSON(fiber.Map{
		"message": "评论删除成功",
	})
}