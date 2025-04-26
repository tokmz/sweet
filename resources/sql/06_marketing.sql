-- Sweet社交电商分销系统 - 营销模块表结构
-- 创建时间: 2023-06-01
-- 表前缀: sw_

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 优惠券表
-- ----------------------------
DROP TABLE IF EXISTS `sw_coupon`;
CREATE TABLE `sw_coupon` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) NOT NULL COMMENT '优惠券名称',
  `type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '优惠券类型：1满减券 2折扣券 3无门槛券',
  `code` varchar(20) NOT NULL DEFAULT '' COMMENT '优惠码',
  `value` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '面值/折扣',
  `min_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '最低消费金额',
  `total` int(11) NOT NULL DEFAULT 0 COMMENT '发行数量，0表示不限量',
  `used` int(11) NOT NULL DEFAULT 0 COMMENT '已使用数量',
  `limit_per_user` int(11) NOT NULL DEFAULT 1 COMMENT '每人限领数量',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `valid_days` int(11) NOT NULL DEFAULT 0 COMMENT '有效天数，0表示不限制',
  `product_ids` text COMMENT '适用商品ID，为空表示全场通用',
  `category_ids` text COMMENT '适用分类ID，为空表示全场通用',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '使用说明',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `is_visible` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否可见：1是 2否',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_end_time` (`end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券表';

-- ----------------------------
-- 用户优惠券表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user_coupon`;
CREATE TABLE `sw_user_coupon` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `coupon_id` bigint(20) UNSIGNED NOT NULL COMMENT '优惠券ID',
  `coupon_name` varchar(50) NOT NULL COMMENT '优惠券名称',
  `coupon_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '优惠券类型：1满减券 2折扣券 3无门槛券',
  `value` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '面值/折扣',
  `min_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '最低消费金额',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1未使用 2已使用 3已过期',
  `use_time` datetime DEFAULT NULL COMMENT '使用时间',
  `order_id` bigint(20) UNSIGNED DEFAULT NULL COMMENT '订单ID',
  `order_no` varchar(50) DEFAULT NULL COMMENT '订单编号',
  `source` varchar(20) NOT NULL DEFAULT '' COMMENT '获取来源：receive/send/register/activity',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_coupon_id` (`coupon_id`),
  KEY `idx_status` (`status`),
  KEY `idx_end_time` (`end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户优惠券表';

-- ----------------------------
-- 秒杀活动表
-- ----------------------------
DROP TABLE IF EXISTS `sw_seckill_activity`;
CREATE TABLE `sw_seckill_activity` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) NOT NULL COMMENT '活动名称',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '活动描述',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1未开始 2进行中 3已结束 4已取消',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `banner` varchar(255) NOT NULL DEFAULT '' COMMENT '活动banner',
  `rules` text COMMENT '活动规则',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_end_time` (`end_time`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀活动表';

-- ----------------------------
-- 秒杀商品表
-- ----------------------------
DROP TABLE IF EXISTS `sw_seckill_product`;
CREATE TABLE `sw_seckill_product` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `activity_id` bigint(20) UNSIGNED NOT NULL COMMENT '活动ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `sku_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'SKU ID',
  `seckill_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '秒杀价格',
  `original_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '原价',
  `stock` int(11) NOT NULL DEFAULT 0 COMMENT '秒杀库存',
  `sales` int(11) NOT NULL DEFAULT 0 COMMENT '销量',
  `limit_num` int(11) NOT NULL DEFAULT 1 COMMENT '限购数量',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_activity_id` (`activity_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_sku_id` (`sku_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀商品表';

-- ----------------------------
-- 拼团活动表
-- ----------------------------
DROP TABLE IF EXISTS `sw_group_activity`;
CREATE TABLE `sw_group_activity` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) NOT NULL COMMENT '活动名称',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '活动描述',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `group_num` int(11) NOT NULL DEFAULT 2 COMMENT '成团人数',
  `duration` int(11) NOT NULL DEFAULT 24 COMMENT '团有效时间(小时)',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1未开始 2进行中 3已结束 4已取消',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `banner` varchar(255) NOT NULL DEFAULT '' COMMENT '活动banner',
  `rules` text COMMENT '活动规则',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_end_time` (`end_time`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='拼团活动表';

-- ----------------------------
-- 拼团商品表
-- ----------------------------
DROP TABLE IF EXISTS `sw_group_product`;
CREATE TABLE `sw_group_product` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `activity_id` bigint(20) UNSIGNED NOT NULL COMMENT '活动ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `sku_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'SKU ID',
  `group_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '拼团价格',
  `original_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '原价',
  `stock` int(11) NOT NULL DEFAULT 0 COMMENT '拼团库存',
  `sales` int(11) NOT NULL DEFAULT 0 COMMENT '销量',
  `limit_num` int(11) NOT NULL DEFAULT 1 COMMENT '限购数量',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_activity_id` (`activity_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_sku_id` (`sku_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='拼团商品表';

-- ----------------------------
-- 拼团记录表
-- ----------------------------
DROP TABLE IF EXISTS `sw_group_record`;
CREATE TABLE `sw_group_record` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `group_no` varchar(50) NOT NULL COMMENT '团编号',
  `activity_id` bigint(20) UNSIGNED NOT NULL COMMENT '活动ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `sku_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'SKU ID',
  `leader_user_id` bigint(20) UNSIGNED NOT NULL COMMENT '团长用户ID',
  `leader_order_id` bigint(20) UNSIGNED NOT NULL COMMENT '团长订单ID',
  `leader_order_no` varchar(50) NOT NULL COMMENT '团长订单编号',
  `required_num` int(11) NOT NULL DEFAULT 2 COMMENT '成团人数',
  `current_num` int(11) NOT NULL DEFAULT 1 COMMENT '当前人数',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1进行中 2已成团 3未成团',
  `expire_time` datetime NOT NULL COMMENT '过期时间',
  `success_time` datetime DEFAULT NULL COMMENT '成团时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_group_no` (`group_no`),
  KEY `idx_activity_id` (`activity_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_sku_id` (`sku_id`),
  KEY `idx_leader_user_id` (`leader_user_id`),
  KEY `idx_leader_order_id` (`leader_order_id`),
  KEY `idx_status` (`status`),
  KEY `idx_expire_time` (`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='拼团记录表';

-- ----------------------------
-- 拼团成员表
-- ----------------------------
DROP TABLE IF EXISTS `sw_group_member`;
CREATE TABLE `sw_group_member` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `group_id` bigint(20) UNSIGNED NOT NULL COMMENT '拼团ID',
  `group_no` varchar(50) NOT NULL COMMENT '团编号',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `order_id` bigint(20) UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no` varchar(50) NOT NULL COMMENT '订单编号',
  `is_leader` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否团长：1否 2是',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1待支付 2已支付 3已取消',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_group_no` (`group_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_is_leader` (`is_leader`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='拼团成员表';

SET FOREIGN_KEY_CHECKS = 1;