package services

import (
	"awesomeProject/kafa-job/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"os"
	"time"
	"github.com/olivere/elastic"
	"encoding/json"
)

const (
	FILENAME   = "../sequence.log"
	DATEFORMAT = "2006-01-02 15:04:05 +0800 CST"
)

func getItem(curPage, pageSize int, startDate time.Time, endDate time.Time) (result []models.Item) {
	o := orm.NewOrm()

	setter := o.Raw("select * from t_cm_item where modify_time between ? and ? limit ?,?", startDate, endDate, (curPage-1)*pageSize, pageSize)
	_, _ = setter.QueryRows(&result)
	return result
}
func ScheduleImport() {
	defaultDate, _ := time.ParseInLocation(DATEFORMAT, "2016-07-01 00:00:00 +0800 CST", time.Local)
	beego.Info(fmt.Sprintf("defaultDate: %s", defaultDate))
	var lastDate time.Time
	var _ *os.File
	_, err := os.Stat(FILENAME)
	if err != nil {
		if os.IsNotExist(err) {
			//文件不存在
			_, err = os.Create(FILENAME)
			lastDate = defaultDate
		}
	} else {
		b, err := ioutil.ReadFile(FILENAME)
		if err != nil {
			panic(err)
		}
		if b == nil || len(b) == 0 {
			lastDate = defaultDate
		} else {
			beego.Info(fmt.Sprintf("bstring: %s", string(b)))
			fileDate, _ := time.ParseInLocation(DATEFORMAT, string(b), time.Local)
			beego.Info(fmt.Sprintf("fileDate: %s", fileDate))
			lastDate = fileDate
		}

	}
	pageSize := 10
	now := time.Now().Local()
	beego.Info(fmt.Sprintf("now: %s", now))
	for {
		if lastDate.After(now) {
			beego.Info("上一次同步的时间在当前之间之后，退出")
			break
		}
		curPage := 1
		endDate := lastDate.AddDate(0, 0, 1)
		beego.Info(fmt.Sprintf("开始同步商品: %s--%s", lastDate, endDate))
		for {
			result := getItem(curPage, pageSize, lastDate, endDate)
			offset := (curPage - 1) * pageSize

			beego.Info(fmt.Sprintf("开始同步商品: %s--%s ,从第  %d 条---到第 %d 条", lastDate, endDate, offset, offset+pageSize))
			if result == nil || len(result) == 0 {
				break
			}
			elasticParams(&result)
			time.Sleep(1*time.Second)
			curPage++

		}
		lastDate = endDate
	}
	ioutil.WriteFile(FILENAME, []byte(now.Format(DATEFORMAT)), 0644)

}
func elasticParams(results *[]models.Item)  {
	for _, v := range *results {

		completion:=&elastic.SuggestField{}
		completion.Input(
			v.ItemName,
			v.ItemName+v.ItemUnit+v.ItemSize,
			v.ItemName+v.ItemSize+v.ItemUnit,
			v.ItemUnit+v.ItemSize+v.ItemName,
			v.ItemSize+v.ItemUnit+v.ItemName,
			v.Brand+v.ItemName+v.ItemUnit+v.ItemSize,
			v.Brand+v.ItemName+v.ItemSize+v.ItemUnit,
			v.Brand+v.ItemUnit+v.ItemSize+v.ItemName,
			v.Brand+v.ItemSize+v.ItemUnit+v.ItemName)
		data,_:=json.Marshal(completion)
		beego.Info(fmt.Sprintf("封装过后的值：%s",string(data)))
		tags:=[]string{v.Brand,v.ItemName,v.ItemUnit,v.ItemSize,v.FirstCategory,v.SecondCategory}

		v:=models.ElasticItem{
			Item: v,
			CompletionSnapShoot: completion,
			TermSnapShoot:tags}

		asyncProducer(&v)
	}
}