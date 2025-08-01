# Sweet 数据库设计文档

## 概述

本目录包含 Sweet 社交电商分销系统的数据库设计文件。

## 文件说明

### sweet.sql

完整的数据库结构定义文件，包含以下模块：

#### 系统管理模块
- `sw_sys_user` - 系统用户表
- `sw_sys_role` - 角色表
- `sw_sys_menu` - 菜单表
- `sw_sys_menu_config` - 菜单配置表
- `sw_sys_role_menu` - 角色菜单关联表
- `sw_sys_api` - API接口表
- `sw_sys_api_group` - API分组表
- `sw_sys_role_api` - 角色API关联表

#### 组织架构模块
- `sw_sys_dept` - 部门表
- `sw_sys_post` - 岗位表

#### 日志模块
- `sw_sys_login_log` - 登录日志表
- `sw_sys_operation_log` - 操作日志表

#### 文件管理模块
- `sw_sys_file` - 文件表（简化版本）

## 数据库特性

### 设计原则
- 遵循三范式设计
- 使用软删除机制（deleted_at字段）
- 统一的字符集：utf8mb4
- 合理的索引设计
- 完整的外键约束

### 字段规范
- 主键：使用 `bigint unsigned` 自增ID
- 时间字段：`created_at`、`updated_at`、`deleted_at`
- 状态字段：使用 `tinyint unsigned`，1=启用/正常，2=禁用/删除
- 排序字段：使用 `int unsigned`，数值越大排序越靠前

### 索引策略
- 主键索引：所有表都有自增主键
- 唯一索引：用于保证数据唯一性（如用户名、编码等）
- 普通索引：用于提高查询性能
- 复合索引：用于多字段组合查询
- 软删除索引：deleted_at字段索引

## 使用说明

### 导入数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE sweet CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入数据结构
mysql -u root -p sweet < sweet.sql
```

### 环境要求
- MySQL 8.0+
- 字符集：utf8mb4
- 排序规则：utf8mb4_unicode_ci

### 注意事项
1. 导入前请确保数据库版本兼容
2. 建议在测试环境先验证SQL文件
3. 生产环境导入前请备份现有数据
4. 外键约束已启用，删除数据时注意关联关系

## 版本历史

### v1.0.0
- 初始版本，包含完整的系统管理功能
- 支持用户、角色、菜单、API权限管理
- 包含组织架构和日志功能
- 简化的文件管理功能

## 开发说明

### GORM 代码生成

本项目使用 GORM Gen 自动生成模型代码，相关配置文件：
- `scripts/gorm_gen.go` - GORM 代码生成器配置

生成模型代码：
```bash
go run scripts/gorm_gen.go
```

### 数据库连接配置

请参考项目根目录的配置文件示例进行数据库连接配置。

## 联系方式

如有问题或建议，请联系开发团队。