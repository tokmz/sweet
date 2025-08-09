package system

import (
	"context"
	"errors"
	"sweet/internal/global"
	"sweet/internal/models"
	systemDTO "sweet/internal/models/dto/system"
	"sweet/internal/models/entity"
	"sweet/internal/models/query"
	"sweet/pkg/crypto"
	"sweet/pkg/errs"
	"sweet/pkg/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() IUserService {
	return &UserService{}
}

func (s *UserService) CreateUser(ctx context.Context, req *systemDTO.CreateUserReq) error {
	return global.Query.Transaction(func(tx *query.Query) error {
		dao := tx.SysUser
		do := dao.WithContext(ctx)

		// 检查用户名、手机号、邮箱是否已存在
		do = do.Where(dao.Username.Eq(req.Username))
		if req.Email != nil {
			do = do.Where(dao.Email.Eq(*req.Email))
		}
		if req.Phone != nil {
			do = do.Where(dao.Phone.Eq(*req.Phone))
		}

		if user, err := do.First(); err == nil {
			// 判断是哪个字段重复
			if user.Username == req.Username {
				return errs.ErrUserExists
			}
			if user.Email == req.Email {
				return errs.ErrEmailExists
			}
			return errs.ErrPhoneExists
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error(
				"查询用户失败",
				zap.String("Username", req.Username),
				zap.String("Email", *req.Email),
				zap.String("Phone", *req.Phone),
				zap.Error(err),
			)
			return errs.ErrServer
		} else {
			// 获取salt
			salt := crypto.Salt()
			// 创建用户
			userEntity := entity.SysUser{
				Username: req.Username,
				Password: crypto.MD5(req.Password + salt),
				Salt:     salt,
				Realname: req.Realname,
				Nickname: req.Nickname,
				Avatar:   req.Avatar,
				Email:    req.Email,
				Phone:    req.Phone,
				Status:   req.Status,
				RoleID:   req.RoleID,
				DeptID:   req.DeptID,
				PostID:   req.PostID,
				Remark:   req.Remark,
			}
			if err = dao.WithContext(ctx).Create(&userEntity); err != nil {
				// 创建用户失败
				global.Logger.Error(
					"创建用户失败",
					zap.Any("entity", userEntity),
					zap.Error(err),
				)
				return errs.ErrServer
			}
			return nil
		}
	})

}

