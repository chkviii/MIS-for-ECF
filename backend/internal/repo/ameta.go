package repo

import (
	"gorm.io/gorm"
)

type AMetaRepository struct {
	db *gorm.DB
}

func NewAMetaRepository(db *gorm.DB) *AMetaRepository {
	return &AMetaRepository{db: db}
}

func (r *AMetaRepository) Create(meta *ArticleMeta) error {
	return r.db.Create(meta).Error
}

func (r *AMetaRepository) Update(id uint, meta *ArticleMeta) error {
	return r.db.Model(&ArticleMeta{}).Where("id = ?", id).Updates(meta).Error
}

func (r *AMetaRepository) GetByID(id uint) (*ArticleMeta, error) {
	var meta ArticleMeta
	err := r.db.First(&meta, id).Error
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func (r *AMetaRepository) FindAll() ([]*ArticleMeta, error) {
	var metas []*ArticleMeta
	err := r.db.Find(&metas).Error
	if err != nil {
		return nil, err
	}
	return metas, nil
}

func (r *AMetaRepository) FindByAuthorID(authorID string) ([]*ArticleMeta, error) {
	var metas []*ArticleMeta
	err := r.db.Where("author_id = ?", authorID).Find(&metas).Error
	if err != nil {
		return nil, err
	}
	return metas, nil
}

func (r *AMetaRepository) FetchByTitle(title string, limit, offset int) ([]*ArticleMeta, error) {
	var metas []*ArticleMeta
	err := r.db.Where("LOWER(title) LIKE LOWER(?)", "%"+title+"%").
		Limit(limit).Offset(offset).
		Find(&metas).Error
	if err != nil {
		return nil, err
	}
	return metas, nil
}

func (r *AMetaRepository) FetchByTag(tag string, limit, offset int) ([]*ArticleMeta, error) {
	var metas []*ArticleMeta
	err := r.db.Where("LOWER(tags) LIKE LOWER(?)", "%"+tag+"%").
		Limit(limit).Offset(offset).
		Find(&metas).Error
	if err != nil {
		return nil, err
	}
	return metas, nil
}
