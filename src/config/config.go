/*
Package config default config
*/
package config

import (
	"io/ioutil"
	"log"
	"os"

	jsoniter "github.com/json-iterator/go"
)

// Config 配置
type Config struct {
	AppInfo appInfo `json:"AppInfo"`
	Log     logConf `json:"Log"`
	DB      db      `json:"DB"`
	Redis   redis   `json:"Redis"`
	Wechat  wechat  `json:"Wechat"`
}

type appInfo struct {
	Env  string `json:"Env"` // example: local, dev, prod
	Addr string `json:"Addr"`
}

type logConf struct {
	LogBasePath string `json:"LogBasePath"`
	LogFileName string `json:"LogFileName"`
}

type wechat struct {
	AppID          string `json:"AppID"`
	AppSecret      string `json:"AppSecret"`
	Token          string `json:"Token"`
	EncodingAESKey string `json:"EncodingAESKey"`
}

type db struct {
	DriverName  string `json:"DriverName"`
	Host        string `json:"Host"`
	Port        string `json:"Port"`
	DBName      string `json:"DBName"`
	User        string `json:"User"`
	PW          string `json:"PW"`
	AdminDBName string `json:"AdminDBName"`
}

type redis struct {
	Host string `json:"Host"`
	Port string `json:"Port"`
	PW   string `json:"PW"`
}

// Conf 配置
var Conf *Config

var filePrefix = "/app/config/"

func init() {
	log.Println("begin init all configs")
	initConf()
	log.Println("over init all configs")
}

func initConf() {
	log.Println("begin init default config")

	Conf = &Config{}
	fileName := "default.json"

	if v, ok := os.LookupEnv("CONFIG_PATH_PREFIX"); ok {
		filePrefix = v
	}
	// read default config
	data, err := ioutil.ReadFile(filePrefix + fileName)
	if err != nil {
		log.Println("config-initConf: read default.json error")
		log.Panic(err)
		return
	}
	err = jsoniter.Unmarshal(data, Conf)
	if err != nil {
		log.Println("config-initConf: unmarshal default.json error")
		log.Panic(err)
		return
	}

	// read env and config path
	if v, ok := os.LookupEnv("ENV"); ok {
		fileName = v + ".json"
	}
	if fileName != "default.json" {
		// read env config
		data, err = ioutil.ReadFile(filePrefix + fileName)
		if err != nil {
			log.Println("config-initConf: read [env].json error")
			log.Panic(err)
			return
		}
		err = jsoniter.Unmarshal(data, Conf)
		if err != nil {
			log.Println("config-initConf: unmarshal [env].json error")
			log.Panic(err)
			return
		}
	}

	if v, ok := os.LookupEnv("WeixinAppID"); ok {
		Conf.Wechat.AppID = v
	}
	if v, ok := os.LookupEnv("WeixinAppSecret"); ok {
		Conf.Wechat.AppSecret = v
	}
	if v, ok := os.LookupEnv("WeixinToken"); ok {
		Conf.Wechat.Token = v
	}
	if v, ok := os.LookupEnv("WeixinEncodingAESKey"); ok {
		Conf.Wechat.EncodingAESKey = v
	}

	if v, ok := os.LookupEnv("MONGO_INITDB_ROOT_USERNAME"); ok {
		Conf.DB.User = v
	}
	if v, ok := os.LookupEnv("MONGO_INITDB_ROOT_PASSWORD"); ok {
		Conf.DB.PW = v
	}
	if v, ok := os.LookupEnv("MONGO_INITDB_DATABASE"); ok {
		Conf.DB.DBName = v
	}
	if v, ok := os.LookupEnv("RedisPass"); ok {
		Conf.Redis.PW = v
	}

	log.Println("over init default config")
}
