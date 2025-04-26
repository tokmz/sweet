-- Sweet社交电商分销系统 - 订单模块表结构
-- 创建时间: 2023-06-01
-- 表前缀: sw_

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 订单表
-- ----------------------------
DROP TABLE IF EXISTS `sw_order`;
CREATE TABLE `sw_order` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` varchar(50) NOT NULL COMMENT '订单编号',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `order_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '订单类型：1普通订单 2秒杀订单 3拼团订单',
  `order_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '订单状态：1待付款 2待发货 3待收货 4已完成 5已取消 6已退款',
  `pay_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '支付状态：1未支付 2已支付 3部分支付 4已退款',
  `ship_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '发货状态：1未发货 2已发货 3部分发货 4已收货',
  `pay_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '支付方式：1未选择 2微信支付 3支付宝 4余额支付',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `ship_time` datetime DEFAULT NULL COMMENT '发货时间',
  `confirm_time` datetime DEFAULT NULL COMMENT '确认收货时间',
  `cancel_time` datetime DEFAULT NULL COMMENT '取消时间',
  `cancel_reason` varchar(255) NOT NULL DEFAULT '' COMMENT '取消原因',
  `finish_time` datetime DEFAULT NULL COMMENT '完成时间',
  `total_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '订单总金额',
  `product_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '商品总金额',
  `shipping_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '运费',
  `discount_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '优惠金额',
  `coupon_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '优惠券抵扣金额',
  `pay_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '实付金额',
  `refund_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '退款金额',
  `buyer_message` varchar(255) NOT NULL DEFAULT '' COMMENT '买家留言',
  `seller_message` varchar(255) NOT NULL DEFAULT '' COMMENT '卖家备注',
  `is_comment` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否已评价：1否 2是',
  `is_settlement` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否已结算：1否 2是',
  `is_distribution` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否分销订单：1否 2是',
  `distribution_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '分销用户ID',
  `distribution_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '分销佣金',
  `source` varchar(20) NOT NULL DEFAULT '' COMMENT '订单来源：app/h5/mini/pc',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT '下单IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_order_status` (`order_status`),
  KEY `idx_pay_status` (`pay_status`),
  KEY `idx_ship_status` (`ship_status`),
  KEY `idx_is_distribution` (`is_distribution`),
  KEY `idx_distribution_user_id` (`distribution_user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- ----------------------------
-- 订单商品表
-- ----------------------------
DROP TABLE IF EXISTS `sw_order_item`;
CREATE TABLE `sw_order_item` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单项ID',
  `order_id` bigint(20) UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no` varchar(50) NOT NULL COMMENT '订单编号',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `product_name` varchar(100) NOT NULL COMMENT '商品名称',
  `product_image` varchar(255) NOT NULL DEFAULT '' COMMENT '商品图片',
  `sku_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'SKU ID',
  `sku_code` varchar(50) NOT NULL DEFAULT '' COMMENT 'SKU编码',
  `spec_value_str` varchar(255) NOT NULL DEFAULT '' COMMENT '规格值描述',
  `price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '商品单价',
  `cost_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '成本价',
  `quantity` int(11) NOT NULL DEFAULT 0 COMMENT '购买数量',
  `total_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '总金额',
  `discount_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '优惠金额',
  `pay_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '实付金额',
  `refund_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '退款金额',
  `is_comment` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否已评价：1否 2是',
  `is_settlement` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否已结算：1否 2是',
  `is_distribution` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否参与分销：0否 1是',
  `distribution_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '分销佣金',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_sku_id` (`sku_id`),
  KEY `idx_is_comment` (`is_comment`),
  KEY `idx_is_settlement` (`is_settlement`),
  KEY `idx_is_distribution` (`is_distribution`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单商品表';

-- ----------------------------
-- 订单收货地址表
-- ----------------------------
DROP TABLE IF EXISTS `sw_order_address`;
CREATE TABLE `sw_order_address` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `order_id` bigint(20) UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no` varchar(50) NOT NULL COMMENT '订单编号',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '收货人姓名',
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '收货人手机号',
  `province` varchar(20) NOT NULL DEFAULT '' COMMENT '省份',
  `city` varchar(20) NOT NULL DEFAULT '' COMMENT '城市',
  `district` varchar(20) NOT NULL DEFAULT '' COMMENT '区/县',
  `detail` varchar(255) NOT NULL DEFAULT '' COMMENT '详细地址',
  `post_code` varchar(10) NOT NULL DEFAULT '' COMMENT '邮政编码',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单收货地址表';

-- ----------------------------
-- 订单支付表
-- ----------------------------
DROP TABLE IF EXISTS `sw_order_payment`;
CREATE TABLE `sw_order_payment` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `order_id` bigint(20) UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no` varchar(50) NOT NULL COMMENT '订单编号',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `payment_no` varchar(100) NOT NULL DEFAULT '' COMMENT '支付单号',
  `transaction_id` varchar(100) NOT NULL DEFAULT '' COMMENT '第三方支付流水号',
  `pay_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '支付方式：1未选择 2微信支付 3支付宝 4余额支付',
  `pay_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '支付金额',
  `pay_status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '支付状态：0未支付 1已支付 2支付失败 3已退款',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `refund_time` datetime DEFAULT NULL COMMENT '退款时间',
  `refund_amount` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '退款金额',
  `refund_reason` varchar(255) NOT NULL DEFAULT '' COMMENT '退款原因',
  `pay_channel` varchar(20) NOT NULL DEFAULT '' COMMENT '支付渠道',
  `client_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '客户端IP',
  `notify_data` text COMMENT '回调通知数据',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_payment_no` (`payment_no`),
  KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_pay_status` (`pay_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单支付表';

-- ----------------------------
-- 订单物流表
-- ----------------------------
DROP TABLE IF EXISTS `sw_order_shipping`;
CREATE TABLE `sw_order_shipping` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `order_id` bigint(20) UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no` varchar(50) NOT NULL COMMENT '订单编号',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `shipping_code` varchar(50) NOT NULL DEFAULT '' COMMENT '物流单号',
  `shipping_company` varchar(50) NOT NULL DEFAULT '' COMMENT '物流公司',
  `shipping_time` datetime DEFAULT NULL COMMENT '发货时间',
  `confirm_time` datetime DEFAULT NULL COMMENT '确认收货时间',
  `shipping_status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '物流状态：0未发货 1已发货 2已收货',
  `shipping_data` text COMMENT '物流跟踪数据',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_shipping_code` (`shipping_code`),
  KEY `idx_shipping_status` (`shipping_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单物流表';

-- ----------------------------
-- 购物车表
-- ----------------------------
DROP TABLE IF EXISTS `sw_cart`;
CREATE TABLE `sw_cart` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `sku_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'SKU ID',
  `quantity` int(11) NOT NULL DEFAULT 0 COMMENT '购买数量',
  `is_selected` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否选中：0否 1是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_sku` (`user_id`, `sku_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_sku_id` (`sku_id`),
  KEY `idx_is_selected` (`is_selected`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';

SET FOREIGN_KEY_CHECKS = 1;