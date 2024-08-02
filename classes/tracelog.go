package classes

import (
	"chanel/lib"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type TraceLog struct {
	tools *lib.Tools

	// 必要資訊
	service string
	podName string

	// 選填資訊
	topic       string
	url         string
	method      string
	args        interface{}
	headers     interface{}
	domain      string
	clientIP    string
	err         error
	code        int32
	response    interface{}
	extraInfo   interface{}
	requestTime float64
	traceID     string
}

func TraceLogInit(tools *lib.Tools) *TraceLog {
	return &TraceLog{
		tools:   tools,
		service: getServiceName(),
		podName: getPodName(),
	}
}

func getPodName() (podName string) {
	return os.Getenv("HOSTNAME")
}

func getServiceName() (sevice string) {
	return os.Getenv("SERVICE_NAME")
}

func (tl *TraceLog) SetTopic(topic string) {
	tl.topic = topic
}

func (tl *TraceLog) SetUrl(url string) {
	tl.url = url
}

func (tl *TraceLog) SetMethod(method string) {
	tl.method = method
}

func (tl *TraceLog) SetArgs(args interface{}) {
	tl.args = args
}

func (tl *TraceLog) SetHeaders(headers interface{}) {
	tl.headers = headers
}

func (tl *TraceLog) SetDomain(domain string) {
	tl.domain = domain
}

func (tl *TraceLog) SetClientIP(ip string) {
	tl.clientIP = ip
}

func (tl *TraceLog) SetCode(code int32) {
	tl.code = code
}

func (tl *TraceLog) SetResponse(response interface{}) {
	tl.response = response
}

func (tl *TraceLog) SetExtraInfo(extra interface{}) {
	tl.extraInfo = extra
}

func (tl *TraceLog) SetRequestTime(requestTime float64) {
	tl.requestTime = requestTime
}

func (tl *TraceLog) SetTraceID(traceID string) {
	tl.traceID = traceID
}

func getlogfields(tl *TraceLog) logrus.Fields {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	})

	logField := logrus.Fields{}

	// 固定資訊
	logField["service"] = tl.service
	logField["podName"] = tl.podName
	logField["url"] = tl.url
	logField["method"] = tl.method
	logField["headers"] = tl.headers
	logField["requestTime"] = tl.requestTime
	logField["traceID"] = tl.traceID

	// 非固定資訊
	if tl.topic != "" {
		logField["topic"] = tl.topic
	}
	if tl.args != nil {
		logField["args"] = tl.args
	}
	if tl.domain != "" {
		logField["domain"] = tl.domain
	}
	if tl.clientIP != "" {
		logField["clientIP"] = tl.clientIP
	}
	if tl.err != nil {
		logField["error"] = tl.err
	}
	if tl.code != 0 {
		logField["code"] = tl.code
	}
	if tl.response != nil {
		logField["response"] = tl.response
	}
	if tl.extraInfo != nil {
		logField["extraInfo"] = tl.extraInfo
	}

	return logField
}

func (tl *TraceLog) PrintInfo(msg string) {
	lf := getlogfields(tl)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.WithFields(lf).Info(msg)
	tl.clear()
}

func (tl *TraceLog) PrintError(msg string, err error) {
	tl.err = err
	lf := getlogfields(tl)
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.WithFields(lf).Error(msg)
	tl.clear()
}

func (tl *TraceLog) PrintWarn(msg string, err error) {
	tl.err = err
	lf := getlogfields(tl)
	logrus.SetLevel(logrus.WarnLevel)
	logrus.WithFields(lf).Warn(msg)
	tl.clear()
}

func (tl *TraceLog) clear() {
	tl.topic = ""
	tl.url = ""
	tl.method = ""
	tl.args = nil
	tl.headers = nil
	tl.domain = ""
	tl.clientIP = ""
	tl.err = nil
	tl.code = 0
	tl.response = nil
	tl.extraInfo = nil
	tl.requestTime = 0
	tl.traceID = ""
}
