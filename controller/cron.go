package controller

import (
	"chanel/classes"
	"chanel/structs"

	"github.com/gin-gonic/gin"
)

// @Summary 新增排程
// @description execute格式為: * * * * * * * (秒 分 時 日 月 週 年)
// @Tags 新增
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
