package db

type Side string

const (
	BUY  Side = "BUY"
	SELL Side = "SELL"
)

type PriceType string

const (
	LIMIT  PriceType = "LIMIT"
	MARKET PriceType = "MARKET"
)
