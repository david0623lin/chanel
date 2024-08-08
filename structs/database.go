package structs

// * 資料庫模組結構
// * -----------------
// * 結構命名規則： DB名稱 + Model前綴 + 資料表名稱

type ChanelModelTasks struct {
	ID         int32  `gorm:"primaryKey; autoIncrement;" json:"id"`
	Topic      string `gorm:"index; type:varchar(50); default:''; not null; comment:'主題'" json:"topic"`
	Protocol   string `gorm:"type:varchar(10); default:''; not null; comment:'協議'" json:"protocol"`
	Domain     string `gorm:"index; type:varchar(50); default:''; not null; comment:'網域'" json:"domain"`
	Path       string `gorm:"index; type:varchar(50); default:''; not null; comment:'路徑'" json:"path"`
	Port       string `gorm:"type:varchar(10); default:''; not null; comment:'埠號'" json:"port"`
	Method     string `gorm:"index; type:varchar(10); default:''; not null; comment:'方法'" json:"method"`
	Args       string `gorm:"type:varchar(255); default:''; not null; comment:'參數'" json:"args"`
	Headers    string `gorm:"type:varchar(255); default:''; not null; comment:'表頭'" json:"headers"`
	Execute    int64  `gorm:"type:int(10); default:0; not null; comment:'執行時間'" json:"execute"`
	Status     int32  `gorm:"index; type:int(10); default:0; not null; comment:'狀態 1:未執行,2:已執行,3:異常'" json:"status"`
	Remark     string `gorm:"type:varchar(255); default:''; not null; comment:'備註'" json:"remark"`
	CreateTime int64  `gorm:"type:int(10); default:0; not null; comment:'建立時間'" json:"create_time"`
	UpdateTime int64  `gorm:"type:int(10); default:0; not null; comment:'更新時間'" json:"update_time"`
}

func (ChanelModelTasks) TableName() string {
	return "tasks"
}

type ChanelModelTaskRecords struct {
	ID         int32  `gorm:"primaryKey; autoIncrement;" json:"id"`
	TaskID     int32  `gorm:"unique; type:int(10); comment:'任務ID'" json:"task_id"`
	Status     int32  `gorm:"index; type:int(10); default:0; not null; comment:'狀態 1:成功,2:失敗'" json:"status"`
	Result     string `gorm:"type:varchar(255); default:''; not null; comment:'執行結果'" json:"result"`
	Error      string `gorm:"type:varchar(255); default:''; not null; comment:'錯誤資訊'" json:"error"`
	CreateTime int64  `gorm:"type:int(10); default:0; not null; comment:'建立時間'" json:"create_time"`
}

func (ChanelModelTaskRecords) TableName() string {
	return "task_records"
}

type ChanelModelCrons struct {
	ID         int32  `gorm:"primaryKey; autoIncrement;" json:"id"`
	Protocol   string `gorm:"type:varchar(10); default:''; not null; comment:'協議'" json:"protocol"`
	Domain     string `gorm:"index; type:varchar(50); default:''; not null; comment:'網域'" json:"domain"`
	Path       string `gorm:"index; type:varchar(50); default:''; not null; comment:'路徑'" json:"path"`
	Port       string `gorm:"type:varchar(10); default:''; not null; comment:'埠號'" json:"port"`
	Method     string `gorm:"index; type:varchar(10); default:''; not null; comment:'方法'" json:"method"`
	Args       string `gorm:"type:varchar(255); default:''; not null; comment:'參數'" json:"args"`
	Headers    string `gorm:"type:varchar(255); default:''; not null; comment:'表頭'" json:"headers"`
	Execute    string `gorm:"type:varchar(20); default:''; not null; comment:'執行時間'" json:"execute"`
	Status     int32  `gorm:"index; type:int(10); default:0; not null; comment:'狀態 1:啟用,2:停用'" json:"status"`
	Remark     string `gorm:"type:varchar(255); default:''; not null; comment:'備註'" json:"remark"`
	CreateTime int64  `gorm:"type:int(10); default:0; not null; comment:'建立時間'" json:"create_time"`
	UpdateTime int64  `gorm:"type:int(10); default:0; not null; comment:'更新時間'" json:"update_time"`
}

func (ChanelModelCrons) TableName() string {
	return "crons"
}

type ChanelModelCronRecords struct {
	ID         int32  `gorm:"primaryKey; autoIncrement;" json:"id"`
	CronID     int32  `gorm:"index; type:int(10); default:0; not null; comment:'排程ID'" json:"task_id"`
	Status     int32  `gorm:"index; type:int(10); default:0; not null; comment:'狀態 1:成功,2:失敗'" json:"status"`
	Result     string `gorm:"type:varchar(255); default:''; not null; comment:'執行結果'" json:"result"`
	Error      string `gorm:"type:varchar(255); default:''; not null; comment:'錯誤資訊'" json:"error"`
	CreateTime int64  `gorm:"type:int(10); default:0; not null; comment:'建立時間'" json:"create_time"`
}

func (ChanelModelCronRecords) TableName() string {
	return "cron_records"
}

type ChanelModelAdmins struct {
	ID         int32  `gorm:"primaryKey; autoIncrement;" json:"id"`
	Uuid       string `gorm:"unique; type:varchar(50); default:''; not null; comment:'使用者ID'" json:"uuid"`
	Account    string `gorm:"unique; type:varchar(20); default:''; not null; comment:'帳號'" json:"account"`
	Password   string `gorm:"type:varchar(50); default:''; not null; comment:'密碼'" json:"password"`
	Status     int32  `gorm:"index; type:int(10); default:0; not null; comment:'狀態 1:啟用,2:停用'" json:"status"`
	Remark     string `gorm:"type:varchar(255); default:''; not null; comment:'備註'" json:"remark"`
	CreateTime int64  `gorm:"type:int(10); default:0; not null; comment:'建立時間'" json:"create_time"`
	UpdateTime int64  `gorm:"type:int(10); default:0; not null; comment:'更新時間'" json:"update_time"`
}

func (ChanelModelAdmins) TableName() string {
	return "admins"
}
