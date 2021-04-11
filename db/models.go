package db

import (
	"github.com/shopspring/decimal"
	"time"
)

// Order 订单记录
type Order struct {
	ID          uint
	OrderId     string
	Side        Side
	Vol         decimal.Decimal
	Price       decimal.Decimal
	Amount      decimal.Decimal
	Symbol      string
	TimeInForce string
	Type        PriceType
	Timestamp   time.Time
	Strategy    string
	CurrentFund decimal.Decimal
	Commission  decimal.Decimal
}

// Fund 资金记录
type Fund struct {
	ID        uint
	TotalFund decimal.Decimal
	Strategy  string `gorm:"uniqueindex;type:varchar(50);"`
}

// KLine 分钟K线
type KLine struct {
	ID        int64
	Symbol    string `gorm:"index:idx_symbol_time"`
	OpenTime  time.Time
	CloseTime time.Time `gorm:"index:idx_symbol_time"`
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
