# Sweet 工具包

这个包提供了一系列常用的工具函数，用于简化开发过程中的常见操作。

## 指针工具 (ptr.go)

为了解决 Go 语言中指针类型的安全处理问题，提供了一系列指针操作辅助函数。

### 从指针获取值

这些函数可以安全地从可能为 nil 的指针中获取值，如果指针为 nil，则返回零值。

```go
// 从 *int64 获取 int64 值，nil 返回 0
value := utils.SafeInt64(pointerValue)

// 从 *int 获取 int 值，nil 返回 0
value := utils.SafeInt(pointerValue)

// 从 *string 获取 string 值，nil 返回 ""
value := utils.SafeString(pointerValue)

// 从 *bool 获取 bool 值，nil 返回 false
value := utils.SafeBool(pointerValue)
```

### 类型转换

这些函数用于在指针类型之间进行转换，特别是数据库中常用的整型转布尔值。

```go
// 从 *int64 类型转为 bool 值，1 为 true，其他为 false
value := utils.SafeBoolFromInt64(pointerValue)

// 从 *int 类型转为 bool 值，1 为 true，其他为 false
value := utils.SafeBoolFromInt(pointerValue)
```

### 创建指针

将值转换为对应的指针类型。

```go
// 创建指向值的指针
strPtr := utils.ToPtr("string value")
intPtr := utils.ToPtr(123)
boolPtr := utils.ToPtr(true)
```

## 切片工具 (slice.go)

这个包扩展了 Go 1.21+ 的标准库 `slices` 包，提供了一些标准库中没有的切片操作函数。对于标准库已有的函数，为了兼容性我们提供了简单的包装。

> **注意**: 如果你的项目使用 Go 1.21 或更高版本，建议优先使用标准库 `slices` 包中的函数。

### 与标准库重叠的函数

以下函数是对标准库的简单封装，为了向后兼容性而保留：

```go
// 检查切片是否包含指定元素（使用标准库 slices.Contains）
found := utils.Contains([]string{"a", "b", "c"}, "b")

// 使用自定义函数检查切片是否包含符合条件的元素（使用标准库 slices.ContainsFunc）
found := utils.ContainsFunc(users, func(u User) bool {
    return u.Age > 18
})
```

### 标准库中不存在的函数

以下函数是对标准库的补充，提供了更多实用的操作：

#### 转换和过滤

```go
// 将切片中的每个元素转换为另一个类型
names := utils.Map(users, func(u User) string {
    return u.Name
})

// 过滤切片中符合条件的元素
adults := utils.Filter(users, func(u User) bool {
    return u.Age >= 18
})
```

#### 聚合和去重

```go
// 累加切片中的所有数字
sum := utils.Reduce([]int{1, 2, 3, 4}, 0, func(acc, val int) int {
    return acc + val
})

// 去除切片中的重复元素
unique := utils.Unique([]int{1, 2, 2, 3, 3, 3})
```

#### 分割和扁平化

```go
// 将切片分割成多个子切片，每个子切片最多包含3个元素
chunks := utils.Chunk([]int{1, 2, 3, 4, 5, 6, 7}, 3)
// 结果: [[1 2 3] [4 5 6] [7]]

// 将二维切片扁平化为一维切片
flat := utils.Flatten([][]int{{1, 2}, {3, 4}, {5, 6}})
// 结果: [1 2 3 4 5 6]
```

## 字符串工具 (string.go)

提供了一系列字符串处理的实用工具函数，包括判断、转换、生成和操作等功能。

### 判断和检查

```go
// 检查字符串是否为空
isEmpty := utils.IsEmpty("") // true

// 检查字符串是否为空白
isBlank := utils.IsBlank("   ") // true

// 检查字符串是否包含任意一个子字符串
hasAny := utils.ContainsAny("hello world", "foo", "world") // true

// 检查字符串是否包含所有子字符串
hasAll := utils.ContainsAll("hello world", "hello", "world") // true
```

### 命名风格转换

```go
// 驼峰命名转换为蛇形命名
snake := utils.ToSnakeCase("HelloWorld") // "hello_world"

// 蛇形命名转换为驼峰命名
camel := utils.ToCamelCase("hello_world") // "helloWorld"

// 蛇形命名转换为帕斯卡命名
pascal := utils.ToPascalCase("hello_world") // "HelloWorld"
```

### 字符串操作

```go
// 将字符串截断到指定长度
truncated := utils.Truncate("Hello World", 5) // "Hello"

// 移除字符串中的多余空格
trimmed := utils.RemoveExtraSpaces("  Hello   World  ") // "Hello World"

// 反转字符串
reversed := utils.ReverseString("Hello") // "olleH"
```

### 随机字符串生成

```go
// 生成指定长度的随机字符串
random := utils.RandomString(8) // 例如: "a1b2C3D4"

// 使用指定字符集生成随机字符串
randomNum := utils.RandomStringWithCharset(6, "0123456789") // 例如: "123456"
```

## 类型转换工具 (conv.go)

提供了一系列类型转换函数，用于在不同数据类型之间进行安全转换。

### 转为字符串

```go
// 将任意类型转为字符串
str := utils.ToString(123) // "123"
str = utils.ToString(true) // "true"
str = utils.ToString(3.14) // "3.14"
```

### 转为数值

```go
// 将任意类型转为int
num := utils.ToInt("123") // 123
num = utils.ToInt(true) // 1
num = utils.ToInt(3.14) // 3

// 将任意类型转为int64
num64 := utils.ToInt64("123") // 123
num64 = utils.ToInt64(true) // 1
num64 = utils.ToInt64(3.14) // 3
```

### 转为布尔值

```go
// 将任意类型转为bool
b := utils.ToBool(1) // true
b = utils.ToBool("true") // true
b = utils.ToBool(0) // false
```

## 使用场景

1. 数据库模型字段处理：当数据库模型中使用指针类型字段时，可以安全地访问这些字段。
2. API 响应构建：当需要构建 API 响应时，可以安全地从可能为 nil 的指针中获取值。
3. 配置参数处理：处理配置参数时，可以安全地获取可能不存在的配置项的默认值。
4. 数据转换和处理：使用切片工具函数，可以轻松地进行数据转换、过滤和聚合。
5. 字符串处理：在处理用户输入、格式转换、命名风格变换等场景中使用字符串工具函数。
6. 类型安全转换：在处理不同数据源或格式时，使用类型转换函数安全地转换数据。 