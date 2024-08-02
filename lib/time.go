package lib

import (
	"time"
)

func (tools *Tools) Locat() *time.Location {
	return tools.loc
}

func (tools *Tools) NowDateTime() string {
	return time.Now().In(tools.loc).Format("2006-01-02 15:04:05")
}

func (tools *Tools) NowDate() string {
	return time.Now().In(tools.loc).Format("2006-01-02")
}

func (tools *Tools) NowTime() string {
	return time.Now().In(tools.loc).Format("15:04:05")
}

func (tools *Tools) NowYear() string {
	return time.Now().In(tools.loc).Format("2006")
}

func (tools *Tools) NowMonth() string {
	return time.Now().In(tools.loc).Format("01")
}

func (tools *Tools) NowDay() string {
	return time.Now().In(tools.loc).Format("02")
}

func (tools *Tools) NowHour() string {
	return time.Now().In(tools.loc).Format("15")
}

func (tools *Tools) NowMin() string {
	return time.Now().In(tools.loc).Format("04")
}

func (tools *Tools) NowSec() string {
	return time.Now().In(tools.loc).Format("05")
}

func (tools *Tools) NowUnix() int64 {
	return time.Now().In(tools.loc).Unix()
}

func (tools *Tools) Now() time.Time {
	return time.Now().In(tools.loc)
}

// 指定時間戳轉當日整天區間日期時間
func (tools *Tools) RangeDateTime(unix int64) []string {
	day := time.Unix(unix, 0)

	return []string{
		time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, tools.loc).Format("2006-01-02 15:04:05"),
		time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, tools.loc).Format("2006-01-02 15:04:05"),
	}
}

// 指定時間戳轉當日整天區間時間戳
func (tools *Tools) RangeUnix(unix int64) []int64 {
	day := time.Unix(unix, 0)

	return []int64{
		time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, tools.loc).Unix(),
		time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, tools.loc).Unix(),
	}
}

// 指定日期時間轉時間戳
func (tools *Tools) DateTimeToUnix(DateTime string) int64 {
	t, err := time.Parse("2006-01-02 15:04:05", DateTime)

	if err != nil {
		return 0
	}
	return t.In(tools.loc).Unix()
}

// 指定時間戳轉日期時間
func (tools *Tools) UnixToDateTime(unix int64) string {
	return time.Unix(unix, 0).In(tools.loc).Format("2006-01-02 15:04:05")
}

// 指定時間戳增減時間, d 給正的是增加, 負的是減少
func (tools *Tools) UnixAdd(unix int64, d time.Duration) int64 {
	return time.Unix(unix, 0).In(tools.loc).Add(d).Unix()
}
