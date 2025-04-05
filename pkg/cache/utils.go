package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// 检查错误是否为键不存在
func isNilError(err error) bool {
	return err == redis.Nil
}

// 处理命令执行超时
func withTimeout(ctx context.Context, timeout int) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	}
	return ctx, func() {}
}
