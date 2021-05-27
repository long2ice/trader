package main

import (
	"github.com/long2ice/trader/conf"
	"github.com/long2ice/trader/engine"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/server"
	"github.com/long2ice/trader/strategy"
	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	eng := (*engine.GetEngine(exchange.BinanceSpot, conf.BinanceApiKey, conf.BinanceApiSecret)).(*engine.Engine)
	client, _ := gorm.Open(mysql.Open("mysql://"), &gorm.Config{})
	eng.InitConfig("config.yml")
	eng.SetDb(client)
	s := &UpDownRate{
		KLineLimit: 10,
		Rate:       decimal.NewFromInt(6).Div(decimal.NewFromInt(4)),
		Base: strategy.NewStrategy(
			"ETH",
			"USDT",
			eng.Exchange,
			[]string{"ethusdt@kline_1m", "ethusdt@miniTicker"},
			decimal.NewFromFloat(1),
			decimal.NewFromFloat(0.06),
			decimal.NewFromFloat(0.02)),
	}
	s.RegisterStreamCallback("ethusdt@kline_1m", s.On1mKLine)
	s.RegisterStreamCallback("ethusdt@miniTicker", s.OnTicker)

	eng.RegisterStrategy(s)
	eng.Start(false)
	server.Start()
}
