package basic

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"sweet/internal/global"
	basicDto "sweet/internal/models/dto/basic"
	"sweet/internal/models/entity"
)

// FileService 文件服务实现
type FileService struct{}

// NewFileService 创建文件服务实例
func NewFileService() IFileService {
	return &FileService{}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(ctx context.Context, req *basicDto.UploadFileReq) (*basicDto.UploadFileRes, error) {
	if req.File == nil {
		return nil, errors.New("文件不能为空")
	}

	// 打开上传的文件
	src, err := req.File.Open()
	if err != nil {
		global.Logger.Error("打开上传文件失败", zap.Error(err))
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 计算文件MD5
	md5Hash, err := s.calculateFileMD5(src)
	if err != nil {
		global.Logger.Error("计算文件MD5失败", zap.Error(err))
		return nil, fmt.Errorf("计算文件MD5失败: %v", err)
	}

	// 检查文件是否已存在（去重）
	existingFile, err := s.GetFileByMD5(ctx, md5Hash)
	if err != nil {
		global.Logger.Error("检查文件MD5失败", zap.Error(err))
		return nil, fmt.Errorf("检查文件是否存在失败: %v", err)
	}
	if existingFile != nil {
		// 文件已存在，返回已存在的文件信息
		return &basicDto.UploadFileRes{
			ID:           existingFile.ID,
			Name:         existingFile.Name,
			OriginalName: existingFile.OriginalName,
			FileURL:      *existingFile.FileURL,
			FileSize:     existingFile.FileSize,
			FileType:     existingFile.FileType,
			FileExt:      existingFile.FileExt,
			Md5:          existingFile.Md5,
			StorageType:  *existingFile.StorageType,
		}, nil
	}

	// 重新定位到文件开头
	src.Seek(0, io.SeekStart)

	// 生成文件名和路径
	originalName := req.File.Filename
	fileExt := filepath.Ext(originalName)
	fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), md5Hash[:8], fileExt)

	// 确定存储类型
	storageType := int64(1) // 默认本地存储
	if req.StorageType != nil {
		storageType = *req.StorageType
	}

	// 根据存储类型保存文件
	var filePath, fileURL string
	switch storageType {
	case 1: // 本地存储
		filePath, fileURL, err = s.saveToLocal(src, fileName)
	case 2: // 阿里云OSS
		filePath, fileURL, err = s.saveToAliOSS(src, fileName)
	case 3: // 腾讯云COS
		filePath, fileURL, err = s.saveToTencentCOS(src, fileName)
	case 4: // 七牛云
		filePath, fileURL, err = s.saveToQiniu(src, fileName)
	default:
		return nil, errors.New("不支持的存储类型")
	}

	if err != nil {
		global.Logger.Error("保存文件失败", zap.Int64("storage_type", storageType), zap.Error(err))
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 获取用户ID（从上下文中获取）
	var uploadUserID *int64
	if userID := ctx.Value("user_id"); userID != nil {
		if uid, ok := userID.(int64); ok {
			uploadUserID = &uid
		}
	}

	// 保存文件信息到数据库
	fileEntity := &entity.SysFile{
		Name:         fileName,
		OriginalName: originalName,
		FilePath:     filePath,
		FileURL:      &fileURL,
		FileSize:     req.File.Size,
		FileType:     req.File.Header.Get("Content-Type"),
		FileExt:      strings.TrimPrefix(fileExt, "."),
		Md5:          md5Hash,
		StorageType:  &storageType,
		UploadUserID: uploadUserID,
		Status:       func() *int64 { s := int64(1); return &s }(), // 正常状态
	}

	if err := global.Query.SysFile.WithContext(ctx).Create(fileEntity); err != nil {
		global.Logger.Error("保存文件信息到数据库失败", zap.Error(err))
		return nil, fmt.Errorf("保存文件信息失败: %v", err)
	}

	global.Logger.Info("文件上传成功", zap.Int64("file_id", fileEntity.ID), zap.String("file_name", fileName))

	return &basicDto.UploadFileRes{
		ID:           fileEntity.ID,
		Name:         fileEntity.Name,
		OriginalName: fileEntity.OriginalName,
		FileURL:      fileURL,
		FileSize:     fileEntity.FileSize,
		FileType:     fileEntity.FileType,
		FileExt:      fileEntity.FileExt,
		Md5:          fileEntity.Md5,
		StorageType:  storageType,
	}, nil
}

// DownloadFile 下载文件
func (s *FileService) DownloadFile(ctx context.Context, req *basicDto.DownloadFileReq) (filePath string, fileName string, err error) {
	// 查询文件信息
	fileEntity, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(req.ID), global.Query.SysFile.Status.Eq(1)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("文件不存在")
		}
		global.Logger.Error("查询文件信息失败", zap.Int64("file_id", req.ID), zap.Error(err))
		return "", "", fmt.Errorf("查询文件信息失败: %v", err)
	}

	// 检查文件是否存在
	if err := s.ValidateFile(ctx, req.ID); err != nil {
		return "", "", err
	}

	// 根据存储类型返回文件路径
	switch *fileEntity.StorageType {
	case 1: // 本地存储
		return fileEntity.FilePath, fileEntity.OriginalName, nil
	default: // 云存储
		if fileEntity.FileURL != nil {
			return *fileEntity.FileURL, fileEntity.OriginalName, nil
		}
		return fileEntity.FilePath, fileEntity.OriginalName, nil
	}
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(ctx context.Context, req *basicDto.DeleteFileReq) error {
	// 检查文件是否存在
	_, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(req.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文件不存在")
		}
		global.Logger.Error("查询文件信息失败", zap.Int64("file_id", req.ID), zap.Error(err))
		return fmt.Errorf("查询文件信息失败: %v", err)
	}

	// 软删除文件记录
	if _, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(req.ID)).Delete(); err != nil {
		global.Logger.Error("删除文件记录失败", zap.Int64("file_id", req.ID), zap.Error(err))
		return fmt.Errorf("删除文件失败: %v", err)
	}

	global.Logger.Info("文件删除成功", zap.Int64("file_id", req.ID))
	return nil
}

