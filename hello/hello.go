package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func main() {
	type User struct {
		UserId   int    `json:"user_id" bson:"user_id"`
		UserName string `json:"user_name" bson:"user_name"`
	}
	// 输出json格式
	u := &User{UserId: 1, UserName: "tony"}
	j, _ := json.Marshal(u)
	fmt.Println(string(j))
	// 输出内容：{"user_id":1,"user_name":"tony"}

	// 获取tag中的内容
	t := reflect.TypeOf(u)
	field := t.Elem().Field(0)
	fmt.Println(field.Tag.Get("json"))
	// 输出：user_id
	fmt.Println(field.Tag.Get("bson"))
	// 输出：user_id
}
