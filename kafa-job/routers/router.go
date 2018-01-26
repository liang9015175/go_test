package routers

import (
	"awesomeProject/kafa-job/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