// BatchDeleteFile 批量删除文件
func (s *FileService) BatchDeleteFile(ctx context.Context, req *basicDto.BatchDeleteFileReq) error {
	if len(req.IDs) == 0 {
		return errors.New("文件ID列表不能为空")
	}

	// 批量软删除
	result, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.In(req.IDs...)).Delete()
	if err != nil {
		global.Logger.Error("批量删除文件失败", zap.Any("file_ids", req.IDs), zap.Error(err))
		return fmt.Errorf("批量删除文件失败: %v", err)
	}

	global.Logger.Info("批量删除文件成功", zap.Any("file_ids", req.IDs), zap.Int64("affected_rows", result.RowsAffected))
	return nil
}

// GetFile 获取文件详情
func (s *FileService) GetFile(ctx context.Context, req *basicDto.GetFileReq) (*basicDto.FileDetailRes, error) {
	fileEntity, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(req.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文件不存在")
		}
		global.Logger.Error("查询文件详情失败", zap.Int64("file_id", req.ID), zap.Error(err))
		return nil, fmt.Errorf("查询文件详情失败: %v", err)
	}

	return &basicDto.FileDetailRes{
		ID:           fileEntity.ID,
		Name:         fileEntity.Name,
		OriginalName: fileEntity.OriginalName,
		FilePath:     fileEntity.FilePath,
		FileURL:      fileEntity.FileURL,
		FileSize:     fileEntity.FileSize,
		FileType:     fileEntity.FileType,
		FileExt:      fileEntity.FileExt,
		Md5:          fileEntity.Md5,
		StorageType:  fileEntity.StorageType,
		UploadUserID: fileEntity.UploadUserID,
		Status:       fileEntity.Status,
		CreatedAt:    fileEntity.CreatedAt,
		UpdatedAt:    fileEntity.UpdatedAt,
	}, nil
}

