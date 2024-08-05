package request

import (
	"chanel/classes"
	"chanel/structs"
	"net/http"

	"github.com/gorhill/cronexpr"
)

func (request *Request) GetCrons(params structs.GetCronsRequest, response structs.Response) structs.Response {
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

func (request *Request) CreateCron(params structs.CreateCronRequest, response structs.Response) structs.Response {
	// 必要參數
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
	} else {
		_, err := cronexpr.Parse(params.Execute)

		if err != nil {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Execute")
			return response
		}
	}

	return response
}

func (request *Request) UpdateCron(params structs.UpdateCronRequest, response structs.Response) structs.Response {
	if !request.tools.Request(params.CronID) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "CronID")
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

func (request *Request) DeleteCron(params structs.DeleteCronRequest, response structs.Response) structs.Response {
	if !request.tools.Request(params.CronID) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "CronID")
		return response
	}

	return response
}
