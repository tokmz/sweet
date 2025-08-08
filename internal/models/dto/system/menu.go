package system

import (
	"sweet/internal/models"
	"time"
)

// CreateMenuReq 创建菜单
type CreateMenuReq struct {
	Uid           int64   `json:"uid"`             // 用户ID
	ParentID      *int64  `json:"parent_id"`       // 父菜单ID
	Name          string  `json:"name"`            // 组件名称/路由名称
	Title         string  `json:"title"`           // 菜单名称
	Path          *string `json:"path"`            // 路由地址
	Component     *string `json:"component"`       // 组件地址
	MenuType      *int64  `json:"menu_type"`       // 菜单类型（1 目录 2 菜单 3 按钮）
	Status        *int64  `json:"status"`          // 菜单状态(1 正常 2 停用)
	Perms         *string `json:"perms"`           // 权限标识
	Icon          *string `json:"icon"`            // 菜单图标
	Order_        *int64  `json:"order"`           // 显示顺序 从大到小
	Remark        *string `json:"remark"`          // 备注
	Query         *string `json:"query"`           // 路由参数
	IsFrame       *int64  `json:"is_frame"`        // 是否外联（1 是 2 否）
	ShowBadge     *int64  `json:"show_badge"`      // 是否显示徽章（1 是 2 否）
	ShowTextBadge *string `json:"show_text_badge"` // 文本徽章内容
	IsHide        *int64  `json:"is_hide"`         // 是否在菜单中隐藏（1 是 2 否）
	IsHideTab     *int64  `json:"is_hide_tab"`     // 是否在标签页中隐藏（1 是 2 否）
	Link          *string `json:"link"`            // 外链地址
	IsIframe      *int64  `json:"is_iframe"`       // 是否iframe（1 是 2 否）
	KeepAlive     *int64  `json:"keep_alive"`      // 是否缓存页面（1 是 2否）
	FixedTab      *int64  `json:"fixed_tab"`       // 是否固定标签页（1 是 2 否）
	IsFirstLevel  *int64  `json:"is_first_level"`  // 是否为一级菜单
	ActivePath    *string `json:"active_path"`     // 激活菜单路径
}

// DeleteMenuReq 删除菜单
type DeleteMenuReq models.IDReq

// UpdateMenuReq 更新菜单
type UpdateMenuReq struct {
	models.IDReq
	Uid           int64   `json:"uid"`             // 用户ID
	ParentID      *int64  `json:"parent_id"`       // 父菜单ID
	Name          string  `json:"name"`            // 组件名称/路由名称
	Title         string  `json:"title"`           // 菜单名称
	Path          *string `json:"path"`            // 路由地址
	Component     *string `json:"component"`       // 组件地址
	MenuType      *int64  `json:"menu_type"`       // 菜单类型（1 目录 2 菜单 3 按钮）
	Status        *int64  `json:"status"`          // 菜单状态(1 正常 2 停用)
	Perms         *string `json:"perms"`           // 权限标识
	Icon          *string `json:"icon"`            // 菜单图标
	Order_        *int64  `json:"order"`           // 显示顺序 从大到小
	Remark        *string `json:"remark"`          // 备注
	Query         *string `json:"query"`           // 路由参数
	IsFrame       *int64  `json:"is_frame"`        // 是否外联（1 是 2 否）
	ShowBadge     *int64  `json:"show_badge"`      // 是否显示徽章（1 是 2 否）
	ShowTextBadge *string `json:"show_text_badge"` // 文本徽章内容
	IsHide        *int64  `json:"is_hide"`         // 是否在菜单中隐藏（1 是 2 否）
	IsHideTab     *int64  `json:"is_hide_tab"`     // 是否在标签页中隐藏（1 是 2 否）
	Link          *string `json:"link"`            // 外链地址
	IsIframe      *int64  `json:"is_iframe"`       // 是否iframe（1 是 2 否）
	KeepAlive     *int64  `json:"keep_alive"`      // 是否缓存页面（1 是 2否）
	FixedTab      *int64  `json:"fixed_tab"`       // 是否固定标签页（1 是 2 否）
	IsFirstLevel  *int64  `json:"is_first_level"`  // 是否为一级菜单
	ActivePath    *string `json:"active_path"`     // 激活菜单路径
}

// MenuTreeReq 菜单树列表请求
type MenuTreeReq struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

