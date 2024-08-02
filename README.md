
# Chanel 後端開發框架模組 📝

![Static Badge](https://img.shields.io/badge/License-Golang-red)

![Static Badge](https://img.shields.io/badge/Package-Gin__Swagger-blue)

![Static Badge](https://img.shields.io/badge/Database-Mysql-green)

![Static Badge](https://img.shields.io/badge/Version-1.0-yellow)

---

## 相關指令 🚀
* 環境變數 (依照本地環境調整內容)
```
cp .env.local .env
```
* 啟用服務
```
go run *.go -e .env
```
* 檢查 coding style
```
golangci-lint run --timeout 10m
```
* 產生 Swagger 文件 (需先安裝套件, 往下有教學)
```
swag init
```

---

## 參數命名、請求表頭、回傳結構 📚
### 參數命名
* <font style="color: rgb(200,100,100)">小寫駝峰式</font> 命名法
### 請求表頭
```
Content-Type: application/json    # 統一使用 json（必填）
Sid: 加密密文（string）       # 各專案自定義加密方法（選填, 以API需求為主）
```
### 回傳結構
* <font style="color: rgb(200,100,100)">小寫駝峰式</font> 命名法
```
{
    "code": 0,                    # 狀態碼, 不等於0就代表異常
    "message": "Success",         # 訊息
    "result": {},                 # 資料, 沒資料時不顯示
    "error":nil                   # 錯誤, 沒錯誤時不顯示
}
```
---

## 結構樹與說明 🔥

### 結構樹
```
├── .env.local                    # 本地環境設定檔
├── .gitignore                    # 排除特定檔案不受GIT版控
├── api                           # 外服務 API
│   ├── basic.go                      # 外服務 API（初始化）
│   ├── beckham.go                    # Http 外服務 beckham
│   ├── camila.go                     # gRPC 外服務 camila
│   └── proto                         # gRPC 外服務存放 proto 檔的資料夾
│       ├── camila                    # 每個gRPC 外服務的 proto 編譯結果資料夾
│       │   ├── camila.pb.go
│       │   └── camila_grpc.pb.go
│       └── camila.proto
├── classes                       # 類別
│   ├── curl.go                       # Curl 物件（初始化）
│   └── error.go                      # 錯誤代碼、訊息（初始化）
├── config                        # 設定檔
│   └── config.go                     # 載入環境設定
├── controller                    # 控制器, API 參數處理、檢查
│   ├── basic.go                      # 控制器（初始化）
│   ├── public.go                     # 未歸類的 API 群組
│   └── user.go                       # user API 群組
├── cron                          # 排程
│   ├── basic.go                  # 排程（初始化）
│   └── cron.go                   # 要執行的排程任務
├── database                      # 資料庫
│   ├── mysql.go                      # mysql 連線（初始化）
│   └── redis.go                      # redis 連線（初始化）
├── lib                           # 工具包
│   ├── basic.go                      # 自定義工具 (初始化)
│   ├── crypto.go                     # 加解密處理相關
│   ├── decimal.go                    # 精度處理相關
│   ├── time.go                       # 日期時間處理相關
│   └── type.go                       # 類型處理相關
├── main.go                       # 主程序
├── repository                    # 資料表存儲庫
│   └── user.go                       # user 資料表
├── request                       # 請求資料
│   ├── basic.go                      # 請求資料檢查物件（初始化）
│   ├── public.go                     # 未歸類的 API 請求參數檢查
│   └── user.go                       # user API 請求參數檢查
├── service                       # 商業邏輯, API 處理邏輯、流程
│   ├── basic.go                      # 邏輯（初始化）
│   ├── public.go                     # 未歸類的 API 群組邏輯
│   └── user.go                       # user API 群組邏輯
├── server                        # 主服務
│   ├── basic.go                      # 主服務（初始化）
│   ├── cron.go                       # 排程設定
│   ├── middleware.go                 # 中介層設定
│   └── route.go                      # 路由設定
├── structs                       # 自定義結構
│   ├── api.go                        # 外服務 API 接口請求回傳結構
│   ├── database.go                   # 資料表結構
│   └── server.go                     # 主服務 API 相關的結構
└── tracelog                      # Log
    ├── traceid.go                    # 產生唯一ID
    └── tracelog.go                   # 自定義 log 資訊（初始化）
```

### 說明

* basic.go 檔
    * 依照目錄的分類去存放 "初始化" 的方法

* public.go 檔
    * 依照路由存放 "未歸類" API 的方法
    * 基本上只有 Controller、Service 會有此檔

* Route 命名規則：
    * 不可包含 GET | POST | PUT | DELETE 保留字（直接定義在 Method 就好）

* Controller 命名規則：
    * 檔案：對應 Router 群組名稱
    * 方法：API動作敘述 + API完整名稱

* Service 命名規則：
    * 檔案：對應 Controller 群組名稱
    * 方法：API動作敘述 + API完整名稱 (最好是跟對應的 Controller 方法名稱相同)

* Model 命名規則：
    * Model前綴 + 資料表名稱

* Repository 命名規則：
    * 檔案：對應 Model 名稱
    * 方法：自定義（看得懂就行但盡量詳細）

* Server API 結構 命名規則：
    * 方法：API動作敘述 + API完整名稱 + 類別（最好是跟對應的 Controller 方法名稱相同）

* 外服務 API 結構 命名規則：
    * 檔案：外服務名稱
    * 結構：外服務名稱 + API方法 + API名稱 + 類型（Request、Response）

### 範例

* Route
    * 方法：GET
    * 路徑：/chanel/user/data
* Controller
    * 路徑：/controller/user.go
    * 方法：GetUserData()
* Service
    * 路徑：/service/user.go
    * 方法：GetUserData()
* API 結構
    * 路徑：/structs/server.go
    * 方法：GetUserDataRequest、GetUserDataResponse

---

## Git ⛅️
### Branch 命名

| 用途 | 範例 |
| -------- | -------- |
| 開發新功能 | add/xxx_xxx |
| 調整、修正功能 | modify/xxx_xxx |
| 移除功能 | delete/xxx_xxx |

### Commit 註解

| 用途 | 範例 |
| -------- | -------- |
| 新增功能 | git commit -m "<font style="color: rgb(200,100,100)">Add</font>: 新增什麼功能" |
| 調整功能 | git commit -m "<font style="color: rgb(200,100,100)">Changed</font>: 調整什麼功能" |
| 修正功能 | git commit -m "<font style="color: rgb(200,100,100)">Fixed</font>: 修正什麼功能" |
| 刪除功能 | git commit -m "<font style="color: rgb(200,100,100)">Delete</font>: 刪除什麼功能" |
| 衝突 | git commit -m "衝突處理" |

---

## 本地安裝 Swagger ✨
1. 安裝相關套件 Package
```
go install github.com/swaggo/swag/cmd/swag@latest
```
```
export PATH=$HOME/go/bin:$PATH
```
```
swag init
```
```
go get -u github.com/swaggo/gin-swagger
```
```
go get -u github.com/swaggo/files
```
2. 初始化 gin 的路由 .go 檔案新增以下
```
"github.com/gin-gonic/gin"
swaggerFiles "github.com/swaggo/files"
ginSwagger "github.com/swaggo/gin-swagger"

r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```
3. main.go 需要加上以下這段載入 (chanel要替換成vendor名稱)
```
_ "chanel/docs"
```
---