package system

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	dto "sweet/internal/models/dto/system"
	"sweet/internal/service/system"
	"sweet/pkg/response"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService system.IUserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService system.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "用户信息"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Router /api/v1/system/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建用户失败", err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body dto.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Router /api/v1/system/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", "无效的用户ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	req.ID = id
	user, err := h.userService.UpdateUser(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新用户失败", err.Error())
		return
	}

	response.Success(c, user)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户
// @Tags 用户管理
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/system/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", "无效的用户ID")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除用户失败", err.Error())
		return
	}

	response.Success(c, nil)
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据ID获取用户详情
// @Tags 用户管理
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Router /api/v1/system/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", "无效的用户ID")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户失败", err.Error())
		return
	}

	response.Success(c, user)
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param username query string false "用户名"
// @Param realname query string false "真实姓名"
// @Param status query int false "状态"
// @Success 200 {object} response.Response{data=dto.UserListResponse}
// @Router /api/v1/system/users [get]
func (h *UserHandler) GetUserList(c *gin.Context) {
	var req dto.UserQueryRequest

	// 解析查询参数
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			req.Page = p
		}
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if size := c.Query("size"); size != "" {
		if s, err := strconv.Atoi(size); err == nil {
			req.Size = s
		}
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	req.Username = c.Query("username")
	req.Realname = c.Query("realname")

	if status := c.Query("status"); status != "" {
		if s, err := strconv.ParseInt(status, 10, 64); err == nil {
			req.Status = &s
		}
	}

	result, err := h.userService.GetUserList(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户列表失败", err.Error())
		return
	}

	response.Success(c, result)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param password body dto.ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Router /api/v1/system/users/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, "修改密码失败", err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 启用或禁用用户
// @Tags 用户管理
// @Param id path int true "用户ID"
// @Param status query int true "状态(1:启用 0:禁用)"
// @Success 200 {object} response.Response
// @Router /api/v1/system/users/{id}/status [put]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", "无效的用户ID")
		return
	}

	status, err := strconv.ParseInt(c.Query("status"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", "无效的状态值")
		return
	}

	if err := h.userService.UpdateUserStatus(c.Request.Context(), id, status); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新用户状态失败", err.Error())
		return
	}

	response.Success(c, nil)
}