package controller

import (
	"chanel/classes"
	"chanel/structs"

	"github.com/gin-gonic/gin"
)

// @Summary 登入
// @description
// @Tags Admin
// @produce application/json
// @param param body structs.AdminLoginRequest true "參數"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/admin/login [post]
func (ctl *Controoller) Login(c *gin.Context) {
	var (
		response = structs.Response{}
		params   = structs.AdminLoginRequest{}
	)

	defer func() {
		c.Set("Response", response)
	}()
	// 取得請求上下文
	ctx := c.Request.Context()

	// 參數處理
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Code = classes.ParserRequestBodyError
		response.Message = ctl.tools.FormatMsg(ctl.myErr.Msg(classes.ParserRequestBodyError), "")
		response.Error = ctl.tools.FormatErr(ctl.myErr.Msg(classes.ParserRequestBodyError), "Login.ShouldBindJSON", err)
		return
	}
	// 參數檢查
	if response = ctl.request.Login(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.Login(params, ctx)
}

// @Summary 註冊帳號
// @description
// @Tags Admin
// @produce application/json
// @param param body structs.AdminRegisterRequest true "參數"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/admin/register [post]
func (ctl *Controoller) Register(c *gin.Context) {
	var (
		response = structs.Response{}
		params   = structs.AdminRegisterRequest{}
	)

	defer func() {
		c.Set("Response", response)
	}()
	// 取得請求上下文
	ctx := c.Request.Context()

	// 參數處理
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Code = classes.ParserRequestBodyError
		response.Message = ctl.tools.FormatMsg(ctl.myErr.Msg(classes.ParserRequestBodyError), "")
		response.Error = ctl.tools.FormatErr(ctl.myErr.Msg(classes.ParserRequestBodyError), "Register.ShouldBindJSON", err)
		return
	}
	// 參數檢查
	if response = ctl.request.Register(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.Register(params, ctx)
}
