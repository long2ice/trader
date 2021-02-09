package db

type Side string

const (
	BUY  Side = "buy"
	SELL Side = "sell"
)

type PriceType string

const (
	Limit  PriceType = "LIMIT"
	Market PriceType = "MARKET"
)
