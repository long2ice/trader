# Trader

![kline](https://raw.githubusercontent.com/long2ice/trader/master/images/kline.png)

## Introduction

`Trader` is a framework that automated cryptocurrency exchange with strategy.

**Disclaimer: The profit and loss is yourself when use this framework!**

## Features

- Current support spot/future of [binance](https://www.binance.com/), more exchange work in progress.
- Easy write your own strategy.
- Back test engine support to test your strategy.

## Quick Start

### Write Strategy

Just see [examples](https://github.com/long2ice/trader/tree/main/examples)

#### Inherit `strategy.Base`

```go
package strategy

type UpDownRate struct {
	strategy.Base
}
```

#### Implement `OnConnect` and `On1mKline`, or `OnTicker`.

```go
package strategy

func (s *UpDownRate) OnConnect() {
	// you can do something initialization here.
}

func (s *UpDownRate) On1mKline() {
	// here you can buy or sell order depending on 1m kline.
}
```

#### Back test strategy

The result will store in `orders` table in database.

```go
package tests

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
```

#### Run strategy

After test strategy, and it's effective, that's the time to run it.

```go
package main

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
	eng.Start(true)
}
```

### Back test

`Trader` has a back test engine, for that you can test your strategy.

![backtest](https://raw.githubusercontent.com/long2ice/trader/master/images/backtest.png?raw=true)

### Config

Copy `config.example.yaml` to `config.yaml` then update it, and current use `MySQL`.

You can get other custom config by `viper.GetString()` etc, see [viper](https://github.com/spf13/viper) framework for
more details.

### Run trader

Also see [examples](https://github.com/long2ice/trader/tree/main/examples).

```shell
INFO[0001] Get account balances success                  balances="[{BTC 0.00000092 0} {ETH 0.0000053 0} {USDT 111.51295105 0} {BUSD 0.00914978 0}]"
INFO[0001] Register strategy success                     strategy=UpDownRate
INFO[0002] Subscribe account success                    
INFO[0005] Subscribe market data success                 streams="[ethusdt@kline_1m ethusdt@miniTicker]"
INFO[0005] Start trader success   
```

You can earn `BTC` and `ETH` when sleep now!

## License

This project is licensed under the
[Apache-2.0](https://github.com/long2ice/trader/blob/master/LICENSE) License.