package api

import (
	"chanel/api/proto/camila"
	"chanel/classes"
	"chanel/config"
	"chanel/lib"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Camila struct {
	client camila.CamilaServiceClient
	myErr  *classes.MyErr
}

func CamilaInit(config *config.Config, gRpcSet grpc.DialOption, myErr *classes.MyErr) *Camila {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("camila error -> %v", lib.PanicParser(err)))
		}
	}()

	conn, err := grpc.Dial(
		config.CamilaDomain+":"+config.CamilaPort,
		grpc.WithTransportCredentials(credentials.NewTLS(nil)),
		gRpcSet,
	)

	if err != nil {
		panic(fmt.Sprintf("gRPC 連線錯誤, ERR: %s", err.Error()))
	}

	return &Camila{
		client: camila.NewCamilaServiceClient(conn),
		myErr:  myErr,
	}
}
