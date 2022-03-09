package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
)

type Database struct {
	Type string
	User string
	Password string
	Host string
	Name string
	TablePrefix string
}

var DatabaseSetting = &Database{}


var cfg *ini.File

func Setup() {
	var err error
	cfg,err = ini.Load("conf/app.ini")
	if err != nil{
		log.Fatalf("读取配置文件失败")
	}
	Mapto("database",DatabaseSetting)
	fmt.Println("走过来了")
}

func Mapto(section string,v interface{})  {
	err := cfg.Section(section).MapTo(v)
	if err != nil{
		log.Fatalf("cfg.mapto")
	}
}

