# Sweet社交电商分销系统 - 数据库设计

## 概述

本目录包含Sweet社交电商分销系统的数据库表结构SQL文件。所有表均使用`sw_`作为表前缀，按照不同的功能模块进行了拆分，便于管理和维护。

## 文件说明

数据库SQL文件按照模块进行了拆分，每个文件对应一个功能模块：

1. **01_user.sql**: 用户模块，包含用户基本信息、地址、第三方授权、积分等表结构
2. **02_product.sql**: 商品模块，包含商品、分类、规格、SKU、评价等表结构
3. **03_order.sql**: 订单模块，包含订单、订单商品、订单地址、支付、物流、购物车等表结构
4. **04_payment.sql**: 支付模块，包含支付记录、退款记录、用户余额、充值等表结构
5. **05_distribution.sql**: 分销模块，包含分销商、分销关系、佣金记录、提现等表结构
6. **06_marketing.sql**: 营销模块，包含优惠券、秒杀、拼团等表结构
7. **07_system.sql**: 系统模块，包含管理员、角色权限、系统配置、日志等表结构
8. **08_content.sql**: 内容模块，包含文章、广告、通知等表结构
9. **09_community.sql**: 社区圈子模块，包含圈子、帖子、评论、关注等表结构

## 表前缀说明

所有表均使用`sw_`作为表前缀，表示Sweet系统。这样设计的好处：

1. 避免与其他系统的表名冲突
2. 便于识别表所属系统
3. 方便数据库管理和维护

## 数据库设计规范

1. 所有表必须包含`id`主键，类型为`bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT`
2. 所有表必须包含`created_at`和`updated_at`字段，用于记录创建和更新时间
3. 需要软删除功能的表增加`deleted_at`字段
4. 所有表必须使用`utf8mb4`字符集，以支持完整的Unicode字符
5. 所有表必须添加适当的索引，提高查询效率
6. 所有字段必须有注释，说明字段用途

## 使用方法

### 初始化数据库

1. 创建数据库

```sql
CREATE DATABASE sweet DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

2. 导入SQL文件

可以按照模块顺序依次导入SQL文件，也可以使用以下命令一次性导入所有SQL文件：

```bash
# 方法一：使用MySQL命令行工具
mysql -u用户名 -p密码 sweet < 01_user.sql
mysql -u用户名 -p密码 sweet < 02_product.sql
# ... 依次导入其他SQL文件

# 方法二：使用cat命令合并后导入
cat *.sql > all.sql
mysql -u用户名 -p密码 sweet < all.sql
```

### 数据库配置

在应用中配置数据库连接信息：

```go
config := database.Config{
    Master:         "user:pass@tcp(127.0.0.1:3306)/sweet?charset=utf8mb4&parseTime=True&loc=Local",
    Slaves:         []string{"user:pass@tcp(127.0.0.1:3307)/sweet?charset=utf8mb4&parseTime=True&loc=Local"},
    EnableLog:      true,
    EnableTrace:    true,
    SlowThreshold:  500, // 慢查询阈值500毫秒
    DriverType:     "mysql",
    QueryTimeout:   3000, // 3秒超时
    MaxRetries:     3,
    RetryDelay:     100, // 100毫秒
    Pool: database.ConnPool{
        MaxIdleConns:    10,
        MaxOpenConns:    100,
        ConnMaxLifetime: 3600, // 1小时
    },
}

