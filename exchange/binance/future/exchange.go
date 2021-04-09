package future

import (
	"crypto/tls"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/exchange/binance"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"strings"
	"time"
)

const (
	//wsAddr = "wss://stream.goodusahost.cn"
	//wsAddr = "wss://stream.shyqxxy.com"
	wsAddr       = "wss://fstream.binance.com"
	wsMarketAddr = wsAddr + "/stream"
)

type Future struct {
	Api Api
	binance.Exchange
	//余额信息
	balances []exchange.Balance
}

type CancelOrderService struct {
	binance.CancelOrderService
}
type CreateOrderService struct {
	binance.CreateOrderService
	PositionSide string
}

func (service *CancelOrderService) Collect() map[string]interface{} {
	params := make(map[string]interface{})
	params["symbol"] = service.Symbol
	if service.OrigClientOrderId != nil {
		params["origClientOrderId"] = *service.OrigClientOrderId
	}
	if service.OrderId != nil {
		params["orderId"] = *service.OrderId
	}
	return params
}
func (service *CreateOrderService) Collect() map[string]interface{} {
	params := make(map[string]interface{})
	params["symbol"] = service.Symbol
	params["side"] = service.Side
	params["positionSide"] = service.PositionSide
	params["type"] = service.Type
	if service.Type == db.LIMIT {
		params["timeInForce"] = *service.TimeInForce
		params["price"] = *service.Price
		params["quantity"] = *service.Quantity
	} else if service.Type == db.MARKET {
		params["quantity"] = *service.Quantity
	}
	if service.NewClientOrderId != nil {
		params["newClientOrderId"] = *service.NewClientOrderId
	}
	if service.StopPrice != nil {
		params["stopPrice"] = *service.StopPrice
	}
	if service.NewOrderRespType != nil {
		params["newOrderRespType"] = *service.NewOrderRespType
	}
	return params
}
func init() {
	exchange.RegisterExchange(exchange.BinanceFuture, &Future{})
}

func (future *Future) SubscribeMarketData(symbols []string, callback func(map[string]interface{})) error {
	addr := wsMarketAddr + "?streams=" + strings.Join(symbols, "/")
	wsMarketDataClient, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.WithField("err", err).WithField("symbols", symbols).Error("订阅行情失败")
		return err
	}
	go func() {
		for {
			var message map[string]interface{}
			err := wsMarketDataClient.ReadJSON(&message)
			if err != nil {
				log.WithField("err", err).Warning("解析行情消息错误，重新连接")
				wsMarketDataClient, _, err = websocket.DefaultDialer.Dial(wsMarketAddr, nil)
				if err != nil {
					log.WithField("err", err).Warning("重新连接失败，2秒后重试...")
					time.Sleep(time.Second * 2)
				} else {
					wsMarketDataClient, _, err = websocket.DefaultDialer.Dial(addr, nil)
					if err != nil {
						log.WithField("err", err).WithField("symbols", symbols).Error("重新连接失败")
					} else {
						log.Info("重新连接成功")
					}
				}
				continue
			}
			_, ok := message["error"]
			if ok {
				log.WithField("error", message["error"]).Error()
				continue
			}
			callback(message)
		}
	}()

	go func() {
		for range time.NewTicker(time.Second * 60 * 10).C {
			err := wsMarketDataClient.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(10*time.Second))
			if err != nil {
				log.WithField("err", err).Error("行情pong失败")
			}
		}
	}()
	return nil
}

func (future *Future) SubscribeAccount(callback func(map[string]interface{})) error {
	listenKey, ok := future.Api.CreateSpotListenKey()
	if !ok {
		return errors.New("createSpotListenKey失败")
	}
	wsUrl := wsAddr + "/stream?streams=" + listenKey
	c, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		log.WithField("err", err).WithField("wsUrl", wsUrl).Error("连接websocket失败")
	}
	go func() {
		for {
			var message map[string]interface{}
			err := c.ReadJSON(&message)
			if err != nil {
				log.WithField("err", err).Warning("解析账号消息错误，重新连接")
				c, _, err = websocket.DefaultDialer.Dial(wsUrl, nil)
				if err != nil {
					log.WithField("err", err).Warning("重新连接失败，2秒后重试...")
					time.Sleep(time.Second * 2)
				} else {
					log.Info("重新连接成功")
				}
				continue
			}
			log.WithField("message", message).Debug()
			callback(message)
		}
	}()
	//每30分钟延期SpotListenKey
	go func() {
		for range time.NewTicker(time.Second * 60 * 30).C {
			_, ok := future.Api.CreateSpotListenKey()
			if !ok {
				log.Error("延期SpotListenKey失败")
			}
		}
	}()
	//每10分钟ping
	go func() {
		for range time.NewTicker(time.Second * 60 * 10).C {
			err = c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
			if err != nil {
				log.WithField("err", err).Error("账户信息ping失败")
			}
		}
	}()
	return nil
}
func (future *Future) AddOrder(order db.Order) (map[string]interface{}, error) {
	var positionSide string
	if order.Vol.GreaterThan(decimal.Zero) {
		positionSide = "LONG"
	} else {
		positionSide = "SHORT"
	}
	service := CreateOrderService{
		binance.CreateOrderService{
			Symbol:      order.Symbol,
			Side:        order.Side,
			Type:        order.Type,
			Price:       &order.Price,
			Quantity:    &order.Vol,
			TimeInForce: &order.TimeInForce,
			Api:         &future.Api,
		},
		positionSide,
	}
	return service.Do()
}

func (future *Future) CancelOrder(symbol string, orderId string) (map[string]interface{}, error) {
	service := CancelOrderService{
		binance.CancelOrderService{
			Symbol:  symbol,
			OrderId: &orderId,
			Api:     &future.Api,
		},
	}
	return service.Do()
}
func (future *Future) NewKLineService() exchange.IKLineService {
	var p exchange.IKLineService
	p = &exchange.KLineService{
		Api: &future.Api,
	}
	return p
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
