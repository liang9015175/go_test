package controllers

import (
	"awesomeProject/quickstart/models"
	"awesomeProject/quickstart/services"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
func Show(){

	curPage := 1
	pageSize := 10
	for {

		ret, err := services.SelectSchedule(curPage, pageSize, time.Now().Unix()*1000)
		beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】当前页: %d,可能需要批量自留的任务列表长度: %d", curPage, len(ret)))
		curPage++
		if err == nil && len(ret) == 0 {
			break
		}
	label:
		for _, v := range ret {
			beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】开始处理任务ID: %d", v.Id))
			endTime := v.EndTime
			selfCount := v.SelfCount
			intervalCount := v.IntervalCount
			maxAuctionId := v.MaxAuctionId
			maxAuctionStatus := v.MaxAuctionStatus
			maxAuctionUser := v.MaxAuctionUser
			maxSelfId := v.MaxSelfId
			if endTime < time.Now().Unix()*1000 {
				beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID: %d 已经过期", v.Id))
				if maxAuctionId == int64(0) {
					beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID : %d 没有任何竞拍记录，设置默认竞拍人为自拍人:%d,更新为成交", v.Id, v.BuyUserId))
					record(v.BuyUserId, endTime, v.Id, 2)
				} else {
					if maxAuctionStatus == 2 {
						beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID : %d 已经成功有人中标，返回", v.Id))
						continue label
					}
					if v.BuyUserId == maxAuctionUser {
						beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID: %d 存在竞拍记录，最后一次竞拍人为自拍人:%d ，更新为成交", v.Id, maxAuctionUser))
						if maxAuctionStatus != 2 {
							bid := new(models.AuctionOlBid)
							bid.Id = v.MaxAuctionId
							bid.Status = 2
							updateBidRecord(bid)
						}
					} else {
						beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID: %d 存在竞拍记录，最后一次竞拍人为其他人:%d ，插入自拍记录，并成交", v.Id, maxAuctionUser))
						record(v.BuyUserId, endTime, v.Id, 2)
					}
				}

			} else {
				beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID: %d 未过期", v.Id))
				if maxAuctionId == 0 {
					beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID : %d 没有任何竞拍记录，设置默认竞拍人为自拍人:%d", v.Id, v.BuyUserId))
					record(v.BuyUserId, time.Now().Unix()*1000, v.Id, 1)
				} else {
					if maxAuctionStatus == 2 {
						beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID : %d 已经成功有人中标，返回", v.Id))
						continue label
					}
					if v.BuyUserId == maxAuctionUser {
						beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID: %d 存在竞拍记录，最后一次竞拍人为自拍人:%d ，继续等待", v.Id, maxAuctionUser))
						continue label
					} else {
						if maxSelfId == 0 {
							beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID : %d 存在竞拍记录,但自拍人:%d 为首次出价", v.Id, v.BuyUserId))
							record(v.BuyUserId, time.Now().Unix()*1000, v.Id, 1)
						} else {
							beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID: %d 存在竞拍记录，最后一次竞拍人为其他人:%d ，继续出价", v.Id, maxAuctionUser))
							if intervalCount >= selfCount+1 {
								beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID: %d 他人间隔出价次数:%d , 开始第%d次出价", v.Id, intervalCount, selfCount+1))
								record(v.BuyUserId, time.Now().Unix()*1000, v.Id, 1)
							} else {
								beego.Info(fmt.Sprintf("【拍卖管理-定时更新自留竞拍记录】任务ID: %d 他人间隔出价次数:%d ,本次为第%d次出价，继续等待 ", v.Id, intervalCount, selfCount+1))
							}

						}

					}

				}
			}
		}
	}
	return

}

func record(buyUserId int64, bidTime int64, productId int64, status int) {
	auctionOlId := new(models.AuctionOlBid)
	auctionOlId.BidTime = bidTime
	auctionOlId.UserId = buyUserId
	auctionOlId.Status = status
	auctionOlId.AuctionProductId = productId
	auctionOlId.Note = "机器人自动出价"

	id, _ := services.InsertBidRecord(auctionOlId)
	////最新竞拍记录
	model := models.AuctionOlBid{Id: int64(id)}
	services.SelectAuctionBid(&model)
	////竞拍的商品
	product := models.AuctionOlProduct{Id: productId}
	services.SelectAuctionProduct(&product)
	////更新当前价格
	rowAffect, _ := services.UpdateProduct(productId, product.CurrentPrice, model.BidPrice)
	if rowAffect == 0 {
		services.UpdateProduct(productId, product.CurrentPrice, model.BidPrice)
	}

}
func updateBidRecord(bid *models.AuctionOlBid) {
	services.UpdateBid(bid)
}
func init() {
	cron := cron.New()
	spec := "0 1/2 * * * ?"
	cron.AddFunc(spec, Show)
	cron.Start()
}