func (s *UserService) DeleteUser(ctx context.Context, req *systemDTO.DeleteUserReq) error {
	dao := global.Query.SysUser
	if _, err := dao.WithContext(ctx).Where(dao.ID.In(req.Ids...)).Delete(); err != nil {
		global.Logger.Error(
			"删除用户失败",
			zap.Any("ids", req.Ids),
			zap.Error(err),
		)
		return errs.ErrServer
	}
	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *systemDTO.UpdateUserReq) error {
	return global.Query.Transaction(func(tx *query.Query) error {
		dao := tx.SysUser

		// 检查用户是否存在
		if _, err := dao.WithContext(ctx).Where(dao.ID.Eq(req.ID)).First(); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errs.ErrUserNotFound
			}
			global.Logger.Error(
				"查询用户失败",
				zap.Int64("id", req.ID),
				zap.Error(err),
			)
			return errs.ErrServer
		}

		// 检查邮箱重复性（如果提供了邮箱且不为空）
		if req.Email != nil && *req.Email != "" {
			if existingUser, err := dao.WithContext(ctx).Where(
				dao.ID.Neq(req.ID),
				dao.Email.Eq(*req.Email),
			).First(); err == nil && existingUser != nil {
				return errs.ErrEmailExists
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				global.Logger.Error(
					"检查邮箱重复性失败",
					zap.Int64("id", req.ID),
					zap.String("email", *req.Email),
					zap.Error(err),
				)
				return errs.ErrServer
			}
		}

		// 检查手机号重复性（如果提供了手机号且不为空）
		if req.Phone != nil && *req.Phone != "" {
			if existingUser, err := dao.WithContext(ctx).Where(
				dao.ID.Neq(req.ID),
				dao.Phone.Eq(*req.Phone),
			).First(); err == nil && existingUser != nil {
				return errs.ErrPhoneExists
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				global.Logger.Error(
					"检查手机号重复性失败",
					zap.Int64("id", req.ID),
					zap.String("phone", *req.Phone),
					zap.Error(err),
				)
				return errs.ErrServer
			}
		}

		// 构建更新实体
		updateEntity := entity.SysUser{
			Realname: req.Realname,
			Nickname: req.Nickname,
			Avatar:   req.Avatar,
			Email:    req.Email,
			Phone:    req.Phone,
			Status:   req.Status,
			RoleID:   req.RoleID,
			DeptID:   req.DeptID,
			PostID:   req.PostID,
			Remark:   req.Remark,
		}

		// 如果需要更新密码
		if req.Password != "" {
			salt := crypto.Salt()
			updateEntity.Password = crypto.MD5(req.Password + salt)
			updateEntity.Salt = salt
		}

		// 执行更新
		if _, err := dao.WithContext(ctx).Where(dao.ID.Eq(req.ID)).Updates(updateEntity); err != nil {
			global.Logger.Error(
				"更新用户失败",
				zap.Int64("id", req.ID),
				zap.Any("entity", updateEntity),
				zap.Error(err),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (s *UserService) ListUser(ctx context.Context, req *systemDTO.ListUserReq) (*systemDTO.ListUserRes, error) {
	dao := global.Query.SysUser
	do := dao.WithContext(ctx)

	// 构建查询条件
	if req.Username != "" {
		do = do.Where(dao.Username.Like("%" + req.Username + "%"))
	}
	if req.RoleID > 0 {
		do = do.Where(dao.RoleID.Eq(req.RoleID))
	}
	if req.DeptID > 0 {
		do = do.Where(dao.DeptID.Eq(req.DeptID))
	}
	if req.PostID > 0 {
		do = do.Where(dao.PostID.Eq(req.PostID))
	}
	if req.Status > 0 {
		do = do.Where(dao.Status.Eq(req.Status))
	}

	// 时间范围查询
	if req.StartTime > 0 {
		do = do.Where(dao.CreatedAt.Gte(utils.UnixToTime(req.StartTime)))
	}
	if req.EndTime > 0 {
		do = do.Where(dao.CreatedAt.Lte(utils.UnixToTime(req.EndTime)))
	}

	// 排序
	if req.Field != "" {
		if orderExpr, ok := dao.GetFieldByName(req.Field); ok {
			if req.Order == "asc" {
				do = do.Order(orderExpr)
			} else {
				do = do.Order(orderExpr.Desc())
			}
		}
	} else {
		// 默认按创建时间倒序
		do = do.Order(dao.CreatedAt.Desc())
	}

	// 关联查询角色信息
	do = do.Preload(dao.Role)

	// 分页查询
	offset := (req.Page - 1) * req.Size
	users, count, err := do.FindByPage(offset, req.Size)
	if err != nil {
		global.Logger.Error(
			"查询用户列表失败",
			zap.Any("req", req),
			zap.Error(err),
		)
		return nil, errs.ErrServer
	}

	// 转换为DTO
	list := make([]*systemDTO.ListUserItem, 0, len(users))
	for _, user := range users {
		item := &systemDTO.ListUserItem{
			ID:        user.ID,
			Username:  user.Username,
			Realname:  user.Realname,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Status:    user.Status,
			RoleID:    utils.Deref(user.RoleID),
			CreatedAt: user.CreatedAt,
		}
		// 设置角色名称
		if user.Role != nil {
			item.RoleName = user.Role.Name
		}
		list = append(list, item)
	}

	return &systemDTO.ListUserRes{
		Total: count,
		List:  list,
	}, nil
}

func (s *UserService) GetUserDetail(ctx context.Context, req *models.IDReq) (*systemDTO.UserDetailRes, error) {
	dao := global.Query.SysUser
	do := dao.WithContext(ctx)

	// 关联查询角色、部门、岗位信息
	user, err := do.Where(dao.ID.Eq(req.ID)).Preload(dao.Role).Preload(dao.Dept).Preload(dao.Post).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		global.Logger.Error(
			"查询用户详情失败",
			zap.Int64("id", req.ID),
			zap.Error(err),
		)
		return nil, errs.ErrServer
	}

	// 转换为DTO
	detail := &systemDTO.UserDetailRes{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		Salt:      user.Salt,
		Realname:  user.Realname,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
		RoleID:    utils.Deref(user.RoleID),
		DeptID:    utils.Deref(user.DeptID),
		PostID:    utils.Deref(user.PostID),
		Remark:    user.Remark,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// 设置关联信息
	if user.Role != nil {
		detail.RoleName = user.Role.Name
	}
	if user.Dept != nil {
		detail.DeptName = user.Dept.Name
	}
	if user.Post != nil {
		detail.PostName = user.Post.Name
	}

	return detail, nil
}
