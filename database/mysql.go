package database

import (
	"chanel/config"
	"chanel/lib"
	"chanel/repository"
	"chanel/structs"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Mysql struct {
	config   *config.Config
	ChanelDB chanelDB
}

type conn struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

type chanelDB struct {
	WDB        *gorm.DB
	RDB        *gorm.DB
	Repository *repository.Chanel // 資料表
}

func MysqlInit(config *config.Config) *Mysql {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("mysql error -> %v", lib.PanicParser(err)))
		}
	}()

	return &Mysql{
		config: config,
	}
}

func (m *Mysql) Start() {
	// Chanel DB Write
	wDB, err := m.write(conn{
		User:     m.config.MysqlWriteUser,
		Password: m.config.MysqlWritePassword,
		Host:     m.config.MysqlWriteHost,
		Port:     m.config.MysqlWritePort,
		DB:       m.config.MysqlChanelDB,
	})

	if err != nil {
		panic(fmt.Sprintf("mysql start write 錯誤, ERR: %s", err.Error()))
	}

	// Chanel DB Read
	rDB, err := m.read(conn{
		User:     m.config.MysqlReadUser,
		Password: m.config.MysqlReadPassword,
		Host:     m.config.MysqlReadHost,
		Port:     m.config.MysqlReadPort,
		DB:       m.config.MysqlChanelDB,
	})

	if err != nil {
		panic(fmt.Sprintf("mysql start read 錯誤, ERR: %s", err.Error()))
	}

	m.ChanelDB.WDB = wDB
	m.ChanelDB.RDB = rDB
	m.ChanelDB.Repository = repository.ChanelInit(wDB, rDB)
}

func (m *Mysql) write(conn conn) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conn.User,
		conn.Password,
		conn.Host,
		conn.Port,
		conn.DB,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 關閉複數名稱轉換
		},
	})

	if err != nil {
		return db, err
	}

	// 自動產生資料表
	switch conn.DB {
	case "chanel":
		if err = db.AutoMigrate(
			&structs.ChanelModelTasks{},
			&structs.ChanelModelTaskRecords{},
			&structs.ChanelModelCrons{},
			&structs.ChanelModelCronRecords{},
		); err != nil {
			return db, err
		}
		// Todo 由此往下新增資料庫
	}

	// 設定寫入的連線池設定
	sqlDB, err := db.DB()

	if err != nil {
		return db, err
	}
	// 設置閒置連線上限
	sqlDB.SetMaxIdleConns(50)
	// 設置開放連線上限
	sqlDB.SetMaxOpenConns(100)
	// 設置連線的最大閒置時間為 5 分鐘
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}

func (m *Mysql) read(conn conn) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conn.User,
		conn.Password,
		conn.Host,
		conn.Port,
		conn.DB,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 關閉複數名稱轉換
		},
	})

	if err != nil {
		return db, err
	}

	return db, nil
}
