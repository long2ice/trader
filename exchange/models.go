package exchange

import (
	"github.com/shopspring/decimal"
	"time"
)

type Balance struct {
	Asset  string
	Free   decimal.Decimal
	Locked decimal.Decimal
}
type KLine struct {
	Open      decimal.Decimal
	Close     decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Amount    decimal.Decimal
	Volume    decimal.Decimal
	Finish    bool
	CloseTime time.Time
}
type Ticker struct {
	EventType             string          `json:"e"` // 事件类型
	Time                  float64         `json:"E"` // 事件时间
	Symbol                string          `json:"s"` // 交易对
	Change24Price         decimal.Decimal `json:"p"` // 24小时价格变化
	Change24Percent       decimal.Decimal `json:"P"` // 24小时价格变化(百分比)
	AvgPrice              decimal.Decimal `json:"w"` // 平均价格
	First24PricePre       decimal.Decimal `json:"x"` // 整整24小时之前，向前数的最后一次成交价格
	LatestPrice           decimal.Decimal `json:"c"` // 最新成交价格
	LatestVol             decimal.Decimal `json:"Q"` // 最新成交交易的成交量
	LatestHighPriceBuy    decimal.Decimal `json:"b"` // 目前最高买单价
	LatestHighPriceBuyVol decimal.Decimal `json:"B"` // 目前最高买单价的挂单量
	LatestLowPriceSell    decimal.Decimal `json:"a"` // 目前最低卖单价
	LatestLowPriceSellVol decimal.Decimal `json:"A"` // 目前最低卖单价的挂单量
	First24PriceAft       decimal.Decimal `json:"o"` // 整整24小时前，向后数的第一次成交价格
	High24Price           decimal.Decimal `json:"h"` // 24小时内最高成交价
	Low24Price            decimal.Decimal `json:"l"` // 24小时内最低成交价
	Vol                   decimal.Decimal `json:"v"` // 24小时内成交量
	Amount                decimal.Decimal `json:"q"` // 24小时内成交额
	StartTime             float64         `json:"O"` // 统计开始时间
	EndTime               float64         `json:"C"` // 统计结束时间
	FirstOrderId          int             `json:"F"` // 24小时内第一笔成交交易ID
	EndOrderId            int             `json:"L"` // 24小时内最后一笔成交交易ID
	OrderNum              int             `json:"n"` // 24小时内成交数
}
