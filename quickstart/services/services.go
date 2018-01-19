package services

import (
	"awesomeProject/quickstart/models"
	"fmt"
	"github.com/astaxie/beego/orm"
)

func SelectSchedule(curPage, pageSize int, now int64) (results []models.AuctionOlSchedule, err error) {
	o := orm.NewOrm()
	r := o.Raw("SELECT temp.*,(SELECT count(id) FROM auction_ol_bid WHERE auction_product_id=temp.id AND id> temp.max_self_id) AS interval_count,(\n"+
		"SELECT user_id FROM `auction_ol_bid` WHERE id=temp.max_auction_id) AS max_auction_user,(\n"+
		"SELECT STATUS FROM `auction_ol_bid` WHERE id=temp.max_auction_id) AS max_auction_status FROM (\n"+
		"SELECT aop.id,aop.`buy_user_id`,aop.`price`,aop.end_time,count(aob.id) AS self_count,(\n"+
		"SELECT max(id) FROM `auction_ol_bid` WHERE `auction_product_id`=aop.id) AS max_auction_id,max(aob.id) AS max_self_id FROM auction_ol_product aop LEFT JOIN `auction_ol_bid` aob ON aop.id=aob.`auction_product_id` AND aop.buy_user_id=aob.user_id WHERE unsaleable=1 AND check_status=1 AND delete_auction=2 AND begin_time<=? GROUP BY aop.id ORDER BY aop.id) AS temp limit ?,?",
		now, (curPage-1)*pageSize, pageSize)
	num, err := r.QueryRows(&results)
	if err != nil {
		panic(err)
	}
	fmt.Println("user nums: ", num)
	return results, err
}
func InsertBidRecord(bid *models.AuctionOlBid) (id int64, err error) {
	o := orm.NewOrm()
	return o.Insert(bid)
}
func SelectAuctionProduct(product *models.AuctionOlProduct) models.AuctionOlProduct {
	o := orm.NewOrm()
	o.Read(product)
	return *product
}
func SelectAuctionBid(bid *models.AuctionOlBid) models.AuctionOlBid {
	o := orm.NewOrm()
	o.Read(bid)
	return *bid
}
func UpdateProduct(id int64, oldPrice float64, currentPrice float64) (int64, error) {
	o := orm.NewOrm()
	setter := o.Raw("UPDATE ecloud_marketing.`auction_ol_product` SET current_price=? WHERE id=? and current_price=?", currentPrice, id, oldPrice)
	ret, _ := setter.Exec()
	return ret.RowsAffected()
}
func UpdateBid(bid *models.AuctionOlBid) {
	o := orm.NewOrm()
	o.Update(bid, "STATUS")
}
