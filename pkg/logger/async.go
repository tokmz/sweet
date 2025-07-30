package logger

import (
	"sync"
	"time"

	"go.uber.org/zap/zapcore"
)

// asyncCore 异步日志核心
type asyncCore struct {
	zapcore.Core
	buffer        chan zapcore.Entry
	bufferSize    int
	flushInterval time.Duration
	wg            sync.WaitGroup
	mu            sync.RWMutex
	closed        bool
	stopCh        chan struct{}
}

// newAsyncCore 创建异步核心
func newAsyncCore(core zapcore.Core, bufferSize int, flushInterval time.Duration) *asyncCore {
	async := &asyncCore{
		Core:          core,
		buffer:        make(chan zapcore.Entry, bufferSize),
		bufferSize:    bufferSize,
		flushInterval: flushInterval,
		stopCh:        make(chan struct{}),
	}

	// 启动后台写入goroutine
	async.wg.Add(1)
	go async.writeLoop()

	return async
}

// writeLoop 写入循环
func (a *asyncCore) writeLoop() {
	defer a.wg.Done()

	ticker := time.NewTicker(a.flushInterval)
	defer ticker.Stop()

	var entries []zapcore.Entry
	batchSize := 100 // 批量写入大小

	for {
		select {
		case entry := <-a.buffer:
			entries = append(entries, entry)
			// 如果达到批量大小，立即写入
			if len(entries) >= batchSize {
				a.flushEntries(entries)
				entries = entries[:0] // 重置切片
			}

		case <-ticker.C:
			// 定时刷新
			if len(entries) > 0 {
				a.flushEntries(entries)
				entries = entries[:0]
			}

		case <-a.stopCh:
			// 关闭时刷新剩余条目
			if len(entries) > 0 {
				a.flushEntries(entries)
			}
			// 处理缓冲区中剩余的条目
			for {
				select {
				case entry := <-a.buffer:
					a.flushEntry(entry)
				default:
					return
				}
			}
		}
	}
}

// flushEntries 批量刷新条目
func (a *asyncCore) flushEntries(entries []zapcore.Entry) {
	for _, entry := range entries {
		a.flushEntry(entry)
	}
}

// flushEntry 刷新单个条目
func (a *asyncCore) flushEntry(entry zapcore.Entry) {
	// 这里需要重新构造CheckedEntry来写入
	if ce := a.Core.Check(entry, nil); ce != nil {
		ce.Write()
	}
}

// Write 写入日志条目（异步）
func (a *asyncCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.closed {
		// 如果已关闭，直接同步写入
		return a.Core.Write(entry, fields)
	}

	// 创建条目副本（包含字段）
	entryCopy := entry
	entryCopy.Stack = entry.Stack // 复制堆栈信息

	// 尝试异步写入
	select {
	case a.buffer <- entryCopy:
		return nil
	default:
		// 缓冲区满时，直接同步写入
		return a.Core.Write(entry, fields)
	}
}

// Sync 同步刷新
func (a *asyncCore) Sync() error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.closed {
		return a.Core.Sync()
	}

	// 等待缓冲区清空
	for len(a.buffer) > 0 {
		time.Sleep(10 * time.Millisecond)
	}

	return a.Core.Sync()
}

// Close 关闭异步核心
func (a *asyncCore) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return nil
	}

	a.closed = true
	close(a.stopCh)
	a.wg.Wait()

	return a.Core.Sync()
}

// Check 检查日志级别
func (a *asyncCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return a.Core.Check(entry, ce)
}

// With 添加字段
func (a *asyncCore) With(fields []zapcore.Field) zapcore.Core {
	return &asyncCore{
		Core:          a.Core.With(fields),
		buffer:        a.buffer,
		bufferSize:    a.bufferSize,
		flushInterval: a.flushInterval,
		stopCh:        a.stopCh,
	}
}

// Enabled 检查是否启用
func (a *asyncCore) Enabled(level zapcore.Level) bool {
	return a.Core.Enabled(level)
}