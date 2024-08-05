package service

import (
	"chanel/classes"
	"chanel/schedule"
	"chanel/structs"
	"context"
	"encoding/json"
)

func (service *Service) GetCrons(params structs.GetCronsRequest, ctx context.Context) structs.Response {
	var (
		response = structs.Response{}
		result   = []structs.GetCronsResponse{}
	)

	// 查詢
	crons, err := service.mysql.ChanelDB.Repository.Crons.GetCrons(params)

	if err != nil {
		response.Code = classes.MysqlSearchError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.MysqlSearchError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlSearchError), "Crons.GetCrons", err)
		return response
	}

	for _, cron := range crons {
		var args map[string]interface{}

		if cron.Args != "" {
			err := json.Unmarshal([]byte(cron.Args), &args)

			if err != nil {
				response.Code = classes.JsonUnmarshalError
				response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.JsonUnmarshalError), "")
				response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonUnmarshalError), "Args.Unmarshal", err)
				return response
			}
		}

		var headers map[string]string

		if cron.Headers != "" {
			err := json.Unmarshal([]byte(cron.Headers), &headers)

			if err != nil {
				response.Code = classes.JsonUnmarshalError
				response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.JsonUnmarshalError), "")
				response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonUnmarshalError), "Headers.Unmarshal", err)
				return response
			}
		}

		result = append(result, structs.GetCronsResponse{
			ID:         cron.ID,
			Protocol:   cron.Protocol,
			Domain:     cron.Domain,
			Path:       cron.Path,
			Port:       cron.Port,
			Method:     cron.Method,
			Args:       args,
			Headers:    headers,
			Execute:    cron.Execute,
			Status:     cron.Status,
			Remark:     cron.Remark,
			CreateTime: cron.CreateTime,
			UpdateTime: cron.UpdateTime,
		})
	}

	response.Result = result
	return response
}

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

	// 參數 Json加密
	args, err := json.Marshal(params.Args)

	if err != nil {
		response.Code = classes.JsonMarshalError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.JsonMarshalError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonMarshalError), "Args.Marshal", err)
		return response
	}

	// 表頭 Json加密
	headers, err := json.Marshal(params.Headers)

	if err != nil {
		response.Code = classes.JsonMarshalError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.JsonMarshalError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonMarshalError), "Headers.Marshal", err)
		return response
	}

	datas := structs.ChanelModelCrons{
		Protocol:   params.Protocol,
		Domain:     params.Domain,
		Path:       params.Path,
		Port:       params.Port,
		Method:     params.Method,
		Args:       string(args),
		Headers:    string(headers),
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
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.MysqlInsertError), "")
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
		Args:        params.Args,
		Headers:     params.Headers,
		Mode:        schedule.Cron,
		ExecuteCron: datas.Execute,
	}

	// 回傳
	return response
}

