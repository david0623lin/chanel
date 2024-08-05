package classes

type MyErr struct{}

func ErrorInit() *MyErr {
	return &MyErr{}
}

const (
	// 共用
	SystemError    = -999
	RequestSuccess = 0

	// 請求檢查參數用錯誤
	ParserRequestBodyError = 1
	MissingRequireParams   = 2
	ParamsInvalid          = 3

	// 資料庫、快取用錯誤
	MysqlSearchError = 4
	MysqlInsertError = 5
	MysqlUpdateError = 6
	MysqlDeleteError = 7
	CacheError       = 8

	// Json 處理用錯誤
	JsonUnmarshalError = 9
	JsonMarshalError   = 10

	// 身份驗證錯誤
	MissingSession   = 11
	SessionNotFound  = 12
	PermissionDenied = 13

	// 流程回傳代碼
	SystemMaintain         = 14
	CallApiError           = 15
	RouteError             = 16
	OperatingTooFrequently = 17
	CreateTraceIdError     = 18
	ParseUrlParamsError    = 19
	IoReadBodyError        = 20
)

func (e *MyErr) Msg(code int32) string {
	return e.result(code)
}

func (e *MyErr) result(code int32) string {
	errCode := map[int32]string{
		// 共用
		SystemError:    "系統錯誤",
		RequestSuccess: "請求成功",

		// 請求檢查參數用錯誤
		ParserRequestBodyError: "解析請求Body資料錯誤",
		MissingRequireParams:   "缺少必要參數",
		ParamsInvalid:          "參數不合規定",

		// 資料庫、快取用錯誤
		MysqlSearchError: "資料庫查詢錯誤",
		MysqlInsertError: "資料庫新增錯誤",
		MysqlUpdateError: "資料庫修改錯誤",
		MysqlDeleteError: "資料庫刪除錯誤",
		CacheError:       "Redis 操作錯誤",

		// Json 處理用錯誤
		JsonUnmarshalError: "Json Decode 錯誤",
		JsonMarshalError:   "Json Encode 錯誤",

		// 身份驗證錯誤
		MissingSession:   "缺少Session",
		SessionNotFound:  "Session錯誤",
		PermissionDenied: "權限不足",

		// 流程回傳代碼
		SystemMaintain:         "系統維護中",
		CallApiError:           "呼叫外服務API錯誤",
		RouteError:             "路由錯誤",
		OperatingTooFrequently: "操作太頻繁",
		CreateTraceIdError:     "建立TraceID錯誤",
		ParseUrlParamsError:    "解析路由參數錯誤",
		IoReadBodyError:        "取得請求Body資料錯誤",
	}
	return errCode[code]
}
