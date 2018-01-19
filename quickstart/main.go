package main

import (
	_ "awesomeProject/quickstart/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	beego.Run()
}
func init() {
	orm.RegisterDataBase("default", "mysql", "root:iloveyou2016,./@tcp(192.168.0.252:3306)/ecloud_marketing?charset=utf8&parseTime=true&loc=Local")
	beego.SetLogger(logs.AdapterFile, `{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	beego.SetLogFuncCall(true)
	orm.Debug = true

}
