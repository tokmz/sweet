# Crypto 加密工具包

这是一个提供常用加密和哈希功能的工具包，包含 MD5、SHA 系列哈希算法、Base64 编解码以及安全随机盐值生成功能。

## 功能特性

### 🔐 哈希算法
- **MD5**: 快速哈希算法，适用于非安全场景
- **SHA256**: 安全哈希算法，推荐用于密码哈希
- **SHA384**: 更高安全性的哈希算法
- **SHA512**: 最高安全性的哈希算法

### 🔑 编解码
- **Base64 编码**: 将字符串编码为 Base64 格式
- **Base64 解码**: 将 Base64 字符串解码为原始字符串

### 🧂 随机盐值生成
- **安全随机**: 使用加密安全的随机数生成器
- **可定制长度**: 支持自定义盐值长度
- **字符集丰富**: 包含大小写字母和数字

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "zero/pkg/crypto"
)

func main() {
    text := "Hello, World!"
    
    // MD5 哈希
    md5Hash := crypto.MD5(text)
    fmt.Printf("MD5: %s\n", md5Hash)
    
    // SHA256 哈希
    sha256Hash := crypto.Hash256(text)
    fmt.Printf("SHA256: %s\n", sha256Hash)
    
    // Base64 编码
    encoded := crypto.Base64Encode(text)
    fmt.Printf("Base64 Encoded: %s\n", encoded)
    
    // Base64 解码
    decoded, err := crypto.Base64Decode(encoded)
    if err != nil {
        fmt.Printf("Decode error: %v\n", err)
    } else {
        fmt.Printf("Base64 Decoded: %s\n", decoded)
    }
    
    // 生成随机盐值
    salt := crypto.Salt()
    fmt.Printf("Random Salt: %s\n", salt)
}
```

## API 参考

### 哈希函数

#### MD5(str string) string
计算字符串的 MD5 哈希值。

```go
hash := crypto.MD5("password123")
// 输出: 482c811da5d5b4bc6d497ffa98491e38
```

**注意**: MD5 已被认为不够安全，不建议用于密码哈希，仅适用于数据校验等非安全场景。

#### Hash256(str string) string
计算字符串的 SHA256 哈希值，推荐用于密码哈希。

```go
hash := crypto.Hash256("password123")
// 输出: ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f
```

#### Hash384(str string) string
计算字符串的 SHA384 哈希值。

```go
hash := crypto.Hash384("password123")
// 输出: 9bb58f26192e4ba00f01e2e7b136bbd8f4c5b8a3b0c9c0d5e6f7a8b9c0d1e2f3...
```

#### Hash512(str string) string
计算字符串的 SHA512 哈希值，提供最高级别的安全性。

```go
hash := crypto.Hash512("password123")
// 输出: b109f3bbbc244eb82441917ed06d618b9008dd09b3befd1b5e07394c706a8bb9...
```

### Base64 编解码

#### Base64Encode(str string) string
将字符串编码为 Base64 格式。

```go
encoded := crypto.Base64Encode("Hello, 世界!")
// 输出: SGVsbG8sIOS4lueVjCE=
```

#### Base64Decode(str string) (string, error)
将 Base64 编码的字符串解码为原始字符串。

```go
decoded, err := crypto.Base64Decode("SGVsbG8sIOS4lueVjCE=")
if err != nil {
    // 处理解码错误
}
// decoded: "Hello, 世界!"
```

### 随机盐值生成

#### Salt() string
生成默认长度（16位）的随机盐值。

```go
salt := crypto.Salt()
// 输出示例: "aB3dE7gH9jK2mN5p"
```

#### SaltWithLength(length int) string
生成指定长度的随机盐值。

```go
// 生成8位盐值
salt8 := crypto.SaltWithLength(8)
// 输出示例: "aB3dE7gH"

// 生成32位盐值
salt32 := crypto.SaltWithLength(32)
// 输出示例: "aB3dE7gH9jK2mN5pqR6sT8uV1wX4yZ7c"
```

**参数说明**:
- `length`: 盐值长度，建议不少于8位
- 字符集包含: `a-z`, `A-Z`, `0-9`
- 使用加密安全的随机数生成器

## 使用场景

### 1. 密码哈希

```go
package main

import (
    "fmt"
    "zero/pkg/crypto"
)

// 用户注册时的密码处理
func hashPassword(password string) (string, string) {
    // 生成随机盐值
    salt := crypto.Salt()
    
    // 密码 + 盐值进行哈希
    hashedPassword := crypto.Hash256(password + salt)
    
    return hashedPassword, salt
}

