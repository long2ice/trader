package future

import (
	"crypto/tls"
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

const (
	//wsAddr = "wss://stream.goodusahost.cn"
	//wsAddr = "wss://stream.shyqxxy.com"
	wsAddr       = "wss://fstream.binance.com"
	wsMarketAddr = wsAddr + "/stream"
)

type Future struct {
	exchange.BaseExchange
	Api Api
	//余额信息
	balances []exchange.Balance
}

type CancelOrderService struct {
	exchange.CancelOrderService
}
type CreateOrderService struct {
	exchange.CreateOrderService
}

func init() {
	exchange.RegisterExchange(exchange.BinanceFuture, &Future{})
}
func (future *Future) SubscribeMarketData(streams []string, callback func(map[string]interface{})) error {
	panic("implement me")
}

func (future *Future) SubscribeAccount(callback func(map[string]interface{})) error {
	panic("implement me")
}
func (future *Future) AddOrder(order db.Order) (map[string]interface{}, error) {
	return future.Api.AddOrder(&exchange.CreateOrderService{
		Symbol: order.Symbol,
		Side:   order.Side,
		Type:   order.Type,
		Price:  order.Price,
		Api:    &future.Api,
	})
}

func (future *Future) CancelOrder(symbol string, orderId string) (map[string]interface{}, error) {
	return future.Api.CancelOrder(&exchange.CancelOrderService{
		Symbol:  symbol,
		OrderId: orderId,
		Api:     &future.Api,
	})
}

func (future *Future) RefreshAccount() {
	//初始化账号信息
	balances, err := future.Api.AccountInfo()
	if err != nil {
		log.WithField("err", err).Error("获取账号信息失败")
	} else {
		future.Balances = balances
	}
}

func (future *Future) NewExchange(apiKey string, apiSecret string) exchange.IExchange {
	b := &Future{
		Api: Api{exchange.BaseApi{
			ApiKey:      apiKey,
			ApiSecret:   apiSecret,
			RestyClient: resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).SetHeader("X-MBX-APIKEY", apiKey),
		}},
	}
	b.RefreshAccount()
	var iExchange exchange.IExchange
	iExchange = b
	return iExchange
}
