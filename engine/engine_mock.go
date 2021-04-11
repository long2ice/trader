package engine

import (
	"github.com/long2ice/trader/db"
	_ "github.com/long2ice/trader/exchange/mock"
	"github.com/long2ice/trader/strategy"
	"github.com/long2ice/trader/utils"
	"strings"
)

type Mock struct {
	Base
}

func (e *Mock) Start(block bool) {
	for _, s := range e.Strategies {
		db.Client.Where("strategy = ?", utils.GetTypeName(s)).Where("symbol = ?", s.GetSymbol()).Unscoped().Delete(&db.Order{})
		s.OnConnect()
		err := e.SubscribeMarketData(s)
		if err != nil {
			e.GetLogger().WithField("err", err).Error("Fail to subscribe market data")
		}
	}
	e.GetLogger().Info("Mock finished")
}
func (e *Mock) SubscribeMarketData(strategy strategy.IStrategy) error {
	streams := strategy.GetStreams()
	err := e.Exchange.SubscribeMarketData(streams, func(message map[string]interface{}) {
		stream := message["stream"].(string)
		callbacks := strategy.GetStreamCallback(strings.ToLower(stream))
		for _, callback := range callbacks {
			callback(message["data"].(map[string]interface{}))
		}
	})
	if err != nil {
		e.GetLogger().WithField("err", err).WithField("streams", streams).Error("Failed to subscribe market data")
	}
	return nil
}
