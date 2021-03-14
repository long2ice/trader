package engine

import (
	"github.com/long2ice/trader/db"
	_ "github.com/long2ice/trader/exchange/binance/future"
	_ "github.com/long2ice/trader/exchange/binance/spot"
	"github.com/long2ice/trader/strategy"
	"os"
	"os/signal"
)

type Engine struct {
	Base
}

func (e *Engine) Start(block bool) {
	db.Init()
	e.SubscribeAccount()
	for _, s := range e.Strategies {
		//订阅行情
		s.OnConnect()
		err := e.SubscribeMarketData(s)
		if err != nil {
			e.GetLogger().WithField("err", err).Error("Subscribe market data fail")
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
			callback := strategy.GetStreamCallback(stream_)
			callback(message["data"].(map[string]interface{}))
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
		for _, s := range e.Strategies {
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
