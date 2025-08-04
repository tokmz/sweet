package basic

import (
	"mime/multipart"
	"time"
)

// UploadFileReq 上传文件请求
type UploadFileReq struct {
	File        *multipart.FileHeader `form:"file" binding:"required" json:"-"` // 上传的文件
	StorageType *int64                `form:"storage_type" json:"storage_type"` // 存储类型：1=本地存储，2=阿里云OSS，3=腾讯云COS，4=七牛云
	Category    string                `form:"category" json:"category"`         // 文件分类（可选）
	Description string                `form:"description" json:"description"`   // 文件描述（可选）
}

// UploadFileRes 上传文件响应
type UploadFileRes struct {
	ID           int64  `json:"id"`            // 文件ID
	Name         string `json:"name"`          // 文件名称
	OriginalName string `json:"original_name"` // 原始文件名
	FileURL      string `json:"file_url"`      // 文件访问URL
	FileSize     int64  `json:"file_size"`     // 文件大小（字节）
	FileType     string `json:"file_type"`     // 文件类型/MIME类型
	FileExt      string `json:"file_ext"`      // 文件扩展名
	Md5          string `json:"md5"`           // 文件MD5值
	StorageType  int64  `json:"storage_type"`  // 存储类型
}

// DeleteFileReq 删除文件请求
type DeleteFileReq struct {
	ID int64 `uri:"id" binding:"required" json:"id"` // 文件ID
}

// BatchDeleteFileReq 批量删除文件请求
type BatchDeleteFileReq struct {
	IDs []int64 `json:"ids" binding:"required,min=1" validate:"required,min=1"` // 文件ID列表
}

// GetFileReq 获取文件详情请求
type GetFileReq struct {
	ID int64 `uri:"id" binding:"required" json:"id"` // 文件ID
}

// ListFileReq 文件列表查询请求
type ListFileReq struct {
	Page        int    `form:"page" json:"page"`                                                      // 页码
	PageSize    int    `form:"page_size" json:"page_size"`                                            // 每页数量
	Name        string `form:"name" json:"name"`                                                      // 文件名称（模糊查询）
	FileType    string `form:"file_type" json:"file_type"`                                            // 文件类型
	FileExt     string `form:"file_ext" json:"file_ext"`                                              // 文件扩展名
	StorageType *int64 `form:"storage_type" json:"storage_type"`                                      // 存储类型
	UserID      *int64 `form:"user_id" json:"user_id"`                                                // 上传用户ID
	Status      *int64 `form:"status" json:"status"`                                                  // 状态
	StartTime   string `form:"start_time" json:"start_time" validate:"omitempty,datetime=2006-01-02"` // 开始时间
	EndTime     string `form:"end_time" json:"end_time" validate:"omitempty,datetime=2006-01-02"`     // 结束时间
	OrderBy     string `form:"order_by" json:"order_by"`                                              // 排序字段
	OrderType   string `form:"order_type" json:"order_type"`                                          // 排序方式：asc/desc
}

// ListFileItem 文件列表项
type ListFileItem struct {
	ID           int64      `json:"id"`             // 文件ID
	Name         string     `json:"name"`           // 文件名称
	OriginalName string     `json:"original_name"`  // 原始文件名
	FileURL      *string    `json:"file_url"`       // 文件访问URL
	FileSize     int64      `json:"file_size"`      // 文件大小（字节）
	FileType     string     `json:"file_type"`      // 文件类型/MIME类型
	FileExt      string     `json:"file_ext"`       // 文件扩展名
	StorageType  *int64     `json:"storage_type"`   // 存储类型
	UploadUserID *int64     `json:"upload_user_id"` // 上传用户ID
	Status       *int64     `json:"status"`         // 状态
	CreatedAt    *time.Time `json:"created_at"`     // 创建时间
	UpdatedAt    *time.Time `json:"updated_at"`     // 更新时间
}

// ListFileRes 文件列表响应
type ListFileRes struct {
	List  []ListFileItem `json:"list"`  // 文件列表
	Total int64          `json:"total"` // 总数
}

// FileDetailRes 文件详情响应
type FileDetailRes struct {
	ID           int64      `json:"id"`             // 文件ID
	Name         string     `json:"name"`           // 文件名称
	OriginalName string     `json:"original_name"`  // 原始文件名
	FilePath     string     `json:"file_path"`      // 文件路径
	FileURL      *string    `json:"file_url"`       // 文件访问URL
	FileSize     int64      `json:"file_size"`      // 文件大小（字节）
	FileType     string     `json:"file_type"`      // 文件类型/MIME类型
	FileExt      string     `json:"file_ext"`       // 文件扩展名
	Md5          string     `json:"md5"`            // 文件MD5值
	StorageType  *int64     `json:"storage_type"`   // 存储类型
	UploadUserID *int64     `json:"upload_user_id"` // 上传用户ID
	Status       *int64     `json:"status"`         // 状态
	CreatedAt    *time.Time `json:"created_at"`     // 创建时间
	UpdatedAt    *time.Time `json:"updated_at"`     // 更新时间
}

// DownloadFileReq 下载文件请求
type DownloadFileReq struct {
	ID int64 `uri:"id" binding:"required" json:"id"` // 文件ID
}

// UpdateFileReq 更新文件信息请求
type UpdateFileReq struct {
	ID          int64  `uri:"id" binding:"required" json:"id"`           // 文件ID
	Name        string `json:"name" validate:"omitempty,max=64"`         // 文件名称
	Description string `json:"description" validate:"omitempty,max=255"` // 文件描述
	Status      *int64 `json:"status" validate:"omitempty,oneof=1 2"`    // 状态：1=正常，2=禁用
}

// FileStatisticsRes 文件统计响应
type FileStatisticsRes struct {
	TotalFiles    int64 `json:"total_files"`    // 总文件数
	TotalSize     int64 `json:"total_size"`     // 总文件大小（字节）
	LocalFiles    int64 `json:"local_files"`    // 本地存储文件数
	CloudFiles    int64 `json:"cloud_files"`    // 云存储文件数
	ActiveFiles   int64 `json:"active_files"`   // 正常状态文件数
	DisabledFiles int64 `json:"disabled_files"` // 禁用状态文件数
}
