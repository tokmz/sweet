# Sweet 后台管理系统

## 项目简介

Sweet 是一个基于 Go 语言开发的现代化后台管理系统框架，目前处于开发初期阶段。项目提供了完整的基础工具包和架构设计，为后续业务功能开发奠定了坚实的基础。

## 核心特性

### 🚀 技术特性
- **高性能架构**: 基于 Go 语言和 Gin 框架，支持高并发处理
- **微服务设计**: 模块化架构，易于扩展和维护
- **云原生支持**: 支持 Docker 容器化部署和 Kubernetes 编排
- **完善的监控**: 集成日志、链路追踪和性能监控
- **安全可靠**: JWT 认证、数据加密、权限控制

### 💼 项目特性
- **基础工具包**: 完整的认证、缓存、数据库、日志等工具包
- **架构设计**: 清晰的分层架构和项目结构
- **开发规范**: 完善的代码规范和开发指南
- **扩展性**: 模块化设计，便于后续功能扩展
- **文档完善**: 详细的使用文档和示例代码

## 技术架构

### 技术栈

**核心技术**
- **语言**: Go 1.21+
- **框架**: Gin (HTTP 框架)
- **数据库**: 支持 MySQL / PostgreSQL (通过GORM)
- **缓存**: Redis (完整封装)
- **配置管理**: Viper (多格式支持)
- **认证**: JWT (完整实现)
- **日志**: Zap (异步日志)
- **加密**: 多种加密算法支持

**开发工具**
- **链路追踪**: OpenTelemetry 集成
- **错误处理**: 统一错误处理机制
- **工具函数**: 丰富的工具函数库

### 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    应用层 (待开发)                           │
│                  API接口 + 业务逻辑                          │
├─────────────────────────────────────────────────────────────┤
│                    基础工具包 (已完成)                        │
│  认证JWT  │  缓存Redis │  数据库GORM │  日志Zap  │  配置Viper │
├─────────────────────────────────────────────────────────────┤
│                    基础设施层                                │
│    数据库    │     缓存     │     配置     │     日志      │
└─────────────────────────────────────────────────────────────┘
```

## 项目结构

```
sweet/
├── cmd/                    # 应用程序入口
│   └── server.go           # 服务器启动文件
├── internal/               # 私有应用代码
│   ├── api/               # API 处理器
│   ├── service/           # 业务服务层
│   ├── router/            # 路由配置
│   ├── middleware/        # 中间件
│   ├── models/            # 数据模型
│   │   ├── dto/          # 数据传输对象
│   │   └── http.go       # HTTP 响应模型
│   └── global/           # 全局变量
├── pkg/                   # 可复用的库代码
│   ├── auth/             # 认证授权
│   ├── cache/            # 缓存操作
│   ├── config/           # 配置管理
│   ├── crypto/           # 加密工具
│   ├── database/         # 数据库操作
│   ├── logger/           # 日志系统
│   ├── errs/             # 错误处理
│   └── utils/            # 工具函数
├── common/                # 共享组件
│   └── gin.go            # Gin 框架配置
├── resource/              # 资源文件
├── scripts/               # 脚本文件
├── test/                  # 测试文件
├── go.mod                 # Go 模块文件
├── go.sum                 # 依赖校验文件
└── README.md              # 项目文档
```

## 快速开始

### 环境要求

- Go 1.21 或更高版本
- MySQL 8.0+ 或 PostgreSQL 13+
- Redis 6.0+
- Docker (可选)

### 安装依赖

```bash
# 克隆项目
git clone https://github.com/your-org/sweet.git
cd sweet

# 下载依赖
go mod download
```

### 配置文件

项目使用 Viper 进行配置管理，支持多种配置格式：

```bash
# 查看配置示例
cat pkg/config/example.yaml
cat pkg/config/example.json

# 创建本地配置文件
cp pkg/config/example.yaml config.yaml
# 根据需要修改配置
vim config.yaml
```

### 运行示例

```bash
# 运行配置管理示例
go run pkg/config/example.go

# 运行缓存操作示例
go run pkg/cache/example.go

# 运行数据库示例
go run pkg/database/example.go

# 运行日志示例
go run pkg/logger/example.go

# 运行认证示例
go run pkg/auth/README.md # 查看认证使用示例
```

### 开发环境搭建

```bash
# 安装开发工具
go install github.com/cosmtrek/air@latest
go install github.com/swaggo/swag/cmd/swag@latest

