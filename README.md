# Sweet（蜜趣商城）- 社交电商分销系统

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-1.9.0+-6C86E4?style=flat-square&logo=gin)](https://github.com/gin-gonic/gin)
[![gRPC](https://img.shields.io/badge/gRPC-1.60.0+-4285F4?style=flat-square&logo=google)](https://grpc.io/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

<p align="center">
  <img src="https://via.placeholder.com/200x200?text=Sweet+Logo" alt="Sweet Logo" width="200"/>
</p>

> 专注社交电商场景的分销商城系统，赋能企业快速构建数字化商业生态

## 📑 目录

- [概述](#概述)
- [核心功能](#核心功能)
- [技术架构](#技术架构)
- [项目结构](#项目结构)
- [开发环境](#开发环境)
- [快速开始](#快速开始)
- [API文档](#api文档)
- [数据模型](#数据模型)
- [安全策略](#安全策略)
- [扩展性设计](#扩展性设计)
- [性能优化](#性能优化)
- [国际化支持](#国际化支持)
- [灾备与监控](#灾备与监控)
- [开发示例](#开发示例)
- [贡献指南](#贡献指南)
- [开发规范](#开发规范)
- [版本计划](#版本计划)
- [许可证](#许可证)

## 📋 概述

**Sweet** 是一款专注于社交电商场景的分销商城系统，帮助企业快速搭建多层级分销体系，实现商品销售、裂变推广与团队管理的全链路数字化。系统采用模块化设计，提供开箱即用的功能模块与深度定制能力。

<details>
<summary><strong>🚀 为什么选择 Sweet?</strong></summary>

- **快速部署**：开箱即用的服务，分钟级完成商城搭建
- **灵活定制**：模块化设计，支持按需扩展与品牌定制
- **高性能**：优化的架构设计，支持大规模并发
- **数据驱动**：强大的数据分析能力，助力业务决策
- **生态开放**：丰富的 API 接口，轻松对接第三方系统
- **社交裂变**：内置社交分享与裂变工具，实现营销增长
- **安全可靠**：严格的数据安全措施和容灾备份方案
</details>

---

## 🔥 核心功能

### 1. 分销体系引擎

<table>
  <tr>
    <td>
      <strong>多级分销模式</strong><br/>
      支持 1-3 级分佣比例配置，灵活适应不同行业规则
    </td>
    <td>
      <strong>自动分账系统</strong><br/>
      订单成交后实时计算佣金，支持微信/支付宝自动提现
    </td>
  </tr>
  <tr>
    <td>
      <strong>团队裂变工具</strong><br/>
      邀请码、专属海报、裂变红包等多场景推广组件
    </td>
    <td>
      <strong>分销商等级管理</strong><br/>
      按业绩自动升级，差异化佣金比例与权益
    </td>
  </tr>
</table>

### 2. 商城管理

<table>
  <tr>
    <td>
      <strong>商品中心</strong><br/>
      支持实物/虚拟商品、多规格SKU、库存预警
    </td>
    <td>
      <strong>营销插件</strong><br/>
      拼团、秒杀、优惠券、积分兑换等 20+ 营销工具
    </td>
  </tr>
  <tr>
    <td>
      <strong>多端适配</strong><br/>
      自动生成 H5/小程序/PC 商城，支持品牌自定义UI
    </td>
    <td>
      <strong>订单中心</strong><br/>
      全流程订单跟踪，支持退款/售后/电子发票
    </td>
  </tr>
</table>

### 3. 数据赋能

<table>
  <tr>
    <td>
      <strong>实时数据看板</strong><br/>
      分销商业绩、商品热力图、用户行为分析
    </td>
    <td>
      <strong>智能风控系统</strong><br/>
      自动检测异常订单与违规分销行为
    </td>
  </tr>
  <tr>
    <td>
      <strong>BI 分析工具</strong><br/>
      客户画像生成、销售预测模型、ROI 分析报表
    </td>
    <td>
      <strong>决策助手</strong><br/>
      基于数据分析的营销决策建议与商品推荐
    </td>
  </tr>
</table>

### 4. 系统管理

<table>
  <tr>
    <td>
      <strong>权限管理</strong><br/>
      RBAC 角色权限控制，支持组织架构树形管理
    </td>
    <td>
      <strong>系统配置</strong><br/>
      灵活的系统参数配置，支持动态调整
    </td>
  </tr>
  <tr>
    <td>
      <strong>OpenAPI 接口</strong><br/>
      与 ERP/CRM/支付系统无缝对接
    </td>
    <td>
      <strong>数据字典</strong><br/>
      统一的数据字典管理，保障数据一致性
    </td>
  </tr>
</table>

---

## 🔧 技术架构

### 后端技术栈

<table>
  <tr>
    <th>类别</th>
    <th>技术选型</th>
    <th>说明</th>
  </tr>
  <tr>
    <td>编程语言</td>
    <td>Go 1.24+</td>
    <td>高性能、并发友好的语言</td>
  </tr>
  <tr>
    <td>微服务框架</td>
    <td>gRPC + Protocol Buffers</td>
    <td>高效的RPC框架与数据序列化</td>
  </tr>
  <tr>
    <td>API网关</td>
    <td>Gin</td>
    <td>轻量高效的HTTP框架，处理API请求</td>
  </tr>
  <tr>
    <td>ORM框架</td>
    <td>GORM</td>
    <td>功能强大的ORM库，简化数据库操作</td>
  </tr>
  <tr>
    <td rowspan="2">数据存储</td>
    <td>MySQL 8.0+</td>
    <td>主数据库</td>
  </tr>
  <tr>
    <td>Redis 6.0+</td>
    <td>缓存、分布式锁、会话管理</td>
  </tr>
  <tr>
    <td>消息队列</td>
    <td>Kafka</td>
    <td>异步处理订单、支付通知等</td>
  </tr>
  <tr>
    <td>服务发现</td>
    <td>etcd</td>
    <td>服务注册与发现</td>
  </tr>
  <tr>
    <td>认证授权</td>
    <td>JWT + Casbin</td>
    <td>用户认证与权限管理</td>
  </tr>
  <tr>
    <td>配置中心</td>
    <td>Nacos</td>
    <td>统一配置管理，动态配置下发</td>
  </tr>
  <tr>
    <td>熔断降级</td>
    <td>Sentinel</td>
    <td>自适应熔断策略，流量控制</td>
  </tr>
  <tr>
    <td>API网关</td>
    <td>APISIX</td>
    <td>统一API入口，鉴权、限流、路由</td>
  </tr>
  <tr>
    <td>日志收集</td>
    <td>ELK Stack</td>
    <td>日志收集、分析与可视化</td>
  </tr>
  <tr>
    <td>文档工具</td>
    <td>Swagger Docs</td>
    <td>API文档自动生成</td>
  </tr>
  <tr>
    <td rowspan="2">监控系统</td>
    <td>Prometheus</td>
    <td>指标收集</td>
  </tr>
  <tr>
    <td>Grafana</td>
    <td>可视化监控</td>
  </tr>
  <tr>
    <td rowspan="2">CI/CD</td>
    <td>Docker + Kubernetes</td>
    <td>容器化与编排管理</td>
  </tr>
  <tr>
    <td>GitLab CI/GitHub Actions</td>
    <td>自动化构建部署</td>
  </tr>
</table>

### 前端技术栈

<table>
  <tr>
    <th>应用场景</th>
    <th>技术选型</th>
    <th>说明</th>
  </tr>
  <tr>
    <td>管理后台</td>
    <td>Vue3 + Element Plus</td>
    <td>高效企业级中后台解决方案</td>
  </tr>
  <tr>
    <td>小程序</td>
    <td>uni-app</td>
    <td>多端统一开发框架</td>
  </tr>
  <tr>
    <td>H5商城</td>
    <td>Vue3 + Vant</td>
    <td>移动端UI组件库</td>
  </tr>
  <tr>
    <td>PC端</td>
    <td>Next.js</td>
    <td>支持SSR的React框架</td>
  </tr>
</table>

### 系统架构图

```
+----------------------------------+
|            API网关层              |
|                                  |
+----------------------------------+
               |
+----------------------------------+
|           微服务集群               |
|                                  |
|  +-------------+  +------------+ |
|  | 用户服务集群 |  | 商品服务集群 |  |
|  +-------------+  +------------+ |
|                                  |
|  +-------------+  +------------+ |
|  | 订单服务集群 |  | 支付服务集群 | |
|  +-------------+  +------------+ |
|                                  |
|  +-------------+  +------------+ |
|  | 分销服务集群 |  | 营销服务集群 | |
|  +-------------+  +------------+ |
|                                  |
+----------------------------------+
               |
+----------------------------------+
|           中间件服务层            |
| +------+ +----+ +-----+ +------+ |
| |Redis | |etcd | |Kafka| |Trace | |
| +------+ +----+ +-----+ +------+ |
+----------------------------------+
               |
+----------------------------------+
|            数据存储层             |
|  +--------+          +-------+   |
|  | MySQL  |          |OSS/S3 |   |
|  +--------+          +-------+   |
+----------------------------------+
```

---

## 📖 API文档

API文档由swagger自动生成，启动服务后访问：
```
http://localhost:8888/docs
```

## 🗄️ 数据模型

我们采用了基于DDD（领域驱动设计）的分层架构，并结合微服务框架特性:

<table>
  <tr>
    <th>层次</th>
    <th>职责</th>
    <th>对应组件</th>
  </tr>
  <tr>
    <td>表示层</td>
    <td>处理HTTP请求响应</td>
    <td>API服务</td>
  </tr>
  <tr>
    <td>应用层</td>
    <td>协调业务流程，编排服务</td>
    <td>API Handler、API逻辑</td>
  </tr>
  <tr>
    <td>领域层</td>
    <td>核心业务逻辑与规则</td>
    <td>RPC服务逻辑、领域服务</td>
  </tr>
  <tr>
    <td>基础设施层</td>
    <td>提供技术支持</td>
    <td>Model层、中间件、工具包</td>
  </tr>
</table>

### 微服务架构设计

Sweet采用基于DDD(领域驱动设计)的微服务架构，使用gRPC作为服务间通信协议，主要服务划分如下：

1. **用户服务 (User Service)**
   - 用户注册、登录、认证
   - 会员管理、用户画像
   - 地址管理
   - 用户积分与等级管理
   - 第三方账号绑定与授权

2. **商品服务 (Product Service)**
   - 商品管理、分类、规格
   - 库存管理、商品审核
   - 商品搜索、推荐
   - 商品评价与问答
   - 商品导入导出

3. **订单服务 (Order Service)**
   - 订单创建、支付、履约
   - 购物车管理
   - 售后、退款管理
   - 发票管理
   - 订单统计与分析

4. **支付服务 (Payment Service)**
   - 支付渠道管理
   - 支付处理、退款
   - 支付回调、对账
   - 资金流水管理
   - 支付安全与风控

5. **分销服务 (Distribution Service)**
   - 分销关系管理
   - 佣金计算与结算
   - 分销商管理
   - 分销营销工具
   - 分销数据分析

6. **营销服务 (Marketing Service)**
   - 优惠券、满减活动
   - 秒杀、拼团活动
   - 积分商城、签到奖励
   - 会员营销
   - 活动数据分析

7. **内容服务 (Content Service)**
   - 文章、通知公告
   - 广告、轮播图管理
   - 评论、用户互动
   - SEO内容管理
   - 富媒体内容存储

8. **搜索服务 (Search Service)**
   - 商品全文检索
   - 智能推荐
   - 搜索热词
   - 搜索结果优化
   - 相关性排序

9. **数据分析服务 (Analytics Service)**
   - 业务数据统计
   - 用户行为分析
   - 营销效果分析
   - 实时数据看板
   - 数据挖掘与预测

10. **系统管理服务 (System Service)**
    - 管理员账户、权限
    - 系统配置、字典管理
    - 日志审计
    - 系统监控
    - 定时任务管理

#### 服务间通信模式

- **同步通信**: 基于gRPC的直接服务调用
- **异步通信**: 基于Kafka的事件驱动架构
- **状态共享**: 基于Redis的分布式缓存

#### 服务治理

- **服务注册与发现**: 使用etcd实现自动服务注册与发现
- **负载均衡**: 客户端负载均衡
- **熔断降级**: 使用Sentinel实现自适应熔断保护
- **限流**: 基于令牌桶算法的API限流
- **链路追踪**: 基于OpenTelemetry的全链路追踪

### 数据库设计

系统共划分为10个主要模块，每个模块对应一组相关功能的数据库表：

<details>
<summary><strong>01 系统管理模块</strong></summary>

系统管理模块包含系统基础配置、用户角色权限相关的表：

- `sys_user` - 系统用户表
- `sys_role` - 角色表
- `sys_menu` - 菜单表
- `sys_user_role` - 用户角色关联表
- `sys_role_menu` - 角色菜单关联表
- `sys_dept` - 部门表
- `sys_dict` - 数据字典表
- `sys_dict_item` - 字典项表
- `sys_login_log` - 登录日志表
- `sys_operation_log` - 操作日志表
- `sys_notice` - 通知公告表
</details>

<details>
<summary><strong>02 用户模块</strong></summary>

用户模块包含会员用户信息、收货地址等相关表：

- `biz_user_profile` - 用户资料表
- `biz_user_account` - 用户账户表
- `biz_user_level` - 用户等级表
- `biz_user_address` - 用户收货地址表
- `biz_user_favorite` - 用户收藏表
- `biz_user_browse_history` - 用户浏览历史表
- `biz_user_point` - 用户积分表
- `biz_user_point_record` - 积分记录表
- `biz_user_tags` - 用户标签表
- `biz_user_tag_relation` - 用户标签关联表
</details>

<details>
<summary><strong>03 商品模块</strong></summary>

商品模块包含商品信息、分类、规格等相关表：

- `biz_product` - 商品表
- `biz_product_category` - 商品分类表
- `biz_product_sku` - 商品SKU表
- `biz_product_attribute` - 商品属性表
- `biz_product_attribute_value` - 属性值表
- `biz_product_spec` - 规格表
- `biz_product_spec_value` - 规格值表
- `biz_product_brand` - 品牌表
- `biz_product_comment` - 商品评价表
- `biz_product_stock` - 库存表
- `biz_product_collection` - 商品专题表
- `biz_product_collection_relation` - 专题商品关联表
</details>

<details>
<summary><strong>04 订单模块</strong></summary>

订单模块包含订单信息、支付、物流等相关表：

- `biz_order` - 订单表
- `biz_order_item` - 订单商品表
- `biz_order_delivery` - 订单物流表
- `biz_order_refund` - 退款申请表
- `biz_pay_transaction` - 支付交易表
- `biz_pay_notify` - 支付回调表
- `biz_pay_refund` - 退款记录表
- `biz_cart` - 购物车表
- `biz_logistics_company` - 物流公司表
- `biz_delivery_template` - 运费模板表
</details>

<details>
<summary><strong>05 营销模块</strong></summary>

营销模块包含优惠券、秒杀、拼团等促销相关表：

- `biz_coupon` - 优惠券表
- `biz_coupon_template` - 优惠券模板表
- `biz_coupon_user` - 用户优惠券表
- `biz_seckill` - 秒杀活动表
- `biz_seckill_product` - 秒杀商品表
- `biz_group_buy` - 拼团活动表
- `biz_group_buy_record` - 拼团记录表
- `biz_discount` - 满减活动表
- `biz_discount_product` - 满减商品关联表
- `biz_gift` - 赠品表
- `biz_gift_rule` - 赠品规则表
</details>

<details>
<summary><strong>06 分销模块</strong></summary>

分销模块包含分销关系、佣金等相关表：

- `biz_distributor` - 分销员表
- `biz_distributor_apply` - 分销员申请表
- `biz_distributor_level` - 分销员等级表
- `biz_distributor_relation` - 分销关系表
- `biz_commission` - 佣金记录表
- `biz_commission_withdraw` - 佣金提现表
- `biz_commission_rule` - 佣金规则表
- `biz_commission_settlement` - 佣金结算表
- `biz_distributor_product` - 分销商品表
- `biz_distributor_statistics` - 分销员统计表
</details>

<details>
<summary><strong>07 内容管理模块</strong></summary>

内容管理模块包含文章、广告、导航等内容相关表：

- `biz_article_category` - 文章分类表
- `biz_article` - 文章表
- `biz_article_tag` - 文章标签表
- `biz_article_tag_relation` - 文章-标签关联表
- `biz_ad_position` - 广告位表
- `biz_ad` - 广告表
- `biz_navigation` - 导航菜单表
- `biz_page_design` - 页面装修表
- `biz_notice` - 通知公告表
- `biz_help_category` - 帮助中心分类表
- `biz_help_article` - 帮助中心文章表
</details>

<details>
<summary><strong>08 社区圈子模块</strong></summary>

社区圈子模块包含社交互动、内容分享相关表：

- `biz_community_group` - 社区圈子表
- `biz_community_category` - 圈子分类表
- `biz_community_group_category` - 圈子-分类关联表
- `biz_community_member` - 圈子成员表
- `biz_community_post` - 圈子帖子表
- `biz_community_post_like` - 帖子点赞表
- `biz_community_post_favorite` - 帖子收藏表
- `biz_community_comment` - 帖子评论表
- `biz_community_comment_like` - 评论点赞表
- `biz_community_follow` - 用户关注表
- `biz_community_fan` - 用户粉丝表
- `biz_community_topic` - 话题表
- `biz_community_post_topic` - 帖子-话题关联表
- `biz_community_topic_follow` - 话题关注表
</details>

<details>
<summary><strong>09 统计分析模块</strong></summary>

统计分析模块包含各类数据统计表：

- `biz_user_statistics` - 用户统计表
- `biz_product_sales_statistics_day` - 商品销售统计表(按日)
- `biz_product_sales_statistics_month` - 商品销售统计表(按月)
- `biz_order_statistics_day` - 订单统计表(按日)
- `biz_order_statistics_month` - 订单统计表(按月)
- `biz_distribution_statistics_day` - 分销统计表(按日)
- `biz_distribution_statistics_month` - 分销统计表(按月)
- `biz_community_statistics_day` - 社区统计表(按日)
- `biz_community_statistics_month` - 社区统计表(按月)
- `biz_visit_statistics_day` - 访问统计表(按日)
- `biz_user_active_statistics` - 用户活跃统计表
- `biz_channel_statistics_day` - 渠道统计表(按日)
</details>

<details>
<summary><strong>10 通用功能模块</strong></summary>

通用功能模块包含系统共用功能相关表：

- `com_file` - 文件信息表
- `com_sms_record` - 短信记录表
- `com_message_template` - 消息通知模板表
- `com_message` - 消息通知表
- `com_user_message` - 用户消息表
- `com_config` - 系统配置表
- `com_job` - 定时任务表
- `com_job_log` - 任务执行日志表
- `com_cache` - 缓存管理表
- `com_region` - 地区表
- `com_friend_link` - 友情链接表
- `com_system_log` - 系统日志表
</details>

---

## 🔒 安全策略

### 数据安全
- **敏感数据加密**：使用 AES-256 算法加密存储敏感信息（支付凭证、个人身份信息等）
- **传输安全**：全站 HTTPS，API 通信采用 TLS 1.3
- **防注入攻击**：参数校验和过滤，防止SQL注入和XSS攻击

### 认证与授权
- **多因素认证**：支持短信/邮箱验证码、TOTP 等二次验证
- **JWT认证**：标准JWT令牌认证机制
- **细粒度权限控制**：基于Casbin的RBAC模型权限系统
- **登录保护**：登录失败次数限制，异地登录提醒，风险操作二次验证

### 服务安全
- **API网关防护**：限流、熔断、黑白名单
- **服务间安全通信**：RPC服务通过TLS加密
- **请求追踪**：全链路请求追踪，异常快速定位

### 合规性
- **数据留存策略**：符合相关法规要求的数据留存与删除机制
- **操作审计**：所有关键操作留存审计日志，支持审计追溯
- **隐私保护**：用户同意管理，数据使用透明化

---

## 🔌 扩展性设计

### 微服务扩展性
- **服务自治**：各微服务独立开发、部署、扩展
- **服务发现**：基于etcd的自动服务发现与注册
- **负载均衡**：客户端负载均衡

### 容量扩展
- **水平扩展**：无状态服务支持多实例水平扩展
- **分库分表**：支持数据分片策略
- **弹性伸缩**：基于K8s的Pod自动扩缩容

### 功能扩展
- **插件化设计**：核心功能支持插件化扩展
- **配置中心**：动态配置，无需重启即可生效
- **服务编排**：基于工作流的业务流程编排

---

## ⚡ 性能优化

### 项目优化
- **Redis缓存**：利用redis做缓存策略
- **自适应熔断**：保护系统免受级联故障
- **限流保护**：API级别的限流策略

### 服务优化
- **本地缓存+分布式缓存**：多级缓存减轻数据库压力
- **批处理请求**：批量API减少网络开销
- **异步处理**：非即时反馈业务异步化处理
- **资源池化**：连接池复用，减少资源创建开销

### 前端优化
- **资源压缩**：JS/CSS 压缩，图片优化
- **懒加载**：图片和组件懒加载
- **CDN 加速**：静态资源 CDN 分发
- **SSR/预渲染**：关键页面服务端渲染

---

## 🌎 国际化支持

### 语言支持
- **多语言系统**：支持简体中文、繁体中文、英语、日语、韩语等多种语言
- **动态语言包**：支持在线更新语言包，无需重新部署
- **本地化UI**：根据用户语言自动调整界面布局和设计风格
- **翻译工具链**：提供翻译管理平台，支持多人协作翻译

### 区域适配
- **多币种支持**：支持人民币、美元、欧元、日元等多种货币结算
- **时区处理**：数据存储 UTC 时间，展示根据用户时区动态转换
- **区域特性**：支持不同地区的税费规则、营销规则
- **本地化内容**：根据用户所在地区推送定制化内容和活动

### 支付适配
- **多支付渠道**：支持微信支付、支付宝、Stripe、PayPal、Apple Pay等全球支付方式
- **结算方式**：支持多种结算周期和方式
- **跨境支付**：支持跨境电商场景下的支付解决方案
- **分账处理**：支持多商户分账、跨境分账处理

### 合规适配
- **GDPR合规**：符合欧盟通用数据保护条例要求
- **地区法规**：自动适配不同地区的数据隐私和电商法规
- **税务处理**：支持自动计算不同国家/地区的税费
- **出口合规**：支持跨境电商的出口清关、关税处理

---

## 🔄 灾备与监控

### 容灾设计
- **多可用区部署**：跨可用区高可用架构
- **数据备份**：自动定时备份 + 手动备份，支持时间点恢复
- **故障转移**：自动检测与故障转移机制

### 全链路监控
- **业务监控**：核心业务指标实时监控（订单量、支付成功率等）
- **系统监控**：基础设施监控（CPU、内存、磁盘等）
- **API 监控**：接口响应时间、错误率监控
- **用户体验监控**：前端性能监控、用户行为分析

### 告警系统
- **多级告警**：基于严重程度的多级告警机制
- **智能告警**：异常检测算法，减少误报
- **通知渠道**：支持短信、邮件、企业微信等多种告警通知渠道

---

## 🤝 贡献指南

我们欢迎所有形式的贡献，无论是新功能、文档改进还是问题修复。

<details>
<summary><strong>贡献流程</strong></summary>

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request
</details>

## 📏 开发规范

- **代码风格**：遵循 [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- **微服务设计**：遵循 [The Twelve-Factor App](https://12factor.net/zh_cn/) 和微服务最佳实践
- **提交信息**：遵循 [Conventional Commits](https://www.conventionalcommits.org/zh-hans/v1.0.0/)
- **API设计**：RESTful API设计规范和gRPC API设计标准
- **文档规范**：API文档使用Swagger，代码注释完整
- **测试规范**：单元测试覆盖率>80%，集成测试覆盖核心流程

---

## 🚀 版本计划

| 版本 | 特性 | 状态 |
|------|------|------|
| v0.1.0 | 基础框架搭建 | 进行中 |
| v0.2.0 | 用户、商品、订单核心功能 | 计划中 |
| v0.3.0 | 分销体系引擎 | 计划中 |
| v0.4.0 | 营销插件系统 | 计划中 |
| v0.5.0 | 数据分析与可视化 | 计划中 |
| v0.6.0 | 社区互动系统 | 计划中 |
| v0.7.0 | 移动端/小程序应用 | 计划中 |
| v0.8.0 | 国际化与多语言 | 计划中 |
| v0.9.0 | 性能优化与安全加固 | 计划中 |
| v1.0.0 | 首个稳定版本发布 | 计划中 |

---

## 📄 许可证

[MIT License](LICENSE)

<p align="center">
  <sub>Built with ❤️ by Sweet Team</sub>
</p>