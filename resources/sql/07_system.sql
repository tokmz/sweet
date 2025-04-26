-- Sweet社交电商分销系统 - 系统模块表结构
-- 创建时间: 2023-06-01
-- 表前缀: sw_

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 管理员表
-- ----------------------------
DROP TABLE IF EXISTS `sw_admin`;
CREATE TABLE `sw_admin` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `salt` varchar(20) NOT NULL COMMENT '密码盐',
  `nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `role_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '角色ID',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_status` (`status`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

-- ----------------------------
-- 角色表
-- ----------------------------
DROP TABLE IF EXISTS `sw_role`;
CREATE TABLE `sw_role` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '角色描述',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- ----------------------------
-- 菜单表
-- ----------------------------
DROP TABLE IF EXISTS `sw_menu`;
CREATE TABLE `sw_menu` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父菜单ID',
  `name` varchar(50) NOT NULL COMMENT '菜单名称',
  `path` varchar(100) NOT NULL DEFAULT '' COMMENT '路由路径',
  `component` varchar(100) NOT NULL DEFAULT '' COMMENT '组件路径',
  `redirect` varchar(100) NOT NULL DEFAULT '' COMMENT '重定向地址',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '菜单标题',
  `icon` varchar(50) NOT NULL DEFAULT '' COMMENT '图标',
  `permission` varchar(100) NOT NULL DEFAULT '' COMMENT '权限标识',
  `type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '类型：1目录 2菜单 3按钮 4接口',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `is_hidden` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否隐藏：1否 2是',
  `is_keep_alive` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否缓存：1否 2是',
  `is_affix` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否固定在标签栏：1否 2是',
  `is_iframe` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否为外链：1否 2是',
  `frame_src` varchar(255) NOT NULL DEFAULT '' COMMENT '内嵌iframe地址',
  `transition_name` varchar(50) NOT NULL DEFAULT '' COMMENT '路由切换动画',
  `ignore_route` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否忽略路由：1否 2是',
  `is_root` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否为根节点：1否 2是',
  `active_menu` varchar(100) NOT NULL DEFAULT '' COMMENT '高亮的菜单',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';

-- ----------------------------
-- 角色菜单关联表
-- ----------------------------
DROP TABLE IF EXISTS `sw_role_menu`;
CREATE TABLE `sw_role_menu` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint(20) UNSIGNED NOT NULL COMMENT '角色ID',
  `menu_id` bigint(20) UNSIGNED NOT NULL COMMENT '菜单ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_menu` (`role_id`, `menu_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单关联表';

-- ----------------------------
-- 系统配置表
-- ----------------------------
DROP TABLE IF EXISTS `sw_system_config`;
CREATE TABLE `sw_system_config` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `group` varchar(50) NOT NULL DEFAULT 'default' COMMENT '分组',
  `name` varchar(50) NOT NULL COMMENT '配置名称',
  `key` varchar(50) NOT NULL COMMENT '配置键',
  `value` text NOT NULL COMMENT '配置值',
  `type` varchar(20) NOT NULL DEFAULT 'string' COMMENT '类型：string/number/boolean/array/object',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '配置描述',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_key` (`key`),
  KEY `idx_group` (`group`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- ----------------------------
-- 操作日志表
-- ----------------------------
DROP TABLE IF EXISTS `sw_operation_log`;
CREATE TABLE `sw_operation_log` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `admin_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
  `admin_name` varchar(50) NOT NULL DEFAULT '' COMMENT '管理员名称',
  `module` varchar(50) NOT NULL DEFAULT '' COMMENT '模块',
  `action` varchar(50) NOT NULL DEFAULT '' COMMENT '操作',
  `method` varchar(10) NOT NULL DEFAULT '' COMMENT '请求方法',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '请求URL',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `user_agent` varchar(255) NOT NULL DEFAULT '' COMMENT '用户代理',
  `request_data` text COMMENT '请求数据',
  `response_data` text COMMENT '响应数据',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1成功 2失败',
  `error_message` varchar(255) NOT NULL DEFAULT '' COMMENT '错误信息',
  `execution_time` int(11) NOT NULL DEFAULT 0 COMMENT '执行时间(ms)',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_module` (`module`),
  KEY `idx_action` (`action`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- ----------------------------
-- 登录日志表
-- ----------------------------
DROP TABLE IF EXISTS `sw_login_log`;
CREATE TABLE `sw_login_log` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `admin_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
  `admin_name` varchar(50) NOT NULL DEFAULT '' COMMENT '管理员名称',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `user_agent` varchar(255) NOT NULL DEFAULT '' COMMENT '用户代理',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1成功 2失败',
  `error_message` varchar(255) NOT NULL DEFAULT '' COMMENT '错误信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';

-- ----------------------------
-- 地区表
-- ----------------------------
DROP TABLE IF EXISTS `sw_region`;
CREATE TABLE `sw_region` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `parent_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父ID',
  `name` varchar(50) NOT NULL COMMENT '名称',
  `code` varchar(20) NOT NULL COMMENT '编码',
  `level` tinyint(1) NOT NULL DEFAULT 1 COMMENT '级别：1省 2市 3区县',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_code` (`code`),
  KEY `idx_level` (`level`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='地区表';

SET FOREIGN_KEY_CHECKS = 1;