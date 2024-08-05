package schedule

import (
	"chanel/classes"
	"chanel/structs"
	"net/http"
	"time"
)

func (schedule *Schedule) doTask(task *Job) {
	// 初始化 Trace log
	traceLog := classes.TraceLogInit(schedule.tools)
	traceLog.SetTopic("schedule")
	traceLog.SetMethod("doTask")
	traceLog.SetArgs(task)

	t := time.Now()
	traceID, err := schedule.tools.NewTraceID()

	if err != nil {
		traceLog.SetRequestTime(schedule.tools.GetDownRunTime(t))
		traceLog.SetCode(classes.CreateTraceIdError)
		traceLog.PrintError(schedule.myErr.Msg(classes.CreateTraceIdError), err)
		return
	}
	traceLog.SetTraceID(traceID)

	// 寫入 traceID 到 header
	if task.Headers == nil {
		task.Headers = make(map[string]string)
	}

	// 建立 curl 物件
	curl := classes.CurlInit(schedule.tools)

	switch task.Protocol {
	case classes.ProtocolHttp:
		curl.SetHttp()
	case classes.ProtocolHttps:
		curl.SetHttps()
	}
	curl.NewRequest(task.Domain, task.Port, task.Path)
	curl.SetHeaders(task.Headers)
	curl.SetTraceID(traceID)

	var result string

	switch task.Method {
	case http.MethodGet:
		result, err = curl.SetQueries(task.Args).Get()
	case http.MethodPost:
		result, err = curl.SetBody(task.Args).Post()
	case http.MethodPut:
		result, err = curl.SetBody(task.Args).Put()
	case http.MethodDelete:
		result, err = curl.SetQueries(task.Args).Delete()
	}

	// 新增 執行紀錄
	var insertDatas structs.ChanelModelTaskRecords
	insertDatas.TaskID = task.ID
	insertDatas.Result = result
	insertDatas.CreateTime = schedule.tools.NowUnix()

	// 如果打 API 錯誤or失敗不用再寫 Log, 因為 Curl時就會寫一筆完整的
	if err != nil {
		insertDatas.Status = 2
		insertDatas.Error = err.Error()
	} else {
		insertDatas.Status = 1
	}
	err = schedule.mysql.ChanelDB.Repository.TaskRecords.Create(insertDatas)

	if err != nil {
		// 紀錄 新增執行紀錄錯誤
		traceLog.SetRequestTime(schedule.tools.GetDownRunTime(t))
		traceLog.SetArgs(insertDatas)
		traceLog.SetCode(classes.MysqlInsertError)
		traceLog.PrintError(schedule.myErr.Msg(classes.MysqlInsertError), err)
	}

	// 更新 任務狀態
	updateDatas := structs.ChanelModelTasks{
		Status:     2, // 已執行
		UpdateTime: schedule.tools.NowUnix(),
	}
	err = schedule.mysql.ChanelDB.Repository.Tasks.UpdateTask(updateDatas, task.ID)

	if err != nil {
		// 紀錄 更新任務狀態錯誤
		traceLog.SetRequestTime(schedule.tools.GetDownRunTime(t))
		traceLog.SetArgs(updateDatas)
		traceLog.SetCode(classes.MysqlUpdateError)
		traceLog.PrintError(schedule.myErr.Msg(classes.MysqlUpdateError), err)
		return
	}

	// 紀錄 執行完成
	traceLog.SetRequestTime(schedule.tools.GetDownRunTime(t))
	traceLog.PrintInfo(Success)
}
