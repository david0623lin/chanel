package controller

import (
	"chanel/api"
	"chanel/classes"
	"chanel/config"
	"chanel/database"
	"chanel/lib"
	"chanel/request"
	"chanel/service"
	"context"
	"fmt"
)

type Controoller struct {
	config  *config.Config
	ctx     context.Context
	mysql   *database.Mysql
	redis   *database.Redis
	tools   *lib.Tools
	myErr   *classes.MyErr
	request *request.Request
	service *service.Service
}

// 初始化 Controller 物件
func ControllerInit(config *config.Config, ctx context.Context, mysql *database.Mysql, redis *database.Redis) *Controoller {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("controller error -> %v", lib.PanicParser(err)))
		}
	}()

	return &Controoller{
		config: config,
		ctx:    ctx,
		mysql:  mysql,
		redis:  redis,
	}
}

func (ctl *Controoller) SetTools(tools *lib.Tools) {
	ctl.tools = tools
}

func (ctl *Controoller) SetError(myErr *classes.MyErr) {
	ctl.myErr = myErr
}

func (ctl *Controoller) SetRequest() {
	ctl.request = request.RequestInit(ctl.tools, ctl.myErr)
}

func (ctl *Controoller) SetService(api *api.Api) {
	service := service.ServiceInit(ctl.config, ctl.ctx, ctl.mysql, ctl.redis)
	service.SetTools(ctl.tools)
	service.SetError(ctl.myErr)
	service.SetApi(api)

	ctl.service = service
}