// MenuTreeItem 菜单树列表响应Item
type MenuTreeItem struct {
	ID        int64      `json:"id"`         // 菜单ID
	ParentID  *int64     `json:"parent_id"`  // 父菜单ID
	Title     string     `json:"title"`      // 菜单名称
	Path      *string    `json:"path"`       // 路由地址
	MenuType  *int64     `json:"menu_type"`  // 菜单类型（1 目录 2 菜单 3 按钮）
	Status    *int64     `json:"status"`     // 菜单状态(1 正常 2 停用)
	Order_    *int64     `json:"order"`      // 显示顺序 从大到小
	IsHide    *int64     `json:"is_hide"`    // 是否在菜单中隐藏（1 是 2 否）
	KeepAlive *int64     `json:"keep_alive"` // 是否缓存页面（1 是 2否）
	Remark    *string    `json:"remark"`     // 备注
	CreatedAt *time.Time `json:"created_at"` // 创建时间
}

// MenuTreeRes 菜单树列表响应
type MenuTreeRes []*MenuTreeItem

// MenuButtonItem 按钮详情Item
type MenuButtonItem struct {
	ID    int64   `json:"id"`
	Title string  `json:"title"`
	Perms *string `json:"perms"`
}

// MenuOptionsItem 菜单选项响应Item
type MenuOptionsItem struct {
	ID       int64              `json:"id"`
	Title    string             `json:"title"`              // 名称
	Type     *int64             `json:"type"`               // 菜单类型（1 目录 2 菜单 3 按钮）
	Button   []*MenuButtonItem  `json:"button,omitempty"`   // 按钮列表
	Children []*MenuOptionsItem `json:"children,omitempty"` // 子菜单
}

// MenuOptionsRes 菜单选项响应
type MenuOptionsRes []*MenuOptionsItem

// MenuButton 按钮详情
type MenuButton struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Perms  *string `json:"perms"`
	Status *int64  `json:"status"`
}

// MenuDetailRes 菜单详情响应
type MenuDetailRes struct {
	ID            int64         `json:"id"`              // 菜单ID
	ParentID      *int64        `json:"parent_id"`       // 父菜单ID
	Name          string        `json:"name"`            // 组件名称/路由名称
	Title         string        `json:"title"`           // 菜单名称
	Path          *string       `json:"path"`            // 路由地址
	Component     *string       `json:"component"`       // 组件地址
	MenuType      *int64        `json:"menu_type"`       // 菜单类型（1 目录 2 菜单 3 按钮）
	Status        *int64        `json:"status"`          // 菜单状态(1 正常 2 停用)
	Perms         *string       `json:"perms"`           // 权限标识
	Icon          *string       `json:"icon"`            // 菜单图标
	Order_        *int64        `json:"order"`           // 显示顺序 从大到小
	Remark        *string       `json:"remark"`          // 备注
	Query         *string       `json:"query"`           // 路由参数
	IsFrame       *int64        `json:"is_frame"`        // 是否外联（1 是 2 否）
	ShowBadge     *int64        `json:"show_badge"`      // 是否显示徽章（1 是 2 否）
	ShowTextBadge *string       `json:"show_text_badge"` // 文本徽章内容
	IsHide        *int64        `json:"is_hide"`         // 是否在菜单中隐藏（1 是 2 否）
	IsHideTab     *int64        `json:"is_hide_tab"`     // 是否在标签页中隐藏（1 是 2 否）
	Link          *string       `json:"link"`            // 外链地址
	IsIframe      *int64        `json:"is_iframe"`       // 是否iframe（1 是 2 否）
	KeepAlive     *int64        `json:"keep_alive"`      // 是否缓存页面（1 是 2否）
	FixedTab      *int64        `json:"fixed_tab"`       // 是否固定标签页（1 是 2 否）
	IsFirstLevel  *int64        `json:"is_first_level"`  // 是否为一级菜单
	ActivePath    *string       `json:"active_path"`     // 激活菜单路径
	CreateBy      *int64        `json:"create_by"`       // 创建者
	UpdateBy      *int64        `json:"update_by"`       // 更新者
	CreatedAt     *time.Time    `json:"created_at"`      // 创建时间
	UpdatedAt     *time.Time    `json:"updated_at"`      // 更新时间
	Buttons       []*MenuButton `json:"buttons"`
}

// 创建按钮
type CreateButtonReq struct {
	models.IDReq        // 菜单ID
	Title        string `json:"title"`  // 按钮名称
	Perms        string `json:"perms"`  // 权限标识
	Status       int64  `json:"status"` // 状态
	Uid          int64  `json:"uid"`    // 用户ID
}

// 删除按钮
type DeleteButtonReq struct {
	models.IDReq  // 菜单ID
	models.IdsReq // 按钮ID列表
}

// 更新按钮
type UpdateButtonReq struct {
	models.IDReq        // 菜单ID
	Uid          int64  `json:"uid"`    // 用户ID
	Title        string `json:"title"`  // 按钮名称
	Perms        string `json:"perms"`  // 权限标识
	Status       int64  `json:"status"` // 状态
}
