package common

import (
	"net/http"
	"sweet/internal/models"
	"sweet/pkg/auth"
	"sweet/pkg/errs"
	"sweet/pkg/logger"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	Gin  *ginUtils
	once sync.Once
)

type ginUtils struct {
	l logger.Logger
}

func NewGinUtils(l logger.Logger) {
	once.Do(func() {
		Gin = &ginUtils{
			l: l,
		}
	})
}

// Bind 绑定参数
func (g *ginUtils) Bind(ctx *gin.Context, obj interface{}) error {
	if err := ctx.ShouldBind(obj); err != nil {
		g.l.Error("Bind error", zap.Error(err))
		return errs.ErrParams
	}
	return nil
}

// BindQuery 绑定Query参数
func (g *ginUtils) BindQuery(ctx *gin.Context, obj interface{}) error {
	if err := ctx.ShouldBindQuery(obj); err != nil {
		g.l.Error("BindQuery error", zap.Error(err))
		return errs.ErrParams
	}
	return nil
}

// BindHeader 绑定Header参数
func (g *ginUtils) BindHeader(ctx *gin.Context, obj interface{}) error {
	if err := ctx.ShouldBindHeader(obj); err != nil {
		g.l.Error("BindHeader error", zap.Error(err))
		return errs.ErrParams
	}
	return nil
}

// BindUri 绑定Uri参数
func (g *ginUtils) BindUri(ctx *gin.Context, obj interface{}) error {
	if err := ctx.ShouldBindUri(obj); err != nil {
		g.l.Error("BindUri error", zap.Error(err))
		return errs.ErrParams
	}
	return nil
}

// Res 响应
func (g *ginUtils) Res(ctx *gin.Context, err error, data ...any) {
	if err != nil {
		// 判断err 类型
		if e, ok := err.(*errs.Error); ok {
			ctx.JSON(http.StatusOK, models.NewResponse(e.Code, e.Msg, nil))
			return
		} else {
			g.l.Error("Res error", zap.Error(err))
			ctx.JSON(http.StatusOK, models.NewResponse(errs.ErrServer.Code, errs.ErrServer.Msg, nil))
			return
		}
	}

	// 判断是否有响应
	if len(data) == 0 {
		ctx.JSON(http.StatusOK, models.NewResponse(200, "success", nil))
		return
	}

	// 响应
	ctx.JSON(http.StatusOK, models.NewResponse(200, "success", data[0]))
}

// GetClaims 从gin.Context中获取完整的Claims信息
func (g *ginUtils) GetClaims(c *gin.Context) (*auth.Claims, bool) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	claims, ok := claimsVal.(*auth.Claims)
	return claims, ok
}

// Uid 从gin.Context中获取用户ID
func (g *ginUtils) Uid(c *gin.Context) (int64, bool) {
	if uid := c.GetInt64("uid"); uid != 0 {
		return uid, true
	}
	return 0, false
}
