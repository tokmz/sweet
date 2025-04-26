-- Sweet社交电商分销系统 - 社区圈子模块表结构
-- 创建时间: 2023-06-01
-- 表前缀: sw_

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 圈子表
-- ----------------------------
DROP TABLE IF EXISTS `sw_community_circle`;
CREATE TABLE `sw_community_circle` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '圈子ID',
  `name` varchar(50) NOT NULL COMMENT '圈子名称',
  `icon` varchar(255) NOT NULL DEFAULT '' COMMENT '圈子图标',
  `cover` varchar(255) NOT NULL DEFAULT '' COMMENT '圈子封面',
  `description` varchar(500) NOT NULL DEFAULT '' COMMENT '圈子描述',
  `notice` varchar(500) NOT NULL DEFAULT '' COMMENT '圈子公告',
  `category` varchar(50) NOT NULL DEFAULT '' COMMENT '圈子分类',
  `type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '圈子类型：1公开 2私密',
  `join_mode` tinyint(1) NOT NULL DEFAULT 1 COMMENT '加入方式：1自由加入 2审核加入 3仅邀请',
  `invitation_expire_days` int(11) NOT NULL DEFAULT 7 COMMENT '邀请有效期(天)，0表示永久有效',
  `creator_user_id` bigint(20) UNSIGNED NOT NULL COMMENT '创建者用户ID',
  `member_count` int(11) NOT NULL DEFAULT 0 COMMENT '成员数量',
  `post_count` int(11) NOT NULL DEFAULT 0 COMMENT '帖子数量',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `is_recommend` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否推荐：1否 2是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_creator_user_id` (`creator_user_id`),
  KEY `idx_category` (`category`),
  KEY `idx_type` (`type`),
  KEY `idx_join_mode` (`join_mode`),
  KEY `idx_status` (`status`),
  KEY `idx_is_recommend` (`is_recommend`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='圈子表';

-- ----------------------------
-- 圈子成员表
-- ----------------------------
DROP TABLE IF EXISTS `sw_circle_member`;
CREATE TABLE `sw_circle_member` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `circle_id` bigint(20) UNSIGNED NOT NULL COMMENT '圈子ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `role` tinyint(1) NOT NULL DEFAULT 0 COMMENT '角色：0普通成员 1管理员 2创建者',
  `nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '圈子内昵称',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁言',
  `join_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_circle_user` (`circle_id`, `user_id`),
  KEY `idx_circle_id` (`circle_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role` (`role`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='圈子成员表';

-- ----------------------------
-- 圈子加入申请表
-- ----------------------------
DROP TABLE IF EXISTS `sw_circle_join_apply`;
CREATE TABLE `sw_circle_join_apply` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `circle_id` bigint(20) UNSIGNED NOT NULL COMMENT '圈子ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `apply_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '申请类型：1主动申请 2邀请加入',
  `reason` varchar(255) NOT NULL DEFAULT '' COMMENT '申请理由',
  `inviter_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '邀请人用户ID，0表示非邀请加入',
  `invitation_code` varchar(32) NOT NULL DEFAULT '' COMMENT '邀请码',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1待审核 2已通过 3已拒绝 4已过期',
  `audit_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '审核人用户ID',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `audit_remark` varchar(255) NOT NULL DEFAULT '' COMMENT '审核备注',
  `expire_time` datetime DEFAULT NULL COMMENT '邀请过期时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_circle_id` (`circle_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_inviter_user_id` (`inviter_user_id`),
  KEY `idx_invitation_code` (`invitation_code`),
  KEY `idx_status` (`status`),
  KEY `idx_apply_type` (`apply_type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='圈子加入申请表';

-- ----------------------------
-- 帖子表
-- ----------------------------
DROP TABLE IF EXISTS `sw_community_post`;
CREATE TABLE `sw_community_post` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '帖子ID',
  `circle_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '圈子ID，0表示不属于任何圈子',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '标题',
  `content` text NOT NULL COMMENT '内容',
  `images` text COMMENT '图片，JSON格式',
  `videos` text COMMENT '视频，JSON格式',
  `location` varchar(100) NOT NULL DEFAULT '' COMMENT '位置',
  `longitude` decimal(10,6) NOT NULL DEFAULT 0.000000 COMMENT '经度',
  `latitude` decimal(10,6) NOT NULL DEFAULT 0.000000 COMMENT '纬度',
  `view_count` int(11) NOT NULL DEFAULT 0 COMMENT '浏览量',
  `like_count` int(11) NOT NULL DEFAULT 0 COMMENT '点赞数',
  `comment_count` int(11) NOT NULL DEFAULT 0 COMMENT '评论数',
  `share_count` int(11) NOT NULL DEFAULT 0 COMMENT '分享数',
  `collect_count` int(11) NOT NULL DEFAULT 0 COMMENT '收藏数',
  `is_top` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否置顶：1否 2是',
  `is_essence` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否精华：1否 2是',
  `is_anonymous` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否匿名：1否 2是',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1显示 2隐藏',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_circle_id` (`circle_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_top` (`is_top`),
  KEY `idx_is_essence` (`is_essence`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子表';

-- ----------------------------
-- 帖子评论表
-- ----------------------------
DROP TABLE IF EXISTS `sw_post_comment`;
CREATE TABLE `sw_post_comment` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `post_id` bigint(20) UNSIGNED NOT NULL COMMENT '帖子ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `parent_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '父评论ID',
  `reply_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '回复用户ID',
  `content` varchar(500) NOT NULL COMMENT '评论内容',
  `images` text COMMENT '图片，JSON格式',
  `like_count` int(11) NOT NULL DEFAULT 0 COMMENT '点赞数',
  `reply_count` int(11) NOT NULL DEFAULT 0 COMMENT '回复数',
  `is_anonymous` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否匿名：1否 2是',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1显示 2隐藏',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_reply_user_id` (`reply_user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子评论表';

-- ----------------------------
-- 用户点赞表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user_like`;
CREATE TABLE `sw_user_like` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `target_id` bigint(20) UNSIGNED NOT NULL COMMENT '目标ID',
  `target_type` tinyint(1) NOT NULL COMMENT '目标类型：1帖子 2评论 3文章',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_target` (`user_id`, `target_id`, `target_type`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_target_id` (`target_id`),
  KEY `idx_target_type` (`target_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户点赞表';

-- ----------------------------
-- 用户收藏表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user_collect`;
CREATE TABLE `sw_user_collect` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `target_id` bigint(20) UNSIGNED NOT NULL COMMENT '目标ID',
  `target_type` tinyint(1) NOT NULL COMMENT '目标类型：1帖子 2文章 3商品',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_target` (`user_id`, `target_id`, `target_type`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_target_id` (`target_id`),
  KEY `idx_target_type` (`target_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收藏表';

-- ----------------------------
-- 用户关注表
-- ----------------------------
DROP TABLE IF EXISTS `sw_user_follow`;
CREATE TABLE `sw_user_follow` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `follow_user_id` bigint(20) UNSIGNED NOT NULL COMMENT '关注的用户ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_follow` (`user_id`, `follow_user_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_follow_user_id` (`follow_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户关注表';

-- ----------------------------
-- 圈子管理员表
-- ----------------------------
DROP TABLE IF EXISTS `sw_circle_admin`;
CREATE TABLE `sw_circle_admin` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `circle_id` bigint(20) UNSIGNED NOT NULL COMMENT '圈子ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '管理员用户ID',
  `remark` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_circle_user` (`circle_id`, `user_id`),
  KEY `idx_circle_id` (`circle_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='圈子管理员表';

-- ----------------------------
-- 话题表
-- ----------------------------
DROP TABLE IF EXISTS `sw_topic`;
CREATE TABLE `sw_topic` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '话题ID',
  `name` varchar(50) NOT NULL COMMENT '话题名称',
  `icon` varchar(255) NOT NULL DEFAULT '' COMMENT '话题图标',
  `description` varchar(500) NOT NULL DEFAULT '' COMMENT '话题描述',
  `post_count` int(11) NOT NULL DEFAULT 0 COMMENT '帖子数量',
  `view_count` int(11) NOT NULL DEFAULT 0 COMMENT '浏览量',
  `follow_count` int(11) NOT NULL DEFAULT 0 COMMENT '关注数',
  `sort` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
  `is_recommend` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否推荐：1否 2是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_status` (`status`),
  KEY `idx_is_recommend` (`is_recommend`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='话题表';

-- ----------------------------
-- 帖子话题关联表
-- ----------------------------
DROP TABLE IF EXISTS `sw_post_topic`;
CREATE TABLE `sw_post_topic` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `post_id` bigint(20) UNSIGNED NOT NULL COMMENT '帖子ID',
  `topic_id` bigint(20) UNSIGNED NOT NULL COMMENT '话题ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_post_topic` (`post_id`, `topic_id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_topic_id` (`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子话题关联表';

SET FOREIGN_KEY_CHECKS = 1;