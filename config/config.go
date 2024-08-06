package config

import (
	"os"
)

type Config struct {
	Env   string
	Locat string

	ServerName string
	ServerPort string
	PwdSalt    string
	WsMd5Salt  string

	MysqlWriteHost     string
	MysqlWritePort     string
	MysqlWriteUser     string
	MysqlWritePassword string
	MysqlReadHost      string
	MysqlReadPort      string
	MysqlReadUser      string
	MysqlReadPassword  string
	MysqlChanelDB      string

	RedisHost     string
	RedisPort     string
	RedisPoolSize string

	BeckhamDomain string
	BeckhamPort   string
	CamilaDomain  string
	CamilaPort    string
}

func NewConfig() *Config {
	return &Config{
		Env:   os.Getenv("ENV"),
		Locat: os.Getenv("LOCAT"),

		ServerName: os.Getenv("SERVER_NAME"),
		ServerPort: os.Getenv("SERVER_PORT"),
		PwdSalt:    os.Getenv("PWD_SALT"),
		WsMd5Salt:  os.Getenv("WS_MD5_SALT"),

		MysqlWriteHost:     os.Getenv("MYSQL_WRITE_HOST"),
		MysqlWritePort:     os.Getenv("MYSQL_WRITE_PORT"),
		MysqlWriteUser:     os.Getenv("MYSQL_WRITE_USER"),
		MysqlWritePassword: os.Getenv("MYSQL_WRITE_PASSWORD"),
		MysqlReadHost:      os.Getenv("Mysql_READ_HOST"),
		MysqlReadPort:      os.Getenv("MYSQL_READ_PORT"),
		MysqlReadUser:      os.Getenv("MYSQL_READ_USER"),
		MysqlReadPassword:  os.Getenv("MYSQL_READ_PASSWORD"),
		MysqlChanelDB:      os.Getenv("MYSQL_SHARECO_DB"),

		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPoolSize: os.Getenv("REDIS_POOlSIZE"),

		BeckhamDomain: os.Getenv("BECKHAM_DOMAIN"),
		BeckhamPort:   os.Getenv("BECKHAM_PORT"),
		CamilaDomain:  os.Getenv("CAMILA_DOMAIN"),
		CamilaPort:    os.Getenv("CAMILA_PORT"),
	}
}
