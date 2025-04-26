-- Sweet社交电商分销系统 - 商品模块表结构
-- 创建时间: 2025-04-25
-- 更新时间: 2025-04-26
-- 表前缀: sw_

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 商品品牌表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_brand`;
CREATE TABLE `sw_product_brand` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '品牌ID',
  `name` varchar(50) NOT NULL COMMENT '品牌名称',
  `logo` varchar(255) NOT NULL DEFAULT '' COMMENT '品牌Logo',
  `description` varchar(500) NOT NULL DEFAULT '' COMMENT '品牌描述',
  `website` varchar(255) NOT NULL DEFAULT '' COMMENT '品牌官网',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `is_recommend` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否推荐：1否 2是',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`),
  KEY `idx_is_recommend` (`is_recommend`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品品牌表';

-- ----------------------------
-- 商品分类表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_category`;
CREATE TABLE `sw_product_category` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `parent_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父分类ID',
  `name` varchar(64) NOT NULL COMMENT '分类名称',
  `icon` varchar(255) NOT NULL DEFAULT '' COMMENT '分类图标',
  `image` varchar(255) NOT NULL DEFAULT '' COMMENT '分类图片',
  `level` tinyint(1) NOT NULL DEFAULT 1 COMMENT '分类级别',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `is_recommend` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否推荐：1否 2是',
  `commission_rate` decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '分类佣金比例',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_is_recommend` (`is_recommend`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品分类表';

-- ----------------------------
-- 商品标签表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_tag`;
CREATE TABLE `sw_product_tag` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '标签ID',
  `name` varchar(30) NOT NULL COMMENT '标签名称',
  `color` varchar(20) NOT NULL DEFAULT '' COMMENT '标签颜色',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品标签表';

-- ----------------------------
-- 商品与标签关联表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_tag_relation`;
CREATE TABLE `sw_product_tag_relation` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `tag_id` bigint(20) UNSIGNED NOT NULL COMMENT '标签ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_product_tag` (`product_id`,`tag_id`),
  KEY `idx_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品与标签关联表';

-- ----------------------------
-- 商品属性表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_attribute`;
CREATE TABLE `sw_product_attribute` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '属性ID',
  `category_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '属性名称',
  `input_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '录入方式：1手工录入 2从选项列表选取',
  `values` text COMMENT '可选值列表，多个以逗号分隔',
  `is_required` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否必填：0否 1是',
  `is_filter` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否支持筛选：0否 1是',
  `is_search` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否支持搜索：0否 1是',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品属性表';

-- ----------------------------
-- 商品属性值表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_attribute_value`;
CREATE TABLE `sw_product_attribute_value` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '属性值ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `attribute_id` bigint(20) UNSIGNED NOT NULL COMMENT '属性ID',
  `value` varchar(255) NOT NULL DEFAULT '' COMMENT '属性值',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_attribute_id` (`attribute_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品属性值表';

-- ----------------------------
-- 商品表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product`;
CREATE TABLE `sw_product` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `category_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '分类ID',
  `brand_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '品牌ID',
  `name` varchar(100) NOT NULL COMMENT '商品名称',
  `subtitle` varchar(200) NOT NULL DEFAULT '' COMMENT '副标题',
  `product_sn` varchar(50) NOT NULL DEFAULT '' COMMENT '商品编号',
  `main_image` varchar(255) NOT NULL DEFAULT '' COMMENT '主图',
  `album` text COMMENT '相册，JSON格式',
  `video` varchar(255) NOT NULL DEFAULT '' COMMENT '视频URL',
  `detail` text COMMENT '商品详情',
  `price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '售价',
  `market_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '市场价',
  `cost_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '成本价',
  `stock` int(11) NOT NULL DEFAULT 0 COMMENT '总库存',
  `stock_warning` int(11) NOT NULL DEFAULT 0 COMMENT '库存预警值',
  `sales` int(11) NOT NULL DEFAULT 0 COMMENT '销量',
  `virtual_sales` int(11) NOT NULL DEFAULT 0 COMMENT '虚拟销量',
  `unit` varchar(20) NOT NULL DEFAULT '' COMMENT '单位',
  `weight` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '重量(kg)',
  `volume` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '体积(m³)',
  `is_on_sale` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否上架：1否 2是',
  `is_new` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否新品：1否 2是',
  `is_hot` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否热销：1否 2是',
  `is_recommend` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否推荐：1否 2是',
  `is_virtual` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否虚拟商品：1否 2是',
  `has_sku` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否有规格：0否 1是',
  `delivery_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '配送方式：1物流 2到店自提 3同城配送',
  `freight_template_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '运费模板ID',
  `distribution_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '分销状态：1不参与分销 2参与分销',
  `distribution_commission_rate` decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '分销佣金比例',
  `audit_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '审核状态：1待审核 2审核通过 3审核拒绝',
  `audit_remark` varchar(255) NOT NULL DEFAULT '' COMMENT '审核备注',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_brand_id` (`brand_id`),
  KEY `idx_product_sn` (`product_sn`),
  KEY `idx_is_on_sale` (`is_on_sale`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`),
  KEY `idx_distribution_status` (`distribution_status`),
  KEY `idx_audit_status` (`audit_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

-- ----------------------------
-- 商品规格表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_spec`;
CREATE TABLE `sw_product_spec` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '规格ID',
  `name` varchar(50) NOT NULL COMMENT '规格名称',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品规格表';

-- ----------------------------
-- 商品规格值表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_spec_value`;
CREATE TABLE `sw_product_spec_value` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '规格值ID',
  `spec_id` bigint(20) UNSIGNED NOT NULL COMMENT '规格ID',
  `value` varchar(50) NOT NULL COMMENT '规格值',
  `image` varchar(255) NOT NULL DEFAULT '' COMMENT '规格图片',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_spec_id` (`spec_id`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品规格值表';

-- ----------------------------
-- 商品规格关联表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_spec_relation`;
CREATE TABLE `sw_product_spec_relation` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `spec_id` bigint(20) UNSIGNED NOT NULL COMMENT '规格ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_product_spec` (`product_id`,`spec_id`),
  KEY `idx_spec_id` (`spec_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品规格关联表';

-- ----------------------------
-- 商品SKU表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_sku`;
CREATE TABLE `sw_product_sku` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `sku_code` varchar(50) NOT NULL DEFAULT '' COMMENT 'SKU编码',
  `spec_value_ids` varchar(255) NOT NULL DEFAULT '' COMMENT '规格值ID，多个以逗号分隔',
  `spec_value_str` varchar(255) NOT NULL DEFAULT '' COMMENT '规格值描述',
  `image` varchar(255) NOT NULL DEFAULT '' COMMENT 'SKU图片',
  `price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '售价',
  `market_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '市场价',
  `cost_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '成本价',
  `stock` int(11) NOT NULL DEFAULT 0 COMMENT '库存',
  `stock_warning` int(11) NOT NULL DEFAULT 0 COMMENT '库存预警值',
  `sales` int(11) NOT NULL DEFAULT 0 COMMENT '销量',
  `weight` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '重量(kg)',
  `volume` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '体积(m³)',
  `barcode` varchar(50) NOT NULL DEFAULT '' COMMENT '条形码',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_sku_code` (`sku_code`),
  KEY `idx_barcode` (`barcode`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品SKU表';

-- ----------------------------
-- 商品库存记录表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_stock_log`;
CREATE TABLE `sw_product_stock_log` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `sku_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'SKU ID',
  `type` tinyint(1) NOT NULL COMMENT '类型：1入库 2出库 3调整',
  `quantity` int(11) NOT NULL COMMENT '数量',
  `before_stock` int(11) NOT NULL COMMENT '变更前库存',
  `after_stock` int(11) NOT NULL COMMENT '变更后库存',
  `order_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单ID',
  `order_item_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单项ID',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `operator_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作人ID',
  `operator_name` varchar(50) NOT NULL DEFAULT '' COMMENT '操作人姓名',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_sku_id` (`sku_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_type` (`type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品库存记录表';

-- ----------------------------
-- 商品评价表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_comment`;
CREATE TABLE `sw_product_comment` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评价ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `order_id` bigint(20) UNSIGNED NOT NULL COMMENT '订单ID',
  `order_item_id` bigint(20) UNSIGNED NOT NULL COMMENT '订单项ID',
  `sku_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'SKU ID',
  `content` text NOT NULL COMMENT '评价内容',
  `images` text COMMENT '评价图片，JSON格式',
  `video` varchar(255) NOT NULL DEFAULT '' COMMENT '评价视频',
  `star` tinyint(1) NOT NULL DEFAULT 5 COMMENT '评分：1-5',
  `is_anonymous` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否匿名：0否 1是',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：0隐藏 1显示',
  `reply_content` text COMMENT '商家回复内容',
  `reply_time` datetime DEFAULT NULL COMMENT '商家回复时间',
  `is_show` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否显示：0否 1是',
  `is_top` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否置顶：0否 1是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_star` (`star`),
  KEY `idx_status` (`status`),
  KEY `idx_is_show` (`is_show`),
  KEY `idx_is_top` (`is_top`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品评价表';

-- ----------------------------
-- 商品收藏表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_favorite`;
CREATE TABLE `sw_product_favorite` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_product` (`user_id`,`product_id`),
  KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品收藏表';

-- ----------------------------
-- 商品浏览记录表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_view_log`;
CREATE TABLE `sw_product_view_log` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `product_id` bigint(20) UNSIGNED NOT NULL COMMENT '商品ID',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `user_agent` varchar(255) NOT NULL DEFAULT '' COMMENT 'User Agent',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品浏览记录表';

-- ----------------------------
-- 商品运费模板表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_freight_template`;
CREATE TABLE `sw_product_freight_template` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) NOT NULL COMMENT '模板名称',
  `charge_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '计费方式：1按件数 2按重量 3按体积',
  `is_free` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否包邮：0否 1是',
  `free_condition` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '包邮条件，满多少包邮',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品运费模板表';

-- ----------------------------
-- 商品运费模板规则表
-- ----------------------------
DROP TABLE IF EXISTS `sw_product_freight_rule`;
CREATE TABLE `sw_product_freight_rule` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `template_id` bigint(20) UNSIGNED NOT NULL COMMENT '模板ID',
  `region_ids` text NOT NULL COMMENT '地区ID，多个以逗号分隔',
  `region_names` text NOT NULL COMMENT '地区名称，多个以逗号分隔',
  `first_quantity` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '首件/首重/首体积',
  `first_fee` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '首件/首重/首体积费用',
  `additional_quantity` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '续件/续重/续体积',
  `additional_fee` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '续件/续重/续体积费用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品运费模板规则表';

SET FOREIGN_KEY_CHECKS = 1;