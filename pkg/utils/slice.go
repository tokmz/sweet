package utils

import "slices"

// Contains 检查切片是否包含元素（优先使用标准库的slices.Contains）
// 为了向后兼容性保留此函数
func Contains[T comparable](slice []T, element T) bool {
	return slices.Contains(slice, element)
}

// ContainsFunc 使用自定义比较函数检查切片是否包含元素
// 为了向后兼容性保留此函数
func ContainsFunc[T any](slice []T, predicate func(T) bool) bool {
	return slices.ContainsFunc(slice, predicate)
}

// Map 对切片中的每个元素应用函数并返回新切片
// 标准库slices中没有这个函数
func Map[T, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// Filter 过滤切片中符合条件的元素并返回新切片
// 标准库slices中没有这个函数
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce 对切片中的元素进行累积操作
// 标准库slices中没有这个函数
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// Unique 去除切片中的重复元素
// 标准库slices中没有这个函数
func Unique[T comparable](slice []T) []T {
	if len(slice) <= 1 {
		return slices.Clone(slice)
	}

	result := make([]T, 0, len(slice))
	seen := make(map[T]struct{}, len(slice))

	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// Chunk 将切片分成指定大小的多个子切片
// 标准库slices中没有这个函数
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slices.Clone(slice[i:end]))
	}
	return chunks
}

// Flatten 将二维切片压平成一维切片
// 标准库slices中没有这个函数
func Flatten[T any](slices [][]T) []T {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]T, 0, totalLen)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}
