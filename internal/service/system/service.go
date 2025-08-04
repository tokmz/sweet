package system

import (
	"gorm.io/gorm"
)

// SystemService 系统服务集合
type SystemService struct {
	UserService IUserService
}

// NewSystemService 创建系统服务实例
func NewSystemService(db *gorm.DB) *SystemService {
	return &SystemService{
		UserService: NewUserService(db),
	}
}