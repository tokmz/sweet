/*
 Navicat Premium Dump SQL

 Source Server         : 本地
 Source Server Type    : MySQL
 Source Server Version : 80041 (8.0.41)
 Source Host           : localhost:3306
 Source Schema         : sweet_common

 Target Server Type    : MySQL
 Target Server Version : 80041 (8.0.41)
 File Encoding         : 65001

 Date: 05/04/2025 17:02:02
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '部门ID',
  `pid` bigint DEFAULT 0 COMMENT '父部门ID',
  `name` varchar(64) NOT NULL COMMENT '部门名称',
  `code` varchar(32) DEFAULT NULL COMMENT '部门编码',
  `ancestors` varchar(255) DEFAULT NULL COMMENT '祖级列表',
  `leader` varchar(64) DEFAULT NULL COMMENT '负责人',
  `phone` varchar(16) DEFAULT NULL COMMENT '联系电话',
  `email` varchar(64) DEFAULT NULL COMMENT '邮箱',
  `sort` int DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
  `created_by` bigint DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint DEFAULT NULL COMMENT '更新人',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `version` int DEFAULT '0' COMMENT '乐观锁版本',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='部门表';

-- ----------------------------
-- Table structure for sys_dict_data
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_data`;
CREATE TABLE `sys_dict_data` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典数据ID',
  `dict_type` varchar(64) NOT NULL COMMENT '字典类型',
  `dict_label` varchar(128) NOT NULL COMMENT '字典标签',
  `dict_value` varchar(128) NOT NULL COMMENT '字典值',
  `dict_sort` int DEFAULT '0' COMMENT '排序',
  `css_class` varchar(128) DEFAULT NULL COMMENT '样式属性',
  `list_class` varchar(128) DEFAULT NULL COMMENT '表格回显样式',
  `is_default` tinyint(1) DEFAULT '2' COMMENT '是否默认 1-是 2-否',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint DEFAULT NULL COMMENT '更新人',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `version` int DEFAULT '0' COMMENT '乐观锁版本',
  PRIMARY KEY (`id`),
  KEY `idx_dict_type` (`dict_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='字典数据表';

-- ----------------------------
-- Table structure for sys_dict_type
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_type`;
CREATE TABLE `sys_dict_type` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典类型ID',
  `name` varchar(64) NOT NULL COMMENT '字典名称',
  `type` varchar(64) NOT NULL COMMENT '字典类型',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint DEFAULT NULL COMMENT '更新人',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `version` int DEFAULT '0' COMMENT '乐观锁版本',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_type` (`type`) COMMENT '字典类型唯一索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='字典类型表';

-- ----------------------------
-- Table structure for sys_login_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_login_log`;
CREATE TABLE `sys_login_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '访问ID',
  `username` varchar(64) DEFAULT NULL COMMENT '用户账号',
  `ip` varchar(50) DEFAULT NULL COMMENT '登录IP地址',
  `location` varchar(255) DEFAULT NULL COMMENT '登录地点',
  `browser` varchar(50) DEFAULT NULL COMMENT '浏览器类型',
  `os` varchar(50) DEFAULT NULL COMMENT '操作系统',
  `device` varchar(50) DEFAULT NULL COMMENT '设备',
  `status` tinyint(1) DEFAULT '1' COMMENT '登录状态（1成功 2失败）',
  `msg` varchar(255) DEFAULT NULL COMMENT '提示消息',
  `login_time` datetime DEFAULT NULL COMMENT '登录时间',
  PRIMARY KEY (`id`),
  KEY `idx_username` (`username`),
  KEY `idx_login_time` (`login_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='登录日志表';

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint DEFAULT '0' COMMENT '父菜单ID',
  `name` varchar(64) NOT NULL COMMENT '菜单名称',
  `permission` varchar(128) DEFAULT NULL COMMENT '权限标识',
  `type` tinyint(1) NOT NULL COMMENT '类型 1-目录 2-菜单 3-按钮',
  `path` varchar(255) DEFAULT NULL COMMENT '路由地址',
  `component` varchar(255) DEFAULT NULL COMMENT '组件路径',
  `redirect` varchar(255) DEFAULT NULL COMMENT '重定向地址',
  `icon` varchar(128) DEFAULT NULL COMMENT '菜单图标',
  `sort` int DEFAULT '0' COMMENT '排序',
  `hidden` tinyint(1) DEFAULT '2' COMMENT '是否隐藏 1-是 2-否',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
  `always_show` tinyint(1) DEFAULT '2' COMMENT '是否总是显示 1-是 2-否',
  `keep_alive` tinyint(1) DEFAULT '2' COMMENT '是否缓存 1-是 2-否',
  `target` varchar(20) DEFAULT '_self' COMMENT '打开方式 _self _blank',
  `title` varchar(64) DEFAULT NULL COMMENT '菜单标题',
  `active_menu` varchar(255) DEFAULT NULL COMMENT '激活菜单',
  `breadcrumb` tinyint(1) DEFAULT '1' COMMENT '是否显示面包屑 1-是 2-否',
  `affix` tinyint(1) DEFAULT '2' COMMENT '是否固定 1-是 2-否',
  `frame_src` varchar(255) DEFAULT NULL COMMENT 'iframe地址',
  `frame_loading` tinyint(1) DEFAULT '1' COMMENT 'iframe加载状态 1-是 2-否',
  `transition` varchar(32) DEFAULT NULL COMMENT '过渡动画',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint DEFAULT NULL COMMENT '更新人',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `version` int DEFAULT '0' COMMENT '乐观锁版本',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='菜单权限表';

-- ----------------------------
-- Table structure for sys_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_operation_log`;
CREATE TABLE `sys_operation_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `title` varchar(64) DEFAULT NULL COMMENT '模块标题',
  `business_type` tinyint(1) DEFAULT '4' COMMENT '业务类型（1新增 2修改 3删除 4其它）',
  `method` varchar(128) DEFAULT NULL COMMENT '方法名称',
  `request_method` varchar(10) DEFAULT NULL COMMENT '请求方式',
  `operator_type` tinyint(1) DEFAULT '3' COMMENT '操作类别（1后台用户 2手机端用户 3其它）',
  `user_id` bigint DEFAULT NULL COMMENT '操作人ID',
  `username` varchar(64) DEFAULT NULL COMMENT '操作人名称',
  `url` varchar(255) DEFAULT NULL COMMENT '请求URL',
  `ip` varchar(50) DEFAULT NULL COMMENT '操作地址',
  `location` varchar(64) DEFAULT NULL COMMENT '操作地点',
  `browser` varchar(50) DEFAULT NULL COMMENT '浏览器类型',
  `os` varchar(50) DEFAULT NULL COMMENT '操作系统',
  `device` varchar(50) DEFAULT NULL COMMENT '设备',
  `param` text COMMENT '请求参数',
  `result` text COMMENT '返回结果',
  `status` tinyint(1) DEFAULT '1' COMMENT '操作状态（1正常 2异常）',
  `error_stack` text COMMENT '错误消息',
  `cost_time` bigint DEFAULT NULL COMMENT '消耗时间',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='操作日志表';

-- ----------------------------
-- Table structure for sys_post
-- ----------------------------
DROP TABLE IF EXISTS `sys_post`;
CREATE TABLE `sys_post` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '岗位ID',
  `name` varchar(64) NOT NULL COMMENT '岗位名称',
  `code` varchar(32) NOT NULL COMMENT '岗位编码',
  `sort` int DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint DEFAULT NULL COMMENT '更新人',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `version` int DEFAULT '0' COMMENT '乐观锁版本',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_code` (`code`) COMMENT '岗位编码唯一索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='岗位表';

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(64) NOT NULL COMMENT '角色名称',
  `code` varchar(32) NOT NULL COMMENT '角色编码',
  `sort` int DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
  `data_scope` tinyint(1) DEFAULT '1' COMMENT '数据范围 1-全部 2-自定义 3-本部门 4-本部门及子部门 5-仅本人',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint DEFAULT NULL COMMENT '更新人',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `version` int DEFAULT '0' COMMENT '乐观锁版本',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_code` (`code`) COMMENT '角色编码唯一索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色表';

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `menu_id` bigint NOT NULL COMMENT '菜单ID',
  `created_by` bigint DEFAULT NULL COMMENT '创建人',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_role_menu` (`role_id`,`menu_id`) COMMENT '角色菜单唯一索引',
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色菜单关联表';

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `password` varchar(128) NOT NULL COMMENT '密码',
  `salt` varchar(16) DEFAULT NULL COMMENT '加密盐',
  `real_name` varchar(64) DEFAULT NULL COMMENT '真实姓名',
  `nick_name` varchar(64) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像地址',
  `email` varchar(128) DEFAULT NULL COMMENT '邮箱',
  `mobile` varchar(20) DEFAULT NULL COMMENT '手机号',
  `gender` tinyint(1) DEFAULT '3' COMMENT '性别 1-男 2-女 3-未知',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `post_id` bigint DEFAULT NULL COMMENT '岗位ID',
  `role_id` bigint DEFAULT NULL COMMENT '角色ID',
  `is_admin` tinyint(1) DEFAULT '2' COMMENT '是否管理员 1-是 2-否',
  `login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态 1-正常 2-禁用',
  `created_by` bigint DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint DEFAULT NULL COMMENT '更新人',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `version` int DEFAULT '0' COMMENT '乐观锁版本',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_username` (`username`) COMMENT '用户名唯一索引',
  KEY `idx_mobile` (`mobile`),
  KEY `idx_email` (`email`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统用户表';

SET FOREIGN_KEY_CHECKS = 1;
