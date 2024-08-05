package controller

import (
	"chanel/classes"
	"chanel/structs"

	"github.com/gin-gonic/gin"
)

// @Summary 取得排程列表
// @description 一些說明
// @Tags Cron
// @produce application/json
// @Param Sid header string true "SessionID"
// @Param Path query string false "路徑"
// @Param Method query string false "方法"
// @Param Status query int false "狀態 1:啟用, 2停用"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/cron/list [get]
func (ctl *Controoller) GetCrons(c *gin.Context) {
	var (
		response = structs.Response{}
	)

	defer func() {
		c.Set("Response", response)
	}()
	// 取得請求上下文
	ctx := c.Request.Context()

	// 參數處理
	params := structs.GetCronsRequest{
		Path:   c.Query("Path"),
		Method: c.Query("Method"),
		Status: ctl.tools.StrToInt32(c.Query("Status")),
	}
	// 參數檢查
	if response = ctl.request.GetCrons(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.GetCrons(params, ctx)
}

// @Summary 新增排程
// @description execute格式為: * * * * * * * (秒 分 時 日 月 週 年)
// @Tags Cron
// @produce application/json
// @Param Sid header string true "SessionID"
// @param param body structs.CreateCronRequest true "參數"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/cron/create [post]
func (ctl *Controoller) CreateCron(c *gin.Context) {
	var (
		response = structs.Response{}
		params   = structs.CreateCronRequest{}
	)

	defer func() {
		c.Set("Response", response)
	}()
	// 取得請求上下文
	ctx := c.Request.Context()

	// 參數處理
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Code = classes.ParserRequestBodyError
		response.Message = ctl.tools.FormatMsg(structs.RequestErrorMsg, "")
		response.Error = ctl.tools.FormatErr(ctl.myErr.Msg(classes.ParserRequestBodyError), "CreateCron.ShouldBindJSON", err)
		return
	}
	// 參數檢查
	if response = ctl.request.CreateCron(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.CreateCron(params, ctx)
}

// @Summary 修改排程
// @description 一些說明
// @Tags Cron
// @produce application/json
// @Param Sid header string true "SessionID"
// @param param body structs.UpdateCronRequest true "參數"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/cron/update [put]
func (ctl *Controoller) UpdateCron(c *gin.Context) {
	var (
		response = structs.Response{}
		params   = structs.UpdateCronRequest{}
	)

	defer func() {
		c.Set("Response", response)
	}()
	// 取得請求上下文
	ctx := c.Request.Context()

	// 參數處理
	if err := c.ShouldBindJSON(&params); err != nil {
		response.Code = classes.ParserRequestBodyError
		response.Message = ctl.tools.FormatMsg(structs.RequestErrorMsg, "")
		response.Error = ctl.tools.FormatErr(ctl.myErr.Msg(classes.ParserRequestBodyError), "UpdateCron.ShouldBindJSON", err)
		return
	}
	// 參數檢查
	if response = ctl.request.UpdateCron(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.UpdateCron(params, ctx)
}

// @Summary 刪除排程
// @description 一些說明
// @Tags Cron
// @produce application/json
// @Param Sid header string true "SessionID"
// @Param CronID query int true "排程ID"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/cron/remove [delete]
func (ctl *Controoller) DeleteCron(c *gin.Context) {
	var (
		response = structs.Response{}
	)

	defer func() {
		c.Set("Response", response)
	}()
	// 取得請求上下文
	ctx := c.Request.Context()

	// 參數處理
	params := structs.DeleteCronRequest{
		CronID: ctl.tools.StrToInt32(c.DefaultQuery("CronID", "")),
	}
	// 參數檢查
	if response = ctl.request.DeleteCron(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.DeleteCron(params, ctx)
}
