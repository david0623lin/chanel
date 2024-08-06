package request

import (
	"chanel/classes"
	"chanel/structs"
	"regexp"
)

func (request *Request) Login(params structs.AdminLoginRequest, response structs.Response) structs.Response {
	// 必要參數
	if !request.tools.Request(params.Account) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "Account")
		return response
	} else {
		// 允許字母、數字、下劃線和破折號，至少3個字符，最多20個字符
		if !regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`).MatchString(params.Account) {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Account")
			return response
		}
	}

	if !request.tools.Request(params.Password) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "Password")
		return response
	} else {
		// 不允許空格，至少8個字符，最多100個字符
		if !regexp.MustCompile(`^[^\s]{8,20}$`).MatchString(params.Password) {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Password")
			return response
		}
	}

	return response
}

func (request *Request) Register(params structs.AdminRegisterRequest, response structs.Response) structs.Response {
	// 必要參數
	if !request.tools.Request(params.Account) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "Account")
		return response
	} else {
		// 允許字母、數字、下劃線和破折號，至少3個字符，最多20個字符
		if !regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`).MatchString(params.Account) {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Account")
			return response
		}
	}

	if !request.tools.Request(params.Password) {
		response.Code = classes.MissingRequireParams
		response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.MissingRequireParams), "Password")
		return response
	} else {
		// 不允許空格，至少8個字符，最多100個字符
		if !regexp.MustCompile(`^[^\s]{8,20}$`).MatchString(params.Password) {
			response.Code = classes.ParamsInvalid
			response.Message = request.tools.FormatMsg(request.myErr.Msg(classes.ParamsInvalid), "Password")
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
