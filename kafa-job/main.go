package main

import (
	_ "awesomeProject/kafa-job/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	_"github.com/go-sql-driver/mysql"
)

func main() {
	beego.Run()
}
func init()  {
	err:=orm.RegisterDataBase("default", "mysql", "root:@tcp(127.0.0.1:3306)/dd_store?charset=utf8&parseTime=true&loc=Local")
	if err!=nil{
		print(err.Error())
	}
	beego.SetLogger(logs.AdapterFile, `{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	beego.SetLogFuncCall(true)
	orm.Debug = true

}
