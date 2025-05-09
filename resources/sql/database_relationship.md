# Sweet社交电商分销系统 - 数据库表关系图

## 概述

本文档提供了Sweet社交电商分销系统各模块表之间的关系图，帮助开发人员更直观地理解数据库结构。

## 用户模块关系图

```
+-------------+       1:n       +------------------+
|   sw_user   |---------------->| sw_user_address  |
+-------------+                 +------------------+
       |                                
       | 1:n                             
       v                                
+----------------+                      
| sw_user_oauth  |                      
+----------------+                      
       ^                                
       | 1:1                             
       |                                
+----------------+       1:1       +----------------+
| sw_user_points |<--------------->| sw_user_balance|
+----------------+                 +----------------+
```

## 商品模块关系图

```
+---------------------+       1:n       +---------------------+
| sw_product_category |---------------->| sw_product_category |
+---------------------+  (parent_id)    +---------------------+
       |                                        
       | 1:n                                     
       v                                        
+-------------+       1:n       +----------------+
| sw_product  |---------------->| sw_product_sku |
+-------------+                 +----------------+
       |                                
       | 1:n                             
       v                                
+-------------------+                   
| sw_product_comment|                   
+-------------------+                   

+----------------+       1:n       +----------------------+
| sw_product_spec|---------------->| sw_product_spec_value|
+----------------+                 +----------------------+
       ^
       |
       | m:n (via sw_product_spec_relation)
       |
+-------------+       1:n       +----------------------------+
| sw_product  |---------------->| sw_product_attribute_value |
+-------------+                 +----------------------------+
       |                                ^
       |                                |
       |                                | 1:n
       |                                |
       |                         +--------------------+
       |                         | sw_product_attribute |
       |                         +--------------------+
       |
       | 1:n
       v
+-------------------+
| sw_product_favorite|
+-------------------+

+-------------------------+       1:n       +----------------------+
| sw_product_freight_template |---------------->| sw_product_freight_rule |
+-------------------------+                 +----------------------+
```

## 订单模块关系图

```
+-------------+       1:n       +----------------+
|   sw_user   |---------------->|    sw_order    |
+-------------+                 +----------------+
                                        |        
                                        | 1:n     
                                        v        
+-------------+       1:n       +----------------+
| sw_product  |---------------->| sw_order_item  |
+-------------+                 +----------------+
                                        ^
                                        |
                                        | 1:n
                                        |
                                +----------------+       1:1       +------------------+
                                |    sw_order   |---------------->| sw_order_address |
                                +----------------+                 +------------------+
```

## 支付模块关系图

```
+-------------+       1:n       +-------------------+
|   sw_order  |---------------->| sw_payment_record |
+-------------+                 +-------------------+
                                        |
                                        | 1:n
                                        v
                                +------------------+
                                | sw_refund_record |
                                +------------------+

+-------------+       1:1       +------------------+       1:n       +----------------+
|   sw_user   |---------------->| sw_user_balance  |---------------->| sw_balance_log |
+-------------+                 +------------------+                 +----------------+
```

## 分销模块关系图

```
+-------------+       1:1       +----------------+       1:n       +----------------+
|   sw_user   |---------------->| sw_distributor |---------------->| sw_distributor |
+-------------+                 +----------------+   (parent_id)   +----------------+
       |                                ^
       | 1:n                             |
       v                                | 1:n
+-------------------------+             |
| sw_distribution_relation|             |
+-------------------------+             |
                                        |
                                +---------------------+
                                | sw_distributor_level|
                                +---------------------+
```

## 营销模块关系图

```
+-------------+                 +------------+
|   sw_user   |                 |  sw_coupon |
+-------------+                 +------------+
       |                               |
       |                               |
       |            +------------------+
       |            |
       +----------->|<--------------+
                    | sw_user_coupon|
                    +---------------+
                    
+--------------------+       1:n       +-------------------+
| sw_seckill_activity|---------------->| sw_seckill_product |
+--------------------+                 +-------------------+
                                               ^
                                               |
                                               | 1:n
                                               |
                                        +-------------+
                                        | sw_product  |
                                        +-------------+
```

## 社区圈子模块关系图

```
+---------------------+                 +-------------+
| sw_community_circle |                 |   sw_user   |
+---------------------+                 +-------------+
       |                                      |
       |                                      |
       |            +------------------------+|
       |            |                         |
       +----------->|<-----------------------+|
                    |   sw_circle_member     ||
                    +------------------------+|
                                              |
                                              | 1:n
                                              v
+---------------------+       1:n       +------------------+       1:n       +----------------+
| sw_community_circle |---------------->| sw_community_post|---------------->| sw_post_comment|
+---------------------+                 +------------------+                 +----------------+
                                               ^
                                               |
                                               | 1:n
                                               |
                                        +-------------+
                                        |   sw_user   |
                                        +-------------+
                                              |
                                              | m:n
                                              v
                                      +----------------+
                                      | sw_user_follow |
                                      +----------------+
```

## 系统模块关系图

```
+----------+       1:n       +---------+
| sw_admin |---------------->| sw_role |
+----------+                 +---------+
                                  |
                                  | m:n (via sw_role_menu)
                                  v
+---------+       1:n       +---------+
| sw_menu |---------------->| sw_menu |
+---------+  (parent_id)    +---------+
```

\n## 说明

- **1:1** 表示一对一关系
- **1:n** 表示一对多关系
- **m:n** 表示多对多关系
- 箭头方向表示引用关系，指向被引用的表

这些关系图与README.md中的表关系说明相对应，提供了更直观的视觉表示，帮助开发人员理解表之间的关联。在实际开发中，可以参考这些关系图来正确实现业务逻辑和数据完整性约束。