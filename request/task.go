package request

import (
	"chanel/classes"
	"chanel/structs"
	"net/http"
)

func (request *Request) GetTasks(params structs.GetTasksRequest, response structs.Response) structs.Response {
	if request.tools.Request(params.StartTime) && request.tools.Request(params.EndTime) {
		if params.StartTime > params.EndTime {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "StartTime Or EndTime")
			return response
		}
	}

	if request.tools.Request(params.Method) {
		if !request.tools.InStrArray(params.Method, []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}) {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Method")
			return response
		}
	}

	if request.tools.Request(params.Status) {
		if !request.tools.InInt32Array(params.Status, []int32{1, 2}) {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Status")
			return response
		}
	}
	return response
}

func (request *Request) CreateTask(params structs.CreateTaskRequest, response structs.Response) structs.Response {
	if !request.tools.Request(params.Protocol) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "Protocol")
		return response
	} else {
		if !request.tools.InStrArray(params.Protocol, []string{classes.ProtocolHttp, classes.ProtocolHttps}) {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Protocol")
			return response
		}
	}

	if !request.tools.Request(params.Domain) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "Domain")
		return response
	}

	if !request.tools.Request(params.Path) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "Path")
		return response
	}

	if !request.tools.Request(params.Method) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "Method")
		return response
	} else {
		if !request.tools.InStrArray(params.Method, []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}) {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Method")
			return response
		}
	}

	if !request.tools.Request(params.Execute) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "Execute")
		return response
	}

	return response
}

func (request *Request) UpdateTask(params structs.UpdateTaskRequest, response structs.Response) structs.Response {
	if !request.tools.Request(params.TaskID) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "TaskID")
		return response
	}

	if request.tools.Request(params.Protocol) {
		if !request.tools.InStrArray(params.Protocol, []string{classes.ProtocolHttp, classes.ProtocolHttps}) {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Protocol")
			return response
		}
	}

	if request.tools.Request(params.Method) {
		if !request.tools.InStrArray(params.Method, []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}) {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Method")
			return response
		}
	}

	return response
}

func (request *Request) DeleteTask(params structs.DeleteTaskRequest, response structs.Response) structs.Response {
	if !request.tools.Request(params.TaskID) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "TaskID")
		return response
	}

	return response
}
