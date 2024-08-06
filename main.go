package main

import (
	"chanel/api"
	"chanel/classes"
	"chanel/config"
	"chanel/database"
	"chanel/lib"
	"chanel/schedule"
	"chanel/server"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "chanel/docs"

	_ "github.com/joho/godotenv/autoload"
)

// @title chanel
// @version 1.0
// @description API Demo 頁面
// @schemes http
func main() {
	// recover 防止因服務 panic 直接關閉
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(PanicFormat(err))
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())

	// 自動載入 .env
	config := config.NewConfig()

	// 初始化 自定義工具
	tools := lib.ToolsInit(config)

	// 初始化 自定義錯誤
	myErr := classes.ErrorInit()

	// 初始化 Mysql
	mysql := database.MysqlInit(config)
	mysql.Start()

	// 初始化 Redis
	redis := database.RedisInit(config, ctx, tools)
	redis.Start()

	// 初始化 Api
	api := api.ApiInit(config, ctx)
	api.SetTools(tools)
	api.SetError(myErr)
	api.SetRepo()

	// 初始化 主服務
	srv := server.ServerInit(config, ctx, mysql, redis)
	srv.SetTools(tools)
	srv.SetError(myErr)
	srv.SetApi(api)
	srv.SetController()
	go srv.Start()

	// 初始化 Schedules
	schedules := schedule.ScheduleInit(config, ctx, mysql, redis)
	schedules.SetTools(tools)
	schedules.SetError(myErr)
	schedules.SetApi(api)
	schedule.JobChan = make(chan *schedule.Job, 10)

	go schedules.LoadJobs()
	go schedules.StartJobs()

	// 初始化 Websocket
	websocket := classes.WebsocketInit(config, ctx, redis)
	websocket.SetTools(tools)
	websocket.SetError(myErr)

	go websocket.Start()
	go websocket.Messages()
	go websocket.Heart()

	// 優雅退出
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	exit := make(chan struct{})

	go func() {
		<-signals
		// 優雅退出要結束的程式寫在這 Ex.關閉連線
		cancel()
		srv.StopServer()
		// close(schedule.TaskChan)
		fmt.Println("Exit.")
		exit <- struct{}{}
	}()
	<-ctx.Done()
	<-exit
}

func PanicFormat(e interface{}) string {
	return fmt.Sprintf("[Panic] %v", e)
}
