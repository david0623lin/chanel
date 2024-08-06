package classes

import (
	"chanel/config"
	"chanel/database"
	"chanel/lib"
	"chanel/structs"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	config   *config.Config
	ctx      context.Context
	redis    *database.Redis
	port     string
	heart    time.Duration
	upgrader websocket.Upgrader
	tools    *lib.Tools
	myErr    *MyErr
}

func WebsocketInit(config *config.Config, ctx context.Context, redis *database.Redis) *Websocket {
	return &Websocket{
		config: config,
		ctx:    ctx,
		redis:  redis,
		port:   "8088",
		heart:  1 * time.Minute,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (ws *Websocket) SetTools(tools *lib.Tools) {
	ws.tools = tools
}

func (ws *Websocket) SetError(myErr *MyErr) {
	ws.myErr = myErr
}

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan []byte)

func (ws *Websocket) Start() {
	http.HandleFunc("/ws", ws.connections)
	log.Fatal(http.ListenAndServe(":"+ws.port, nil))
}

// 客端連線
func (ws *Websocket) connections(w http.ResponseWriter, r *http.Request) {
	// 初始化 Trace log
	traceLog := TraceLogInit(ws.tools)
	traceLog.SetTopic("websocket")
	traceLog.SetMethod("connections")
	traceLog.SetUrl("/ws")

	// 建立連線
	wsUpgrade, err := ws.upgrader.Upgrade(w, r, nil)

	if err != nil {
		traceLog.code = WebsocketUpgradeError
		traceLog.PrintError(ws.myErr.Msg(WebsocketUpgradeError), err)
		return
	}
	defer wsUpgrade.Close()

	// 取得 conn url key資料
	rawQueryList := strings.Split(r.URL.RawQuery, "key=")

	if len(rawQueryList) < 2 || rawQueryList[1] == "" {
		traceLog.args = r.URL.RawQuery
		traceLog.code = MissingRequireParams
		traceLog.PrintError(ws.myErr.Msg(MissingRequireParams), err)
		return
	}

	// url 編碼處理
	rawQueryList[1], err = url.QueryUnescape(rawQueryList[1])

	if err != nil {
		traceLog.args = rawQueryList[1]
		traceLog.code = ParseUrlParamsError
		traceLog.PrintError(ws.myErr.Msg(ParseUrlParamsError), err)
		return
	}

	// 解析 wid Info
	var tokenDecrypted []byte
	tokenDecrypted, err = base64.StdEncoding.DecodeString(rawQueryList[1])

	if err != nil {
		traceLog.args = rawQueryList[1]
		traceLog.code = WebsocketParseKeyError
		traceLog.PrintError(ws.myErr.Msg(WebsocketParseKeyError), err)
		return
	}

	// 取得 widInfo 詳細資料
	decrypted := ws.tools.AesDecryptCBC(tokenDecrypted, []byte(ws.tools.Md5(ws.config.WsMd5Salt)))
	var widInfo structs.Websockets
	err = json.Unmarshal([]byte(decrypted), &widInfo)

	if err != nil {
		traceLog.args = string(decrypted)
		traceLog.code = JsonUnmarshalError
		traceLog.PrintError(ws.myErr.Msg(JsonUnmarshalError), err)
		return
	}

	// 驗證 widInfo 內容
	sid, err := ws.redis.Client.Get(ws.ctx, widInfo.Uuid).Result()

	if err != nil {
		traceLog.args = widInfo
		traceLog.code = CacheError
		traceLog.PrintError(ws.myErr.Msg(CacheError), err)
		return
	}

	// 檢查 wid 是否過期
	if sid != widInfo.Sid {
		traceLog.args = widInfo
		traceLog.code = WebsocketKeyExpired
		traceLog.PrintError(ws.myErr.Msg(WebsocketKeyExpired), err)
		return
	}

	// 新增客戶端連線到連線清單
	clients[wsUpgrade] = widInfo.Uuid

	// 處理來自客戶端的訊息
	for {
		_, msg, err := wsUpgrade.ReadMessage()

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, wsUpgrade)
			break
		}
		// 將訊息傳送給所有客戶端
		broadcast <- msg
	}
}

// 處理客戶端訊息的廣播
func (ws *Websocket) Messages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)

			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// 客端定時心跳
func (ws *Websocket) Heart() {
	for {
		broadcast <- []byte("Ping")
		time.Sleep(ws.heart)
	}
}
