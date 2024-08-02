package server

import (
	"chanel/api"
	"chanel/classes"
	"chanel/config"
	"chanel/controller"
	"chanel/database"
	"chanel/lib"
	"chanel/schedule"
	"context"
	"fmt"
	"sync"
)

type Server struct {
	config            *config.Config
	ctx               context.Context
	mysql             *database.Mysql
	redis             *database.Redis
	tools             *lib.Tools
	myErr             *classes.MyErr
	api               *api.Api
	controller        *controller.Controoller
	schedule          *schedule.Schedule
	GracefulWaitGroup sync.WaitGroup // 優雅退出使用
}

func ServerInit(config *config.Config, ctx context.Context, mysql *database.Mysql, redis *database.Redis) *Server {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("server error -> %v", lib.PanicParser(err)))
		}
	}()

	return &Server{
		config:            config,
		ctx:               ctx,
		mysql:             mysql,
		redis:             redis,
		GracefulWaitGroup: sync.WaitGroup{},
	}
}

func (srv *Server) SetTools(tools *lib.Tools) {
	srv.tools = tools
}

func (srv *Server) SetError(myErr *classes.MyErr) {
	srv.myErr = myErr
}

func (srv *Server) SetApi(api *api.Api) {
	srv.api = api
}

func (srv *Server) SetController() {
	controller := controller.ControllerInit(srv.config, srv.ctx, srv.mysql, srv.redis)
	controller.SetTools(srv.tools)
	controller.SetError(srv.myErr)
	controller.SetRequest()
	controller.SetService(srv.api)

	srv.controller = controller
}
