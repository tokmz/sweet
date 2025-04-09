package repository

import (
	"context"
	"sweet/pkg/errs"

	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"

	"gorm.io/gorm"
)

var (
	ErrAccountNotFound = errs.New(2000, "账号不存在")
	ErrAccountExists   = errs.New(2001, "账号已存在")
	ErrInvalidPassword = errs.New(2002, "密码错误")
	ErrInvalidUserID   = errs.New(2003, "无效ID")
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error
	// Update 更新用户
	Update(ctx context.Context, user *entity.User) error
	// Delete 删除用户
	Delete(ctx context.Context, ids []int64) error
	// FindOne 查询用户
	FindOne(ctx context.Context, id int64) (*entity.User, error)
	// FindOneByUsername 查询用户通过Username
	FindOneByUsername(ctx context.Context, username string) (*entity.User, error)
	// FindOneByEmail 查询用户通过Email
	FindOneByEmail(ctx context.Context, email string) (*entity.User, error)
	// FindOneByMobile 查询用户通过Mobile
	FindOneByMobile(ctx context.Context, mobile string) (*entity.User, error)
	// FindList 查询用户列表
	FindList(ctx context.Context, params *UserListParams) (list []*entity.User, total int64, err error)
	// ScanOne 查询用户
	ScanOne(ctx context.Context, id int64, val any) error
	// ScanOneByUsername 查询用户通过Username
	ScanOneByUsername(ctx context.Context, username string, val any) error
	// ScanOneByEmail 查询用户通过Email
	ScanOneByEmail(ctx context.Context, email string, val any) error
	// ScanOneByMobile 查询用户通过Mobile
	ScanOneByMobile(ctx context.Context, mobile string, val any) error
	// ScanList 查询用户列表
	ScanList(ctx context.Context, params *UserListParams, list any, total *int64) error

	/*
		登录日志相关
	*/

	// FindOneLog 查询登录日志
	FindOneLog(ctx context.Context, id int64) (*entity.LoginLog, error)
	// FindListLog 查询登录日志列表
	FindListLog(ctx context.Context, params *LoginListParams) (list []*entity.LoginLog, total int64, err error)
	// ScanOneLog 查询登录日志
	ScanOneLog(ctx context.Context, id int64, val any) error
	// ScanListLog 查询登录日志列表
	ScanListLog(ctx context.Context, params *LoginListParams, list any, total *int64) error
}

// UserListParams 用户列表查询参数
type UserListParams struct {
	Username string // 用户名
	RealName string // 真实姓名
	Mobile   string // 手机号
	Email    string // 邮箱
	Status   *int64 // 状态
	DeptID   *int64 // 部门ID
	RoleID   *int64 // 角色ID
	PostID   *int64 // 岗位ID
	Page     int    // 页码
	Size     int    // 每页数量
}

type LoginListParams struct {
	Uid    int64  // 用户ID
	Status int64  // 登录状态
	Device string // 登录设备
	Os     string // 操作系统
	Start  string // 开始时间
	End    string // 结束时间
	Page   int
	Size   int
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
	q  *query.Query
}

func (u *userRepository) Create(ctx context.Context, user *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) Update(ctx context.Context, user *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) Delete(ctx context.Context, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindOne(ctx context.Context, id int64) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindOneByUsername(ctx context.Context, username string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindOneByMobile(ctx context.Context, mobile string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindList(ctx context.Context, params *UserListParams) (list []*entity.User, total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) ScanOne(ctx context.Context, id int64, val any) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) ScanOneByUsername(ctx context.Context, username string, val any) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) ScanOneByEmail(ctx context.Context, email string, val any) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) ScanOneByMobile(ctx context.Context, mobile string, val any) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) ScanList(ctx context.Context, params *UserListParams, list any, total *int64) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindOneLog(ctx context.Context, id int64) (*entity.LoginLog, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) FindListLog(ctx context.Context, params *LoginListParams) (list []*entity.LoginLog, total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) ScanOneLog(ctx context.Context, id int64, val any) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) ScanListLog(ctx context.Context, params *LoginListParams, list any, total *int64) error {
	//TODO implement me
	panic("implement me")
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
		q:  query.Use(db),
	}
}
