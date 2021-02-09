# Trader

![backtest](https://raw.githubusercontent.com/long2ice/trader/main/images/kline.png)

## Introduction

`Trader` is a framework that automated cryptocurrency exchange with strategy.

**Important: The profit and loss is yourself when use this framework!**

## Features

- Current support spot of [binance](https://www.binance.com/).
- Easy write your own strategy.
- Back test engine support.

## Quick Start

### Write Strategy

Just see [examples](https://github.com/long2ice/trader/tree/main/examples)

### Back test

`Trader` has a back test engine, for that you can test your strategy.

![backtest](https://raw.githubusercontent.com/long2ice/trader/main/images/backtest.png)

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