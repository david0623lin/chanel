package server

import (
	"chanel/structs"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ! 路由的命名不可包含 GET|POST|PUT|DELETE, 對應的 Controller 方法命名不在此限制

// 執行主服務
func (srv *Server) Start() {
	// 關閉 gin 詳細資訊
	if srv.config.Env == "Prod" {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()

		// 關掉 gin 預設的請求紀錄
		gin.DefaultWriter = io.Discard
	}
	router := gin.Default()

	// PING
	router.GET("/ping", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, structs.ServerReturnJson{
			Message: "PONG",
		})
	})

	// 載入 swagger
	if srv.config.Env != "Prod" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 中介層
	router.Use(srv.Available(router))
	router.Use(srv.Cors)
	router.Use(srv.Maintain)
	router.Use(srv.Session)
	router.Use(srv.Service)

	// 主群組
	chanelGroup := router.Group("chanel")
	{
		// 任務
		taskGroup := chanelGroup.Group("task")
		taskGroup.GET("/list", srv.controller.GetTasks)
		taskGroup.GET("/detail", srv.controller.GetTaskDetail)
		taskGroup.POST("/create", srv.controller.CreateTask)
		taskGroup.PUT("/update", srv.controller.UpdateTask)
		taskGroup.DELETE("/remove", srv.controller.DeleteTask)

		// 排程
		cronGroup := chanelGroup.Group("cron")
		cronGroup.GET("/list", srv.controller.GetCrons)
		cronGroup.POST("/create", srv.controller.CreateCron)
		cronGroup.PUT("/update", srv.controller.UpdateCron)
		cronGroup.DELETE("/remove", srv.controller.DeleteCron)

		// 管理後台
		AdminGroup := chanelGroup.Group("admin")
		AdminGroup.POST("/login", srv.controller.Login)
		AdminGroup.POST("/register", srv.controller.Register)
	}
	// 啟動服務
	if err := router.Run(fmt.Sprintf(":%s", srv.config.ServerPort)); err != nil {
		panic(err)
	}
}

func (srv *Server) StopServer() {
	srv.GracefulWaitGroup.Wait()
	fmt.Println("ApiGateway Graceful Stop")
}
