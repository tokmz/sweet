package dto

// CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	ParentID     *int64  `json:"parent_id"`     // 父菜单ID
	Name         string  `json:"name"`          // 菜单名称
	Permission   *string `json:"permission"`    // 权限标识
	Type         int64   `json:"type"`          // 类型 1-目录 2-菜单 3-按钮
	Path         *string `json:"path"`          // 路由地址
	Component    *string `json:"component"`     // 组件路径
	Redirect     *string `json:"redirect"`      // 重定向地址
	Icon         *string `json:"icon"`          // 菜单图标
	Sort         *int    `json:"sort"`          // 排序
	Hidden       *int64  `json:"hidden"`        // 是否隐藏 1-是 2-否
	Status       *int64  `json:"status"`        // 状态 1-正常 2-禁用
	AlwaysShow   *int64  `json:"always_show"`   // 是否总是显示 1-是 2-否
	KeepAlive    *int64  `json:"keep_alive"`    // 是否缓存 1-是 2-否
	Target       *string `json:"target"`        // 打开方式 _self _blank
	Title        *string `json:"title"`         // 菜单标题
	ActiveMenu   *string `json:"active_menu"`   // 激活菜单
	Breadcrumb   *int64  `json:"breadcrumb"`    // 是否显示面包屑 1-是 2-否
	Affix        *int64  `json:"affix"`         // 是否固定 1-是 2-否
	FrameSrc     *string `json:"frame_src"`     // iframe地址
	FrameLoading *int64  `json:"frame_loading"` // iframe加载状态 1-是 2-否
	Transition   *string `json:"transition"`    // 过渡动画
	Remark       *string `json:"remark"`        // 备注
	CreatedBy    *int64  `json:"created_by"`    // 创建人
}

// UpdateMenuReq 更新菜单请求
type UpdateMenuReq struct {
	ID           int64   `json:"id"`            // 菜单ID
	ParentID     *int64  `json:"parent_id"`     // 父菜单ID
	Name         string  `json:"name"`          // 菜单名称
	Permission   *string `json:"permission"`    // 权限标识
	Type         int64   `json:"type"`          // 类型 1-目录 2-菜单 3-按钮
	Path         *string `json:"path"`          // 路由地址
	Component    *string `json:"component"`     // 组件路径
	Redirect     *string `json:"redirect"`      // 重定向地址
	Icon         *string `json:"icon"`          // 菜单图标
	Sort         *int    `json:"sort"`          // 排序
	Hidden       *int64  `json:"hidden"`        // 是否隐藏 1-是 2-否
	Status       *int64  `json:"status"`        // 状态 1-正常 2-禁用
	AlwaysShow   *int64  `json:"always_show"`   // 是否总是显示 1-是 2-否
	KeepAlive    *int64  `json:"keep_alive"`    // 是否缓存 1-是 2-否
	Target       *string `json:"target"`        // 打开方式 _self _blank
	Title        *string `json:"title"`         // 菜单标题
	ActiveMenu   *string `json:"active_menu"`   // 激活菜单
	Breadcrumb   *int64  `json:"breadcrumb"`    // 是否显示面包屑 1-是 2-否
	Affix        *int64  `json:"affix"`         // 是否固定 1-是 2-否
	FrameSrc     *string `json:"frame_src"`     // iframe地址
	FrameLoading *int64  `json:"frame_loading"` // iframe加载状态 1-是 2-否
	Transition   *string `json:"transition"`    // 过渡动画
	Remark       *string `json:"remark"`        // 备注
	UpdatedBy    *int64  `json:"updated_by"`    // 更新人
}

// FindMenuTreeReq 查询菜单树请求
type FindMenuTreeReq struct {
	ParentID *int64  `json:"parent_id"` // 父菜单ID
	Name     *string `json:"name"`      // 菜单名称
	RoleID   *int64  `json:"role_id"`   // 角色ID
}

// RouteTreeReq 查询路由树请求
type RouteTreeReq struct {
	RoleID int64 `json:"role_id"` // 角色ID
}