# 启动热重载开发
air
```

## 已完成模块

### 1. 认证授权模块 (pkg/auth)
- JWT Token 生成和验证
- 权限检查和角色管理
- 用户认证中间件
- 完整的权限控制体系

### 2. 缓存模块 (pkg/cache)
- Redis 客户端封装
- 字符串、哈希、列表、集合操作
- 分布式锁和计数器
- 链路追踪集成

### 3. 数据库模块 (pkg/database)
- GORM 封装和配置
- 读写分离支持
- 慢查询监控
- 链路追踪集成

### 4. 日志模块 (pkg/logger)
- Zap 日志封装
- 异步日志写入
- 结构化日志输出
- 上下文日志支持

### 5. 配置模块 (pkg/config)
- Viper 配置管理
- 多格式配置文件支持
- 环境变量绑定
- 配置热重载

### 6. 工具模块
- **加密工具** (pkg/crypto): 密码哈希和验证
- **错误处理** (pkg/errs): 统一错误定义和处理
- **工具函数** (pkg/utils): 常用工具函数库
- **HTTP模型** (internal/models): 请求响应结构体

## 开发计划

### 待开发功能

**第一阶段：基础功能**
- 用户管理系统
- 角色权限管理
- 系统配置管理
- 基础API接口

**第二阶段：高级功能**
- 数据统计分析
- 日志审计系统
- 文件管理系统
- 系统监控面板

**第三阶段：扩展功能**
- 工作流引擎
- 消息通知系统
- 数据导入导出
- API文档生成

### API 设计规范

项目将遵循 RESTful API 设计规范：

```http
# 认证相关
POST   /api/v1/auth/login     # 用户登录
POST   /api/v1/auth/refresh   # 刷新Token
POST   /api/v1/auth/logout    # 用户登出

# 用户管理
GET    /api/v1/users          # 获取用户列表
GET    /api/v1/users/:id      # 获取用户详情
POST   /api/v1/users          # 创建用户
PUT    /api/v1/users/:id      # 更新用户
DELETE /api/v1/users/:id      # 删除用户

# 角色权限
GET    /api/v1/roles          # 获取角色列表
POST   /api/v1/roles          # 创建角色
PUT    /api/v1/roles/:id      # 更新角色
DELETE /api/v1/roles/:id      # 删除角色
```

## 开发指南

### 代码规范

- 遵循 Go 语言官方代码规范
- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码质量
- 所有公开函数必须有注释
- 错误处理不能忽略

### 提交规范

```
type(scope): description

type: feat, fix, docs, style, refactor, test, chore
scope: 影响的模块
description: 简短描述

示例:
feat(user): add user registration API
fix(order): fix order status update bug
docs(readme): update installation guide
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 运行基准测试
go test -bench=. ./...
```

## 部署指南

### 生产环境部署

1. **环境准备**
   - 服务器配置：4核8G内存起步
   - 数据库：MySQL 主从配置
   - 缓存：Redis 集群
   - 负载均衡：Nginx

2. **配置优化**
   - 数据库连接池配置
   - Redis 连接池配置
   - 日志级别调整
   - 性能监控配置

3. **安全配置**
   - HTTPS 证书配置
   - 防火墙规则
   - 数据库访问控制
   - API 限流配置

### 监控告警

- **应用监控**: 使用 Prometheus 收集应用指标
- **日志监控**: 使用 ELK 进行日志分析
- **链路追踪**: 使用 Jaeger 进行分布式追踪
- **告警通知**: 集成钉钉、企业微信等通知渠道

## 性能优化

### 数据库优化
- 合理设计索引
- 读写分离
- 分库分表
- 连接池优化

### 缓存策略
- Redis 缓存热点数据
- 本地缓存减少网络开销
- 缓存预热和更新策略
- 缓存穿透和雪崩防护

### 应用优化
- Goroutine 池管理
- 内存复用
- HTTP 连接复用
- 静态资源 CDN 加速

## 常见问题

### Q: 如何配置数据库连接？
A: 使用 `pkg/database` 包，参考 `pkg/database/example.go` 中的配置示例，支持 MySQL 和 PostgreSQL。

### Q: 如何使用缓存功能？
A: 使用 `pkg/cache` 包，参考 `pkg/cache/example.go` 中的使用示例，支持 Redis 的各种数据类型操作。

### Q: 如何配置日志系统？
A: 使用 `pkg/logger` 包，参考 `pkg/logger/example.go` 中的配置示例，支持异步日志和结构化输出。

### Q: 如何实现JWT认证？
A: 使用 `pkg/auth` 包，查看 `pkg/auth/README.md` 了解完整的认证和权限管理功能。

### Q: 如何添加新的工具包？
A: 在 `pkg/` 目录下创建新的包，遵循项目的代码规范和文档标准。

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 版本历史

- **v0.2.0** (计划中)
  - 基础API框架搭建
  - 用户管理功能实现
  - 权限管理功能实现
  - 系统配置功能实现

- **v0.1.0** (当前)
  - 项目初始化和架构设计
  - 认证授权工具包 (pkg/auth)
  - 缓存操作工具包 (pkg/cache)
  - 数据库操作工具包 (pkg/database)
  - 日志系统工具包 (pkg/logger)
  - 配置管理工具包 (pkg/config)
  - 加密工具包 (pkg/crypto)
  - 错误处理工具包 (pkg/errs)
  - 基础HTTP模型 (internal/models)
  - 项目文档和开发规范

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系我们

- 项目主页: https://github.com/your-org/sweet
- 问题反馈: https://github.com/your-org/sweet/issues
- 邮箱: dev@sweet.com
- 微信群: 扫描二维码加入开发者群

---

**Sweet 团队** ❤️ 用心打造优质的后台管理系统框架