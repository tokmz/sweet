-- Sweet社交电商分销系统 - 用户模块表结构
-- 创建时间: 2023-06-01
-- 表前缀: sw_

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 用户表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user`;
CREATE TABLE `sw_user` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(100) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` varchar(20) NOT NULL DEFAULT '' COMMENT '密码盐',
  `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像URL',
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
  `gender` tinyint(1) NOT NULL DEFAULT 1 COMMENT '性别：1未知 2男 3女',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `user_level` tinyint(1) NOT NULL DEFAULT 0 COMMENT '用户等级',
  `is_distributor` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否分销商：1否 2是',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `register_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '注册IP',
  `register_source` varchar(20) NOT NULL DEFAULT '' COMMENT '注册来源',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_mobile` (`mobile`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_is_distributor` (`is_distributor`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- ----------------------------
-- 用户地址表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user_address`;
CREATE TABLE `sw_user_address` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '地址ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '收货人姓名',
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '收货人手机号',
  `province` varchar(20) NOT NULL DEFAULT '' COMMENT '省份',
  `city` varchar(20) NOT NULL DEFAULT '' COMMENT '城市',
  `district` varchar(20) NOT NULL DEFAULT '' COMMENT '区/县',
  `detail` varchar(255) NOT NULL DEFAULT '' COMMENT '详细地址',
  `post_code` varchar(10) NOT NULL DEFAULT '' COMMENT '邮政编码',
  `is_default` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否默认地址：1否 2是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户地址表';

-- ----------------------------
-- 用户第三方授权表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user_oauth`;
CREATE TABLE `sw_user_oauth` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `oauth_type` varchar(20) NOT NULL COMMENT '第三方类型：weixin/qq/weibo等',
  `oauth_id` varchar(128) NOT NULL COMMENT '第三方唯一ID',
  `union_id` varchar(128) NOT NULL DEFAULT '' COMMENT '第三方联合ID',
  `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT '第三方昵称',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '第三方头像',
  `access_token` varchar(255) NOT NULL DEFAULT '' COMMENT '访问令牌',
  `refresh_token` varchar(255) NOT NULL DEFAULT '' COMMENT '刷新令牌',
  `expires_in` int(11) NOT NULL DEFAULT 0 COMMENT '过期时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_oauth_type_id` (`oauth_type`, `oauth_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户第三方授权表';

-- ----------------------------
-- 用户积分表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user_points`;
CREATE TABLE `sw_user_points` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `points` int(11) NOT NULL DEFAULT 0 COMMENT '积分余额',
  `total_points` int(11) NOT NULL DEFAULT 0 COMMENT '累计获得积分',
  `used_points` int(11) NOT NULL DEFAULT 0 COMMENT '已使用积分',
  `frozen_points` int(11) NOT NULL DEFAULT 0 COMMENT '冻结积分',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分表';

-- ----------------------------
-- 用户积分记录表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user_points_log`;
CREATE TABLE `sw_user_points_log` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `change_points` int(11) NOT NULL COMMENT '变动积分',
  `before_points` int(11) NOT NULL COMMENT '变动前积分',
  `after_points` int(11) NOT NULL COMMENT '变动后积分',
  `change_type` tinyint(1) NOT NULL COMMENT '变动类型：1增加 2减少',
  `source_type` varchar(30) NOT NULL COMMENT '来源类型：register/order/sign/admin等',
  `source_id` varchar(50) NOT NULL DEFAULT '' COMMENT '来源ID',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_source_type` (`source_type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分记录表';

SET FOREIGN_KEY_CHECKS = 1;