package repository

import (
	"github.com/nhutphat1203/hestia-backend/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByAccount(account string) (model.User, error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &UserRepo{
		db: db,
	}
}
func (r *UserRepo) GetUserByAccount(account string) (model.User, error) {
	var user model.User
	result := r.db.Where("account = ?", account).First(&user)
	return user, result.Error
}