// 用户登录时的密码验证
func verifyPassword(password, hashedPassword, salt string) bool {
    // 使用相同的盐值对输入密码进行哈希
    inputHash := crypto.Hash256(password + salt)
    
    // 比较哈希值
    return inputHash == hashedPassword
}

func main() {
    password := "mySecretPassword"
    
    // 注册
    hashedPwd, salt := hashPassword(password)
    fmt.Printf("Hashed Password: %s\n", hashedPwd)
    fmt.Printf("Salt: %s\n", salt)
    
    // 登录验证
    isValid := verifyPassword(password, hashedPwd, salt)
    fmt.Printf("Password Valid: %t\n", isValid)
}
```

### 2. 数据完整性校验

```go
// 文件完整性校验
func checkFileIntegrity(content string, expectedHash string) bool {
    actualHash := crypto.Hash256(content)
    return actualHash == expectedHash
}

// API 签名验证
func generateAPISignature(data, secretKey string) string {
    return crypto.Hash256(data + secretKey)
}
```

### 3. 数据编码传输

```go
// 敏感数据编码传输
func encodeForTransmission(data string) string {
    return crypto.Base64Encode(data)
}

// 接收端解码
func decodeFromTransmission(encodedData string) (string, error) {
    return crypto.Base64Decode(encodedData)
}
```

### 4. 会话令牌生成

```go
// 生成会话令牌
func generateSessionToken(userID string) string {
    // 用户ID + 随机盐值 + 时间戳
    timestamp := fmt.Sprintf("%d", time.Now().Unix())
    salt := crypto.Salt()
    
    tokenData := userID + salt + timestamp
    return crypto.Base64Encode(crypto.Hash256(tokenData))
}
```

## 安全建议

### 1. 密码哈希最佳实践

```go
// ✅ 推荐：使用盐值 + SHA256
func securePasswordHash(password string) (string, string) {
    salt := crypto.SaltWithLength(16) // 至少16位盐值
    hash := crypto.Hash256(password + salt)
    return hash, salt
}

// ❌ 不推荐：直接使用MD5
func insecurePasswordHash(password string) string {
    return crypto.MD5(password) // 不安全
}
```

### 2. 盐值管理

```go
// ✅ 推荐：每个密码使用唯一盐值
func createUser(username, password string) {
    hash, salt := securePasswordHash(password)
    // 将 hash 和 salt 分别存储到数据库
    saveToDatabase(username, hash, salt)
}

// ❌ 不推荐：使用固定盐值
const FIXED_SALT = "mysalt123" // 不安全
```

### 3. 算法选择

- **MD5**: 仅用于非安全场景（如数据校验）
- **SHA256**: 推荐用于一般安全需求
- **SHA384/SHA512**: 用于高安全性要求的场景

## 性能考虑

### 哈希算法性能对比

| 算法 | 速度 | 安全性 | 推荐场景 |
|------|------|--------|----------|
| MD5 | 最快 | 低 | 数据校验 |
| SHA256 | 快 | 高 | 密码哈希 |
| SHA384 | 中等 | 很高 | 高安全需求 |
| SHA512 | 较慢 | 最高 | 极高安全需求 |

### 性能优化建议

```go
// 对于高频操作，可以考虑缓存哈希结果
var hashCache = make(map[string]string)

func cachedHash(input string) string {
    if hash, exists := hashCache[input]; exists {
        return hash
    }
    
    hash := crypto.Hash256(input)
    hashCache[input] = hash
    return hash
}
```

## 注意事项

1. **MD5 安全性**: MD5 已被证明存在碰撞漏洞，不应用于安全敏感场景
2. **盐值存储**: 盐值应与哈希值分开存储，但不需要加密
3. **随机数安全**: 本包使用 `crypto/rand` 提供加密安全的随机数
4. **错误处理**: Base64 解码可能失败，务必检查错误返回值
5. **字符编码**: 所有函数都假设输入为 UTF-8 编码的字符串

## 常见问题

### Q: 为什么需要使用盐值？
A: 盐值可以防止彩虹表攻击，即使两个用户使用相同密码，加盐后的哈希值也不同。

### Q: 盐值需要保密吗？
A: 盐值不需要保密，但应该为每个密码生成唯一的盐值。

### Q: 如何选择合适的哈希算法？
A: 对于密码哈希，推荐使用 SHA256；对于数据校验，可以使用 MD5；对于高安全需求，使用 SHA512。

### Q: Base64 编码是加密吗？
A: Base64 是编码而非加密，任何人都可以轻易解码，不应用于敏感数据保护。

## 许可证

MIT License