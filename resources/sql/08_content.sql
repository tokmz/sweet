-- Sweet社交电商分销系统 - 内容模块表结构
-- 创建时间: 2023-06-01
-- 表前缀: sw_

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 文章分类表
-- ----------------------------
DROP TABLE IF EXISTS `sw_article_category`;
CREATE TABLE `sw_article_category` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `parent_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `code` varchar(50) NOT NULL DEFAULT '' COMMENT '分类编码',
  `icon` varchar(255) NOT NULL DEFAULT '' COMMENT '分类图标',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章分类表';

-- ----------------------------
-- 文章表
-- ----------------------------
DROP TABLE IF EXISTS `sw_article`;
CREATE TABLE `sw_article` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '文章ID',
  `category_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '分类ID',
  `title` varchar(100) NOT NULL COMMENT '文章标题',
  `subtitle` varchar(200) NOT NULL DEFAULT '' COMMENT '副标题',
  `cover` varchar(255) NOT NULL DEFAULT '' COMMENT '封面图',
  `summary` varchar(500) NOT NULL DEFAULT '' COMMENT '摘要',
  `content` text NOT NULL COMMENT '内容',
  `author` varchar(50) NOT NULL DEFAULT '' COMMENT '作者',
  `source` varchar(50) NOT NULL DEFAULT '' COMMENT '来源',
  `view_count` int(11) NOT NULL DEFAULT 0 COMMENT '浏览量',
  `like_count` int(11) NOT NULL DEFAULT 0 COMMENT '点赞数',
  `comment_count` int(11) NOT NULL DEFAULT 0 COMMENT '评论数',
  `is_top` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否置顶：1否 2是',
  `is_recommend` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否推荐：1否 2是',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1已发布 2已下架 3草稿',
  `seo_title` varchar(100) NOT NULL DEFAULT '' COMMENT 'SEO标题',
  `seo_keywords` varchar(200) NOT NULL DEFAULT '' COMMENT 'SEO关键词',
  `seo_description` varchar(500) NOT NULL DEFAULT '' COMMENT 'SEO描述',
  `publish_time` datetime DEFAULT NULL COMMENT '发布时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_top` (`is_top`),
  KEY `idx_is_recommend` (`is_recommend`),
  KEY `idx_publish_time` (`publish_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章表';

-- ----------------------------
-- 文章评论表
-- ----------------------------
DROP TABLE IF EXISTS `sw_article_comment`;
CREATE TABLE `sw_article_comment` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `article_id` bigint(20) UNSIGNED NOT NULL COMMENT '文章ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `parent_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父评论ID',
  `content` varchar(500) NOT NULL COMMENT '评论内容',
  `like_count` int(11) NOT NULL DEFAULT 0 COMMENT '点赞数',
  `reply_count` int(11) NOT NULL DEFAULT 0 COMMENT '回复数',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1显示 2隐藏',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_article_id` (`article_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章评论表';

-- ----------------------------
-- 广告位表
-- ----------------------------
DROP TABLE IF EXISTS `sw_ad_position`;
CREATE TABLE `sw_ad_position` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '广告位ID',
  `name` varchar(50) NOT NULL COMMENT '广告位名称',
  `code` varchar(50) NOT NULL COMMENT '广告位编码',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '广告位描述',
  `width` int(11) NOT NULL DEFAULT 0 COMMENT '宽度',
  `height` int(11) NOT NULL DEFAULT 0 COMMENT '高度',
  `template` text COMMENT '广告位模板',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='广告位表';

-- ----------------------------
-- 广告表
-- ----------------------------
DROP TABLE IF EXISTS `sw_ad`;
CREATE TABLE `sw_ad` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '广告ID',
  `position_id` bigint(20) UNSIGNED NOT NULL COMMENT '广告位ID',
  `name` varchar(50) NOT NULL COMMENT '广告名称',
  `media_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '媒体类型：1图片 2视频 3文字',
  `image` varchar(255) NOT NULL DEFAULT '' COMMENT '图片地址',
  `video` varchar(255) NOT NULL DEFAULT '' COMMENT '视频地址',
  `text` varchar(500) NOT NULL DEFAULT '' COMMENT '文字内容',
  `link` varchar(255) NOT NULL DEFAULT '' COMMENT '链接地址',
  `link_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '链接类型：1站内 2站外',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_position_id` (`position_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_end_time` (`end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='广告表';

-- ----------------------------
-- 通知公告表
-- ----------------------------
DROP TABLE IF EXISTS `sw_notice`;
CREATE TABLE `sw_notice` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '通知ID',
  `title` varchar(100) NOT NULL COMMENT '通知标题',
  `content` text NOT NULL COMMENT '通知内容',
  `type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '通知类型：1系统通知 2活动通知 3更新通知',
  `target` tinyint(1) NOT NULL DEFAULT 1 COMMENT '目标对象：1全部用户 2指定用户 3指定用户组',
  `target_ids` text COMMENT '目标ID，JSON格式',
  `is_top` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否置顶：1否 2是',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1已发布 2已下架 3草稿',
  `publish_time` datetime DEFAULT NULL COMMENT '发布时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_target` (`target`),
  KEY `idx_status` (`status`),
  KEY `idx_is_top` (`is_top`),
  KEY `idx_publish_time` (`publish_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知公告表';

-- ----------------------------
-- 用户通知表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user_notice`;
CREATE TABLE `sw_user_notice` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `notice_id` bigint(20) UNSIGNED NOT NULL COMMENT '通知ID',
  `is_read` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否已读：1否 2是',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_notice` (`user_id`, `notice_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_notice_id` (`notice_id`),
  KEY `idx_is_read` (`is_read`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户通知表';

-- ----------------------------
-- 消息模板表
-- ----------------------------
DROP TABLE IF EXISTS `sw_message_template`;
CREATE TABLE `sw_message_template` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `code` varchar(50) NOT NULL COMMENT '模板编码',
  `name` varchar(50) NOT NULL COMMENT '模板名称',
  `type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '模板类型：1短信 2邮件 3站内信 4推送',
  `content` text NOT NULL COMMENT '模板内容',
  `params` varchar(255) NOT NULL DEFAULT '' COMMENT '参数列表，JSON格式',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_type` (`code`, `type`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息模板表';

-- ----------------------------
-- 消息记录表
-- ----------------------------
DROP TABLE IF EXISTS `sw_message_log`;
CREATE TABLE `sw_message_log` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `template_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '模板ID',
  `template_code` varchar(50) NOT NULL DEFAULT '' COMMENT '模板编码',
  `type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '消息类型：1短信 2邮件 3站内信 4推送',
  `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
  `receiver` varchar(100) NOT NULL DEFAULT '' COMMENT '接收者',
  `content` text NOT NULL COMMENT '消息内容',
  `params` varchar(500) NOT NULL DEFAULT '' COMMENT '参数内容，JSON格式',
  `send_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '发送状态：1待发送 2发送成功 3发送失败',
  `send_time` datetime DEFAULT NULL COMMENT '发送时间',
  `error_msg` varchar(255) NOT NULL DEFAULT '' COMMENT '错误信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_template_code` (`template_code`),
  KEY `idx_type` (`type`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_send_status` (`send_status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息记录表';

SET FOREIGN_KEY_CHECKS = 1;