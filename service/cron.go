package service

import (
	"chanel/classes"
	"chanel/schedule"
	"chanel/structs"
	"context"
	"encoding/json"
	"fmt"
)

func (service *Service) CreateCron(params structs.CreateCronRequest, ctx context.Context) structs.Response {
	var (
		response = structs.Response{}
	)

	if params.Port == "" {
		if params.Protocol == classes.ProtocolHttps {
			params.Port = "443"
		} else {
			params.Port = "80"
		}
	}

	// 解析 參數
	var args map[string]interface{}

	if params.Args != "" {
		err := json.Unmarshal([]byte(params.Args), &args)

		if err != nil {
			response.Code = classes.JsonUnmarshalError
			response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
			response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonUnmarshalError), "Args.Unmarshal", err)
			return response
		}
	}

	// 解析 表頭
	var header = make(map[string]string)

	if params.Headers != "" {
		var headers map[string]interface{}
		err := json.Unmarshal([]byte(params.Headers), &headers)

		if err != nil {
			response.Code = classes.JsonUnmarshalError
			response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
			response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonUnmarshalError), "Headers.Unmarshal", err)
			return response
		}

		for k, v := range headers {
			header[k] = fmt.Sprintf("%v", v)
		}
	}

	datas := structs.ChanelModelCrons{
		Protocol:   params.Protocol,
		Domain:     params.Domain,
		Path:       params.Path,
		Port:       params.Port,
		Method:     params.Method,
		Args:       params.Args,
		Headers:    params.Headers,
		Execute:    params.Execute,
		Status:     1, // 預設啟用
		Remark:     params.Remark,
		CreateTime: service.tools.NowUnix(),
		UpdateTime: service.tools.NowUnix(),
	}

	// 資料新增
	cronID, err := service.mysql.ChanelDB.Repository.Crons.CreateCron(datas)

	if err != nil {
		response.Code = classes.MysqlInsertError
		response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlInsertError), "Crons.CreateCron", err)
		return response
	}

	// 寫入頻道
	schedule.JobChan <- &schedule.Job{
		ID:          cronID,
		Protocol:    datas.Protocol,
		Domain:      datas.Domain,
		Path:        datas.Path,
		Port:        datas.Port,
		Method:      datas.Method,
		Args:        args,
		Headers:     header,
		Mode:        schedule.Cron,
		ExecuteCron: datas.Execute,
	}

	// 回傳
	return response
}
