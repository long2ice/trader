package engine

import (
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	_ "github.com/long2ice/trader/exchange/binance/future"
	_ "github.com/long2ice/trader/exchange/binance/spot"
	"github.com/long2ice/trader/strategy"
	"github.com/long2ice/trader/utils"
	"github.com/shopspring/decimal"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Engine struct {
	engineBase
}

func (e *Engine) Start(block bool) {
	db.Init()
	e.SubscribeAccount()
	for _, s := range e.strategies {
		//订阅行情
		s.OnConnect()
		err := e.SubscribeMarketData(s)
		if err != nil {
			e.GetLogger().WithField("err", err).Fatal("Subscribe market data fail")
		}
	}
	e.GetLogger().Info("Start engine success")
	if block {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, os.Kill)
		<-interrupt
	}
}
func (e *Engine) SubscribeMarketData(strategy strategy.IStrategy) error {
	streams := strategy.GetStreams()
	err := e.Exchange.SubscribeMarketData(streams, func(message map[string]interface{}) {
		if stream, ok := message["stream"]; ok {
			stream_ := stream.(string)
			if utils.Contains(streams, stream_) {
				stream_ := strings.ToLower(stream_)
				data := message["data"].(map[string]interface{})
				if strings.Contains(stream_, "ticker") { //ticker
					c, _ := data["c"]
					v, _ := data["v"]
					q, _ := data["q"]
					tc, _ := decimal.NewFromString(c.(string))
					tv, _ := decimal.NewFromString(v.(string))
					tq, _ := decimal.NewFromString(q.(string))
					ticker := exchange.Ticker{
						LatestPrice: tc,
						Volume:      tv,
						Amount:      tq,
					}
					go strategy.OnTicker(ticker)
				} else if strings.Contains(stream_, "@kline_1m") { //kline
					k, _ := data["k"].(map[string]interface{})
					h, _ := k["h"]
					kh, _ := decimal.NewFromString(h.(string))
					l, _ := k["l"]
					kl, _ := decimal.NewFromString(l.(string))
					o, _ := k["o"]
					ko, _ := decimal.NewFromString(o.(string))
					c, _ := k["c"]
					kc, _ := decimal.NewFromString(c.(string))
					v, _ := k["v"]
					kv, _ := decimal.NewFromString(v.(string))
					q, _ := k["q"]
					kq, _ := decimal.NewFromString(q.(string))
					x, _ := k["x"]
					kx := x.(bool)
					t, _ := k["T"]
					kt := t.(float64)
					kLine := exchange.KLine{
						Open:      ko,
						Close:     kc,
						High:      kh,
						Low:       kl,
						Amount:    kq,
						Volume:    kv,
						Finish:    kx,
						CloseTime: time.Unix(int64(kt/1000), 0),
					}
					go strategy.On1mKline(kLine)
				}
			}
		}
	})
	if err != nil {
		e.GetLogger().WithField("err", err).WithField("streams", streams).Error("Failed to subscribe market data")
	} else {
		e.GetLogger().WithField("streams", streams).Info("Subscribe market data success")
	}
	return nil
}
func (e *Engine) SubscribeAccount() {
	err := e.Exchange.SubscribeAccount(func(message map[string]interface{}) {
		e.GetLogger().WithField("message", message).Info("Account data")
		eventType, _ := message["e"].(string)
		for _, s := range e.strategies {
			switch eventType {
			case "outboundAccountPosition": //账户更新
				go s.OnAccount(message)
			case "executionReport": //订单更新
				go s.OnOrderUpdate(message)
			}
		}
	})
	if err != nil {
		e.GetLogger().WithField("err", err).Error("Failed to subscribe account")
	} else {
		e.GetLogger().Info("Subscribe account success")
	}
}
