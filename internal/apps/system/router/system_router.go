package router

import "github.com/gin-gonic/gin"

func RegisterSystemV1(group gin.RouterGroup) {
	v1 := group.Group("/system/v1")
	{
		v1.GET("/heart", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{})
		})
	}
}
