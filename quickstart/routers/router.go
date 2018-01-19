package routers

import (
	"awesomeProject/quickstart/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Get("/hello", func(context *context.Context) {
		context.Output.Body([]byte("songliang"))
	})
}
