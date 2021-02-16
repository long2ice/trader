package main

import (
	"errors"
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/strategy"
	"github.com/long2ice/trader/utils"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DEBUG = viper.GetBool("common.debug")

/**
判断前若干次交易的升降比，当大于设定值的时候执行买和卖，注意：这只是个简单的策略，并不能盈利
**/
type UpDownRate struct {
	strategy.Base
	Rate         decimal.Decimal //升/降>Rate 卖，降/升>Rate买
	side         db.Side         //BUY or SELL
	KLineLimit   int             //存多少个K线
	lastBuyPrice decimal.Decimal //上次买入价
	priceWindow  *priceWindow
}

func (s *UpDownRate) OnConnect() {
	var order db.Order
	result := db.Client.Where("symbol = ?", s.GetSymbol()).Where("strategy = ?", utils.GetTypeName(s)).Last(&order)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		s.side = db.SELL
	} else {
		if order.Side == db.BUY {
			s.lastBuyPrice = order.Price
		}
		s.side = order.Side
	}
	ex := s.Exchange
	//加载最近价格
	s.priceWindow = newPriceWindow()
	kLines, err := ex.NewKLineService().SetSymbol(s.GetSymbol()).SetInterval("1m").SetLimit(s.KLineLimit).Do()
	if err != nil {
		s.GetLogger().WithField("err", err).WithField("symbol", s.GetSymbol()).Fatal("Get latest kline error")
	} else {
		s.priceWindow.addKLines(kLines)
	}
}
func (s *UpDownRate) On1mKline(kLine exchange.KLine) {
	if !kLine.Finish {
		return
	}
	//更新priceWindow
	s.priceWindow.addKline(kLine)
	isDown, upTimes, downTimes := s.priceWindow.getUpDown()
	if upTimes.Equal(decimal.Zero) {
		return
	}
	var gains decimal.Decimal
	if s.lastBuyPrice.GreaterThan(decimal.Zero) {
		gains = kLine.Close.Sub(s.lastBuyPrice).Div(s.lastBuyPrice)
	} else {
		gains = decimal.Zero
	}
	isGains := gains.Neg().GreaterThan(s.StopLoss)
	//达到止盈或止损
	if (gains.GreaterThanOrEqual(s.StopProfit) || isGains) && s.side == db.BUY {
		//涨幅超过止盈或跌幅达到止损卖出
		s.GetLogger().WithField("symbol", s.GetSymbol()).WithField("涨跌幅", gains).WithField("交易量", kLine.Volume).WithField("当前最新价", s.LatestPrice).Info("达到止盈止损，卖出")
		//执行卖出
		order := db.Order{
			Side:      s.side,
			Symbol:    s.GetSymbol(),
			Timestamp: kLine.CloseTime,
			Strategy:  utils.GetTypeName(s),
		}
		price := kLine.Close
		var orderId string
		vol := s.Exchange.GetBalance(s.BaseAsset).Free.Truncate(5)
		if !DEBUG && vol.GreaterThan(decimal.Zero) {
			ret, err := s.Exchange.AddOrder(order)
			if err != nil {
				s.GetLogger().WithField("err", err).WithField("symbol", s.GetSymbol()).Error("创建卖出订单失败")
				return
			}
			vol, _ = decimal.NewFromString(ret["executedQty"].(string))
			price, _ = decimal.NewFromString(ret["price"].(string))
			orderId_, _ := ret["orderId"].(float64)
			orderId = utils.FloatToString(orderId_)
			s.Exchange.RefreshAccount()
		}
		s.side = db.SELL
		order.Price = price
		order.OrderId = orderId

		//数据库记录
		db.Client.Create(&order)
	}
	//满足策略并且降价幅度大于指定值，或者止损后
	if ((isDown && downTimes.Div(upTimes).GreaterThanOrEqual(s.Rate)) || isGains) && s.side == db.SELL {
		s.GetLogger().WithField("symbol", s.GetSymbol()).WithField("交易量", kLine.Volume).WithField("当前最新价", s.LatestPrice).Info("买入")
		price := kLine.Close
		var orderId string
		vol := decimal.Zero
		free := s.Exchange.GetBalance(s.QuoteAsset).Free.Truncate(8)
		if free.GreaterThanOrEqual(s.Fund.TotalFund) {
			free = free.Mul(s.FundRatio).Truncate(8)
		}
		s.side = db.BUY
		order := db.Order{
			Side:      s.side,
			Vol:       vol,
			Price:     price,
			OrderId:   orderId,
			Symbol:    s.GetSymbol(),
			Timestamp: kLine.CloseTime,
			Strategy:  utils.GetTypeName(s),
		}
		//余额大于10美元执行买入
		if !DEBUG && free.GreaterThan(decimal.NewFromInt(10)) {
			ret, err := s.Exchange.AddOrder(order)
			if err != nil {
				s.GetLogger().WithField("symbol", s.GetSymbol()).WithField("err", err).Error("创建购买订单失败")
				return
			}
			vol, _ := decimal.NewFromString(ret["executedQty"].(string))
			price = free.Div(vol)
			orderId_, _ := ret["orderId"].(float64)
			orderId = utils.FloatToString(orderId_)
			s.Exchange.RefreshAccount()
		}
		s.lastBuyPrice = price
		order.Price = price
		//数据库记录
		db.Client.Create(&order)
	}
}
