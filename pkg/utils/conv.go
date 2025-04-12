package utils

import (
	"fmt"
	"strconv"
)

// ToString 将任意类型转为字符串
func ToString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch value := v.(type) {
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// ToInt 将任意类型转为int
func ToInt(v interface{}) int {
	if v == nil {
		return 0
	}

	switch value := v.(type) {
	case int:
		return value
	case int8:
		return int(value)
	case int16:
		return int(value)
	case int32:
		return int(value)
	case int64:
		return int(value)
	case uint:
		return int(value)
	case uint8:
		return int(value)
	case uint16:
		return int(value)
	case uint32:
		return int(value)
	case uint64:
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	case string:
		i, _ := strconv.Atoi(value)
		return i
	case bool:
		if value {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// ToInt64 将任意类型转为int64
func ToInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}

	switch value := v.(type) {
	case int:
		return int64(value)
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value
	case uint:
		return int64(value)
	case uint8:
		return int64(value)
	case uint16:
		return int64(value)
	case uint32:
		return int64(value)
	case uint64:
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	case string:
		i, _ := strconv.ParseInt(value, 10, 64)
		return i
	case bool:
		if value {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// ToBool 将任意类型转为bool
func ToBool(v interface{}) bool {
	if v == nil {
		return false
	}

	switch value := v.(type) {
	case bool:
		return value
	case int:
		return value > 0
	case int8:
		return value > 0
	case int16:
		return value > 0
	case int32:
		return value > 0
	case int64:
		return value > 0
	case uint:
		return value > 0
	case uint8:
		return value > 0
	case uint16:
		return value > 0
	case uint32:
		return value > 0
	case uint64:
		return value > 0
	case float32:
		return value > 0
	case float64:
		return value > 0
	case string:
		b, _ := strconv.ParseBool(value)
		return b
	default:
		return false
	}
}
