package spot

import (
	"github.com/long2ice/trader/conf"
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBinanceExchange_RefreshAccount(t *testing.T) {
	Assert := assert.New(t)
	ex, err := exchange.NewExchange(exchange.BinanceSpot, conf.BinanceApiKey, conf.BinanceApiSecret)
	if err != nil {
		log.WithField("err", err).Fatal("创建变交易所失败")
	}
	ex.RefreshAccount()

	btc := ex.GetBalance("BTC")
	usdt := ex.GetBalance("USDT")

	log.WithField("btc", btc).WithField("usdt", usdt).Info()

	Assert.NotEqual(btc, decimal.NewFromInt(0))
	Assert.NotEqual(usdt, decimal.NewFromInt(0))

}
func TestBinanceExchange_AddOrder(t *testing.T) {
	ex, err := exchange.NewExchange(exchange.BinanceSpot, conf.BinanceApiKey, conf.BinanceApiSecret)
	if err != nil {
		log.WithField("err", err).Fatal("创建交易所失败")
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
