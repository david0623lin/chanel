package schedule

import (
	"chanel/classes"
	"chanel/lib"
	"chanel/structs"
	"net/http"
	"time"
)

type cronNextInfo struct {
	t time.Time
	s string
}

func (schedule *Schedule) doCron(cron *Job) {
	// 初始化 Trace log
	traceLog := classes.TraceLogInit(schedule.tools)
	traceLog.SetTopic("schedule")
	traceLog.SetMethod("doCron")
	traceLog.SetArgs(cron)
	t := time.Now()

	defer func() {
		if err := recover(); err != nil {
			traceLog.SetRequestTime(schedule.tools.GetDownRunTime(t))
			traceLog.SetCode(structs.SystemErrorCode)
			traceLog.PrintError(structs.SystemErrorMsg, schedule.tools.FormatErr(structs.SystemErrorMsg, "doCron.Panic", lib.PanicParser(err)))
			return
		}
	}()
	traceID, err := schedule.tools.NewTraceID()

	if err != nil {
		traceLog.SetRequestTime(schedule.tools.GetDownRunTime(t))
		traceLog.SetCode(classes.CreateTraceIdError)
		traceLog.PrintError(schedule.myErr.Msg(classes.CreateTraceIdError), err)
		return
	}
	traceLog.SetTraceID(traceID)

	// 寫入 traceID 到 header
	if cron.Headers == nil {
		cron.Headers = make(map[string]string)
	}
	cron.Headers[structs.TraceID] = traceID

	// 建立 curl 物件
	curl := classes.CurlInit(schedule.tools)

	switch cron.Protocol {
	case classes.ProtocolHttp:
		curl.SetHttp()
	case classes.ProtocolHttps:
		curl.SetHttps()
	}
	curl.NewRequest(cron.Domain, cron.Port, cron.Path)
	curl.SetHeaders(cron.Headers)
	curl.SetTraceID(traceID)

	var result string

	switch cron.Method {
	case http.MethodGet:
		result, err = curl.SetQueries(cron.Args).Get()
	case http.MethodPost:
		result, err = curl.SetBody(cron.Args).Post()
	case http.MethodPut:
		result, err = curl.SetBody(cron.Args).Put()
	case http.MethodDelete:
		result, err = curl.SetQueries(cron.Args).Delete()
	}

	// 新增 執行紀錄
	var insertDatas structs.ChanelModelCronRecords
	insertDatas.CronID = cron.ID
	insertDatas.Result = result
	insertDatas.CreateTime = schedule.tools.NowUnix()

	// 如果打 API 錯誤or失敗不用再寫 Log, 因為 Curl時就會寫一筆完整的
	if err != nil {
		insertDatas.Status = 2
		insertDatas.Error = err.Error()
	} else {
		insertDatas.Status = 1
	}
	err = schedule.mysql.ChanelDB.Repository.CronRecords.Create(insertDatas)

	if err != nil {
		// 紀錄 新增執行紀錄錯誤
		traceLog.SetRequestTime(schedule.tools.GetDownRunTime(t))
		traceLog.SetArgs(insertDatas)
		traceLog.SetCode(classes.MysqlInsertError)
		traceLog.PrintError(schedule.myErr.Msg(classes.MysqlInsertError), err)
		return
	}

	// 紀錄 執行完成
	traceLog.SetRequestTime(schedule.tools.GetDownRunTime(t))
	traceLog.PrintInfo(Success)
}
