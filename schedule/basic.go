package schedule

import (
	"chanel/api"
	"chanel/classes"
	"chanel/config"
	"chanel/database"
	"chanel/lib"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

type Schedule struct {
	config   *config.Config
	ctx      context.Context
	mysql    *database.Mysql
	redis    *database.Redis
	tools    *lib.Tools
	myErr    *classes.MyErr
	traceLog *classes.TraceLog

	api *api.Api
}

// 任務頻道
var JobChan chan *Job

type Job struct {
	ID          int32
	Protocol    string
	Domain      string
	Path        string
	Port        string
	Method      string
	Args        map[string]interface{}
	Headers     map[string]string
	Mode        int32 // 1:任務, 2:排程
	ExecuteTask int64
	ExecuteCron string
}

const (
	Success = "SUCCESS"
	JobGap  = 100 * time.Millisecond // 間隔時間
	Task    = 1                      // 任務
	Cron    = 2                      // 排程
)

func ScheduleInit(config *config.Config, ctx context.Context, mysql *database.Mysql, redis *database.Redis) *Schedule {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("schedule error -> %v", lib.PanicParser(err)))
		}
	}()

	return &Schedule{
		config: config,
		ctx:    ctx,
		mysql:  mysql,
		redis:  redis,
	}
}

func (schedule *Schedule) SetTools(tools *lib.Tools) {
	schedule.tools = tools
}

func (schedule *Schedule) SetTraceLog() {
	// 初始化 Trace log
	traceLog := classes.TraceLogInit(schedule.tools)
	traceLog.SetTopic("schedule")

	schedule.traceLog = traceLog
}

func (schedule *Schedule) SetError(myErr *classes.MyErr) {
	schedule.myErr = myErr
}

func (schedule *Schedule) SetApi(api *api.Api) {
	schedule.api = api
}

// 服務啟動時載入 Job
func (schedule *Schedule) LoadJobs() {
	var jobs []*Job

	// 取得未執行的任務
	taskLists, err := schedule.mysql.ChanelDB.Repository.Tasks.GetNotCompletedTasks()

	if err != nil {
		panic(fmt.Sprintf("Schedule LoadJobs 取得未執行的任務錯誤, ERR: %s", err.Error()))
	}

	for _, taskItem := range taskLists {
		// 把取得的任務都解析後寫入任務列表, 後續一次新增, 避免新增後執行到一半遇到解析異常 Panic 中斷情況
		jobs = append(jobs, &Job{
			ID:          taskItem.ID,
			Protocol:    taskItem.Protocol,
			Domain:      taskItem.Domain,
			Path:        taskItem.Path,
			Port:        taskItem.Port,
			Method:      taskItem.Method,
			Args:        schedule.parserArgs(taskItem.ID, taskItem.Args),
			Headers:     schedule.parserHeaders(taskItem.ID, taskItem.Headers),
			Mode:        Task,
			ExecuteTask: taskItem.Execute,
		})
	}

	// 取得啟用中的排程
	cronLists, err := schedule.mysql.ChanelDB.Repository.Crons.GetCrons()

	if err != nil {
		panic(fmt.Sprintf("Schedule LoadJobs 取得啟用中的排程錯誤, ERR: %s", err.Error()))
	}

	for _, cronItem := range cronLists {
		// 把取得的任務都解析後寫入任務列表, 後續一次新增, 避免新增後執行到一半遇到解析異常 Panic 中斷情況
		jobs = append(jobs, &Job{
			ID:          cronItem.ID,
			Protocol:    cronItem.Protocol,
			Domain:      cronItem.Domain,
			Path:        cronItem.Path,
			Port:        cronItem.Port,
			Method:      cronItem.Method,
			Args:        schedule.parserArgs(cronItem.ID, cronItem.Args),
			Headers:     schedule.parserHeaders(cronItem.ID, cronItem.Headers),
			Mode:        Cron,
			ExecuteCron: cronItem.Execute,
		})
	}

	// 一次新增所有解析正確完成的任務到 TaskChan
	for _, job := range jobs {
		JobChan <- job
	}
}

func (schedule *Schedule) StartJobs() {
	// 排程紀錄下一個執行時間用
	cronNext := make(map[int32]time.Time)

	for {
		select {
		case Job := <-JobChan:
			switch Job.Mode {
			case Task:
				if schedule.tools.NowUnix() < Job.ExecuteTask {
					JobChan <- Job
				} else {
					go schedule.doTask(Job)
				}
			case Cron:
				// 解析 排程字串 (新增時已經檢查過就不再檢查 err)
				cronExpr, _ := cronexpr.Parse(Job.ExecuteCron)

				if _, exists := cronNext[Job.ID]; !exists {
					// 寫入 下一次排程時間
					cronNext[Job.ID] = cronExpr.Next(schedule.tools.Now())
				}

				if schedule.tools.NowUnix() == cronNext[Job.ID].Unix() {
					go schedule.doCron(Job)

					// 更新 下一次排程時間
					cronExpr, _ = cronexpr.Parse(Job.ExecuteCron)
					cronNext[Job.ID] = cronExpr.Next(schedule.tools.Now())
				}
				// 一定要把 Cron 寫回去 Job, 非常重要 !!
				JobChan <- Job
			}
		default:
		}

		time.Sleep(JobGap)
	}
}

// 解析 參數
func (schedule *Schedule) parserArgs(id int32, args string) map[string]interface{} {
	var resp map[string]interface{}

	if args != "" {
		err := json.Unmarshal([]byte(args), &resp)

		if err != nil {
			panic(fmt.Sprintf("Schedule LoadJobs ID[%d] 解析參數錯誤, ERR: %s", id, err.Error()))
		}
	}
	return resp
}

// 解析 表頭
func (schedule *Schedule) parserHeaders(id int32, headers string) map[string]string {
	var resp = make(map[string]string)

	if headers != "" {
		var header map[string]interface{}
		err := json.Unmarshal([]byte(headers), &header)

		if err != nil {
			panic(fmt.Sprintf("Schedule LoadJobs ID[%d] 解析表頭錯誤, ERR: %s", id, err.Error()))
		}

		for k, v := range header {
			resp[k] = fmt.Sprintf("%v", v)
		}
	}
	return resp
}
