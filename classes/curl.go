package classes

import (
	"bytes"
	"chanel/lib"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	// Protocol
	ProtocolHttp  = "http"
	ProtocolHttps = "https"

	// headers Key
	HeadersConnection  = "Connection"
	HeadersUserAgent   = "User-Agent"
	HeadersContentType = "Content-Type"

	// ContentType Value
	ContentTypeJson          = "application/json"
	ContentTypeJsonUTF8      = "application/json;charset=utf-8"
	ContentTypeFormUrlEncode = "application/x-www-form-urlencoded"
	ContentTypeFormData      = "multipart/form-data"
)

type Curl struct {
	protocol string
	client   *http.Client
	request  *http.Request
	method   string
	path     string
	url      string
	port     string
	timeout  time.Duration
	headers  map[string]string
	cookies  map[string]string
	queries  map[string]interface{}
	body     map[string]interface{}
	traceLog *TraceLog
	tools    *lib.Tools
}

func CurlInit(tools *lib.Tools) *Curl {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("curl error -> %v", lib.PanicParser(err)))
		}
	}()

	return &Curl{
		timeout:  10 * time.Second, // 預設超時
		headers:  make(map[string]string),
		cookies:  make(map[string]string),
		queries:  make(map[string]interface{}),
		body:     make(map[string]interface{}),
		traceLog: TraceLogInit(tools), // 初始化 traceLog
		tools:    tools,
	}
}

// 新增請求
func (r *Curl) NewRequest(domain, port, path string) *Curl {
	r.client = http.DefaultClient
	r.port = port
	r.path = path
	r.url = fmt.Sprintf("%s://%s:%s%s", r.protocol, domain, port, path)

	r.traceLog.SetDomain(domain)
	return r
}

func (r *Curl) SetTraceID(traceID string) *Curl {
	r.traceLog.SetTraceID(traceID)
	return r
}

// 設定請求超時時間
func (r *Curl) SetTimeOut(timeout time.Duration) *Curl {
	if timeout > 0 && timeout < 30*time.Second {
		r.timeout = timeout
	}
	return r
}

// 設定 protocol HTTP 協議
func (r *Curl) SetHttp() *Curl {
	r.protocol = ProtocolHttp
	return r
}

// 設定 protocol HTTPs 協議
func (r *Curl) SetHttps() *Curl {
	r.protocol = ProtocolHttps
	return r
}

// 設定 headers
func (r *Curl) SetHeaders(headers map[string]string) *Curl {
	r.headers = headers
	return r
}

// 設定 cookies
func (r *Curl) SetCookies(cookies map[string]string) *Curl {
	r.cookies = cookies
	return r
}

// 設定 url 查詢參數
func (r *Curl) SetQueries(queries map[string]interface{}) *Curl {
	r.queries = queries
	return r
}

// 設定請求 body 資訊
func (r *Curl) SetBody(bodyData map[string]interface{}) *Curl {
	r.body = bodyData
	return r
}

func (r *Curl) Get() (string, error) {
	return r.setMethod(http.MethodGet).send()
}

func (r *Curl) Post() (string, error) {
	return r.setMethod(http.MethodPost).send()
}

func (r *Curl) Put() (string, error) {
	return r.setMethod(http.MethodPut).send()
}

func (r *Curl) Delete() (string, error) {
	return r.setMethod(http.MethodDelete).send()
}

