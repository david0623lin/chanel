package structs

// * 外服務 API 相關結構
// * -----------------
// * 檔案命名規則：專案名稱
// * 結構命名規則：專案名稱 + API方法 + API名稱 + 類型（請求 or 回傳）

type BeckhamGetVersionRequest struct {
	TraceID string
	Sid     string
	HallId  int32
	Device  string
}

type BeckhamGetVersionResponse struct {
	Code    int32
	Message string
	Result  *struct {
		Version string
		Url     string
	}
}

type BeckhamPostAdminLoginRequest struct {
	TraceID   string
	Sid       string
	Account   string
	Pwd       string
	Equipment string
}

type BeckhamPostAdminLoginResponse struct {
	Code    int32
	Message string
	Result  *struct {
		Version string
		Url     string
	}
}

type BeckhamPutVersionRequest struct {
	TraceID string
	Sid     string
	Device  string
	HallId  int32
	Url     string
	Version string
}

type BeckhamPutVersionResponse struct {
	Code    int32
	Message string
	Result  *struct {
		Version string
		Url     string
	}
}

type BeckhamDeleteEmergencyDataRequest struct {
	TraceID string
	Sid     string
	HallId  int32
}

type BeckhamDeleteEmergencyDataResponse struct {
	Code    int32
	Message string
	Result  *struct {
		Version string
		Url     string
	}
}
