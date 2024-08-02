package lib

import (
	"math"

	"github.com/shopspring/decimal"
)

// 精度處裡工具

// 精度相加
func (tools *Tools) GetFloatBYDecimalAdd(num1, num2 float64) float64 {
	dec1 := decimal.NewFromFloat(num1)
	dec2 := decimal.NewFromFloat(num2)
	res := dec1.Add(dec2)
	res = res.Round(6)
	resFloat, _ := res.Float64()
	return resFloat
}

// 精度相減
func (tools *Tools) GetFloatBYDecimalSub(num1, num2 float64) float64 {
	dec1 := decimal.NewFromFloat(num1)
	dec2 := decimal.NewFromFloat(num2)
	res := dec1.Sub(dec2)
	res = res.Round(6)
	resFloat, _ := res.Float64()
	return resFloat
}

// 精度相乘
func (tools *Tools) GetFloatBYDecimalMul(num1, num2 float64) float64 {
	dec1 := decimal.NewFromFloat(num1)
	dec2 := decimal.NewFromFloat(num2)
	res := dec1.Mul(dec2)
	res = res.Round(6)
	resFloat, _ := res.Float64()
	return resFloat
}

// 精度相除
func (tools *Tools) GetFloatBYDecimalDiv(num1, num2 float64) float64 {
	dec1 := decimal.NewFromFloat(num1)
	dec2 := decimal.NewFromFloat(num2)
	if dec2.Equal(decimal.NewFromInt(0)) {
		return 0
	}
	res := dec1.Div(dec2)
	res = res.Round(6)
	resFloat, _ := res.Float64()
	return resFloat
}

// 四捨五入到指定第幾位
func (tools *Tools) DecimalPlaces(x float64, places int) float64 {
	// 計算要到第幾位
	decimal := math.Pow(10, float64(places))
	return math.Round(x*decimal) / decimal
}
