package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/olivere/elastic"
	"time"
)

type Item struct {
	Id               int		`json:"id" `
	ItemId           string    `json:"item_id"`            //keyword
	ItemName         string    `json:"item_name"`          //keyword  itemName.ikPinyin=>text
	ItemNo           string    `json:"item_no"`            //keyword
	BarCode          string    `json:"bar_code"`           //keyword
	FirstCategoryId  string    `json:"first_category_id"`  //keyword
	FirstCategory    string    `json:"first_category"`     //keyword
	SecondCategoryId string    `json:"second_category_id"` //keyword
	SecondCategory   string    `json:"second_category"`    //keyword
	ItemUnit         string    `json:"item_unit"`          //keyword
	ItemSize         string    `json:"item_size"`          //keyword
	ImageUrl         string    `json:"image_url"`          //keyword
	StockSize        int       `json:"stock_size"`
	Brand            string    `json:"brand"` //keyword
	CreateTime       time.Time `json:"create_time"`
	ModifyTime       time.Time `json:"modify_time"`
	ImageMiddleUrl   string    `json:"image_middle_url"`
}
type ElasticItem struct {
	Item
	/**
		itemName =>洗发水
		itemName itemUint itemSize =>洗发水瓶装350ML
		itemName itemSize itemUnit=> 洗发水350ML瓶装
		itemUnit itemName itemSize=> 瓶装洗发水350ML
		itemSize itemUnit itemName=> 350ML瓶装洗发水

	    brand itemName itemUint itemSize =>潘婷洗发水瓶装350ML
		brand itemName itemSize itemUnit=> 潘婷洗发水350ML瓶装
		brand itemUnit itemName itemSize=> 潘婷瓶装洗发水350ML
		brand itemSize itemUnit itemName=> 潘婷350ML瓶装洗发水

	*/
	CompletionSnapShoot *elastic.SuggestField `json:"completion_snap_shoot", omitempty` //completion  ik_pinyin_analyzer
	/**
	brand itemName itemUint itemSize
	*/
	TermSnapShoot []string `json:"term_snap_shoot",omitempty` //text ik_pinyin_analyzer
}

func init() {
	orm.RegisterModel(new(Item))
}
