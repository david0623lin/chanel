package service

import (
	"chanel/api"
	"chanel/classes"
	"chanel/config"
	"chanel/database"
	"chanel/lib"
	"context"
	"fmt"
)

type Service struct {
	config *config.Config
	ctx    context.Context
	mysql  *database.Mysql
	redis  *database.Redis
	tools  *lib.Tools
	myErr  *classes.MyErr
	api    *api.Api
}

// 初始化 Service 物件
func ServiceInit(config *config.Config, ctx context.Context, mysql *database.Mysql, redis *database.Redis) *Service {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("service error -> %v", lib.PanicParser(err)))
		}
	}()

	return &Service{
		config: config,
		ctx:    ctx,
		mysql:  mysql,
		redis:  redis,
	}
}

func (service *Service) SetTools(tools *lib.Tools) {
	service.tools = tools
}

func (service *Service) SetError(myErr *classes.MyErr) {
	service.myErr = myErr
}

func (service *Service) SetApi(api *api.Api) {
	service.api = api
}
