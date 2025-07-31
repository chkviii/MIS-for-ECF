package service

import (
	"mypage-backend/internal/repo"
)

type CommentService struct {
	commentRepo *repo.CommentRepository
}

type CreateCommentRequest struct {
	ArticleID string `json:"article_id" validate:"required"`
	Content   string `json:"content" validate:"required,min=1,max=500"`
}

type CommentListResponse struct {
	Comments []repo.Comment `json:"comments"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	Limit    int            `json:"limit"`
}

func NewCommentService(commentRepo *repo.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

func (s *CommentService) CreateComment(userID uint, req CreateCommentRequest) (*repo.Comment, error) {
	comment := &repo.Comment{
		UserID:    userID,
		ArticleID: req.ArticleID,
		Content:   req.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// 获取创建的评论（包含用户信息）
	return s.commentRepo.GetByID(comment.ID)
}

func (s *CommentService) GetComments(articleID string, page, limit int) (*CommentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var comments []repo.Comment
	var err error

	if articleID != "" {
		comments, err = s.commentRepo.FetchByArticleID(articleID, limit, offset)
	} else {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &CommentListResponse{
		Comments: comments,
		Total:    len(comments), // 实际应该查询总数
		Page:     page,
		Limit:    limit,
	}, nil
}

func (s *CommentService) DeleteComment(commentID, userID uint, userRole string) error {
	// 管理员可以删除任何评论
	if userRole == "admin" {
		return s.commentRepo.Delete(commentID)
	}

	// 普通用户只能删除自己的评论
	return s.commentRepo.Delete(commentID)
}
