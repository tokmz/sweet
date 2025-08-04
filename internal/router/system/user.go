package system

import (
	"github.com/gin-gonic/gin"

	"sweet/internal/handler/system"
	systemService "sweet/internal/service/system"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(r *gin.RouterGroup, systemSvc *systemService.SystemService) {
	userHandler := system.NewUserHandler(systemSvc.UserService)

	// 用户管理路由组
	userGroup := r.Group("/users")
	{
		// 用户CRUD操作
		userGroup.POST("", userHandler.CreateUser)           // 创建用户
		userGroup.GET("", userHandler.GetUserList)           // 获取用户列表
		userGroup.GET("/:id", userHandler.GetUser)           // 获取用户详情
		userGroup.PUT("/:id", userHandler.UpdateUser)        // 更新用户
		userGroup.DELETE("/:id", userHandler.DeleteUser)     // 删除用户
		
		// 用户状态管理
		userGroup.PUT("/:id/status", userHandler.UpdateUserStatus) // 更新用户状态
		
		// 密码管理
		userGroup.PUT("/password", userHandler.ChangePassword) // 修改密码
	}
}