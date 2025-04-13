package vo

// RouteTree 路由树结构
type RouteTree struct {
	ID        int64        `json:"id"`        // 菜单ID
	ParentID  int64        `json:"parentId"`  // 父菜单ID
	Name      string       `json:"name"`      // 路由名称
	Path      string       `json:"path"`      // 路由地址
	Component string       `json:"component"` // 组件路径
	Redirect  string       `json:"redirect"`  // 重定向地址
	Meta      RouteMeta    `json:"meta"`      // 路由元信息
	Children  []*RouteTree `json:"children"`  // 子路由
}

// RouteMeta 路由元信息
type RouteMeta struct {
	Title        string `json:"title"`        // 菜单标题
	Icon         string `json:"icon"`         // 菜单图标
	Hidden       bool   `json:"hidden"`       // 是否隐藏
	KeepAlive    bool   `json:"keepAlive"`    // 是否缓存
	AlwaysShow   bool   `json:"alwaysShow"`   // 是否总是显示
	Target       string `json:"target"`       // 打开方式
	ActiveMenu   string `json:"activeMenu"`   // 激活菜单
	Breadcrumb   bool   `json:"breadcrumb"`   // 是否显示面包屑
	Affix        bool   `json:"affix"`        // 是否固定
	FrameSrc     string `json:"frameSrc"`     // iframe地址
	FrameLoading bool   `json:"frameLoading"` // iframe加载状态
	Transition   string `json:"transition"`   // 过渡动画
	Permission   string `json:"permission"`   // 权限标识
}

// ItemTree 菜单树结构(用于前端分配权限)
type ItemTree struct {
	ID       int64       `json:"id"`       // ID
	Name     string      `json:"name"`     // 名称
	Checked  bool        `json:"checked"`  // 是否选中
	Children []*ItemTree `json:"children"` // 子节点
}
