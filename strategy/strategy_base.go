package strategy

import (
	"fmt"
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/utils"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type IStrategy interface {
	//call when market connected
	OnConnect()
	//kline msg
	On1mKline(kLine exchange.KLine)
	//ticker msg
	OnTicker(ticker exchange.Ticker)
	OnAccount(message map[string]interface{})
	OnOrderUpdate(message map[string]interface{})
	GetStreams() []string
	GetSymbol() string
	GetBaseAsset() string
	GetQuoteAsset() string
	GetFundRatio() decimal.Decimal
	GetFund() decimal.Decimal
	GetStopLoss() decimal.Decimal
	GetStopProfit() decimal.Decimal
	GetLatestPrice() decimal.Decimal
	GetLogger() *log.Entry
	GetAvailableFunds() decimal.Decimal
}

type Base struct {
	IStrategy
	//资产币
	BaseAsset string
	//交易币，通常为USDT
	QuoteAsset string
	//交易所
	Exchange exchange.IExchange
	//需要订阅的行情
	Streams []string
	//使用资金比例
	FundRatio decimal.Decimal
	//资金
	Fund db.Fund
	//止损
	StopLoss decimal.Decimal
	//止盈
	StopProfit decimal.Decimal
	//当前最新价
	LatestPrice decimal.Decimal
}

// 获取交易对
func (strategy *Base) GetSymbol() string {
	return fmt.Sprintf("%s%s", strategy.BaseAsset, strategy.QuoteAsset)
}

//获取行情streams
func (strategy *Base) GetStreams() []string {
	return strategy.Streams
}

//获取可用资金
func (strategy *Base) GetAvailableFunds() decimal.Decimal {
	return strategy.FundRatio.Mul(strategy.Fund.TotalFund)
}

//响应ticker
func (strategy *Base) OnTicker(ticker exchange.Ticker) {
	strategy.LatestPrice = ticker.LatestPrice
}

//响应account
func (strategy *Base) OnAccount(message map[string]interface{}) {
	go strategy.Exchange.RefreshAccount()
}
func (strategy *Base) GetLogger() *log.Entry {
	return log.WithField("strategy", utils.GetTypeName(strategy))
}
func (strategy *Base) GetBaseAsset() string {
	return strategy.BaseAsset
}
func (strategy *Base) GetQuoteAsset() string {
	return strategy.QuoteAsset

}
func (strategy *Base) GetFundRatio() decimal.Decimal {
	return strategy.FundRatio

}
func (strategy *Base) GetFund() decimal.Decimal {
	return strategy.Fund.TotalFund

}
func (strategy *Base) GetStopLoss() decimal.Decimal {
	return strategy.StopLoss
}
func (strategy *Base) GetStopProfit() decimal.Decimal {
	return strategy.StopProfit

}
func (strategy *Base) GetLatestPrice() decimal.Decimal {
	return strategy.LatestPrice

}
