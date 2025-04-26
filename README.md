# Sweet社交电商分销系统

## 项目简介

Sweet是一款基于Go语言开发的社交电商分销系统，旨在为企业提供强大、高效、可扩展的社交分销解决方案。该系统整合了社交裂变、分销管理、商品管理、订单处理、支付集成等多种功能，帮助企业快速搭建自己的社交电商平台。

## 技术架构

### 技术栈

- **后端框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **认证授权**: JWT + Casbin
- **日志**: Zap
- **配置管理**: Viper
- **监控追踪**: OpenTelemetry

### 系统架构

Sweet采用现代化的分层架构设计，遵循领域驱动设计(DDD)原则：

- **接口层(API)**: 负责处理HTTP请求，包括参数验证、路由分发等
- **应用层(Service)**: 实现核心业务逻辑，协调领域对象
- **领域层(Domain)**: 定义业务实体和核心业务规则
- **基础设施层(Infrastructure)**: 提供技术支持，如数据库访问、缓存、消息队列等

## 核心功能

### 用户管理
- 用户注册、登录、认证
- 用户信息管理
- 角色权限控制

### 分销管理
- 分销关系建立与维护
- 分销层级管理
- 佣金计算与结算
- 分销员审核与管理

### 商品管理
- 商品发布与管理
- 商品分类与标签
- 商品库存管理
- 商品定价策略

### 订单管理
- 订单创建与处理
- 订单状态追踪
- 售后服务管理
- 物流信息集成

### 营销活动
- 优惠券管理
- 满减活动
- 限时折扣
- 社交裂变活动

### 数据分析
- 销售数据统计
- 用户行为分析
- 分销效果评估
- 业绩报表生成

## 项目目录结构

```
sweet/
├── cmd/                 # 应用入口点
├── common/              # 通用工具和辅助函数
├── internal/            # 私有应用代码
│   ├── api/             # API处理器
│   ├── middleware/      # HTTP中间件
│   ├── models/          # 数据模型
│   │   ├── entity/      # 数据库实体
│   │   ├── vo/          # 视图对象
│   │   └── dto/         # 数据传输对象
│   ├── repo/            # 数据访问层
│   └── service/         # 业务逻辑层
├── pkg/                 # 可被外部应用程序使用的库代码
└── resources/           # 配置文件和静态资源
```

## 如何开始

### 环境要求

- Go 1.24+
- MySQL 8.0+
- Redis 6.0+

### 安装和运行

1. 克隆代码库
```bash
git clone https://github.com/yourusername/sweet.git
cd sweet
```

2. 安装依赖
```bash
go mod download
```

3. 配置数据库
- 创建数据库
- 修改`resources/config/`下的配置文件

4. 启动服务
```bash
go run cmd/sweet/main.go
```

## 开发规范

### 代码规范
- 遵循Go语言官方规范
- 使用驼峰命名法
- 公开函数必须有注释
- 错误处理必须明确

### API设计
- 遵循RESTful设计原则
- 统一的响应格式
- 版本控制

### 数据库规范
- 表设计遵循三范式
- 使用软删除
- 合理设计索引

## 贡献指南

欢迎贡献代码，请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情见 [LICENSE](LICENSE) 文件

## 联系我们

如有任何问题或建议，请通过以下方式联系我们：

- 邮箱: support@sweet.example.com
- 问题追踪: 请在GitHub上提交Issue 