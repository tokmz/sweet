/*
 Navicat Premium Dump SQL

 Source Server         : 本地
 Source Server Type    : MySQL
 Source Server Version : 80041 (8.0.41)
 Source Host           : localhost:3306
 Source Schema         : sweet

 Target Server Type    : MySQL
 Target Server Version : 80041 (8.0.41)
 File Encoding         : 65001

 Date: 30/07/2025 23:48:11
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sw_sys_dept
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_dept`;
CREATE TABLE `sw_sys_dept` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '部门ID',
  `parent_id` bigint unsigned DEFAULT '0' COMMENT '父部门ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '部门名称',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '部门编码',
  `sort` int unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态: 1-启用, 2-禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_parent_status` (`parent_id`,`status`),
  KEY `idx_status_sort` (`status`,`sort`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

-- ----------------------------
-- Table structure for sw_sys_post
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_post`;
CREATE TABLE `sw_sys_post` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '岗位ID',
  `dept_id` bigint unsigned NOT NULL COMMENT '所属部门',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '岗位名称',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '岗位编码',
  `sort` int unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态: 1-启用, 2-禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status_sort` (`status`,`sort`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_dept_status` (`dept_id`,`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_post_dept` FOREIGN KEY (`dept_id`) REFERENCES `sw_sys_dept` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='岗位表';

-- ----------------------------
-- Table structure for sw_sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_role`;
CREATE TABLE `sw_sys_role` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色标识',
  `sort` int unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `is_system` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否系统内置：1=是，2否',
  `is_super` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '是否超级管理员：1=是，2否',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态：1=正常，2=禁用',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_is_system` (`is_system`),
  KEY `idx_is_super` (`is_super`),
  KEY `idx_status_sort` (`status`,`sort`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统角色表';

-- ----------------------------
-- Table structure for sw_sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_menu`;
CREATE TABLE `sw_sys_menu` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '组件名称/路由名称',
  `title` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '菜单名称',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父菜单ID',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '路由地址',
  `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '组件地址',
  `menu_type` tinyint(1) DEFAULT '1' COMMENT '菜单类型（1 目录 2 菜单 3 按钮）',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '菜单状态(1 正常 2 停用)',
  `perms` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '权限标识',
  `icon` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '菜单图标',
  `order` int DEFAULT '0' COMMENT '显示顺序 从大到小',
  `remark` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建者',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新者',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_menu_type` (`menu_type`),
  KEY `idx_status` (`status`),
  KEY `idx_parent_status` (`parent_id`,`status`),
  KEY `idx_menu_type_status` (`menu_type`,`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_menu_parent` FOREIGN KEY (`parent_id`) REFERENCES `sw_sys_menu` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单表';

-- ----------------------------
-- Table structure for sw_sys_menu_config
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_menu_config`;
CREATE TABLE `sw_sys_menu_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `menu_id` bigint unsigned NOT NULL COMMENT '菜单ID',
  `query` json DEFAULT NULL COMMENT '路由参数',
  `is_frame` tinyint(1) DEFAULT '2' COMMENT '是否外联（1 是 2 否）',
  `show_badge` tinyint(1) DEFAULT '2' COMMENT '是否显示徽章（1 是 2 否）',
  `show_text_badge` varchar(24) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '文本徽章内容',
  `is_hide` tinyint(1) DEFAULT '2' COMMENT '是否在菜单中隐藏（1 是 2 否）',
  `is_hide_tab` tinyint(1) DEFAULT '2' COMMENT '是否在标签页中隐藏（1 是 2 否）',
  `link` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '外链地址',
  `is_iframe` tinyint(1) DEFAULT '2' COMMENT '是否iframe（1 是 2 否）',
  `keep_alive` tinyint(1) DEFAULT '2' COMMENT '是否缓存页面（1 是 2否）',
  `fixed_tab` tinyint(1) DEFAULT '2' COMMENT '是否固定标签页（1 是 2 否）',
  `is_first_level` tinyint(1) DEFAULT '2' COMMENT '是否为一级菜单',
  `active_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '激活菜单路径',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_menu_id` (`menu_id`),
  KEY `idx_is_hide` (`is_hide`),
  KEY `idx_keep_alive` (`keep_alive`),
  CONSTRAINT `fk_menu_config_menu` FOREIGN KEY (`menu_id`) REFERENCES `sw_sys_menu` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单配置表';

-- ----------------------------
-- Table structure for sw_sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_role_menu`;
CREATE TABLE `sw_sys_role_menu` (
  `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
  `menu_id` bigint unsigned NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`role_id`,`menu_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_menu_id` (`menu_id`),
  CONSTRAINT `fk_role_menu_role` FOREIGN KEY (`role_id`) REFERENCES `sw_sys_role` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_role_menu_menu` FOREIGN KEY (`menu_id`) REFERENCES `sw_sys_menu` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色菜单关联表';

-- ----------------------------
-- Table structure for sw_sys_api_group
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_api_group`;
CREATE TABLE `sw_sys_api_group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '分组ID',
  `name` varchar(50) NOT NULL COMMENT '分组名称',
  `code` varchar(50) NOT NULL COMMENT '分组编码',
  `description` varchar(200) DEFAULT NULL COMMENT '分组描述',
  `sort` int NOT NULL DEFAULT '0' COMMENT '显示顺序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '分组状态（1正常 2停用）',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建者',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新者',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='API分组表';

-- ----------------------------
-- Table structure for sw_sys_api
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_api`;
CREATE TABLE `sw_sys_api` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'API ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'API名称',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'API路径',
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'HTTP方法',
  `group` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'API分组',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'API描述',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT 'API状态（1正常 2停用）',
  `is_auth` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否需要认证（1需要 2不需要）',
  `create_by` bigint unsigned DEFAULT NULL COMMENT '创建者',
  `update_by` bigint unsigned DEFAULT NULL COMMENT '更新者',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_path_method` (`path`,`method`),
  KEY `idx_group` (`group`),
  KEY `idx_status` (`status`),
  KEY `idx_is_auth` (`is_auth`),
  KEY `idx_group_status` (`group`,`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_api_group` FOREIGN KEY (`group`) REFERENCES `sw_sys_api_group` (`code`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='API接口表';

-- 角色API关联表----------------------------
-- Table structure for sw_sys_role_api
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_role_api`;
CREATE TABLE `sw_sys_role_api` (
  `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
  `api_id` bigint unsigned NOT NULL COMMENT 'API ID',
  PRIMARY KEY (`role_id`,`api_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_api_id` (`api_id`),
  CONSTRAINT `fk_role_api_role` FOREIGN KEY (`role_id`) REFERENCES `sw_sys_role` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_role_api_api` FOREIGN KEY (`api_id`) REFERENCES `sw_sys_api` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色API关联表';

-- ----------------------------
-- Table structure for sw_sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_user`;
CREATE TABLE `sw_sys_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
  `username` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录用户名',
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录密码',
  `salt` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码盐',
  `realname` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '真实姓名',
  `nickname` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '头像',
  `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '邮箱',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '手机号',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态：1=正常，2=禁用',
  `role_id` bigint unsigned DEFAULT NULL COMMENT '角色ID',
  `dept_id` bigint unsigned DEFAULT NULL COMMENT '部门ID',
  `post_id` bigint unsigned DEFAULT NULL COMMENT '岗位ID',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`) COMMENT '邮箱唯一约束（非空时）',
  UNIQUE KEY `uk_phone` (`phone`) COMMENT '手机号唯一约束（非空时）',
  KEY `idx_status` (`status`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_dept` (`dept_id`),
  KEY `idx_post` (`post_id`),
  KEY `idx_status_dept` (`status`,`dept_id`),
  KEY `idx_status_role` (`status`,`role_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_user_role` FOREIGN KEY (`role_id`) REFERENCES `sw_sys_role` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_user_dept` FOREIGN KEY (`dept_id`) REFERENCES `sw_sys_dept` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_user_post` FOREIGN KEY (`post_id`) REFERENCES `sw_sys_post` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统管理员表';

-- ----------------------------
-- Table structure for sw_sys_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_operation_log`;
CREATE TABLE `sw_sys_operation_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '操作用户ID',
  `module` varchar(50) NOT NULL COMMENT '操作模块',
  `operation` varchar(50) NOT NULL COMMENT '操作类型',
  `method` varchar(10) NOT NULL COMMENT 'HTTP方法',
  `url` varchar(255) NOT NULL COMMENT '请求URL',
  `ip` varchar(45) NOT NULL COMMENT '操作IP',
  `location` varchar(100) DEFAULT NULL COMMENT 'IP归属地',
  `user_agent` varchar(500) DEFAULT NULL COMMENT '用户代理',
  `request_params` text COMMENT '请求参数',
  `response_data` text COMMENT '响应数据',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '操作状态（1成功 2失败）',
  `error_msg` varchar(500) DEFAULT NULL COMMENT '错误信息',
  `cost_time` int unsigned DEFAULT NULL COMMENT '耗时（毫秒）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_module` (`module`),
  KEY `idx_operation` (`operation`),
  KEY `idx_status` (`status`),
  KEY `idx_ip` (`ip`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_user_module` (`user_id`,`module`),
  KEY `idx_module_operation` (`module`,`operation`),
  CONSTRAINT `fk_operation_log_user` FOREIGN KEY (`user_id`) REFERENCES `sw_sys_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='系统操作日志表';

-- ----------------------------
-- Table structure for sw_sys_login_log
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_login_log`;
CREATE TABLE `sw_sys_login_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '用户ID',
  `username` varchar(64) NOT NULL COMMENT '登录用户名',
  `login_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '登录类型（1账号密码 2手机验证码 3邮箱验证码 4第三方登录 5微信 6QQ 7支付宝）',
  `client_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '客户端类型（1Web 2移动端 3小程序 4API 5管理后台）',
  `ip` varchar(45) NOT NULL COMMENT '登录IP',
  `location` varchar(64) DEFAULT NULL COMMENT 'IP归属地',
  `user_agent` text DEFAULT NULL COMMENT '用户代理',
  `device_info` varchar(64) DEFAULT NULL COMMENT '设备信息',
  `browser` varchar(100) DEFAULT NULL COMMENT '浏览器',
  `os` varchar(100) DEFAULT NULL COMMENT '操作系统',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '登录状态（1成功 2失败 3异常）',
  `fail_reason` varchar(64) DEFAULT NULL COMMENT '失败原因',
  `session_id` varchar(128) DEFAULT NULL COMMENT '会话ID',
  `login_duration` int unsigned DEFAULT NULL COMMENT '登录持续时间（秒）',
  `logout_type` tinyint(1) DEFAULT NULL COMMENT '退出类型（1主动退出 2超时退出 3强制退出）',
  `risk_level` tinyint(1) NOT NULL DEFAULT '1' COMMENT '风险等级（1低风险 2中风险 3高风险）',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除（0否 1是）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_ip` (`ip`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_login_type` (`login_type`),
  KEY `idx_client_type` (`client_type`),
  KEY `idx_risk_level` (`risk_level`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_composite_query` (`user_id`, `status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='系统登录日志表';

-- ----------------------------
-- Table structure for sw_sys_file
-- ----------------------------
DROP TABLE IF EXISTS `sw_sys_file`;
CREATE TABLE `sw_sys_file` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名称',
  `original_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '原始文件名',
  `file_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件路径',
  `file_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件访问URL',
  `file_size` bigint unsigned NOT NULL COMMENT '文件大小（字节）',
  `file_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件类型/MIME类型',
  `file_ext` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件扩展名',
  `md5` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件MD5值',
  `storage_type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '存储类型：1=本地存储，2=阿里云OSS，3=腾讯云COS，4=七牛云',
  `upload_user_id` bigint unsigned DEFAULT NULL COMMENT '上传用户ID',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态：1=正常，2=禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_md5_size` (`md5`, `file_size`),
  KEY `idx_file_type` (`file_type`),
  KEY `idx_upload_user` (`upload_user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_file_upload_user` FOREIGN KEY (`upload_user_id`) REFERENCES `sw_sys_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统文件表';

SET FOREIGN_KEY_CHECKS = 1;
