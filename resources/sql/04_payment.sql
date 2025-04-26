-- Sweet社交电商分销系统 - 支付模块表结构
-- 创建时间: 2023-06-01
-- 表前缀: sw_

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 支付记录表
-- ----------------------------
DROP TABLE IF EXISTS `sw_payment_record`;
CREATE TABLE `sw_payment_record` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `payment_no` varchar(100) NOT NULL COMMENT '支付单号',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `order_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单ID',
  `order_no` varchar(50) NOT NULL DEFAULT '' COMMENT '订单编号',
  `business_type` varchar(30) NOT NULL COMMENT '业务类型：order/recharge/vip等',
  `payment_type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '支付方式：0未选择 1微信支付 2支付宝 3余额支付',
  `payment_method` varchar(20) NOT NULL DEFAULT '' COMMENT '支付方法：wechat_jsapi/wechat_h5/alipay_app等',
  `transaction_id` varchar(100) NOT NULL DEFAULT '' COMMENT '第三方支付流水号',
  `amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '支付金额',
  `currency` varchar(10) NOT NULL DEFAULT 'CNY' COMMENT '货币类型',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '支付状态：1未支付 2已支付 3支付失败 4已退款 5部分退款',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `expire_time` datetime DEFAULT NULL COMMENT '过期时间',
  `client_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '客户端IP',
  `device` varchar(50) NOT NULL DEFAULT '' COMMENT '设备信息',
  `payment_params` text COMMENT '支付参数',
  `callback_data` text COMMENT '回调数据',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_payment_no` (`payment_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_status` (`status`),
  KEY `idx_business_type` (`business_type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付记录表';

-- ----------------------------
-- 退款记录表
-- ----------------------------
DROP TABLE IF EXISTS `sw_refund_record`;
CREATE TABLE `sw_refund_record` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `refund_no` varchar(100) NOT NULL COMMENT '退款单号',
  `payment_no` varchar(100) NOT NULL COMMENT '支付单号',
  `transaction_id` varchar(100) NOT NULL DEFAULT '' COMMENT '第三方支付流水号',
  `refund_id` varchar(100) NOT NULL DEFAULT '' COMMENT '第三方退款单号',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `order_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单ID',
  `order_no` varchar(50) NOT NULL DEFAULT '' COMMENT '订单编号',
  `business_type` varchar(30) NOT NULL COMMENT '业务类型：order/recharge/vip等',
  `payment_type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '支付方式：0未选择 1微信支付 2支付宝 3余额支付',
  `amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '退款金额',
  `total_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '订单总金额',
  `reason` varchar(255) NOT NULL DEFAULT '' COMMENT '退款原因',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '退款状态：1处理中 2成功 3失败',
  `refund_time` datetime DEFAULT NULL COMMENT '退款时间',
  `callback_data` text COMMENT '回调数据',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_refund_no` (`refund_no`),
  KEY `idx_payment_no` (`payment_no`),
  KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_refund_id` (`refund_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款记录表';

-- ----------------------------
-- 用户余额表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user_balance`;
CREATE TABLE `sw_user_balance` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `balance` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '余额',
  `frozen_balance` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '冻结余额',
  `total_recharge` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '累计充值',
  `total_consume` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '累计消费',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户余额表';

-- ----------------------------
-- 余额变动记录表
-- ----------------------------
DROP TABLE IF EXISTS `sw_balance_log`;
CREATE TABLE `sw_balance_log` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `change_amount` decimal(10,2) NOT NULL COMMENT '变动金额',
  `before_balance` decimal(10,2) NOT NULL COMMENT '变动前余额',
  `after_balance` decimal(10,2) NOT NULL COMMENT '变动后余额',
  `change_type` tinyint(1) NOT NULL COMMENT '变动类型：1增加 2减少',
  `source_type` varchar(30) NOT NULL COMMENT '来源类型：recharge/order/refund/commission/admin等',
  `source_id` varchar(50) NOT NULL DEFAULT '' COMMENT '来源ID',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_source_type` (`source_type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='余额变动记录表';

-- ----------------------------
-- 充值套餐表
-- ----------------------------
DROP TABLE IF EXISTS `sw_recharge_plan`;
CREATE TABLE `sw_recharge_plan` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) NOT NULL COMMENT '套餐名称',
  `amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '充值金额',
  `gift_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '赠送金额',
  `total_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '总金额',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='充值套餐表';

-- ----------------------------
-- 充值订单表
-- ----------------------------
DROP TABLE IF EXISTS `sw_recharge_order`;
CREATE TABLE `sw_recharge_order` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `order_no` varchar(50) NOT NULL COMMENT '订单编号',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `plan_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '套餐ID',
  `amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '充值金额',
  `gift_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '赠送金额',
  `total_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '总金额',
  `pay_type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '支付方式：0未选择 1微信支付 2支付宝',
  `payment_no` varchar(100) NOT NULL DEFAULT '' COMMENT '支付单号',
  `transaction_id` varchar(100) NOT NULL DEFAULT '' COMMENT '第三方支付流水号',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1未支付 2已支付 3已取消',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `cancel_time` datetime DEFAULT NULL COMMENT '取消时间',
  `client_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '客户端IP',
  `source` varchar(20) NOT NULL DEFAULT '' COMMENT '来源：app/h5/mini/pc',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_plan_id` (`plan_id`),
  KEY `idx_payment_no` (`payment_no`),
  KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='充值订单表';

SET FOREIGN_KEY_CHECKS = 1;