package controllers

import (
	"github.com/astaxie/beego"
	"awesomeProject/quickstart/models"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
func (c *MainController) Show() (u models.User)  {
	 o:=orm.NewOrm()
	 category:=new(models.Category)
	 category.Id=5
	 err:=o.Read(category)
	 if err!=nil{
	 	beego.Informational("发生异常,%s",err)
	 }
	 //c.Data["json"]=map[string]interface{}{"name":"sngliang", "age":27, "tel":18680558310}
	 c.Data["json"]=category
	 c.ServeJSON()
	 return
}
