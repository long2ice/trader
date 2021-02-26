package future

import (
	"github.com/long2ice/trader/conf"
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	conf.InitConfig("config.yml")
	code := m.Run()
	os.Exit(code)
}
func TestBinanceExchange_AddOrder(t *testing.T) {
	ex, err := exchange.NewExchange(exchange.BinanceFuture, conf.BinanceApiKey, conf.BinanceApiSecret)
	if err != nil {
		log.WithField("err", err).Error("创建交易所失败")
	}
	ret, err := ex.AddOrder(db.Order{
		Side:   "BUY",
		Vol:    decimal.NewFromFloat(0.001),
		Price:  decimal.NewFromInt(20000),
		Symbol: "BTCUSDT",
		Type:   "LIMIT",
	})
	log.WithField("account", ex.GetBalance("USDT")).WithField("ret", ret).Info()
}
