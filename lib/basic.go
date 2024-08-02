package lib

import (
	"chanel/config"
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
)

// 自訂工具
type Tools struct {
	config *config.Config
	loc    *time.Location
}

func ToolsInit(config *config.Config) *Tools {
	loc, err := time.LoadLocation(config.Locat)

	if err != nil {
		panic(fmt.Sprintf("初始化 自定義工具 錯誤, ERR: %s", err.Error()))
	}

	return &Tools{
		config: config,
		loc:    loc,
	}
}

func PanicParser(err interface{}) error {
	switch r := err.(type) {
	case error:
		return r
	default:
		return fmt.Errorf("%v", r)
	}
}

// 參數檢查
func (tools *Tools) Request(param any) bool {
	switch v := param.(type) {
	case int:
		return v != 0
	case int16:
		return v != int16(0)
	case int32:
		return v != int32(0)
	case int64:
		return v != int64(0)
	case float32:
		return v != float32(0)
	case float64:
		return v != float64(0)
	case string:
		return v != ""
	case []interface{}:
		return len(v) != 0
	case []string:
		return len(v) != 0
	case []int:
		return len(v) != 0
	case []int16:
		return len(v) != 0
	case []int32:
		return len(v) != 0
	case []int64:
		return len(v) != 0
	case []float32:
		return len(v) != 0
	case []float64:
		return len(v) != 0
	case map[interface{}]interface{}:
		return len(v) != 0
	case map[string]interface{}:
		return len(v) != 0
	case map[int]interface{}:
		return len(v) != 0
	case map[int16]interface{}:
		return len(v) != 0
	case map[int32]interface{}:
		return len(v) != 0
	case map[int64]interface{}:
		return len(v) != 0
	case map[float32]interface{}:
		return len(v) != 0
	case map[float64]interface{}:
		return len(v) != 0
	default:
		// 未定義的檢查
		return false
	}
}

// 格式化回傳訊息
func (tools *Tools) FormatMsg(message, remark string) string {
	if remark == "" {
		return message
	}
	return fmt.Sprintf("%s, [%s]", message, remark)
}

// 格式化回傳錯誤訊息
func (tools *Tools) FormatErr(message, remark string, err error) error {
	if remark == "" {
		return fmt.Errorf("%s, Err [%s]", message, err.Error())
	}
	return fmt.Errorf("%s, Debug [%s], Err [%s]", message, remark, err.Error())
}

// 取得執行完成的總時間
func (tools *Tools) GetDownRunTime(startTime time.Time) float64 {
	return float64(time.Since(startTime).Milliseconds()) / 1000
}

// 產生唯一 traceID
func (tools *Tools) NewTraceID() (traceID string, err error) {
	node, err := snowflake.NewNode(1)

	if err != nil {
		return "", err
	}
	traceID = node.Generate().String()

	return traceID, nil
}

func (tools *Tools) RandCharset(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
