package service

import (
	"chanel/classes"
	"chanel/structs"
	"context"
	"encoding/base64"
	"encoding/json"
	"time"
)

func (service *Service) Login(params structs.AdminLoginRequest, ctx context.Context) structs.Response {
	var (
		response = structs.Response{}
	)

	// 取得帳號
	admin, err := service.mysql.ChanelDB.Repository.Admins.GetAdmin(params)

	if err != nil {
		if err.Error() == "RecordNotFound" {
			response.Code = classes.AccountNotFound
			response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.AccountNotFound), "")
			return response
		}
		response.Code = classes.MysqlSearchError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.MysqlSearchError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlSearchError), "Admins.GetAdmin", err)
		return response
	}

	// 驗證密碼
	if admin.Password != service.tools.PwdEncode(params.Password) {
		response.Code = classes.PasswordError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.PasswordError), "")
		return response
	}

	// 檢查已停用
	if admin.Status == 2 {
		response.Code = classes.AccountDisable
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.AccountDisable), "")
		return response
	}

	// 產生 sid
	var sessionID string
	sessionID, err = service.tools.NewSessionID(admin.Uuid)

	if err != nil {
		response.Code = classes.MakeSessionError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.MakeSessionError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MakeSessionError), "NewSessionID", err)
		return response
	}

	// 寫入登入快取
	err = service.redis.Client.Set(service.ctx, admin.Uuid, sessionID, 1*time.Hour).Err()

	if err != nil {
		response.Code = classes.CacheError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.CacheError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.CacheError), "Set.sessionID", err)
		return response
	}

	// 產生 wid
	var encryptData []byte
	encryptData, err = json.Marshal(structs.Websockets{
		Sid:  sessionID,
		Uuid: admin.Uuid,
	})

	if err != nil {
		response.Code = classes.JsonMarshalError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.JsonMarshalError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.JsonMarshalError), "Marshal", err)
		return response
	}

	// 加密
	key := []byte(service.tools.Md5(service.config.WsMd5Salt))
	encryptText := service.tools.AesEncryptCBC(encryptData, key)
	websocketID := base64.StdEncoding.EncodeToString(encryptText)

	var result = structs.AdminLoginResponse{
		Account: params.Account,
		Sid:     sessionID,
		Wid:     websocketID,
	}

	// 回傳
	response.Result = result
	return response
}

func (service *Service) Register(params structs.AdminRegisterRequest, ctx context.Context) structs.Response {
	var (
		response = structs.Response{}
	)

	// 加密密碼
	params.Password = service.tools.PwdEncode(params.Password)

	// 帳號狀態預設啟用
	if params.Status == 0 {
		params.Status = 1
	}

	// 取得帳號
	err := service.mysql.ChanelDB.Repository.Admins.CreateAdmin(params)

	if err != nil {
		response.Code = classes.MysqlInsertError
		response.Message = service.tools.FormatMsg(service.myErr.Msg(classes.MysqlInsertError), "")
		response.Error = service.tools.FormatErr(service.myErr.Msg(classes.MysqlInsertError), "Admins.CreateAdmin", err)
		return response
	}
	return response
}