// ListFile 获取文件列表
func (s *FileService) ListFile(ctx context.Context, req *basicDto.ListFileReq) (*basicDto.ListFileRes, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 构建查询条件
	query := global.Query.SysFile.WithContext(ctx)

	// 文件名模糊查询
	if req.Name != "" {
		query = query.Where(global.Query.SysFile.Name.Like("%" + req.Name + "%")).Or(global.Query.SysFile.OriginalName.Like("%" + req.Name + "%"))
	}

	// 文件类型筛选
	if req.FileType != "" {
		query = query.Where(global.Query.SysFile.FileType.Eq(req.FileType))
	}

	// 文件扩展名筛选
	if req.FileExt != "" {
		query = query.Where(global.Query.SysFile.FileExt.Eq(req.FileExt))
	}

	// 存储类型筛选
	if req.StorageType != nil {
		query = query.Where(global.Query.SysFile.StorageType.Eq(*req.StorageType))
	}

	// 上传用户筛选
	if req.UserID != nil {
		query = query.Where(global.Query.SysFile.UploadUserID.Eq(*req.UserID))
	}

	// 状态筛选
	if req.Status != nil {
		query = query.Where(global.Query.SysFile.Status.Eq(*req.Status))
	}

	// 时间范围筛选
	if req.StartTime != "" {
		startTime, parseErr := time.Parse("2006-01-02", req.StartTime)
		if parseErr == nil {
			query = query.Where(global.Query.SysFile.CreatedAt.Gte(startTime))
		}
	}
	if req.EndTime != "" {
		endTime, parseErr := time.Parse("2006-01-02", req.EndTime)
		if parseErr == nil {
			query = query.Where(global.Query.SysFile.CreatedAt.Lte(endTime.Add(24 * time.Hour)))
		}
	}

	// 统计总数
	total, err := query.Count()
	if err != nil {
		global.Logger.Error("统计文件总数失败", zap.Error(err))
		return nil, fmt.Errorf("统计文件总数失败: %v", err)
	}

	// 排序
	orderBy := "created_at"
	orderType := "desc"
	if req.OrderBy != "" {
		orderBy = req.OrderBy
	}
	if req.OrderType != "" && (req.OrderType == "asc" || req.OrderType == "desc") {
		orderType = req.OrderType
	}
	if orderType == "desc" {
		switch orderBy {
		case "created_at":
			query = query.Order(global.Query.SysFile.CreatedAt.Desc())
		case "updated_at":
			query = query.Order(global.Query.SysFile.UpdatedAt.Desc())
		case "file_size":
			query = query.Order(global.Query.SysFile.FileSize.Desc())
		default:
			query = query.Order(global.Query.SysFile.CreatedAt.Desc())
		}
	} else {
		switch orderBy {
		case "created_at":
			query = query.Order(global.Query.SysFile.CreatedAt)
		case "updated_at":
			query = query.Order(global.Query.SysFile.UpdatedAt)
		case "file_size":
			query = query.Order(global.Query.SysFile.FileSize)
		default:
			query = query.Order(global.Query.SysFile.CreatedAt)
		}
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	files, err := query.Offset(offset).Limit(req.PageSize).Find()
	if err != nil {
		global.Logger.Error("查询文件列表失败", zap.Error(err))
		return nil, fmt.Errorf("查询文件列表失败: %v", err)
	}

	// 转换为响应格式
	list := make([]basicDto.ListFileItem, len(files))
	for i, file := range files {
		list[i] = basicDto.ListFileItem{
			ID:           file.ID,
			Name:         file.Name,
			OriginalName: file.OriginalName,
			FileURL:      file.FileURL,
			FileSize:     file.FileSize,
			FileType:     file.FileType,
			FileExt:      file.FileExt,
			StorageType:  file.StorageType,
			UploadUserID: file.UploadUserID,
			Status:       file.Status,
			CreatedAt:    file.CreatedAt,
			UpdatedAt:    file.UpdatedAt,
		}
	}

	return &basicDto.ListFileRes{
		List:  list,
		Total: total,
	}, nil
}

// UpdateFile 更新文件信息
func (s *FileService) UpdateFile(ctx context.Context, req *basicDto.UpdateFileReq) error {
	// 构建更新数据
	updateData := make(map[string]interface{})
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Description != "" {
		// 注意：这里假设数据库表有description字段，如果没有需要添加
		updateData["description"] = req.Description
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}

	if len(updateData) == 0 {
		return errors.New("没有需要更新的字段")
	}

	// 执行更新
	result, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(req.ID)).Updates(updateData)
	if err != nil {
		global.Logger.Error("更新文件信息失败", zap.Int64("file_id", req.ID), zap.Error(err))
		return fmt.Errorf("更新文件信息失败: %v", err)
	}

	if result.RowsAffected == 0 {
		return errors.New("文件不存在")
	}

	global.Logger.Info("更新文件信息成功", zap.Int64("file_id", req.ID))
	return nil
}

