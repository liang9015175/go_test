package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Name string
	Age  int
	Tel  int64
}
type Category struct {
	Id           int64  `json:"id";orm:auto`
	CategoryName string `json:"category_name"`
	status       int
	createAt     time.Time
	createBy     string
	updateAt     time.Time
	updateBy     string
}
type AuctionOlBid struct {
	Id               int64   `json:"id";orm:auto`
	UserId           int64   `json:"user_id"`
	BidPrice         float64 `json:"bid_price"`
	AuctionProductId int64   `json:"auction_product_id"`
	Note             string  `json:"note"`
	Status           int     `json:"status"`
	BidTime          int64   `json:"bid_time"`
	Waiver           int     `json:"waiver"`
}
type AuctionOlSchedule struct {
	Id               int64   `json:"id"`
	BuyUserId        int64   `json:"buy_user_id"`
	Price            float64 `json:"price"`
	EndTime          int64   `json:"end_time"`
	SelfCount        int64   `json:"self_count"`
	MaxAuctionId     int64   `json:"max_auction_id"`
	MaxAuctionUser   int64   `json:"max_auction_user"`
	MaxSelfId        int64   `json:"max_self_id"`
	MaxAuctionStatus int     `json:"max_auction_status"`
	IntervalCount    int64   `json:"interval_count"`
}
type AuctionOlProduct struct {
	Id int64 `orm:auto`

	CurrentPrice float64
}

func init() {
	orm.RegisterModel(new(Category), new(AuctionOlBid), new(AuctionOlSchedule), new(AuctionOlProduct))
}
