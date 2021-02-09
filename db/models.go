package db

import (
	"github.com/shopspring/decimal"
	"time"
)

//订单记录
type Order struct {
	ID          uint
	OrderId     string
	Side        Side `gorm:"type:enum('buy', 'sell');"`
	Vol         decimal.Decimal
	Price       decimal.Decimal
	Symbol      string
	Type        PriceType `gorm:"type:enum('LIMIT', 'MARKET');"`
	Timestamp   time.Time
	Strategy    string
	CurrentFund decimal.Decimal
	Commission  decimal.Decimal
}

//资金记录
type Fund struct {
	ID        uint
	TotalFund decimal.Decimal
	Strategy  string `gorm:"uniqueindex;type:varchar(50);"`
}

//分钟K线
type KLine struct {
	ID        int64
	Symbol    string
	OpenTime  time.Time
	CloseTime time.Time
	Open      decimal.Decimal
	Close     decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Vol       decimal.Decimal
	Amount    decimal.Decimal
	Num       decimal.Decimal
	BuyVolume decimal.Decimal
	BuyAmount decimal.Decimal
}
