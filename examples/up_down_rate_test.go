package main

import (
	"github.com/long2ice/trader/conf"
	"github.com/long2ice/trader/engine"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/exchange/mock"
	"github.com/long2ice/trader/strategy"
	"github.com/shopspring/decimal"
	"strings"
	"testing"
	"time"
)

func TestUpDownRate(t *testing.T) {
	conf.InitConfig("config.yml")
	eng := (*engine.GetEngine(exchange.Mock, conf.BinanceApiKey, conf.BinanceApiSecret)).(*engine.Mock)
	BaseAsset := "ETH"
	QuoteAsset := "USDT"
	symbol := BaseAsset + QuoteAsset
	h, _ := time.ParseDuration("-48h")
	eng.Exchange.(*mock.Mock).StartTime = time.Now().Add(h)
	eng.Exchange.(*mock.Mock).EndTime = time.Now()
	eng.Exchange.(*mock.Mock).Symbol = symbol

	eng.RegisterStrategy(&UpDownRate{
		KLineLimit: 10,
		Rate:       decimal.NewFromInt(6).Div(decimal.NewFromInt(4)),
		Base: strategy.Base{
			Streams:    []string{strings.ToLower(symbol) + "@kline_1m"},
			FundRatio:  decimal.NewFromFloat(1),
			StopProfit: decimal.NewFromFloat(0.02),
			StopLoss:   decimal.NewFromFloat(0.06),
			BaseAsset:  BaseAsset,
			QuoteAsset: QuoteAsset,
			Exchange:   eng.Exchange,
		}},
	)
	eng.Start(false)
}