db, err := database.Init(config)
if err != nil {
    panic(err)
}
```

## 表结构说明

每个SQL文件的开头都包含了该模块的表结构说明，包括表名、字段名、字段类型、字段说明等信息。可以通过查看SQL文件了解详细的表结构设计。

## 表关系说明

> 为了更直观地展示表之间的关系，我们提供了[数据库表关系图](./database_relationship.md)，建议与下面的文字说明结合阅读。

### 用户模块关系

1. **用户与地址关系**
   - 关系类型：一对多
   - 主表：`sw_user`（用户表）
   - 从表：`sw_user_address`（用户地址表）
   - 关联字段：`user_id`
   - 说明：一个用户可以有多个收货地址

2. **用户与第三方授权关系**
   - 关系类型：一对多
   - 主表：`sw_user`（用户表）
   - 从表：`sw_user_oauth`（用户第三方授权表）
   - 关联字段：`user_id`
   - 说明：一个用户可以绑定多个第三方账号

3. **用户与积分关系**
   - 关系类型：一对一
   - 主表：`sw_user`（用户表）
   - 从表：`sw_user_points`（用户积分表）
   - 关联字段：`user_id`
   - 说明：一个用户对应一条积分记录

4. **用户与余额关系**
   - 关系类型：一对一
   - 主表：`sw_user`（用户表）
   - 从表：`sw_user_balance`（用户余额表）
   - 关联字段：`user_id`
   - 说明：一个用户对应一条余额记录

### 商品模块关系

1. **商品分类层级关系**
   - 关系类型：自引用一对多
   - 表：`sw_product_category`（商品分类表）
   - 关联字段：`parent_id`
   - 说明：一个分类可以有多个子分类，形成树状结构

2. **商品与分类关系**
   - 关系类型：多对一
   - 主表：`sw_product_category`（商品分类表）
   - 从表：`sw_product`（商品表）
   - 关联字段：`category_id`
   - 说明：一个分类下可以有多个商品

3. **商品与SKU关系**
   - 关系类型：一对多
   - 主表：`sw_product`（商品表）
   - 从表：`sw_product_sku`（商品SKU表）
   - 关联字段：`product_id`
   - 说明：一个商品可以有多个SKU（库存单位）

4. **规格与规格值关系**
   - 关系类型：一对多
   - 主表：`sw_product_spec`（商品规格表）
   - 从表：`sw_product_spec_value`（商品规格值表）
   - 关联字段：`spec_id`
   - 说明：一个规格可以有多个规格值

5. **商品与规格关系**
   - 关系类型：多对多
   - 中间表：`sw_product_spec_relation`（商品规格关联表）
   - 关联字段：`product_id`、`spec_id`
   - 说明：一个商品可以有多个规格，一个规格可以应用于多个商品

6. **商品与属性值关系**
   - 关系类型：一对多
   - 主表：`sw_product`（商品表）
   - 从表：`sw_product_attribute_value`（商品属性值表）
   - 关联字段：`product_id`
   - 说明：一个商品可以有多个属性值

7. **属性与属性值关系**
   - 关系类型：一对多
   - 主表：`sw_product_attribute`（商品属性表）
   - 从表：`sw_product_attribute_value`（商品属性值表）
   - 关联字段：`attribute_id`
   - 说明：一个属性可以有多个属性值

8. **商品与评价关系**
   - 关系类型：一对多
   - 主表：`sw_product`（商品表）
   - 从表：`sw_product_comment`（商品评价表）
   - 关联字段：`product_id`
   - 说明：一个商品可以有多条评价

9. **商品与收藏关系**
   - 关系类型：一对多
   - 主表：`sw_product`（商品表）
   - 从表：`sw_product_favorite`（商品收藏表）
   - 关联字段：`product_id`
   - 说明：一个商品可以被多个用户收藏

10. **运费模板与规则关系**
    - 关系类型：一对多
    - 主表：`sw_product_freight_template`（商品运费模板表）
    - 从表：`sw_product_freight_rule`（商品运费模板规则表）
    - 关联字段：`template_id`
    - 说明：一个运费模板可以有多条运费规则

### 订单模块关系

1. **订单与订单商品关系**
   - 关系类型：一对多
   - 主表：`sw_order`（订单表）
   - 从表：`sw_order_item`（订单商品表）
   - 关联字段：`order_id`、`order_no`
   - 说明：一个订单包含多个商品项

2. **订单与收货地址关系**
   - 关系类型：一对一
   - 主表：`sw_order`（订单表）
   - 从表：`sw_order_address`（订单收货地址表）
   - 关联字段：`order_id`、`order_no`
   - 说明：一个订单对应一个收货地址

3. **用户与订单关系**
   - 关系类型：一对多
   - 主表：`sw_user`（用户表）
   - 从表：`sw_order`（订单表）
   - 关联字段：`user_id`
   - 说明：一个用户可以有多个订单

4. **商品与订单商品关系**
   - 关系类型：一对多
   - 主表：`sw_product`（商品表）
   - 从表：`sw_order_item`（订单商品表）
   - 关联字段：`product_id`
   - 说明：一个商品可以出现在多个订单中

### 支付模块关系

1. **订单与支付记录关系**
   - 关系类型：一对多
   - 主表：`sw_order`（订单表）
   - 从表：`sw_payment_record`（支付记录表）
   - 关联字段：`order_id`、`order_no`
   - 说明：一个订单可能有多条支付记录（如支付失败后重试）

2. **支付记录与退款记录关系**
   - 关系类型：一对多
   - 主表：`sw_payment_record`（支付记录表）
   - 从表：`sw_refund_record`（退款记录表）
   - 关联字段：`payment_no`
   - 说明：一笔支付可能有多笔退款（如部分退款）

3. **用户与余额变动记录关系**
   - 关系类型：一对多
   - 主表：`sw_user`（用户表）
   - 从表：`sw_balance_log`（余额变动记录表）
   - 关联字段：`user_id`
   - 说明：一个用户可以有多条余额变动记录

### 分销模块关系

1. **用户与分销商关系**
   - 关系类型：一对一
   - 主表：`sw_user`（用户表）
   - 从表：`sw_distributor`（分销商表）
   - 关联字段：`user_id`
   - 说明：一个用户可以成为一个分销商

2. **分销商层级关系**
   - 关系类型：自引用一对多
   - 表：`sw_distributor`（分销商表）
   - 关联字段：`parent_id`
   - 说明：一个分销商可以有多个下级分销商

3. **用户分销关系**
   - 关系类型：多对多
   - 表：`sw_distribution_relation`（分销关系表）
   - 关联字段：`user_id`、`parent_user_id`
   - 说明：记录用户之间的分销关系，一个用户可以有多个上级，也可以有多个下级

4. **分销等级与分销商关系**
   - 关系类型：一对多
   - 主表：`sw_distributor_level`（分销等级表）
   - 从表：`sw_distributor`（分销商表）
   - 关联字段：`level`
   - 说明：一个分销等级可以有多个分销商

### 营销模块关系

1. **优惠券与用户优惠券关系**
   - 关系类型：一对多
   - 主表：`sw_coupon`（优惠券表）
   - 从表：`sw_user_coupon`（用户优惠券表）
   - 关联字段：`coupon_id`
   - 说明：一种优惠券可以被多个用户领取

2. **用户与优惠券关系**
   - 关系类型：多对多
   - 中间表：`sw_user_coupon`（用户优惠券表）
   - 关联字段：`user_id`、`coupon_id`
   - 说明：一个用户可以领取多种优惠券，一种优惠券可以被多个用户领取

3. **秒杀活动与秒杀商品关系**
   - 关系类型：一对多
   - 主表：`sw_seckill_activity`（秒杀活动表）
   - 从表：`sw_seckill_product`（秒杀商品表）
   - 关联字段：`activity_id`
   - 说明：一个秒杀活动可以包含多个秒杀商品

4. **商品与秒杀商品关系**
   - 关系类型：一对多
   - 主表：`sw_product`（商品表）
   - 从表：`sw_seckill_product`（秒杀商品表）
   - 关联字段：`product_id`
   - 说明：一个商品可以参与多个秒杀活动

### 社区圈子模块关系

1. **圈子与成员关系**
   - 关系类型：多对多
   - 中间表：`sw_circle_member`（圈子成员表）
   - 关联字段：`circle_id`、`user_id`
   - 说明：一个圈子可以有多个成员，一个用户可以加入多个圈子

2. **圈子与帖子关系**
   - 关系类型：一对多
   - 主表：`sw_community_circle`（圈子表）
   - 从表：`sw_community_post`（帖子表）
   - 关联字段：`circle_id`
   - 说明：一个圈子可以有多个帖子

3. **用户与帖子关系**
   - 关系类型：一对多
   - 主表：`sw_user`（用户表）
   - 从表：`sw_community_post`（帖子表）
   - 关联字段：`user_id`
   - 说明：一个用户可以发布多个帖子

4. **帖子与评论关系**
   - 关系类型：一对多
   - 主表：`sw_community_post`（帖子表）
   - 从表：`sw_post_comment`（帖子评论表）
   - 关联字段：`post_id`
   - 说明：一个帖子可以有多条评论

5. **用户关注关系**
   - 关系类型：多对多
   - 表：`sw_user_follow`（用户关注表）
   - 关联字段：`user_id`、`follow_user_id`
   - 说明：记录用户之间的关注关系，一个用户可以关注多人，也可以被多人关注

## 注意事项

1. 在生产环境中使用前，请根据实际情况调整表结构，如字段长度、索引设计等
2. 建议在测试环境中先进行充分测试，确保表结构满足业务需求
3. 数据库设计会随着业务的发展而调整，请定期检查和优化表结构
4. 在进行数据库变更时，请做好备份工作，避免数据丢失