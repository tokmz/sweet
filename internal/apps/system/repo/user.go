package repo

import (
	"context"
	"errors"
	"sweet/pkg/cache"
	"sweet/pkg/errs"
	"sweet/pkg/logger"
	"sweet/pkg/utils"
	"time"

	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"

	"gorm.io/gorm"
)

var (
	ErrAccountNotFound = errs.New(10001, "账号不存在")
	ErrAccountExists   = errs.New(10002, "账号已存在")
	ErrEmailExists     = errs.New(10003, "邮箱已存在")
	ErrMobileExists    = errs.New(10004, "手机号已存在")
	ErrInvalidPassword = errs.New(10005, "密码错误")
	ErrInvalidUserID   = errs.New(10006, "无效ID")
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
	Status   int64  // 状态
	DeptID   int64  // 部门ID
	RoleID   int64  // 角色ID
	PostID   int64  // 岗位ID
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
	c  *cache.Client
}

func (u *userRepository) Create(ctx context.Context, user *entity.User) error {
	return u.q.Transaction(func(tx *query.Query) error {
		do := tx.User.WithContext(ctx)
		field := tx.User
		if info, err := do.Where(field.Username.Eq(user.Username)).
			Or(field.Mobile.Eq(*user.Mobile)).
			Or(field.Email.Eq(*user.Email)).
			First(); err == nil {
			if info.Username == user.Username {
				return ErrAccountExists
			} else if info.Mobile == user.Mobile {
				return ErrMobileExists
			} else if info.Email == user.Email {
				return ErrEmailExists
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"查询用户失败",
				logger.Err(err),
				logger.String("用户名", user.Username),
				logger.String("手机号", utils.SafeString(user.Mobile)),
				logger.String("邮箱", utils.SafeString(user.Email)),
			)
			return errs.ErrServer
		}
		if err := do.Create(user); err != nil {
			logger.Error(
				"创建用户失败",
				logger.Err(err),
				logger.String("用户名", user.Username),
				logger.String("真实姓名", utils.SafeString(user.RealName)),
				logger.Int64("用户状态", utils.SafeInt64(user.Status)),
			)
			return errs.ErrServer
		}
		return nil
	})
}

