package api

import (
	"chanel/api/proto/camila"
	"chanel/classes"
	"chanel/config"
	"chanel/lib"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Api struct {
	config      *config.Config
	ctx         context.Context
	tools       *lib.Tools
	myErr       *classes.MyErr
	gRpcSetting grpc.DialOption // 設定 gRPC Keepalive 參數

	// TODO 由此往下新增其他外服務物件
	Beckham *Beckham
	Camila  camila.CamilaServiceClient
}

// 初始化 API 物件
func ApiInit(config *config.Config, ctx context.Context) *Api {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("api error -> %v", lib.PanicParser(err)))
		}
	}()

	api := &Api{
		config: config,
		ctx:    ctx,
		gRpcSetting: grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    10 * time.Second, // 客户端发送心跳的频率
			Timeout: 0 * time.Second,  // 连接断开后等待重新连接的超时时间
		}),
	}

	return api
}

func (api *Api) SetTools(tools *lib.Tools) {
	api.tools = tools
}

func (api *Api) SetError(myErr *classes.MyErr) {
	api.myErr = myErr
}

// TODO 由此往下新增其他外服務物件
func (api *Api) SetRepo() {
	api.Beckham = BeckhamInit(api.config, api.tools, api.myErr)
	api.Camila = CamilaInit(api.config, api.gRpcSetting, api.myErr).client
}
