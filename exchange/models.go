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
	EventType    string
	Time         time.Time
	Symbol       string
	LatestPrice  decimal.Decimal
	First24Price decimal.Decimal
	High24Price  decimal.Decimal
	Low24Price   decimal.Decimal
	Volume       decimal.Decimal
	Amount       decimal.Decimal
}