func (u *userRepository) Update(ctx context.Context, user *entity.User) error {
	return u.q.Transaction(func(tx *query.Query) error {
		do := tx.User.WithContext(ctx)
		field := tx.User

		if info, err := do.Where(field.Username.Eq(user.Username)).
			Or(field.Mobile.Eq(*user.Mobile)).
			Or(field.Email.Eq(*user.Email)).
			First(); err == nil {
			if info.Username == user.Username {
				return ErrAccountExists
			} else if info.Mobile == user.Mobile {
				return ErrMobileExists
			} else if info.Email == user.Email {
				return ErrEmailExists
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"查询用户失败",
				logger.Err(err),
				logger.Int64("用户ID", user.ID),
				logger.String("用户名", user.Username),
				logger.String("手机号", utils.SafeString(user.Mobile)),
				logger.String("邮箱", utils.SafeString(user.Email)),
			)
			return errs.ErrServer
		}

		// 更新用户信息
		if _, err := do.Where(field.ID.Eq(user.ID)).Updates(user); err != nil {
			logger.Error(
				"更新用户信息失败",
				logger.Err(err),
				logger.Int64("用户ID", user.ID),
				logger.String("用户名", user.Username),
				logger.Int64("用户状态", utils.SafeInt64(user.Status)),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (u *userRepository) Delete(ctx context.Context, ids []int64) error {
	return u.q.Transaction(func(tx *query.Query) error {
		do := tx.User.WithContext(ctx)
		if _, err := do.Where(tx.User.ID.In(ids...)).Delete(); err != nil {
			logger.Error(
				"删除用户失败",
				logger.Err(err),
				logger.Any("用户ID列表", ids),
				logger.Int("删除数量", len(ids)),
			)
			return errs.ErrServer
		}
		return nil
	})
}

func (u *userRepository) FindOne(ctx context.Context, id int64) (*entity.User, error) {
	if info, err := u.q.User.WithContext(ctx).Where(u.q.User.ID.Eq(id)).First(); err == nil {
		return info, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccountNotFound
	} else {
		logger.Error(
			"查询用户失败",
			logger.Err(err),
			logger.Int64("用户ID", id),
		)
		return nil, errs.ErrServer
	}
}

func (u *userRepository) FindOneByUsername(ctx context.Context, username string) (*entity.User, error) {
	if info, err := u.q.User.WithContext(ctx).Where(u.q.User.Username.Eq(username)).First(); err == nil {
		return info, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccountNotFound
	} else {
		logger.Error(
			"查询用户失败",
			logger.Err(err),
			logger.String("用户名", username),
		)
		return nil, errs.ErrServer
	}
}

func (u *userRepository) FindOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	if info, err := u.q.User.WithContext(ctx).Where(u.q.User.Email.Eq(email)).First(); err == nil {
		return info, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccountNotFound
	} else {
		logger.Error(
			"查询用户失败",
			logger.Err(err),
			logger.String("邮箱", email),
		)
		return nil, errs.ErrServer
	}
}

func (u *userRepository) FindOneByMobile(ctx context.Context, mobile string) (*entity.User, error) {
	if info, err := u.q.User.WithContext(ctx).Where(u.q.User.Mobile.Eq(mobile)).First(); err == nil {
		return info, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccountNotFound
	} else {
		logger.Error(
			"查询用户失败",
			logger.Err(err),
			logger.String("手机号", mobile),
		)
		return nil, errs.ErrServer
	}
}

func (u *userRepository) FindList(ctx context.Context, params *UserListParams) (list []*entity.User, total int64, err error) {
	do := u.q.User.WithContext(ctx)
	field := u.q.User
	if params.Username != "" {
		do = do.Where(field.Username.Like("%" + params.Username + "%"))
	}
	if params.RealName != "" {
		do = do.Where(field.RealName.Like("%" + params.RealName + "%"))
	}
	if params.Mobile != "" {
		do = do.Where(field.Mobile.Eq(params.Mobile))
	}
	if params.Email != "" {
		do = do.Where(field.Email.Eq(params.Email))
	}
	if params.Status != 0 {
		do = do.Where(field.Status.Eq(params.Status))
	}

	if params.DeptID != 0 {
		do = do.Where(field.DeptID.Eq(params.DeptID))
	}
	if params.RoleID != 0 {
		do = do.Where(field.RoleID.Eq(params.RoleID))
	}
	if params.PostID != 0 {
		do = do.Where(field.PostID.Eq(params.PostID))
	}

	if list, total, err := do.FindByPage(params.Page, params.Size); err == nil {
		return list, total, nil
	} else {
		logger.Error(
			"查询用户列表失败",
			logger.Err(err),
			logger.String("用户名", params.Username),
			logger.String("真实姓名", params.RealName),
			logger.Int("页码", params.Page),
			logger.Int("每页数量", params.Size),
		)
		return nil, 0, errs.ErrServer
	}
}

func (u *userRepository) ScanOne(ctx context.Context, id int64, val any) error {
	if err := u.q.User.WithContext(ctx).Where(u.q.User.ID.Eq(id)).Scan(val); err == nil {
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAccountNotFound
	} else {
		logger.Error(
			"查询用户失败",
			logger.Err(err),
			logger.Int64("用户ID", id),
		)
		return errs.ErrServer
	}
}

func (u *userRepository) ScanOneByUsername(ctx context.Context, username string, val any) error {
	if err := u.q.User.WithContext(ctx).Where(u.q.User.Username.Eq(username)).Scan(val); err == nil {
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAccountNotFound
	} else {
		logger.Error(
			"查询用户失败",
			logger.Err(err),
			logger.String("用户名", username),
		)
		return errs.ErrServer
	}
}

func (u *userRepository) ScanOneByEmail(ctx context.Context, email string, val any) error {
	if err := u.q.User.WithContext(ctx).Where(u.q.User.Email.Eq(email)).Scan(val); err == nil {
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAccountNotFound
	} else {
		logger.Error(
			"查询用户失败",
			logger.Err(err),
			logger.String("邮箱", email),
		)
		return errs.ErrServer
	}
}

func (u *userRepository) ScanOneByMobile(ctx context.Context, mobile string, val any) error {
	if err := u.q.User.WithContext(ctx).Where(u.q.User.Mobile.Eq(mobile)).Scan(val); err == nil {
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAccountNotFound
	} else {
		logger.Error(
			"查询用户失败",
			logger.Err(err),
			logger.String("手机号", mobile),
		)
		return errs.ErrServer
	}
}

func (u *userRepository) ScanList(ctx context.Context, params *UserListParams, list any, total *int64) error {
	do := u.q.User.WithContext(ctx)
	field := u.q.User
	if params.Username != "" {
		do = do.Where(field.Username.Like("%" + params.Username + "%"))
	}
	if params.RealName != "" {
		do = do.Where(field.RealName.Like("%" + params.RealName + "%"))
	}
	if params.Mobile != "" {
		do = do.Where(field.Mobile.Eq(params.Mobile))
	}
	if params.Email != "" {
		do = do.Where(field.Email.Eq(params.Email))
	}
	if params.Status != 0 {
		do = do.Where(field.Status.Eq(params.Status))
	}
	if params.DeptID != 0 {
		do = do.Where(field.DeptID.Eq(params.DeptID))
	}
	if params.RoleID != 0 {
		do = do.Where(field.RoleID.Eq(params.RoleID))
	}
	if params.PostID != 0 {
		do = do.Where(field.PostID.Eq(params.PostID))
	}

	if count, err := do.ScanByPage(list, params.Page, params.Size); err == nil {
		*total = count
		return nil
	} else {
		logger.Error(
			"查询用户列表失败",
			logger.Err(err),
			logger.String("用户名", params.Username),
			logger.String("真实姓名", params.RealName),
			logger.Int("页码", params.Page),
			logger.Int("每页数量", params.Size),
		)
		return errs.ErrServer
	}
}

func (u *userRepository) FindOneLog(ctx context.Context, id int64) (*entity.LoginLog, error) {
	if info, err := u.q.LoginLog.WithContext(ctx).Where(u.q.LoginLog.ID.Eq(id)).First(); err == nil {
		return info, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccountNotFound
	} else {
		logger.Error(
			"查询登录日志失败",
			logger.Err(err),
			logger.Int64("日志ID", id),
		)
		return nil, errs.ErrServer
	}
}

func (u *userRepository) FindListLog(ctx context.Context, params *LoginListParams) (list []*entity.LoginLog, total int64, err error) {
	do := u.q.LoginLog.WithContext(ctx)
	field := u.q.LoginLog
	if params.Uid != 0 {
		do = do.Where(field.UID.Eq(params.Uid))
	}
	if params.Status != 0 {
		do = do.Where(field.Status.Eq(params.Status))
	}
	if params.Device != "" {
		do = do.Where(field.Device.Eq(params.Device))
	}
	if params.Os != "" {
		do = do.Where(field.Os.Eq(params.Os))
	}
	if params.Start != "" {
		if start, err := time.Parse(time.DateTime, params.Start); err == nil {
			do = do.Where(field.CreatedAt.Gte(start))
		} else {
			logger.Error(
				"查询登录日志失败",
				logger.Err(err),
				logger.String("错误信息", "时间格式错误"),
				logger.String("开始时间", params.Start),
			)
			return nil, 0, errs.ErrServer
		}
	}
	if params.End != "" {
		if end, err := time.Parse(time.DateTime, params.End); err == nil {
			do = do.Where(field.CreatedAt.Lte(end))
		} else {
			logger.Error(
				"查询登录日志失败",
				logger.Err(err),
				logger.String("错误信息", "时间格式错误"),
				logger.String("结束时间", params.End),
			)
			return nil, 0, errs.ErrServer
		}
	}
	if list, total, err = do.FindByPage(params.Page, params.Size); err == nil {
		return list, total, nil
	} else {
		logger.Error(
			"查询登录日志列表失败",
			logger.Err(err),
			logger.Int64("用户ID", params.Uid),
			logger.Int64("登录状态", params.Status),
			logger.Int("页码", params.Page),
			logger.Int("每页数量", params.Size),
		)
		return nil, 0, errs.ErrServer
	}
}

func (u *userRepository) ScanOneLog(ctx context.Context, id int64, val any) error {
	if err := u.q.LoginLog.WithContext(ctx).Where(u.q.LoginLog.ID.Eq(id)).Scan(val); err == nil {
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAccountNotFound
	} else {
		logger.Error(
			"查询登录日志失败",
			logger.Err(err),
			logger.Int64("日志ID", id),
		)
		return errs.ErrServer
	}
}

func (u *userRepository) ScanListLog(ctx context.Context, params *LoginListParams, list any, total *int64) error {
	do := u.q.LoginLog.WithContext(ctx)
	field := u.q.LoginLog

	// 构建查询条件
	if params.Uid != 0 {
		do = do.Where(field.UID.Eq(params.Uid))
	}
	if params.Status != 0 {
		do = do.Where(field.Status.Eq(params.Status))
	}
	if params.Device != "" {
		do = do.Where(field.Device.Eq(params.Device))
	}
	if params.Os != "" {
		do = do.Where(field.Os.Eq(params.Os))
	}

	// 处理时间范围
	if params.Start != "" {
		if start, err := time.Parse(time.DateTime, params.Start); err == nil {
			do = do.Where(field.CreatedAt.Gte(start))
		} else {
			logger.Error(
				"扫描登录日志失败",
				logger.Err(err),
				logger.String("错误信息", "开始时间格式错误"),
				logger.String("开始时间", params.Start),
			)
			return errs.ErrServer
		}
	}
	if params.End != "" {
		if end, err := time.Parse(time.DateTime, params.End); err == nil {
			do = do.Where(field.CreatedAt.Lte(end))
		} else {
			logger.Error(
				"扫描登录日志失败",
				logger.Err(err),
				logger.String("错误信息", "结束时间格式错误"),
				logger.String("结束时间", params.End),
			)
			return errs.ErrServer
		}
	}

	// 执行分页查询并扫描到指定的结构体
	if count, err := do.ScanByPage(list, params.Page, params.Size); err == nil {
		*total = count
		return nil
	} else {
		logger.Error(
			"扫描登录日志列表失败",
			logger.Err(err),
			logger.Int64("用户ID", params.Uid),
			logger.Int64("登录状态", params.Status),
			logger.Int("页码", params.Page),
			logger.Int("每页数量", params.Size),
		)
		return errs.ErrServer
	}
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB, c *cache.Client) UserRepository {
	return &userRepository{
		db: db,
		q:  query.Use(db),
		c:  c,
	}
}
