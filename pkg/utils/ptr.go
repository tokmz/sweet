package utils

// SafeInt64 安全获取int64指针值
func SafeInt64(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

// SafeInt 安全获取int指针值
func SafeInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

// SafeString 安全获取string指针值
func SafeString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// SafeBool 安全获取bool指针值
func SafeBool(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

// SafeBoolFromInt64 安全获取bool值（从int64转换，1为true，其他为false）
func SafeBoolFromInt64(v *int64) bool {
	if v == nil {
		return false
	}
	return *v == 1
}

// SafeBoolFromInt 安全获取bool值（从int转换，1为true，其他为false）
func SafeBoolFromInt(v *int) bool {
	if v == nil {
		return false
	}
	return *v == 1
}

// ToPtr 将值转换为指针
func ToPtr[T any](v T) *T {
	return &v
}
