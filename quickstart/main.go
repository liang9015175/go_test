package main

import (
	_ "awesomeProject/quickstart/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
)

func main() {

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default","mysql","root:@tcp(127.0.0.1:3306)/liuyanan?charset=utf8")

	beego.Run()
}
func init()  {
	beego.SetLogger(logs.AdapterFile,`{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	beego.SetLogFuncCall(true)
	orm.Debug = true

}
