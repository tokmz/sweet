package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sweet/pkg/logger"
)

func Bind(ctx *gin.Context, obj any) error {
	if err := ctx.ShouldBind(&obj); err != nil {
		logger.Error(
			"入参绑定失败",
			logger.Err(err),
		)
		return ErrInvalidParam
	}
	return nil
}

func Res(ctx *gin.Context, err error, data ...any) {
	tid := ctx.GetString("trace_id")

	if err != nil {
		var e *Error
		if !errors.As(err, &e) {
			ctx.JSON(http.StatusOK, NewResponse(ErrServer.Code, ErrServer.Msg, nil, tid))
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, NewResponse(e.Code, e.Msg, nil, tid))
		return
	}

	if len(data) == 0 {
		ctx.JSON(http.StatusOK, NewResponse(200, "成功", nil, tid))
		return
	}

	ctx.JSON(http.StatusOK, NewResponse(0, "成功", data[0], tid))
}
