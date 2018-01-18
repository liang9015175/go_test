package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Name string
	Age int
	Tel int64
}
type Category struct {
	Id int64 `json:"id";orm:auto`
	CategoryName string `json:"category_name"`
	status int
	createAt time.Time
	createBy string
	updateAt time.Time
	updateBy string
}

func init()  {
	orm.RegisterModel(new(Category))
}


