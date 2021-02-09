package main

import (
	"github.com/long2ice/trader/conf"
	"github.com/long2ice/trader/engine"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/strategy"
	"github.com/shopspring/decimal"
)

func main() {
	conf.InitConfig("config.yml")
	eng := (*engine.GetEngine(exchange.BinanceSpot, conf.BinanceApiKey, conf.BinanceApiSecret)).(*engine.Engine)
	eng.RegisterStrategy(&UpDownRate{
		KLineLimit: 10,
		Rate:       decimal.NewFromInt(6).Div(decimal.NewFromInt(4)),
		Base: strategy.Base{
			FundRatio:  decimal.NewFromFloat(1),
			StopProfit: decimal.NewFromFloat(0.02),
			StopLoss:   decimal.NewFromFloat(0.06),
			BaseAsset:  "ETH",
			QuoteAsset: "USDT",
			Exchange:   eng.Exchange,
			Streams:    []string{"ethusdt@kline_1m", "ethusdt@miniTicker"},
		}},
	)
	eng.Start()
}
