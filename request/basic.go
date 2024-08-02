package request

import (
	"chanel/classes"
	"chanel/lib"
)

type Request struct {
	tools *lib.Tools
	myErr *classes.MyErr
}

func RequestInit(tools *lib.Tools, myErr *classes.MyErr) *Request {
	return &Request{
		tools: tools,
		myErr: myErr,
	}
}