func (service *Service) UpdateCron(params structs.UpdateCronRequest, ctx context.Context) structs.Response {
	var (
		response = structs.Response{}
		result   = structs.UpdateCronResponse{}
	)

	// 取出要修改的 cron
	beforeCron, err := service.mysql.ChanelDB.Repository.Crons.GetCronByID(params.CronID)

	if err != nil {
		response.Code = classes.MysqlSearchError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.MysqlSearchError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlSearchError), "Crons.GetCronByID", err)
		return response
	}

	// 取出要修改的 Job
	newJob := &schedule.Job{}

	for job := range schedule.JobChan {
		if job.ID == params.CronID {
			newJob = job
			break
		}
		schedule.JobChan <- job
	}

	if params.Protocol != "" && params.Protocol != beforeCron.Protocol {
		newJob.Protocol = params.Protocol
		result.Detail = append(result.Detail, structs.UpdateCronDetail{
			Field:  "Protocol",
			Before: beforeCron.Protocol,
			After:  params.Protocol,
		})
	}

	if params.Domain != "" && params.Domain != beforeCron.Domain {
		newJob.Domain = params.Domain
		result.Detail = append(result.Detail, structs.UpdateCronDetail{
			Field:  "Domain",
			Before: beforeCron.Domain,
			After:  params.Domain,
		})
	}

	if params.Path != "" && params.Path != beforeCron.Path {
		newJob.Path = params.Path
		result.Detail = append(result.Detail, structs.UpdateCronDetail{
			Field:  "Path",
			Before: beforeCron.Path,
			After:  params.Path,
		})
	}

	if params.Port != "" && params.Port != beforeCron.Port {
		newJob.Port = params.Port
		result.Detail = append(result.Detail, structs.UpdateCronDetail{
			Field:  "Port",
			Before: beforeCron.Port,
			After:  params.Port,
		})
	}

	if params.Method != "" && params.Method != beforeCron.Method {
		newJob.Method = params.Method
		result.Detail = append(result.Detail, structs.UpdateCronDetail{
			Field:  "Method",
			Before: beforeCron.Method,
			After:  params.Method,
		})
	}

	var args []byte

	if params.Args != nil {
		args, err = json.Marshal(params.Args)

		if err != nil {
			response.Code = classes.JsonMarshalError
			response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.JsonMarshalError), "")
			response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonMarshalError), "Args.Marshal", err)
			return response
		}

		if string(args) != beforeCron.Args {
			newJob.Args = params.Args
			result.Detail = append(result.Detail, structs.UpdateCronDetail{
				Field:  "Args",
				Before: beforeCron.Args,
				After:  params.Args,
			})
		}
	}

	var headers []byte

	if params.Headers != nil {
		headers, err = json.Marshal(params.Headers)

		if err != nil {
			response.Code = classes.JsonMarshalError
			response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.JsonMarshalError), "")
			response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonMarshalError), "Headers.Marshal", err)
			return response
		}

		if string(headers) != beforeCron.Headers {
			newJob.Headers = params.Headers
			result.Detail = append(result.Detail, structs.UpdateCronDetail{
				Field:  "Headers",
				Before: beforeCron.Headers,
				After:  params.Headers,
			})
		}
	}

	if params.Execute != "" && params.Execute != beforeCron.Execute {
		newJob.ExecuteCron = params.Execute
		result.Detail = append(result.Detail, structs.UpdateCronDetail{
			Field:  "Execute",
			Before: beforeCron.Execute,
			After:  params.Execute,
		})
	}

	if params.Remark != "" && params.Remark != beforeCron.Remark {
		result.Detail = append(result.Detail, structs.UpdateCronDetail{
			Field:  "Remark",
			Before: beforeCron.Remark,
			After:  params.Remark,
		})
	}

	// 資料更新
	err = service.mysql.ChanelDB.Repository.Crons.UpdateCron(structs.ChanelModelCrons{
		Protocol:   params.Protocol,
		Domain:     params.Domain,
		Path:       params.Path,
		Port:       params.Port,
		Method:     params.Method,
		Args:       string(args),
		Headers:    string(headers),
		Execute:    params.Execute,
		Status:     1,
		Remark:     params.Remark,
		UpdateTime: service.tools.NowUnix(),
	}, params.CronID)

	if err != nil {
		response.Code = classes.MysqlInsertError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.MysqlInsertError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlInsertError), "Crons.UpdateCron", err)
		return response
	}

	// 寫入 Job Chan
	schedule.JobChan <- newJob

	// 回傳
	response.Result = result
	return response
}

func (service *Service) DeleteCron(params structs.DeleteCronRequest, ctx context.Context) structs.Response {
	var (
		response = structs.Response{}
	)

	// 查詢
	err := service.mysql.ChanelDB.Repository.Crons.DeleteCronByID(params.CronID)

	if err != nil {
		response.Code = classes.MysqlDeleteError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.MysqlDeleteError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlDeleteError), "Crons.DeleteCronByID", err)
		return response
	}

	return response
}
