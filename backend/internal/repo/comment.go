package repo

import (
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) GetByArticleID(articleID string, limit, offset int) ([]Comment, error) {
	var comments []Comment
	err := r.db.Preload("User").
		Where("article_id = ?", articleID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) GetAll(limit, offset int) ([]Comment, error) {
	var comments []Comment
	err := r.db.Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&Comment{}).Error
}

func (r *CommentRepository) DeleteByAdmin(id uint) error {
	return r.db.Delete(&Comment{}, id).Error
}

func (r *CommentRepository) GetByID(id uint) (*Comment, error) {
	var comment Comment
	err := r.db.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}