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

func (r *CommentRepository) GetByID(id uint) (*Comment, error) {
	var comment Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) GetByArticleID(articleID string) ([]Comment, error) {
	var comments []Comment
	err := r.db.Where("article_id = ?", articleID).Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) FetchByArticleID(articleID string, limit, offset int) ([]Comment, error) {
	var comments []Comment
	err := r.db.Where("article_id = ? AND parent_id IS NULL", articleID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&Comment{}, id).Error
}
