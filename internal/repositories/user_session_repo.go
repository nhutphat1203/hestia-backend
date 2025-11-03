package repository

import (
	"github.com/nhutphat1203/hestia-backend/internal/model"
	"gorm.io/gorm"
)

type UserSessionRepository interface {
	GetSessionByToken(hashed_token string) (model.UserSession, error)
	GetSessionByUserIdAndToken(userId uint, hashed_token string) (model.UserSession, error)
	Create(userId uint, hashed_token string) (uint, error)
	Update(session model.UserSession) error
	GetLatestSessionByUserId(userId uint) (model.UserSession, error)
}

type UserSessionRepo struct {
	db *gorm.DB
}

func NewUserSessionRepo(db *gorm.DB) UserSessionRepository {
	return &UserSessionRepo{
		db: db,
	}
}

func (r *UserSessionRepo) GetSessionByToken(hashed_token string) (model.UserSession, error) {
	var session model.UserSession
	result := r.db.Where("hashed_refresh_token = ?", hashed_token).First(&session)
	return session, result.Error
}

func (r *UserSessionRepo) Create(userId uint, hashed_token string) (uint, error) {
	session := model.UserSession{
		UserID:             userId,
		HashedRefreshToken: hashed_token,
	}
	result := r.db.Create(&session)
	return session.ID, result.Error
}

func (r *UserSessionRepo) GetLatestSessionByUserId(userId uint) (model.UserSession, error) {
	var session model.UserSession
	result := r.db.Where("user_id = ?", userId).Order("created_at desc").First(&session)
	return session, result.Error
}

func (r *UserSessionRepo) Update(session model.UserSession) error {
	result := r.db.Save(&session)
	return result.Error
}
func (r *UserSessionRepo) GetSessionByUserIdAndToken(userId uint, hashed_token string) (model.UserSession, error) {
	var session model.UserSession
	result := r.db.Where("user_id = ? AND hashed_refresh_token = ?", userId, hashed_token).First(&session)
	return session, result.Error
}
