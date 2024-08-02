package api

import (
	"chanel/classes"
	"chanel/config"
	"chanel/lib"
	"chanel/structs"
	"encoding/json"
	"fmt"
)

type Beckham struct {
	tools  *lib.Tools
	myErr  *classes.MyErr
	domain string
	port   string
}

func BeckhamInit(config *config.Config, tools *lib.Tools, myErr *classes.MyErr) *Beckham {
	return &Beckham{
		tools:  tools,
		myErr:  myErr,
		domain: config.BeckhamDomain,
		port:   config.BeckhamPort,
	}
}

// 取得版號
func (api *Beckham) GetVersion(params structs.BeckhamGetVersionRequest) (structs.BeckhamGetVersionResponse, error) {
	resp := structs.BeckhamGetVersionResponse{}
	curl := classes.CurlInit(api.tools)

	if result, err := curl.SetHttp().NewRequest(api.domain, api.port, "/usain/version").SetHeaders(map[string]string{
		"sid": params.Sid,
	}).SetQueries(map[string]interface{}{
		"hallId": api.tools.Int32ToStr(params.HallId),
		"device": params.Device,
	}).SetTraceID(params.TraceID).Get(); err != nil {
		return resp, err
	} else {
		// 解析
		err := json.Unmarshal([]byte(result), &resp)

		if err != nil {
			return resp, fmt.Errorf("%s\nString: %s\nErr: %s", api.myErr.Msg(classes.JsonUnmarshalError), result, err.Error())
		}
		return resp, nil
	}
}

func (api *Beckham) PutVersion(params structs.BeckhamPutVersionRequest) (structs.BeckhamPutVersionResponse, error) {
	resp := structs.BeckhamPutVersionResponse{}
	curl := classes.CurlInit(api.tools)

	if result, err := curl.SetHttp().NewRequest(api.domain, api.port, "/usain/version").SetHeaders(map[string]string{
		"sid": params.Sid,
	}).SetBody(map[string]interface{}{
		"device":  params.Device,
		"hallId":  api.tools.Int32ToStr(params.HallId),
		"url":     params.Url,
		"version": params.Version,
	}).SetTraceID(params.TraceID).Put(); err != nil {
		return resp, err
	} else {
		// 解析
		err := json.Unmarshal([]byte(result), &resp)

		if err != nil {
			return resp, fmt.Errorf("%s\nString: %s\nErr: %s", api.myErr.Msg(classes.JsonUnmarshalError), result, err.Error())
		}
		return resp, nil
	}
}

func (api *Beckham) PostAdminLogin(params structs.BeckhamPostAdminLoginRequest) (structs.BeckhamPostAdminLoginResponse, error) {
	resp := structs.BeckhamPostAdminLoginResponse{}
	curl := classes.CurlInit(api.tools)

	if result, err := curl.SetHttp().NewRequest(api.domain, api.port, "/usain/allctl/admin/login").SetHeaders(map[string]string{
		"sid": params.Sid,
	}).SetBody(map[string]interface{}{
		"account":   params.Account,
		"pwd":       params.Pwd,
		"equipment": params.Equipment,
	}).SetTraceID(params.TraceID).Post(); err != nil {
		return resp, err
	} else {
		// 解析
		err := json.Unmarshal([]byte(result), &resp)

		if err != nil {
			return resp, fmt.Errorf("%s\nString: %s\nErr: %s", api.myErr.Msg(classes.JsonUnmarshalError), result, err.Error())
		}
		return resp, nil
	}
}

func (api *Beckham) DeleteEmergencyData(params structs.BeckhamDeleteEmergencyDataRequest) (structs.BeckhamDeleteEmergencyDataResponse, error) {
	resp := structs.BeckhamDeleteEmergencyDataResponse{}
	curl := classes.CurlInit(api.tools)

	if result, err := curl.SetHttp().NewRequest(api.domain, api.port, "/usain/allctl/emergency/data").SetHeaders(map[string]string{
		"sid": params.Sid,
	}).SetQueries(map[string]interface{}{
		"hallId": api.tools.Int32ToStr(params.HallId),
	}).SetTraceID(params.TraceID).Delete(); err != nil {
		return resp, err
	} else {
		// 解析
		err := json.Unmarshal([]byte(result), &resp)

		if err != nil {
			return resp, fmt.Errorf("%s\nString: %s\nErr: %s", api.myErr.Msg(classes.JsonUnmarshalError), result, err.Error())
		}
		return resp, nil
	}
}
