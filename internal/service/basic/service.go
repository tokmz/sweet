package basic

import "sync"

// Service 基础服务
type Service struct {
	// loginLog 登录日志服务
	loginLog ILoginLogService
	// operation 操作日志服务
	operation IOperationLogService
	// file 文件服务
	file IFileService
}

var (
	once *sync.Once
)

// NewService 创建基础服务
func NewService() IBasicService {
	s := &Service{}
	once.Do(func() {
		s.loginLog = NewLoginLogService()
		s.operation = NewOperationLogService()
		s.file = NewFileService()
	})
	return s
}

// LoginLog 获取登录日志服务
func (s *Service) LoginLog() ILoginLogService {
	return s.loginLog
}

// OperationLog 获取操作日志服务
func (s *Service) OperationLog() IOperationLogService {
	return s.operation
}

// File 获取文件服务
func (s *Service) File() IFileService {
	return s.file
}
