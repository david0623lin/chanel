package lib

import (
	"strconv"
)

// 型態轉換、檢查工具

// string 轉 int
func (tools *Tools) StrToInt(str string) int {
	num, err := strconv.Atoi(str)

	if err != nil {
		return 0
	}
	return num
}

// string 轉 int32
func (tools *Tools) StrToInt32(str string) int32 {
	num, err := strconv.Atoi(str)

	if err != nil {
		return 0
	}
	return int32(num)
}

// string 轉 int64
func (tools *Tools) StrToInt64(str string) int64 {
	num, err := strconv.Atoi(str)

	if err != nil {
		return 0
	}
	return int64(num)
}

// string 轉 int32 (位址)
func (tools *Tools) StrToInt32Pointer(s *string) *int32 {
	value, _ := strconv.ParseInt(*s, 10, 32)
	result := int32(value)
	return &result
}

// string 轉 boolean
func (tools *Tools) StrToBool(str string) bool {
	bool, _ := strconv.ParseBool(str)
	return bool
}

// string 轉 float64
func (tools *Tools) StrToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0
	}
	return f
}

// int 轉 string
func (tools *Tools) IntToStr(i int) string {
	return strconv.Itoa(i)
}

// int 轉 int32
func (tools *Tools) IntToInt32(i int) int32 {
	return int32(i)
}

// int32 轉 string
func (tools *Tools) Int32ToStr(i int32) string {
	return strconv.Itoa(int(i))
}

// int64 轉 string
func (tools *Tools) Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

// int64 轉 int32
func (tools *Tools) Int64ToInt32(i int64) int32 {
	return int32(i)
}

// 檢查 []string 中是否包含特定 string
func (tools *Tools) InStrArray(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// 檢查 []int 中是否包含特定 int
func (tools *Tools) InIntArray(value int, array []int) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// 檢查 []int32 中是否包含特定 int32
func (tools *Tools) InInt32Array(value int32, array []int32) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// 陣列刪除指定值
func (tools *Tools) UnsetStrArray(value string, arr []string) []string {
	var result []string

	for _, v := range arr {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}

// 檢查 Map 中是否包含特定 key
func (tools *Tools) MapKeyExist(key string, searchMap map[string]interface{}) bool {
	if _, ok := searchMap[key]; ok {
		return true
	} else {
		return false
	}
}