// GetFileStatistics 获取文件统计信息
func (s *FileService) GetFileStatistics(ctx context.Context) (*basicDto.FileStatisticsRes, error) {
	var stats basicDto.FileStatisticsRes
	var err error

	// 总文件数
	stats.TotalFiles, err = global.Query.SysFile.WithContext(ctx).Count()
	if err != nil {
		global.Logger.Error("统计总文件数失败", zap.Error(err))
		return nil, fmt.Errorf("统计总文件数失败: %v", err)
	}

	// 总文件大小
	var totalSize sql.NullInt64
	if err := global.DBClient.WithContext(ctx).Model(&entity.SysFile{}).Select("SUM(file_size)").Scan(&totalSize).Error; err != nil {
		global.Logger.Error("统计总文件大小失败", zap.Error(err))
		return nil, fmt.Errorf("统计总文件大小失败: %v", err)
	}
	if totalSize.Valid {
		stats.TotalSize = totalSize.Int64
	}

	// 本地存储文件数
	stats.LocalFiles, err = global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.StorageType.Eq(1)).Count()
	if err != nil {
		global.Logger.Error("统计本地文件数失败", zap.Error(err))
		return nil, fmt.Errorf("统计本地文件数失败: %v", err)
	}

	// 云存储文件数
	stats.CloudFiles, err = global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.StorageType.Gt(1)).Count()
	if err != nil {
		global.Logger.Error("统计云存储文件数失败", zap.Error(err))
		return nil, fmt.Errorf("统计云存储文件数失败: %v", err)
	}

	// 正常状态文件数
	stats.ActiveFiles, err = global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.Status.Eq(1)).Count()
	if err != nil {
		global.Logger.Error("统计正常文件数失败", zap.Error(err))
		return nil, fmt.Errorf("统计正常文件数失败: %v", err)
	}

	// 禁用状态文件数
	stats.DisabledFiles, err = global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.Status.Eq(2)).Count()
	if err != nil {
		global.Logger.Error("统计禁用文件数失败", zap.Error(err))
		return nil, fmt.Errorf("统计禁用文件数失败: %v", err)
	}

	return &stats, nil
}

// ValidateFile 验证文件
func (s *FileService) ValidateFile(ctx context.Context, fileID int64) error {
	// 查询文件信息
	fileEntity, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(fileID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文件记录不存在")
		}
		return fmt.Errorf("查询文件信息失败: %v", err)
	}

	// 检查文件状态
	if fileEntity.Status != nil && *fileEntity.Status != 1 {
		return errors.New("文件已被禁用")
	}

	// 检查物理文件是否存在
	if _, err := os.Stat(fileEntity.FilePath); os.IsNotExist(err) {
		return errors.New("物理文件不存在")
	}

	return nil
}

