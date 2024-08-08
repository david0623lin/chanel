package controller

import (
	"chanel/classes"
	"chanel/structs"

	"github.com/gin-gonic/gin"
)

// @Summary 取得任務列表
// @description 一些說明
// @Tags Task
// @produce application/json
// @Param Sid header string true "SessionID"
// @Param StartTime query int false "開始時間"
// @Param EndTime query int false "結束時間"
// @Param Topic query string false "主題"
// @Param Method query string false "方法"
// @Param Status query int false "狀態 1:未執行, 2已執行"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/task/list [get]
func (ctl *Controoller) GetTasks(c *gin.Context) {
	var (
		response = structs.Response{}
	)

	defer func() {
		c.Set("Response", response)
	}()
	// 取得請求上下文
	ctx := c.Request.Context()

	// 參數處理
	params := structs.GetTasksRequest{
		StartTime: ctl.tools.StrToInt64(c.Query("StartTime")),
		EndTime:   ctl.tools.StrToInt64(c.Query("EndTime")),
		Method:    c.Query("Method"),
		Topic:     c.Query("Topic"),
		Status:    ctl.tools.StrToInt32(c.Query("Status")),
	}
	// 參數檢查
	if response = ctl.request.GetTasks(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.GetTasks(params, ctx)
}

// @Summary 取得任務詳細內容
// @description 一些說明
// @Tags Task
// @produce application/json
// @Param Sid header string true "SessionID"
// @Param ID query int true "任務ID"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/task/detail [get]
func (ctl *Controoller) GetTaskDetail(c *gin.Context) {
	var (
		response = structs.Response{}
	)

	defer func() {
		c.Set("Response", response)
	}()
	// 取得請求上下文
	ctx := c.Request.Context()

	// 參數處理
	params := structs.GetTaskDetailRequest{
		ID: ctl.tools.StrToInt32(c.Query("ID")),
	}
	// 參數檢查
	if response = ctl.request.GetTaskDetail(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.GetTaskDetail(params, ctx)
}

// @Summary 新增任務
// @description 一些說明
// @Tags Task
// @produce application/json
// @Param Sid header string true "SessionID"
// @param param body structs.CreateTaskRequest true "參數"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/task/create [post]
func (ctl *Controoller) CreateTask(c *gin.Context) {
	var (
		response = structs.Response{}
		params   = structs.CreateTaskRequest{}
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
		response.Error = ctl.tools.FormatErr(ctl.myErr.Msg(classes.ParserRequestBodyError), "CreateTask.ShouldBindJSON", err)
		return
	}
	// 參數檢查
	if response = ctl.request.CreateTask(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.CreateTask(params, ctx)
}

// @Summary 修改任務
// @description 一些說明
// @Tags Task
// @produce application/json
// @Param Sid header string true "SessionID"
// @param param body structs.UpdateTaskRequest true "參數"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/task/update [put]
func (ctl *Controoller) UpdateTask(c *gin.Context) {
	var (
		response = structs.Response{}
		params   = structs.UpdateTaskRequest{}
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
		response.Error = ctl.tools.FormatErr(ctl.myErr.Msg(classes.ParserRequestBodyError), "UpdateTask.ShouldBindJSON", err)
		return
	}
	// 參數檢查
	if response = ctl.request.UpdateTask(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.UpdateTask(params, ctx)
}

// @Summary 刪除任務
// @description 一些說明
// @Tags Task
// @produce application/json
// @Param Sid header string true "SessionID"
// @Param TaskID query int true "任務ID"
// @Success 200 {object} structs.Response "回傳"
// @Router /chanel/task/remove [delete]
func (ctl *Controoller) DeleteTask(c *gin.Context) {
	var (
		response = structs.Response{}
	)

	defer func() {
		c.Set("Response", response)
	}()
	// 取得請求上下文
	ctx := c.Request.Context()

	// 參數處理
	params := structs.DeleteTaskRequest{
		TaskID: ctl.tools.StrToInt32(c.DefaultQuery("TaskID", "")),
	}
	// 參數檢查
	if response = ctl.request.DeleteTask(params, response); response.Code != 0 {
		return
	}
	// 執行
	response = ctl.service.DeleteTask(params, ctx)
}
