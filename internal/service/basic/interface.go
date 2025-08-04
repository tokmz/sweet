package basic

import (
	"context"
	"sweet/internal/models"
	basicDto "sweet/internal/models/dto/basic"
)

// IBasicService 基础服务接口
type IBasicService interface {
	// LoginLog 获取登录日志服务
	LoginLog() ILoginLogService
	// OperationLog 获取操作日志服务
	OperationLog() IOperationLogService
	// File 获取文件服务
	File() IFileService
}

// ILoginLogService 登录日志服务接口
type ILoginLogService interface {
	// CreateLoginLog 创建登录日志
	CreateLoginLog(ctx context.Context, req *basicDto.CreateLoginLogReq) error
	// DeleteLoginLog 删除登录日志
	DeleteLoginLog(ctx context.Context, req *basicDto.DeleteLoginLogReq) error
	// ClearAllLoginLog 清空所有登录日志
	ClearAllLoginLog(ctx context.Context) error
	// ClearLoginLog 清空自己的登录日志
	ClearLoginLog(ctx context.Context, uid int64) error
	// ListLoginLog 获取登录日志列表
	ListLoginLog(ctx context.Context, req *basicDto.ListLoginLogReq) (*basicDto.ListLoginLogRes, error)
	// GetLoginLog 获取登录日志详情
	GetLoginLog(ctx context.Context, req *models.IDReq) (*basicDto.LoginLogDetailRes, error)
}

// IOperationLogService 操作日志服务接口
type IOperationLogService interface {
	// CreateOperationLog 创建操作日志
	CreateOperationLog(ctx context.Context, req *basicDto.CreateOperationLogReq) error
	// DeleteOperationLog 删除操作日志
	DeleteOperationLog(ctx context.Context, req *basicDto.DeleteOperationLogReq) error
	// ClearAllOperationLog 清空所有操作日志
	ClearAllOperationLog(ctx context.Context) error
	// ClearOperationLog 清空自己的操作日志
	ClearOperationLog(ctx context.Context, uid int64) error
	// ListOperationLog 获取操作日志列表
	ListOperationLog(ctx context.Context, req *basicDto.ListOperationLogReq) (*basicDto.ListOperationLogRes, error)
	// GetOperationLog 获取操作日志详情
	GetOperationLog(ctx context.Context, req *models.IDReq) (*basicDto.OperationLogDetailRes, error)
}

// IFileService 文件服务接口
type IFileService interface {
	// UploadFile 上传文件
	UploadFile(ctx context.Context, req *basicDto.UploadFileReq) (*basicDto.UploadFileRes, error)

	// DownloadFile 下载文件
	DownloadFile(ctx context.Context, req *basicDto.DownloadFileReq) (filePath string, fileName string, err error)

	// DeleteFile 删除文件
	DeleteFile(ctx context.Context, req *basicDto.DeleteFileReq) error

	// BatchDeleteFile 批量删除文件
	BatchDeleteFile(ctx context.Context, req *basicDto.BatchDeleteFileReq) error

	// GetFile 获取文件详情
	GetFile(ctx context.Context, req *basicDto.GetFileReq) (*basicDto.FileDetailRes, error)

	// ListFile 获取文件列表
	ListFile(ctx context.Context, req *basicDto.ListFileReq) (*basicDto.ListFileRes, error)

	// UpdateFile 更新文件信息
	UpdateFile(ctx context.Context, req *basicDto.UpdateFileReq) error

	// GetFileStatistics 获取文件统计信息
	GetFileStatistics(ctx context.Context) (*basicDto.FileStatisticsRes, error)

	// ValidateFile 验证文件（检查文件是否存在、完整性等）
	ValidateFile(ctx context.Context, fileID int64) error

	// GetFileURL 获取文件访问URL（支持临时URL生成）
	GetFileURL(ctx context.Context, fileID int64, expireSeconds ...int64) (string, error)

	// CopyFile 复制文件
	CopyFile(ctx context.Context, sourceFileID int64, targetName string) (*basicDto.UploadFileRes, error)

	// MoveFile 移动文件到不同存储类型
	MoveFile(ctx context.Context, fileID int64, targetStorageType int64) error

	// CleanupExpiredFiles 清理过期文件（软删除的文件物理删除）
	CleanupExpiredFiles(ctx context.Context, expireDays int) (int64, error)

	// GetFileByMD5 根据MD5获取文件（去重检查）
	GetFileByMD5(ctx context.Context, md5 string) (*basicDto.FileDetailRes, error)

	// UpdateFileStatus 更新文件状态
	UpdateFileStatus(ctx context.Context, fileID int64, status int64) error

	// GetUserFiles 获取用户上传的文件列表
	GetUserFiles(ctx context.Context, userID int64, req *basicDto.ListFileReq) (*basicDto.ListFileRes, error)
}