// GetFileURL 获取文件访问URL
func (s *FileService) GetFileURL(ctx context.Context, fileID int64, expireSeconds ...int64) (string, error) {
	// 查询文件信息
	fileEntity, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(fileID), global.Query.SysFile.Status.Eq(1)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("文件不存在")
		}
		return "", fmt.Errorf("查询文件信息失败: %v", err)
	}

	// 如果已有URL直接返回
	if fileEntity.FileURL != nil && *fileEntity.FileURL != "" {
		return *fileEntity.FileURL, nil
	}

	// 根据存储类型生成URL
	var url string
	var urlErr error
	switch *fileEntity.StorageType {
	case 1: // 本地存储
		url = fmt.Sprintf("/api/v1/files/%d/download", fileID)
	case 2: // 阿里云OSS
		url, urlErr = s.generateAliOSSURL(fileEntity.FilePath, expireSeconds...)
	case 3: // 腾讯云COS
		url, urlErr = s.generateTencentCOSURL(fileEntity.FilePath, expireSeconds...)
	case 4: // 七牛云
		url, urlErr = s.generateQiniuURL(fileEntity.FilePath, expireSeconds...)
	default:
		return "", errors.New("不支持的存储类型")
	}

	if urlErr != nil {
		return "", fmt.Errorf("生成文件URL失败: %v", urlErr)
	}

	return url, nil
}

// CopyFile 复制文件
func (s *FileService) CopyFile(ctx context.Context, sourceFileID int64, targetName string) (*basicDto.UploadFileRes, error) {
	// 查询源文件信息
	sourceFile, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(sourceFileID), global.Query.SysFile.Status.Eq(1)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("源文件不存在")
		}
		return nil, fmt.Errorf("查询源文件信息失败: %v", err)
	}

	// 验证源文件
	if validateErr := s.ValidateFile(ctx, sourceFileID); validateErr != nil {
		return nil, fmt.Errorf("源文件验证失败: %v", validateErr)
	}

	// 生成新文件名
	if targetName == "" {
		targetName = fmt.Sprintf("copy_%d_%s", time.Now().Unix(), sourceFile.Name)
	}

	// 复制物理文件
	newFilePath, newFileURL, copyErr := s.copyPhysicalFile(sourceFile.FilePath, targetName, *sourceFile.StorageType)
	if copyErr != nil {
		return nil, fmt.Errorf("复制物理文件失败: %v", copyErr)
	}

	// 获取用户ID
	var uploadUserID *int64
	if userID := ctx.Value("user_id"); userID != nil {
		if uid, ok := userID.(int64); ok {
			uploadUserID = &uid
		}
	}

	// 创建新文件记录
	newFile := &entity.SysFile{
		Name:         targetName,
		OriginalName: sourceFile.OriginalName,
		FilePath:     newFilePath,
		FileURL:      &newFileURL,
		FileSize:     sourceFile.FileSize,
		FileType:     sourceFile.FileType,
		FileExt:      sourceFile.FileExt,
		Md5:          sourceFile.Md5,
		StorageType:  sourceFile.StorageType,
		UploadUserID: uploadUserID,
		Status:       func() *int64 { s := int64(1); return &s }(),
	}

	if createErr := global.Query.SysFile.WithContext(ctx).Create(newFile); createErr != nil {
		return nil, fmt.Errorf("保存复制文件信息失败: %v", createErr)
	}

	return &basicDto.UploadFileRes{
		ID:           newFile.ID,
		Name:         newFile.Name,
		OriginalName: newFile.OriginalName,
		FileURL:      newFileURL,
		FileSize:     newFile.FileSize,
		FileType:     newFile.FileType,
		FileExt:      newFile.FileExt,
		Md5:          newFile.Md5,
		StorageType:  *newFile.StorageType,
	}, nil
}

