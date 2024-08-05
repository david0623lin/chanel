package structs

const (
	// 請求 Headers 固定欄位
	SessionID = "Sid"
	TraceID   = "Tid"
)

// 流程存放要回傳的資訊用
type Response struct {
	Code    int32
	Message string
	Result  interface{}
	Error   error
}

// API 最終回傳用
type ServerReturnJson struct {
	Code    int32
	Message string
	Result  interface{}
	Error   string
}

// 存放請求原始 Headers 資訊
type Headers struct {
	Sid string
	Tid string
}

// 存放解析後的完整 Session 資訊
type Session struct {
}

type contextKey string

const HeadersKey = contextKey("Headers")
const SIDKey = contextKey("SID")
const SessionKey = contextKey("Session")
const TIDKey = contextKey("TID")
const RequestTimeKey = contextKey("RequestTime")
const ArgsKey = contextKey("Args")

type EmptyRequest struct{}

// todo 由此往下新增主服務的 API 結構
// todo -------------------------

type GetTasksRequest struct {
	StartTime int64
	EndTime   int64
	Path      string
	Method    string
	Status    int32
}

type GetTasksResponse struct {
	ID         int32
	Protocol   string
	Domain     string
	Path       string
	Port       string
	Method     string
	Args       map[string]interface{}
	Headers    map[string]string
	Execute    int64
	Status     int32
	Remark     string
	CreateTime int64
	UpdateTime int64
}

type CreateTaskRequest struct {
	Protocol string
	Domain   string
	Path     string
	Port     string
	Method   string
	Args     map[string]interface{}
	Headers  map[string]string
	Execute  int64
	Remark   string
}

type CreateTaskResponse struct {
}

type UpdateTaskRequest struct {
	TaskID   int32
	Protocol string
	Domain   string
	Path     string
	Port     string
	Method   string
	Args     map[string]interface{}
	Headers  map[string]string
	Execute  int64
	Remark   string
}

type UpdateTaskResponse struct {
	Detail []UpdateTaskDetail
}

type UpdateTaskDetail struct {
	Field  string
	Before interface{}
	After  interface{}
}

type DeleteTaskRequest struct {
	TaskID int32
}

type DeleteTaskResponse struct {
}

type GetCronsRequest struct {
	Path   string
	Method string
	Status int32
}

type GetCronsResponse struct {
	ID         int32
	Protocol   string
	Domain     string
	Path       string
	Port       string
	Method     string
	Args       map[string]interface{}
	Headers    map[string]string
	Execute    string
	Status     int32
	Remark     string
	CreateTime int64
	UpdateTime int64
}

type CreateCronRequest struct {
	Protocol string
	Domain   string
	Path     string
	Port     string
	Method   string
	Args     map[string]interface{}
	Headers  map[string]string
	Execute  string
	Remark   string
}

type CreateCronResponse struct {
}

type UpdateCronRequest struct {
	CronID   int32
	Protocol string
	Domain   string
	Path     string
	Port     string
	Method   string
	Args     map[string]interface{}
	Headers  map[string]string
	Execute  string
	Remark   string
}

type UpdateCronResponse struct {
	Detail []UpdateCronDetail
}

type UpdateCronDetail struct {
	Field  string
	Before interface{}
	After  interface{}
}

type DeleteCronRequest struct {
	CronID int32
}

type DeleteCronResponse struct {
}
