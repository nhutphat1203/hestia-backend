package db

import (
	"log"

	"github.com/nhutphat1203/hestia-backend/internal/config"
	"github.com/nhutphat1203/hestia-backend/internal/model"
	hasher "github.com/nhutphat1203/hestia-backend/pkg/hash"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB, cfg *config.Config) {
	// Hash mật khẩu admin từ file .env
	hashedPassword, err := hasher.Hash(cfg.AdminPwd)
	if err != nil {
		log.Fatalf("Không thể hash mật khẩu admin: %v", err)
	}

	admin := model.User{
		Account:        cfg.AdminAcc,
		HashedPassword: hashedPassword,
		Name:           "Admin đẹp trai",
	}

	var existing model.User
	err = db.Where("account = ?", cfg.AdminAcc).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Không thể tạo admin mặc định: %v", err)
		}
		log.Println("✅ Đã seed user admin mặc định thành công")
		return
	}

	if err != nil {
		log.Fatalf("Lỗi khi kiểm tra user admin: %v", err)
	}

	log.Println("ℹ️  Admin mặc định đã tồn tại, bỏ qua seeding")
}
