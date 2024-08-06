package server

import (
	"bytes"
	"chanel/classes"
	"chanel/database"
	"chanel/lib"
	"chanel/structs"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// 路由檢查 中介層, 不符合直接 404 不給進服務
func (srv *Server) Available(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Set("Response", structs.Response{
					Code:    classes.SystemError,
					Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.SystemError), ""),
					Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.SystemError), "Available.Panic", lib.PanicParser(err)),
				})
				srv.response(c)
			}
		}()

		for _, route := range r.Routes() {
			if route.Path == c.Request.URL.Path {
				c.Next()
				break
			}
		}
		c.AbortWithStatus(http.StatusNotFound)
	}
}

// 設定 CORS同源政策 中介層
func (srv *Server) Cors(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.Set("Response", structs.Response{
				Code:    classes.SystemError,
				Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.SystemError), ""),
				Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.SystemError), "Cors.Panic", lib.PanicParser(err)),
			})
			srv.response(c)
		}
	}()

	// 同源設定
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")

	// 後台
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatusJSON(http.StatusNoContent, "")
	}

	// 新增上下文
	ctx := context.Background()
	// 寫入請求開始時間
	ctx = context.WithValue(ctx, structs.RequestTimeKey, time.Now())
	// 寫入請求上下文
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

// 服務維護 中介層
func (srv *Server) Maintain(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.Set("Response", structs.Response{
				Code:    classes.SystemError,
				Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.SystemError), ""),
				Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.SystemError), "Maintain.Panic", lib.PanicParser(err)),
			})
			srv.response(c)
		}
	}()

	exist, err := srv.redis.Client.Exists(srv.ctx, database.Maintain).Result()

	if err != nil {
		c.Set("Response", structs.Response{
			Code:    classes.CacheError,
			Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.CacheError), ""),
			Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.CacheError), "Maintain.Exists", err),
		})
		srv.response(c)
	}

	if exist == 1 {
		c.AbortWithStatusJSON(http.StatusOK, structs.Response{
			Code:    classes.SystemMaintain,
			Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.SystemMaintain), ""),
		})
	}
	c.Next()
}

// 處理 Session 中介層
func (srv *Server) Session(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.Set("Response", structs.Response{
				Code:    classes.SystemError,
				Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.SystemError), ""),
				Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.SystemError), "Session.Panic", lib.PanicParser(err)),
			})
			srv.response(c)
		}
	}()

	// 不需 Sid 白名單（自行依照需求增加）
	white := []string{
		"/chanel/admin/login",
		"/chanel/admin/register", // 後續要拿掉
	}

	if !srv.tools.InStrArray(c.Request.URL.Path, white) && c.Request.Header.Get(structs.SessionID) == "" {
		c.Set("Response", structs.Response{
			Code:    classes.MissingSession,
			Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.MissingSession), ""),
		})
		srv.response(c)
	}

	var (
		headers = structs.Headers{}
		ctx     = c.Request.Context()
	)

	// Sid 處理
	if c.Request.Header.Get(structs.SessionID) != "" {
		headers.Sid = c.Request.Header.Get(structs.SessionID)
		// 寫入來源端的 Sid
		ctx = context.WithValue(ctx, structs.SIDKey, c.Request.Header.Get(structs.SessionID))

		// todo 解析 Sid 取得 Session 資訊
		// srv.api
		session := structs.Session{}
		ctx = context.WithValue(ctx, structs.SessionKey, session)
	}

	// Tid 處理
	var traceID string

	if c.Request.Header.Get(structs.TraceID) == "" {
		var err error
		traceID, err = srv.tools.NewTraceID()

		if err != nil {
			c.Set("Response", structs.Response{
				Code:    classes.CreateTraceIdError,
				Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.CreateTraceIdError), ""),
				Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.CreateTraceIdError), "Session.NewTraceID", err),
			})
			srv.response(c)
		}
	} else {
		// 寫入來源端的 Tid
		headers.Tid = c.Request.Header.Get(structs.TraceID)
		traceID = c.Request.Header.Get(structs.TraceID)
	}
	ctx = context.WithValue(ctx, structs.TIDKey, traceID)

	// 儲存 Header 資訊
	ctx = context.WithValue(ctx, structs.HeadersKey, headers)

	// 寫入請求上下文
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

