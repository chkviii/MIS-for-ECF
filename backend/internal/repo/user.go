package repo

import (
	"mypage-backend/internal/models"
	"strconv"
	"sync"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var userRepoOnce sync.Once
var userRepo *UserRepository

func GetUserRepo() *UserRepository {
	userRepoOnce.Do(func() {
		db := GetDB()
		userRepo = &UserRepository{db: db}
	})
	return userRepo
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(id uint, user *models.User) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *UserRepository) FindAll() ([]*models.User, error) {
	var users []*models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(idstr string) (*models.User, error) {
	// Convert string ID to uint first
	uid, strerr := strconv.ParseUint(idstr, 10, 32)
	if strerr != nil {
		return nil, strerr
	}
	var user models.User
	err := r.db.First(&user, uid).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