// MoveFile 移动文件到不同存储类型
func (s *FileService) MoveFile(ctx context.Context, fileID int64, targetStorageType int64) error {
	// 查询文件信息
	fileEntity, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(fileID), global.Query.SysFile.Status.Eq(1)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文件不存在")
		}
		return fmt.Errorf("查询文件信息失败: %v", err)
	}

	// 检查是否需要移动
	if *fileEntity.StorageType == targetStorageType {
		return errors.New("文件已在目标存储类型中")
	}

	// 验证文件
	if validateErr := s.ValidateFile(ctx, fileID); validateErr != nil {
		return fmt.Errorf("文件验证失败: %v", validateErr)
	}

	// 移动文件到新存储
	newFilePath, newFileURL, moveErr := s.moveFileToStorage(fileEntity.FilePath, fileEntity.Name, *fileEntity.StorageType, targetStorageType)
	if moveErr != nil {
		return fmt.Errorf("移动文件失败: %v", moveErr)
	}

	// 更新数据库记录
	updateData := map[string]interface{}{
		"file_path":    newFilePath,
		"file_url":     newFileURL,
		"storage_type": targetStorageType,
	}

	if _, updateErr := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(fileID)).Updates(updateData); updateErr != nil {
		return fmt.Errorf("更新文件信息失败: %v", updateErr)
	}

	global.Logger.Info("文件移动成功", zap.Int64("file_id", fileID), zap.Int64("from_storage", *fileEntity.StorageType), zap.Int64("to_storage", targetStorageType))
	return nil
}

// CleanupExpiredFiles 清理过期文件
func (s *FileService) CleanupExpiredFiles(ctx context.Context, expireDays int) (int64, error) {
	if expireDays <= 0 {
		return 0, errors.New("过期天数必须大于0")
	}

	// 计算过期时间
	expireTime := time.Now().AddDate(0, 0, -expireDays)

	// 查询过期的软删除文件
	expiredFiles, err := global.Query.SysFile.WithContext(ctx).Unscoped().Where(global.Query.SysFile.DeletedAt.IsNotNull()).Find()
	if err != nil {
		return 0, fmt.Errorf("查询过期文件失败: %v", err)
	}

	// 过滤出真正过期的文件
	var reallyExpiredFiles []*entity.SysFile
	for _, file := range expiredFiles {
		if file.DeletedAt.Valid && file.DeletedAt.Time.Before(expireTime) {
			reallyExpiredFiles = append(reallyExpiredFiles, file)
		}
	}

	var cleanedCount int64
	for _, file := range reallyExpiredFiles {
		// 删除物理文件
		if err := s.deletePhysicalFile(file.FilePath, *file.StorageType); err != nil {
			global.Logger.Error("删除物理文件失败", zap.Int64("file_id", file.ID), zap.String("file_path", file.FilePath), zap.Error(err))
			continue
		}

		// 物理删除数据库记录
		if _, err := global.Query.SysFile.WithContext(ctx).Unscoped().Where(global.Query.SysFile.ID.Eq(file.ID)).Delete(); err != nil {
			global.Logger.Error("物理删除文件记录失败", zap.Int64("file_id", file.ID), zap.Error(err))
			continue
		}

		cleanedCount++
	}

	global.Logger.Info("清理过期文件完成", zap.Int("expire_days", expireDays), zap.Int64("cleaned_count", cleanedCount))
	return cleanedCount, nil
}

// GetFileByMD5 根据MD5获取文件
func (s *FileService) GetFileByMD5(ctx context.Context, md5 string) (*basicDto.FileDetailRes, error) {
	file, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.Md5.Eq(md5), global.Query.SysFile.Status.Eq(1)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 文件不存在，返回nil而不是错误
		}
		return nil, fmt.Errorf("查询文件失败: %v", err)
	}

	return &basicDto.FileDetailRes{
		ID:           file.ID,
		Name:         file.Name,
		OriginalName: file.OriginalName,
		FilePath:     file.FilePath,
		FileURL:      file.FileURL,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		FileExt:      file.FileExt,
		Md5:          file.Md5,
		StorageType:  file.StorageType,
		UploadUserID: file.UploadUserID,
		Status:       file.Status,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
	}, nil
}

