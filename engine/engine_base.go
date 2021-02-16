package engine

import (
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/strategy"
	"github.com/long2ice/trader/utils"
	log "github.com/sirupsen/logrus"
)

type IEngine interface {
	Start(block bool)
	RegisterStrategy(strategy strategy.IStrategy)
	SubscribeMarketData(strategy strategy.IStrategy) error
	SubscribeAccount()
	GetLogger() *log.Entry
}

type engineBase struct {
	IEngine
	ExchangeType exchange.Type
	Exchange     exchange.IExchange
	strategies   []strategy.IStrategy
	apiKey       string
	apiSecret    string
}

func (e *engineBase) GetLogger() *log.Entry {
	return log.WithField("exchange", e.ExchangeType)
}

var engines = make(map[exchange.Type]*IEngine)

func (e *engineBase) RegisterStrategy(s strategy.IStrategy) {
	e.strategies = append(e.strategies, s)
	e.GetLogger().WithField("strategy", utils.GetTypeName(s)).Info("Register strategy success")
}

func GetEngine(exchangeType exchange.Type, apiKey string, apiSecret string) *IEngine {
	if e, ok := engines[exchangeType]; ok {
		return e
	}
	var e IEngine

	ex, err := exchange.NewExchange(exchangeType, apiKey, apiSecret)
	if err != nil {
		e.GetLogger().WithField("err", err).Fatal("New exchange failed")
	}
	eb := engineBase{Exchange: ex, ExchangeType: exchangeType, apiKey: apiKey, apiSecret: apiSecret}
	if exchangeType == exchange.Mock {
		e = &Mock{engineBase: eb}
	} else {
		e = &Engine{engineBase: eb}
	}
	engines[exchangeType] = &e
	return &e
}
