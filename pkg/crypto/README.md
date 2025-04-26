# Crypto 包

本包提供常用的加密和哈希功能。

## 功能

- **密码哈希**: 使用 bcrypt 对密码进行安全的哈希处理和验证。
- **数据加密**: 提供 AES-GCM 和 AES-CBC (PKCS7填充) 对称加密算法。
- **数据解密**: 提供 AES-GCM 和 AES-CBC (PKCS7填充) 对称解密算法。

## 使用方法

### 密码哈希与验证

```go
import "github.com/your-org/sweet/pkg/crypto"

// 哈希密码
hashedPassword, err := crypto.HashPassword("mysecretpassword", crypto.DefaultCost) // 或者指定 cost
if err != nil {
    // 处理错误
}

// 验证密码
isValid := crypto.CheckPasswordHash("mysecretpassword", hashedPassword)
if isValid {
    // 密码正确
} else {
    // 密码错误
}
```

### AES 加密/解密

#### GCM 模式 (提供认证加密)

```go
import "github.com/your-org/sweet/pkg/crypto"

// 创建密钥 (必须是 16, 24, 或 32 字节长度)
key := []byte("a very very very very secret key") // 32字节
plaintext := []byte("exampleplaintext")

// 加密数据
ciphertext, err := crypto.EncryptAES(key, plaintext)
if err != nil {
	// 处理错误
}

// 解密数据
decryptedText, err := crypto.DecryptAES(key, ciphertext)
if err != nil {
	// 处理错误
}

fmt.Printf("Decrypted Text: %s\n", decryptedText)
```

#### CBC 模式 (使用 PKCS7 填充)

```go
import "github.com/your-org/sweet/pkg/crypto"

// 创建密钥 (必须是 16, 24, 或 32 字节长度)
key := []byte("a very very very very secret key") // 32字节
plaintext := []byte("exampleplaintext")

// 加密数据
ciphertext, err := crypto.EncryptAESCBC(key, plaintext)
if err != nil {
	// 处理错误
}

// 解密数据
decryptedText, err := crypto.DecryptAESCBC(key, ciphertext)
if err != nil {
	// 处理错误
}

fmt.Printf("Decrypted Text: %s\n", decryptedText)
```

## 注意事项

- **密钥管理**: 确保安全地管理加密密钥。
- **成本参数**: bcrypt 的 cost 参数影响安全性和性能，请根据应用需求选择合适的值。