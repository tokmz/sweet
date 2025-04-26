-- Sweet社交电商分销系统 - 分销模块表结构
-- 创建时间: 2023-06-01
-- 表前缀: sw_

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 分销商表
-- ----------------------------
DROP TABLE IF EXISTS `sw_distributor`;
CREATE TABLE `sw_distributor` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `level` tinyint(1) NOT NULL DEFAULT 1 COMMENT '分销商等级',
  `parent_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '上级分销商ID',
  `parent_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '上级用户ID',
  `team_count` int(11) NOT NULL DEFAULT 0 COMMENT '团队人数',
  `direct_count` int(11) NOT NULL DEFAULT 0 COMMENT '直推人数',
  `total_commission` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '累计佣金',
  `available_commission` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '可用佣金',
  `frozen_commission` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '冻结佣金',
  `withdrawn_commission` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '已提现佣金',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `apply_time` datetime DEFAULT NULL COMMENT '申请时间',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `audit_remark` varchar(255) NOT NULL DEFAULT '' COMMENT '审核备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_parent_user_id` (`parent_user_id`),
  KEY `idx_level` (`level`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销商表';

-- ----------------------------
-- 分销关系表
-- ----------------------------
DROP TABLE IF EXISTS `sw_distribution_relation`;
CREATE TABLE `sw_distribution_relation` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `parent_user_id` bigint(20) UNSIGNED NOT NULL COMMENT '上级用户ID',
  `level` tinyint(1) NOT NULL COMMENT '关系层级：1直接上级 2间接上级',
  `source` varchar(20) NOT NULL DEFAULT '' COMMENT '绑定来源：register/order/admin',
  `is_locked` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否锁定：1否 2是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_parent_level` (`user_id`, `parent_user_id`, `level`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_user_id` (`parent_user_id`),
  KEY `idx_level` (`level`),
  KEY `idx_is_locked` (`is_locked`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销关系表';

-- ----------------------------
-- 分销等级表
-- ----------------------------
DROP TABLE IF EXISTS `sw_distributor_level`;
CREATE TABLE `sw_distributor_level` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) NOT NULL COMMENT '等级名称',
  `level` tinyint(1) NOT NULL COMMENT '等级值',
  `icon` varchar(255) NOT NULL DEFAULT '' COMMENT '等级图标',
  `upgrade_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '升级方式：1自动升级 2人工审核',
  `upgrade_condition` text COMMENT '升级条件，JSON格式',
  `commission_rate_level1` decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '一级佣金比例',
  `commission_rate_level2` decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '二级佣金比例',
  `commission_rate_level3` decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '三级佣金比例',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '等级描述',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_level` (`level`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销等级表';

-- ----------------------------
-- 分销申请表
-- ----------------------------
DROP TABLE IF EXISTS `sw_distributor_apply`;
CREATE TABLE `sw_distributor_apply` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '姓名',
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
  `id_card` varchar(20) NOT NULL DEFAULT '' COMMENT '身份证号',
  `id_card_front` varchar(255) NOT NULL DEFAULT '' COMMENT '身份证正面',
  `id_card_back` varchar(255) NOT NULL DEFAULT '' COMMENT '身份证反面',
  `reason` varchar(255) NOT NULL DEFAULT '' COMMENT '申请理由',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1待审核 2已通过 3已拒绝',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `audit_remark` varchar(255) NOT NULL DEFAULT '' COMMENT '审核备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销申请表';

-- ----------------------------
-- 分销佣金记录表
-- ----------------------------
DROP TABLE IF EXISTS `sw_commission_log`;
CREATE TABLE `sw_commission_log` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `order_id` bigint(20) UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no` varchar(50) NOT NULL COMMENT '订单编号',
  `order_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '订单金额',
  `commission_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '佣金金额',
  `commission_rate` decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '佣金比例',
  `level` tinyint(1) NOT NULL COMMENT '佣金层级：1一级 2二级 3三级',
  `buyer_user_id` bigint(20) UNSIGNED NOT NULL COMMENT '购买用户ID',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1未结算 2已结算 3已取消',
  `settle_time` datetime DEFAULT NULL COMMENT '结算时间',
  `cancel_time` datetime DEFAULT NULL COMMENT '取消时间',
  `cancel_reason` varchar(255) NOT NULL DEFAULT '' COMMENT '取消原因',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_buyer_user_id` (`buyer_user_id`),
  KEY `idx_level` (`level`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销佣金记录表';

-- ----------------------------
-- 分销提现表
-- ----------------------------
DROP TABLE IF EXISTS `sw_commission_withdraw`;
CREATE TABLE `sw_commission_withdraw` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `withdraw_no` varchar(50) NOT NULL COMMENT '提现单号',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '提现金额',
  `fee` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '手续费',
  `actual_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '实际到账金额',
  `withdraw_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '提现方式：1微信 2支付宝 3银行卡',
  `account_name` varchar(50) NOT NULL DEFAULT '' COMMENT '账户姓名',
  `account_no` varchar(50) NOT NULL DEFAULT '' COMMENT '账户号码',
  `bank_name` varchar(50) NOT NULL DEFAULT '' COMMENT '银行名称',
  `bank_branch` varchar(100) NOT NULL DEFAULT '' COMMENT '开户支行',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1待审核 2审核通过 3审核拒绝 4已打款 5打款失败',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `audit_remark` varchar(255) NOT NULL DEFAULT '' COMMENT '审核备注',
  `payment_time` datetime DEFAULT NULL COMMENT '打款时间',
  `payment_remark` varchar(255) NOT NULL DEFAULT '' COMMENT '打款备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_withdraw_no` (`withdraw_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销提现表';

-- ----------------------------
-- 分销配置表
-- ----------------------------
DROP TABLE IF EXISTS `sw_distribution_config`;
CREATE TABLE `sw_distribution_config` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `key` varchar(50) NOT NULL COMMENT '配置键',
  `value` text NOT NULL COMMENT '配置值',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '配置描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分销配置表';

SET FOREIGN_KEY_CHECKS = 1;