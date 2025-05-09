# Sweet 项目规范

**Sweet** 是一款专注于社交电商场景的分销商城系统，帮助企业快速搭建多层级分销体系，实现商品销售、裂变推广与团队管理的全链路数字化。系统采用**Go+Gin单体架构**设计，提供开箱即用的功能模块与深度定制能力。

## 1. 项目模块

1.  **用户模块**: 用户注册、登录、认证、会员信息、地址管理等。
2.  **商品模块**: 商品管理、分类、规格、库存、评价等。
3.  **订单模块**: 订单创建、支付、履约、购物车、售后退款等。
4.  **支付模块**: 对接支付渠道、处理支付回调、管理资金流水等。
5.  **分销模块**: 分销关系、佣金计算与结算、分销商管理等。
6.  **营销模块**: 优惠券、秒杀、拼团、积分等活动管理。
7.  **内容模块**: 文章、广告、通知等内容管理。
8.  **社区圈子模块**: 圈子管理、帖子发布、评论互动、用户关注等社交功能。
9.  **系统模块**: 管理员、角色权限、系统配置、日志审计等。
10. **通用模块**: 提供文件上传、短信、消息通知等通用服务。

## 2. 技术栈

### 2.1 核心技术

- **语言**: Go 1.24
- **Web框架**: Gin
- **ORM**: GORM
- **缓存**: Redis
- **数据库**: MySQL 8.0+
- **配置管理**: Viper
- **日志**: Zap + Lumberjack
- **链路追踪**: OpenTelemetry
- **消息队列**: RabbitMQ/Kafka
- **搜索引擎**: Elasticsearch

### 2.2 基础组件

- **数据库组件**: 支持主从读写分离、连接池管理、事务管理
- **缓存组件**: 支持多级缓存、分布式锁、过期策略
- **日志组件**: 结构化日志、日志分级、日志轮转
- **配置组件**: 多环境配置、热更新、敏感信息加密
- **认证组件**: JWT认证、OAuth2.0、RBAC权限控制
- **通知组件**: 短信、邮件、站内信、推送

## 3. 项目结构

## 4. 编码规范

### 4.1 命名规范

- **包名**: 使用小写单词，不使用下划线或混合大小写
- **文件名**: 使用小写单词，可使用下划线分隔
- **变量名**: 使用驼峰命名法，局部变量首字母小写，全局变量首字母大写
- **常量名**: 使用全大写，单词间用下划线分隔
- **接口名**: 使用驼峰命名法，通常以"er"结尾
- **结构体名**: 使用驼峰命名法，首字母大写

### 4.2 注释规范

- 所有导出的函数、类型、变量和常量必须有注释
- 包级别的注释应该位于package语句之前
- 使用完整的句子，以被描述的对象开头
- 函数注释应该说明函数的功能、参数和返回值

### 4.3 错误处理

- 错误应该被明确处理，不应该被忽略
- 使用自定义错误类型增强错误信息
- 错误信息应该简洁明了，包含足够的上下文信息
- 使用`errors.Is`和`errors.As`进行错误检查和类型断言

### 4.4 并发处理

- 使用`context.Context`传递截止日期、取消信号和请求范围的值
- 谨慎使用goroutine，确保它们能够正确退出
- 使用互斥锁或读写锁保护共享资源
- 使用通道进行goroutine间的通信

## 5. 数据库规范

### 5.1 表设计

- 表名使用小写，单词间用下划线分隔
- 主键使用`id`，类型为无符号整数或UUID
- 每个表必须包含`created_at`、`updated_at`字段
- 使用软删除，添加`deleted_at`字段
- 字段名使用小写，单词间用下划线分隔

### 5.2 索引规范

- 为经常查询的字段创建索引
- 索引名应该反映其用途和包含的字段
- 避免过度索引，考虑写入性能
- 定期检查索引使用情况，删除未使用的索引

### 5.3 查询优化

- 只查询需要的字段，避免`SELECT *`
- 使用预编译语句防止SQL注入
- 大型查询使用分页
- 使用事务确保数据一致性
- 复杂查询考虑使用存储过程

## 6. API设计规范

### 6.1 RESTful API

- 使用HTTP方法表示操作类型（GET、POST、PUT、DELETE）
- 使用名词复数形式作为资源标识
- 使用嵌套资源表示关系
- 使用HTTP状态码表示请求结果
- 支持分页、排序和过滤

### 6.2 响应格式

```json
{
  "code": 200,        // 业务状态码，200表示成功
  "msg": "成功",      // 状态描述
  "data": {},        // 业务数据
  "trace_id": "xxx"  // 请求ID，用于追踪
}