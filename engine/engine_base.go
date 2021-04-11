package engine

import (
	"github.com/long2ice/trader/conf"
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/strategy"
	"github.com/long2ice/trader/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEngine interface {
	Start(block bool)
	InitConfig(config string)
	SetDb(client *gorm.DB)
	RegisterStrategy(strategy strategy.IStrategy)
	SubscribeMarketData(strategy strategy.IStrategy) error
	SubscribeAccount()
	GetLogger() *log.Entry
}

type Base struct {
	IEngine
	ExchangeType exchange.Type
	Exchange     exchange.IExchange
	Strategies   []strategy.IStrategy
	apiKey       string
	apiSecret    string
}

func (e *Base) GetLogger() *log.Entry {
	return log.WithField("exchange", e.ExchangeType)
}

var engines = make(map[exchange.Type]*IEngine)

func (e *Base) InitConfig(config string) {
	conf.InitConfig(config)
}
func (e *Base) SetDb(client *gorm.DB) {
	db.Init(client)
}
func (e *Base) RegisterStrategy(s strategy.IStrategy) {
	e.Strategies = append(e.Strategies, s)
	e.GetLogger().WithField("symbol", s.GetSymbol()).WithField("strategy", utils.GetTypeName(s)).Info("Register strategy success")
}

func GetEngine(exchangeType exchange.Type, apiKey string, apiSecret string) *IEngine {
	if e, ok := engines[exchangeType]; ok {
		return e
	}
	var e IEngine

	ex, err := exchange.NewExchange(exchangeType, apiKey, apiSecret)
	if err != nil {
		log.WithField("err", err).Error("New exchange failed")
	}
	eb := Base{Exchange: ex, ExchangeType: exchangeType, apiKey: apiKey, apiSecret: apiSecret}
	if exchangeType == exchange.Mock {
		e = &Mock{Base: eb}
	} else {
		e = &Engine{Base: eb}
	}
	engines[exchangeType] = &e
	return &e
}
