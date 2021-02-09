package engine

import (
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/exchange/mock"
	_ "github.com/long2ice/trader/exchange/mock"
	"github.com/long2ice/trader/strategy"
	"github.com/long2ice/trader/utils"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

type Mock struct {
	engineBase
}

func (e *Mock) Start() {
	db.Init()
	for _, s := range e.strategies {
		db.Client.Where("strategy = ?", utils.GetTypeName(s)).Where("symbol = ?", s.GetSymbol()).Unscoped().Delete(&db.Order{})
		s.OnConnect()
		err := e.SubscribeMarketData(s)
		if err != nil {
			log.WithField("err", err).Fatal("Fail to start engine")
		}
		ex := e.Exchange.(*mock.Mock)
		err = exec.Command("python3", "analysis/plot.py", utils.GetTypeName(s), s.GetSymbol(), "100", ex.StartTime.Format("2006-01-02")).Start()
		if err != nil {
			log.WithField("err", err).Error("Failed to plot")
		}
	}
	log.Info("Mock finished")
}
func (e *Mock) SubscribeMarketData(strategy strategy.IStrategy) error {
	streams := strategy.GetStreams()
	err := e.Exchange.SubscribeMarketData(streams, func(message map[string]interface{}) {
		h, _ := message["h"]
		l, _ := message["l"]
		o, _ := message["o"]
		c, _ := message["c"]
		v, _ := message["v"]
		q, _ := message["q"]
		t, _ := message["t"]
		kLine := exchange.KLine{
			Open:      o.(decimal.Decimal),
			Close:     c.(decimal.Decimal),
			High:      h.(decimal.Decimal),
			Low:       l.(decimal.Decimal),
			Amount:    q.(decimal.Decimal),
			Volume:    v.(decimal.Decimal),
			Finish:    true,
			CloseTime: t.(time.Time),
		}
		strategy.OnTicker(exchange.Ticker{
			LatestPrice: kLine.Close,
			Volume:      kLine.Volume,
			Amount:      kLine.Amount,
		})
		strategy.On1mKline(kLine)
	})
	if err != nil {
		log.WithField("err", err).WithField("streams", streams).Error("Failed to subscribe market data")
	}
	return nil
}
