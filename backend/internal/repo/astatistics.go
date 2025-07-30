package repo

import (
	"gorm.io/gorm"
)

type AstatisticsRepository struct {
	db *gorm.DB
}

func NewAstatisticsRepository(db *gorm.DB) *AstatisticsRepository {
	return &AstatisticsRepository{db: db}
}

func (r *AstatisticsRepository) Create(stats *ArticleStatistics) error {
	return r.db.Create(stats).Error
}

func (r *AstatisticsRepository) Update(id uint, stats *ArticleStatistics) error {
	return r.db.Model(&ArticleStatistics{}).Where("id = ?", id).Updates(stats).Error
}

func (r *AstatisticsRepository) GetByID(id uint) (*ArticleStatistics, error) {
	var stats ArticleStatistics
	err := r.db.First(&stats, id).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
