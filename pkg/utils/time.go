package utils

import (
	"fmt"
	"time"
)

// 常用时间格式常量
const (
	// DateTimeFormat 标准日期时间格式
	DateTimeFormat = "2006-01-02 15:04:05"
	// DateFormat 日期格式
	DateFormat = "2006-01-02"
	// TimeFormat 时间格式
	TimeFormat = "15:04:05"
	// DateTimeCompactFormat 紧凑日期时间格式
	DateTimeCompactFormat = "20060102150405"
	// DateCompactFormat 紧凑日期格式
	DateCompactFormat = "20060102"
	// RFC3339Format RFC3339格式
	RFC3339Format = time.RFC3339
	// ISO8601Format ISO8601格式
	ISO8601Format = "2006-01-02T15:04:05Z07:00"
)

// NowUnix 获取当前时间的Unix时间戳（秒）
func NowUnix() int64 {
	return time.Now().Unix()
}

// NowUnixMilli 获取当前时间的Unix时间戳（毫秒）
func NowUnixMilli() int64 {
	return time.Now().UnixMilli()
}

// NowUnixNano 获取当前时间的Unix时间戳（纳秒）
func NowUnixNano() int64 {
	return time.Now().UnixNano()
}

// UnixToTime 将Unix时间戳（秒）转换为time.Time
func UnixToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// UnixMilliToTime 将Unix时间戳（毫秒）转换为time.Time
func UnixMilliToTime(timestamp int64) time.Time {
	return time.UnixMilli(timestamp)
}

// UnixNanoToTime 将Unix时间戳（纳秒）转换为time.Time
func UnixNanoToTime(timestamp int64) time.Time {
	return time.Unix(0, timestamp)
}

// FormatTime 格式化时间为字符串
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// FormatNow 格式化当前时间为字符串
func FormatNow(layout string) string {
	return time.Now().Format(layout)
}

// FormatDateTime 格式化时间为标准日期时间字符串
func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

// FormatDate 格式化时间为日期字符串
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// FormatTimeOnly 格式化时间为时间字符串
func FormatTimeOnly(t time.Time) string {
	return t.Format(TimeFormat)
}

// ParseTime 解析时间字符串
func ParseTime(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// ParseDateTime 解析标准日期时间字符串
func ParseDateTime(value string) (time.Time, error) {
	return time.Parse(DateTimeFormat, value)
}

// ParseDate 解析日期字符串
func ParseDate(value string) (time.Time, error) {
	return time.Parse(DateFormat, value)
}

// ParseTimeOnly 解析时间字符串
func ParseTimeOnly(value string) (time.Time, error) {
	return time.Parse(TimeFormat, value)
}

// IsToday 判断给定时间是否为今天
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}

// IsYesterday 判断给定时间是否为昨天
func IsYesterday(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.Year() == yesterday.Year() && t.YearDay() == yesterday.YearDay()
}

// IsTomorrow 判断给定时间是否为明天
func IsTomorrow(t time.Time) bool {
	tomorrow := time.Now().AddDate(0, 0, 1)
	return t.Year() == tomorrow.Year() && t.YearDay() == tomorrow.YearDay()
}

// IsThisWeek 判断给定时间是否在本周
func IsThisWeek(t time.Time) bool {
	now := time.Now()
	// 获取本周一的日期
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 将周日调整为7
	}
	mondayOffset := weekday - 1
	monday := now.AddDate(0, 0, -mondayOffset)
	mondayStart := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
	sundayEnd := mondayStart.AddDate(0, 0, 7)
	
	return t.After(mondayStart) && t.Before(sundayEnd)
}

// IsThisMonth 判断给定时间是否在本月
func IsThisMonth(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month()
}

// IsThisYear 判断给定时间是否在今年
func IsThisYear(t time.Time) bool {
	return t.Year() == time.Now().Year()
}

// StartOfDay 获取指定日期的开始时间（00:00:00）
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取指定日期的结束时间（23:59:59）
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// StartOfWeek 获取指定日期所在周的开始时间（周一00:00:00）
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // 将周日调整为7
	}
	mondayOffset := weekday - 1
	monday := t.AddDate(0, 0, -mondayOffset)
	return StartOfDay(monday)
}

// EndOfWeek 获取指定日期所在周的结束时间（周日23:59:59）
func EndOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // 将周日调整为7
	}
	sundayOffset := 7 - weekday
	sunday := t.AddDate(0, 0, sundayOffset)
	return EndOfDay(sunday)
}

// StartOfMonth 获取指定日期所在月的开始时间（1号00:00:00）
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 获取指定日期所在月的结束时间（月末23:59:59）
func EndOfMonth(t time.Time) time.Time {
	nextMonth := t.AddDate(0, 1, 0)
	firstOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())
	lastOfMonth := firstOfNextMonth.Add(-time.Nanosecond)
	return lastOfMonth
}

// StartOfYear 获取指定日期所在年的开始时间（1月1日00:00:00）
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 获取指定日期所在年的结束时间（12月31日23:59:59）
func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
}

// DiffDays 计算两个时间相差的天数
func DiffDays(t1, t2 time.Time) int {
	diff := t2.Sub(t1)
	return int(diff.Hours() / 24)
}

// DiffHours 计算两个时间相差的小时数
func DiffHours(t1, t2 time.Time) int {
	diff := t2.Sub(t1)
	return int(diff.Hours())
}

// DiffMinutes 计算两个时间相差的分钟数
func DiffMinutes(t1, t2 time.Time) int {
	diff := t2.Sub(t1)
	return int(diff.Minutes())
}

// DiffSeconds 计算两个时间相差的秒数
func DiffSeconds(t1, t2 time.Time) int {
	diff := t2.Sub(t1)
	return int(diff.Seconds())
}

// AddDays 给时间添加指定天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddHours 给时间添加指定小时数
func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes 给时间添加指定分钟数
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// AddSeconds 给时间添加指定秒数
func AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

// TimeAgo 获取相对时间描述（如：2小时前、3天前等）
func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	
	if diff < time.Minute {
		return "刚刚"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d分钟前", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d小时前", hours)
	} else if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d天前", days)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		return fmt.Sprintf("%d个月前", months)
	} else {
		years := int(diff.Hours() / (24 * 365))
		return fmt.Sprintf("%d年前", years)
	}
}

// IsLeapYear 判断是否为闰年
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DaysInMonth 获取指定年月的天数
func DaysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// GetWeekday 获取中文星期几
func GetWeekday(t time.Time) string {
	weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	return weekdays[t.Weekday()]
}

// GetMonth 获取中文月份
func GetMonth(t time.Time) string {
	months := []string{
		"", "一月", "二月", "三月", "四月", "五月", "六月",
		"七月", "八月", "九月", "十月", "十一月", "十二月",
	}
	return months[t.Month()]
}

// Sleep 休眠指定时间
func Sleep(duration time.Duration) {
	time.Sleep(duration)
}

// SleepSeconds 休眠指定秒数
func SleepSeconds(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

// SleepMilliseconds 休眠指定毫秒数
func SleepMilliseconds(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}