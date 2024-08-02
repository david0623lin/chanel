package service

import (
	"chanel/classes"
	"chanel/schedule"
	"chanel/structs"
	"context"
	"encoding/json"
)

func (service *Service) GetTasks(params structs.GetTasksRequest, ctx context.Context) structs.Response {
	var (
		response = structs.Response{}
		result   = []structs.GetTasksResponse{}
	)

	// 時間處理
	if params.StartTime == 0 || params.EndTime == 0 {
		rangeTime := service.tools.RangeUnix(service.tools.NowUnix())
		params.StartTime = rangeTime[0]
		params.EndTime = rangeTime[1]
	}

	// 查詢
	tasks, err := service.mysql.ChanelDB.Repository.Tasks.GetTasks(params)

	if err != nil {
		response.Code = classes.MysqlSearchError
		response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlSearchError), "Tasks.GetTasks", err)
		return response
	}

	for _, task := range tasks {
		var args map[string]interface{}

		if task.Args != "" {
			err := json.Unmarshal([]byte(task.Args), &args)

			if err != nil {
				response.Code = classes.JsonUnmarshalError
				response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
				response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonUnmarshalError), "Args.Unmarshal", err)
				return response
			}
		}

		var headers map[string]string

		if task.Headers != "" {
			err := json.Unmarshal([]byte(task.Headers), &headers)

			if err != nil {
				response.Code = classes.JsonUnmarshalError
				response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
				response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonUnmarshalError), "Headers.Unmarshal", err)
				return response
			}
		}

		result = append(result, structs.GetTasksResponse{
			Protocol:   task.Protocol,
			Domain:     task.Domain,
			Path:       task.Path,
			Port:       task.Port,
			Method:     task.Method,
			Args:       args,
			Headers:    headers,
			Execute:    task.Execute,
			Status:     task.Status,
			Remark:     task.Remark,
			CreateTime: task.CreateTime,
			UpdateTime: task.UpdateTime,
		})
	}

	response.Result = result
	return response
}

func (service *Service) CreateTask(params structs.CreateTaskRequest, ctx context.Context) structs.Response {
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
		response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonMarshalError), "Args.Marshal", err)
		return response
	}

	// 表頭 Json加密
	headers, err := json.Marshal(params.Headers)

	if err != nil {
		response.Code = classes.JsonMarshalError
		response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonMarshalError), "Headers.Marshal", err)
		return response
	}

	datas := structs.ChanelModelTasks{
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
		CreateTime: service.tools.NowUnix(),
		UpdateTime: service.tools.NowUnix(),
	}

	// 資料新增
	taskID, err := service.mysql.ChanelDB.Repository.Tasks.CreateTask(datas)

	if err != nil {
		response.Code = classes.MysqlInsertError
		response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlInsertError), "Tasks.CreateTask", err)
		return response
	}

	// 寫入頻道
	schedule.JobChan <- &schedule.Job{
		ID:          taskID,
		Protocol:    datas.Protocol,
		Domain:      datas.Domain,
		Path:        datas.Path,
		Port:        datas.Port,
		Method:      datas.Method,
		Args:        params.Args,
		Headers:     params.Headers,
		Mode:        schedule.Task,
		ExecuteTask: datas.Execute,
	}

	// 回傳
	return response
}