// UpdateFileStatus 更新文件状态
func (s *FileService) UpdateFileStatus(ctx context.Context, fileID int64, status int64) error {
	result, err := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.ID.Eq(fileID)).Update(global.Query.SysFile.Status, status)
	if err != nil {
		return fmt.Errorf("更新文件状态失败: %v", err)
	}

	if result.RowsAffected == 0 {
		return errors.New("文件不存在")
	}

	global.Logger.Info("更新文件状态成功", zap.Int64("file_id", fileID), zap.Int64("status", status))
	return nil
}

// GetUserFiles 获取用户上传的文件列表
func (s *FileService) GetUserFiles(ctx context.Context, userID int64, req *basicDto.ListFileReq) (*basicDto.ListFileRes, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 构建查询条件
	query := global.Query.SysFile.WithContext(ctx).Where(global.Query.SysFile.UploadUserID.Eq(userID))

	// 文件名模糊查询
	if req.Name != "" {
		query = query.Where(global.Query.SysFile.Name.Like("%" + req.Name + "%")).Or(global.Query.SysFile.OriginalName.Like("%" + req.Name + "%"))
	}

	// 文件类型筛选
	if req.FileType != "" {
		query = query.Where(global.Query.SysFile.FileType.Eq(req.FileType))
	}

	// 文件扩展名筛选
	if req.FileExt != "" {
		query = query.Where(global.Query.SysFile.FileExt.Eq(req.FileExt))
	}

	// 状态筛选
	if req.Status != nil {
		query = query.Where(global.Query.SysFile.Status.Eq(*req.Status))
	}

	// 时间范围筛选
	if req.StartTime != "" {
		startTime, parseStartErr := time.Parse("2006-01-02", req.StartTime)
		if parseStartErr == nil {
			query = query.Where(global.Query.SysFile.CreatedAt.Gte(startTime))
		}
	}
	if req.EndTime != "" {
		endTime, parseEndErr := time.Parse("2006-01-02", req.EndTime)
		if parseEndErr == nil {
			query = query.Where(global.Query.SysFile.CreatedAt.Lte(endTime.Add(24 * time.Hour)))
		}
	}

	// 统计总数
	total, err := query.Count()
	if err != nil {
		global.Logger.Error("统计用户文件总数失败", zap.Int64("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("统计用户文件总数失败: %v", err)
	}

	// 排序
	orderBy := "created_at"
	orderType := "desc"
	if req.OrderBy != "" {
		orderBy = req.OrderBy
	}
	if req.OrderType != "" && (req.OrderType == "asc" || req.OrderType == "desc") {
		orderType = req.OrderType
	}
	if orderType == "desc" {
		switch orderBy {
		case "created_at":
			query = query.Order(global.Query.SysFile.CreatedAt.Desc())
		case "updated_at":
			query = query.Order(global.Query.SysFile.UpdatedAt.Desc())
		case "file_size":
			query = query.Order(global.Query.SysFile.FileSize.Desc())
		default:
			query = query.Order(global.Query.SysFile.CreatedAt.Desc())
		}
	} else {
		switch orderBy {
		case "created_at":
			query = query.Order(global.Query.SysFile.CreatedAt)
		case "updated_at":
			query = query.Order(global.Query.SysFile.UpdatedAt)
		case "file_size":
			query = query.Order(global.Query.SysFile.FileSize)
		default:
			query = query.Order(global.Query.SysFile.CreatedAt)
		}
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	files, err := query.Offset(offset).Limit(req.PageSize).Find()
	if err != nil {
		global.Logger.Error("查询用户文件列表失败", zap.Int64("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("查询用户文件列表失败: %v", err)
	}

	// 转换为响应格式
	list := make([]basicDto.ListFileItem, len(files))
	for i, file := range files {
		list[i] = basicDto.ListFileItem{
			ID:           file.ID,
			Name:         file.Name,
			OriginalName: file.OriginalName,
			FileURL:      file.FileURL,
			FileSize:     file.FileSize,
			FileType:     file.FileType,
			FileExt:      file.FileExt,
			StorageType:  file.StorageType,
			UploadUserID: file.UploadUserID,
			Status:       file.Status,
			CreatedAt:    file.CreatedAt,
			UpdatedAt:    file.UpdatedAt,
		}
	}

	return &basicDto.ListFileRes{
		List:  list,
		Total: total,
	}, nil
}

// 私有方法：计算文件MD5
func (s *FileService) calculateFileMD5(file multipart.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// 私有方法：保存到本地存储
func (s *FileService) saveToLocal(file multipart.File, fileName string) (filePath, fileURL string, err error) {
	// 创建上传目录
	uploadDir := "uploads/" + time.Now().Format("2006/01/02")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", "", fmt.Errorf("创建上传目录失败: %v", err)
	}

	// 完整文件路径
	filePath = filepath.Join(uploadDir, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		return "", "", fmt.Errorf("保存文件失败: %v", err)
	}

	// 生成访问URL
	fileURL = "/uploads/" + time.Now().Format("2006/01/02") + "/" + fileName

	return filePath, fileURL, nil
}

// 私有方法：保存到阿里云OSS（占位实现）
func (s *FileService) saveToAliOSS(file multipart.File, fileName string) (filePath, fileURL string, err error) {
	// TODO: 实现阿里云OSS上传逻辑
	return "", "", errors.New("阿里云OSS上传功能暂未实现")
}

// 私有方法：保存到腾讯云COS（占位实现）
func (s *FileService) saveToTencentCOS(file multipart.File, fileName string) (filePath, fileURL string, err error) {
	// TODO: 实现腾讯云COS上传逻辑
	return "", "", errors.New("腾讯云COS上传功能暂未实现")
}

// 私有方法：保存到七牛云（占位实现）
func (s *FileService) saveToQiniu(file multipart.File, fileName string) (filePath, fileURL string, err error) {
	// TODO: 实现七牛云上传逻辑
	return "", "", errors.New("七牛云上传功能暂未实现")
}

// 私有方法：生成阿里云OSS URL（占位实现）
func (s *FileService) generateAliOSSURL(filePath string, expireSeconds ...int64) (string, error) {
	// TODO: 实现阿里云OSS URL生成逻辑
	return "", errors.New("阿里云OSS URL生成功能暂未实现")
}

// 私有方法：生成腾讯云COS URL（占位实现）
func (s *FileService) generateTencentCOSURL(filePath string, expireSeconds ...int64) (string, error) {
	// TODO: 实现腾讯云COS URL生成逻辑
	return "", errors.New("腾讯云COS URL生成功能暂未实现")
}

// 私有方法：生成七牛云URL（占位实现）
func (s *FileService) generateQiniuURL(filePath string, expireSeconds ...int64) (string, error) {
	// TODO: 实现七牛云URL生成逻辑
	return "", errors.New("七牛云URL生成功能暂未实现")
}

// 私有方法：复制物理文件（占位实现）
func (s *FileService) copyPhysicalFile(sourcePath, targetName string, storageType int64) (newFilePath, newFileURL string, err error) {
	// TODO: 根据存储类型实现文件复制逻辑
	return "", "", errors.New("文件复制功能暂未实现")
}

// 私有方法：移动文件到新存储（占位实现）
func (s *FileService) moveFileToStorage(sourcePath, fileName string, fromStorage, toStorage int64) (newFilePath, newFileURL string, err error) {
	// TODO: 实现文件在不同存储间移动的逻辑
	return "", "", errors.New("文件移动功能暂未实现")
}

// 私有方法：删除物理文件（占位实现）
func (s *FileService) deletePhysicalFile(filePath string, storageType int64) error {
	switch storageType {
	case 1: // 本地存储
		return os.Remove(filePath)
	case 2, 3, 4: // 云存储
		// TODO: 实现云存储文件删除逻辑
		return errors.New("云存储文件删除功能暂未实现")
	default:
		return errors.New("不支持的存储类型")
	}
}