// 主流程 中介層
func (srv *Server) Service(c *gin.Context) {
	defer func() {
		// 刪除進行中請求
		srv.GracefulWaitGroup.Done()

		if err := recover(); err != nil {
			c.Set("Response", structs.Response{
				Code:    classes.SystemError,
				Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.SystemError), ""),
				Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.SystemError), "Service.Panic", lib.PanicParser(err)),
			})
			srv.response(c)
		}
	}()

	// 新增進行中請求
	srv.GracefulWaitGroup.Add(1)

	// 建立空的 Args
	args := make(map[string]interface{})
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, structs.ArgsKey, args)
	c.Request = c.Request.WithContext(ctx)

	switch c.Request.Method {
	case http.MethodGet, http.MethodDelete:
		// 解析 URL
		parsedURL, err := url.Parse(c.Request.URL.String())

		if err != nil {
			c.Set("Response", structs.Response{
				Code:    classes.ParseUrlParamsError,
				Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.ParseUrlParamsError), ""),
				Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.ParseUrlParamsError), "Service.url.Parse", err),
			})
			srv.response(c)
		}

		for key, values := range parsedURL.Query() {
			if len(values) == 1 {
				args[key] = values[0]
			} else {
				args[key] = values
			}
		}
		// 寫入 Args
		ctx = context.WithValue(ctx, structs.ArgsKey, args)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
		srv.response(c)
	case http.MethodPost, http.MethodPut:
		// 取得 body 資料
		bodyByte, err := io.ReadAll(c.Request.Body)

		if err != nil {
			c.Set("Response", structs.Response{
				Code:    classes.IoReadBodyError,
				Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.IoReadBodyError), ""),
				Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.IoReadBodyError), "Service.io.ReadAll", err),
			})
			srv.response(c)
		}
		err = json.Unmarshal(bodyByte, &args)

		if err != nil {
			c.Set("Response", structs.Response{
				Code:    classes.JsonUnmarshalError,
				Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.JsonUnmarshalError), ""),
				Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.JsonUnmarshalError), "Service.Unmarshal", err),
			})
			srv.response(c)
		}

		// 寫入 Args
		ctx = context.WithValue(ctx, structs.ArgsKey, args)
		c.Request = c.Request.WithContext(ctx)

		// 不需阻擋白名單
		white := []string{}

		if !srv.tools.InStrArray(c.Request.URL.String(), white) {
			body, err := json.Marshal(args)

			if err != nil {
				c.Set("Response", structs.Response{
					Code:    classes.JsonMarshalError,
					Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.JsonMarshalError), ""),
					Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.JsonMarshalError), "Service.Marshal", err),
				})
				srv.response(c)
			}

			// 產生快取 key
			key := c.Request.Method + "_" + c.Request.URL.String() + "_" + string(body)
			cacheKey := srv.tools.Sha256(key)

			// 使用 SETNX 命令設置值
			result, err := srv.redis.Client.SetNX(srv.ctx, cacheKey, true, 1*time.Second).Result()

			if err != nil {
				c.Set("Response", structs.Response{
					Code:    classes.CacheError,
					Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.CacheError), ""),
					Error:   srv.tools.FormatErr(srv.myErr.Msg(classes.CacheError), "Service.SetNX", err),
				})
				srv.response(c)
			}

			// key 已存在
			if !result {
				c.Set("Response", structs.Response{
					Code:    classes.OperatingTooFrequently,
					Message: srv.tools.FormatMsg(srv.myErr.Msg(classes.OperatingTooFrequently), ""),
				})
				srv.response(c)
			}
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyByte))
		c.Next()
		srv.response(c)
	}
}

func (srv *Server) response(c *gin.Context) {
	// 取得控制器執行完成後的回傳
	response := c.MustGet("Response").(structs.Response)

	// 取得請求上下文
	ctx := c.Request.Context()

	// 取得 Args
	args, ok := ctx.Value(structs.ArgsKey).(map[string]interface{})

	if !ok {
		args = make(map[string]interface{})
	}

	// 取得 Headers
	headers, ok := ctx.Value(structs.HeadersKey).(structs.Headers)

	if !ok {
		headers = structs.Headers{}
	}

	// 取得 TraceID
	traceID, ok := ctx.Value(structs.TIDKey).(string)

	if !ok {
		traceID = ""
	}

	// 取得 Session
	session, ok := ctx.Value(structs.SessionKey).(structs.Session)

	if !ok {
		session = structs.Session{}
	}

	// 初始化 TraceLog
	traceLog := classes.TraceLogInit(srv.tools)
	traceLog.SetTopic("server")
	traceLog.SetUrl(c.Request.URL.String())
	traceLog.SetMethod(c.Request.Method)
	traceLog.SetArgs(args)
	traceLog.SetHeaders(headers)
	traceLog.SetDomain(c.Request.Host)
	traceLog.SetClientIP(c.ClientIP())
	traceLog.SetCode(response.Code)
	traceLog.SetTraceID(traceID)
	traceLog.SetRequestTime(srv.tools.GetDownRunTime(ctx.Value(structs.RequestTimeKey).(time.Time)))
	traceLog.SetExtraInfo(session)

	// 失敗
	if response.Code != 0 {
		if response.Error != nil {
			traceLog.PrintError(response.Message, response.Error)

			// 正式站且非內部請求, 不回傳錯誤詳情資訊
			if srv.config.Env == "Prod" && c.Request.Header.Get(structs.TraceID) == "" {
				response.Error = errors.New("")
			}
		} else {
			traceLog.PrintWarn(response.Message, response.Error)
			response.Error = errors.New("")
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, structs.ServerReturnJson{
			Code:    response.Code,
			Message: response.Message,
			Result:  response.Result,
			Error:   response.Error.Error(),
		})
	} else {
		response.Message = srv.tools.FormatMsg(srv.myErr.Msg(classes.Success), "")

		// 如果不需要紀錄回傳結果, 自行註解
		traceLog.SetResponse(response.Result)

		traceLog.PrintInfo(response.Message)
		c.AbortWithStatusJSON(http.StatusOK, structs.ServerReturnJson{
			Code:    response.Code,
			Message: response.Message,
			Result:  response.Result,
		})
	}
}