func (service *Service) UpdateTask(params structs.UpdateTaskRequest, ctx context.Context) structs.Response {
	var (
		response = structs.Response{}
		result   = structs.UpdateTaskResponse{}
	)

	// 取出要修改的 task
	beforeTask, err := service.mysql.ChanelDB.Repository.Tasks.GetTaskByID(params.TaskID)

	if err != nil {
		response.Code = classes.MysqlSearchError
		response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlSearchError), "Tasks.GetTaskByID", err)
		return response
	}

	// 取出要修改的 Job
	newJob := &schedule.Job{}

	for job := range schedule.JobChan {
		if job.ID == params.TaskID {
			newJob = job
			break
		}
		schedule.JobChan <- job
	}

	if params.Protocol != "" && params.Protocol != beforeTask.Protocol {
		newJob.Protocol = params.Protocol
		result.Detail = append(result.Detail, structs.UpdateTaskDetail{
			Field:  "Protocol",
			Before: beforeTask.Protocol,
			After:  params.Protocol,
		})
	}

	if params.Domain != "" && params.Domain != beforeTask.Domain {
		newJob.Domain = params.Domain
		result.Detail = append(result.Detail, structs.UpdateTaskDetail{
			Field:  "Domain",
			Before: beforeTask.Domain,
			After:  params.Domain,
		})
	}

	if params.Path != "" && params.Path != beforeTask.Path {
		newJob.Path = params.Path
		result.Detail = append(result.Detail, structs.UpdateTaskDetail{
			Field:  "Path",
			Before: beforeTask.Path,
			After:  params.Path,
		})
	}

	if params.Port != "" && params.Port != beforeTask.Port {
		newJob.Port = params.Port
		result.Detail = append(result.Detail, structs.UpdateTaskDetail{
			Field:  "Port",
			Before: beforeTask.Port,
			After:  params.Port,
		})
	}

	if params.Method != "" && params.Method != beforeTask.Method {
		newJob.Method = params.Method
		result.Detail = append(result.Detail, structs.UpdateTaskDetail{
			Field:  "Method",
			Before: beforeTask.Method,
			After:  params.Method,
		})
	}

	var args []byte

	if params.Args != nil {
		args, err = json.Marshal(params.Args)

		if err != nil {
			response.Code = classes.JsonMarshalError
			response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
			response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonMarshalError), "Args.Marshal", err)
			return response
		}

		if string(args) != beforeTask.Args {
			newJob.Args = params.Args
			result.Detail = append(result.Detail, structs.UpdateTaskDetail{
				Field:  "Args",
				Before: beforeTask.Args,
				After:  params.Args,
			})
		}
	}

	var headers []byte

	if params.Headers != nil {
		headers, err = json.Marshal(params.Headers)

		if err != nil {
			response.Code = classes.JsonMarshalError
			response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
			response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonMarshalError), "Headers.Marshal", err)
			return response
		}

		if string(headers) != beforeTask.Headers {
			newJob.Headers = params.Headers
			result.Detail = append(result.Detail, structs.UpdateTaskDetail{
				Field:  "Headers",
				Before: beforeTask.Headers,
				After:  params.Headers,
			})
		}
	}

	if params.Execute != 0 && params.Execute != beforeTask.Execute {
		newJob.ExecuteTask = params.Execute
		result.Detail = append(result.Detail, structs.UpdateTaskDetail{
			Field:  "Execute",
			Before: beforeTask.Execute,
			After:  params.Execute,
		})
	}

	if params.Remark != "" && params.Remark != beforeTask.Remark {
		result.Detail = append(result.Detail, structs.UpdateTaskDetail{
			Field:  "Remark",
			Before: beforeTask.Remark,
			After:  params.Remark,
		})
	}

	// 資料更新
	err = service.mysql.ChanelDB.Repository.Tasks.UpdateTask(structs.ChanelModelTasks{
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
	}, params.TaskID)

	if err != nil {
		response.Code = classes.MysqlInsertError
		response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlInsertError), "Tasks.UpdateTask", err)
		return response
	}

	// 寫入 Job Chan
	schedule.JobChan <- newJob

	// 回傳
	response.Result = result
	return response
}

func (service *Service) DeleteTask(params structs.DeleteTaskRequest, ctx context.Context) structs.Response {
	var (
		response = structs.Response{}
	)

	// 查詢
	err := service.mysql.ChanelDB.Repository.Tasks.DeleteTaskByID(params.TaskID)

	if err != nil {
		response.Code = classes.MysqlDeleteError
		response.Message = service.tools.FormatMsg(structs.RequestErrorMsg, "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlDeleteError), "Tasks.DeleteTaskByID", err)
		return response
	}

	return response
}