func (r *Curl) send() (string, error) {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("發送 Curl 請求 error -> %v", lib.PanicParser(err)))
		}
	}()

	var startTime = time.Now()
	var body io.Reader

	if r.tools.InStrArray(r.method, []string{http.MethodPost, http.MethodPut}) && len(r.body) > 0 {
		if contentType, exist := r.headers[HeadersContentType]; exist {
			switch strings.ToLower(contentType) {
			case ContentTypeJson, ContentTypeJsonUTF8:
				if bts, err := json.Marshal(r.body); err != nil {
					return "", err
				} else {
					body = bytes.NewReader(bts)
				}
			case ContentTypeFormUrlEncode:
				formData := url.Values{}

				for k, v := range r.body {
					formData.Add(k, fmt.Sprintf("%v", v))
				}
				body = strings.NewReader(formData.Encode())
			case ContentTypeFormData:
				var b bytes.Buffer
				w := multipart.NewWriter(&b)

				for k, v := range r.body {
					if file, ok := v.(*os.File); ok {
						part, err := w.CreateFormFile(k, file.Name())
						if err != nil {
							return "", err
						}
						_, err = io.Copy(part, file)
						if err != nil {
							return "", err
						}
					} else {
						w.WriteField(k, fmt.Sprintf("%v", v))
					}
				}
				w.Close()
				body = bytes.NewReader(b.Bytes())
				r.headers[HeadersContentType] = w.FormDataContentType()
			}
		} else {
			// 預設 application/json
			if bts, err := json.Marshal(r.body); err != nil {
				return "", err
			} else {
				body = bytes.NewReader(bts)
			}
		}
	}

	if req, err := http.NewRequest(r.method, r.url, body); err != nil {
		return "", err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
		defer cancel()
		r.request = req.WithContext(ctx)
	}
	r.setHeaders().setCookies().setQueries()

	// 寫入 trace log 相關資訊
	r.traceLog.SetTopic("curl")
	r.traceLog.SetUrl(r.path)

	if r.tools.InStrArray(r.method, []string{http.MethodGet, http.MethodDelete}) {
		if r.request.URL.RawQuery != "" {
			r.traceLog.SetUrl(r.path + "?" + r.request.URL.RawQuery)
		} else {
			r.traceLog.SetUrl(r.path)
		}
	}
	r.traceLog.SetMethod(r.method)

	if r.tools.InStrArray(r.method, []string{http.MethodGet, http.MethodDelete}) {
		r.traceLog.SetArgs(r.queries)
	} else {
		r.traceLog.SetArgs(r.body)
	}
	r.traceLog.SetHeaders(r.headers)

	var result []byte
	// 發送
	resp, err := r.client.Do(r.request)

	// 紀錄請求結束時間
	r.traceLog.SetRequestTime(r.tools.GetDownRunTime(startTime))

	if err != nil {
		if resp == nil {
			r.traceLog.PrintError("Response nil", err)
			return string(result), err
		} else {
			r.traceLog.PrintError(resp.Status, err)
			return string(result), err
		}
	}
	// 解析成 json 字串
	result, err = io.ReadAll(resp.Body)

	if err != nil {
		r.traceLog.PrintError("Parser response error", err)
	} else {
		r.traceLog.SetResponse(string(result))
		r.traceLog.PrintInfo("Success")
	}

	return string(result), err
}

func (r *Curl) setMethod(method string) *Curl {
	r.method = method
	return r
}

func (r *Curl) setHeaders() *Curl {
	var foundConnection, foundUserAgent bool

	for k, v := range r.headers {
		r.request.Header.Set(k, v)

		switch k {
		case HeadersConnection:
			foundConnection = false
		case HeadersUserAgent:
			foundUserAgent = false
		}
	}
	// 預設 Connection
	if !foundConnection {
		r.request.Header.Set(HeadersConnection, "close")
		r.headers[HeadersConnection] = "close"
	}

	// 預設 User-Agent
	if !foundUserAgent {
		userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36"
		r.request.Header.Set(HeadersUserAgent, userAgent)
		r.headers[HeadersUserAgent] = userAgent
	}
	return r
}

func (r *Curl) setCookies() *Curl {
	for k, v := range r.cookies {
		r.request.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	return r
}

func (r *Curl) setQueries() *Curl {
	q := r.request.URL.Query()

	for k, v := range r.queries {
		q.Add(k, fmt.Sprintf("%v", v))
	}
	r.request.URL.RawQuery = q.Encode()
	return r
}
