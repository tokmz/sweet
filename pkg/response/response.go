package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
	Total   int64       `json:"total"`   // 总数
	Page    int         `json:"page"`    // 当前页
	Size    int         `json:"size"`    // 每页数量
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, message string, detail interface{}) {
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
		Data:    detail,
	})
}

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string, detail interface{}) {
	Error(c, http.StatusBadRequest, message, detail)
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

// NotFound 404错误响应
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

// InternalServerError 500错误响应
func InternalServerError(c *gin.Context, message string, detail interface{}) {
	Error(c, http.StatusInternalServerError, message, detail)
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    200,
		Message: "success",
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// PageSuccessWithMessage 分页成功响应（自定义消息）
func PageSuccessWithMessage(c *gin.Context, message string, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    200,
		Message: message,
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}