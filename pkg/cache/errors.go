package cache

import (
	"errors"
	"fmt"
)

// 定义错误类型
var (
	// ErrKeyNotExists 键不存在错误
	ErrKeyNotExists = errors.New("key not exists")
	// ErrInvalidMode 无效的Redis模式错误
	ErrInvalidMode = errors.New("invalid redis mode")
	// ErrEmptyAddrs 未提供Redis地址错误
	ErrEmptyAddrs = errors.New("empty redis address")
	// ErrEmptyMasterSet 未提供哨兵主节点名称错误
	ErrEmptyMasterSet = errors.New("empty sentinel master name")
	// ErrConnectionFailed 连接失败错误
	ErrConnectionFailed = errors.New("connection failed")
	// ErrCommandFailed 命令执行失败错误
	ErrCommandFailed = errors.New("command execution failed")
	// ErrTimeout 超时错误
	ErrTimeout = errors.New("operation timeout")
	// ErrInvalidArgument 无效参数错误
	ErrInvalidArgument = errors.New("invalid argument")
)

// ConnectionError 连接错误
type ConnectionError struct {
	Addr string
	Err  error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("connection error to %s: %v", e.Addr, e.Err)
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
}

// CommandError 命令执行错误
type CommandError struct {
	Command string
	Err     error
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("command error [%s]: %v", e.Command, e.Err)
}

func (e *CommandError) Unwrap() error {
	return e.Err
}

// TimeoutError 超时错误
type TimeoutError struct {
	Operation string
	Timeout   int
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("operation [%s] timeout after %dms", e.Operation, e.Timeout)
}

// IsKeyNotExistsError 检查是否为键不存在错误
func IsKeyNotExistsError(err error) bool {
	return errors.Is(err, ErrKeyNotExists)
}

// IsConnectionError 检查是否为连接错误
func IsConnectionError(err error) bool {
	var connErr *ConnectionError
	return errors.As(err, &connErr)
}

// IsCommandError 检查是否为命令执行错误
func IsCommandError(err error) bool {
	var cmdErr *CommandError
	return errors.As(err, &cmdErr)
}

// IsTimeoutError 检查是否为超时错误
func IsTimeoutError(err error) bool {
	var timeoutErr *TimeoutError
	return errors.As(err, &timeoutErr) || errors.Is(err, ErrTimeout)
}
