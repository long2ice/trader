package spot

import (
	"crypto/tls"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"strings"
	"time"
)

const (
	//wsAddr = "wss://stream.goodusahost.cn"
	//wsAddr = "wss://stream.shyqxxy.com"
	wsAddr       = "wss://stream.binance.com:9443"
	wsMarketAddr = wsAddr + "/stream"
)

type Spot struct {
	exchange.BaseExchange
	Api Api
}

func init() {
	exchange.RegisterExchange(exchange.BinanceSpot, &Spot{})
}
func (s *Spot) AddOrder(order db.Order) (map[string]interface{}, error) {
	service := exchange.CreateOrderService{
		Symbol: order.Symbol,
		Side:   order.Side,
		Type:   order.Type,
		Price:  order.Price,
		Api:    &s.Api,
	}
	return s.Api.AddOrder(service.Collect())
}

func (s *Spot) CancelOrder(symbol string, orderId string) (map[string]interface{}, error) {
	service := exchange.CancelOrderService{
		Symbol:  symbol,
		OrderId: orderId,
		Api:     &s.Api,
	}
	return s.Api.CancelOrder(service.Collect())
}
func (s *Spot) NewExchange(apiKey string, apiSecret string) exchange.IExchange {
	b := &Spot{
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
func (s *Spot) RefreshAccount() {
	//初始化账号信息
	balances, err := s.Api.AccountInfo()
	if err != nil {
		log.WithField("err", err).Error("获取账号信息失败")
	} else {
		s.Balances = balances
	}
}
func (s *Spot) NewKLineService() exchange.IKLineService {
	var p exchange.IKLineService
	p = &exchange.KLineService{
		Api: &s.Api,
	}
	return p
}

//订阅账号更新
func (s *Spot) SubscribeAccount(callback func(map[string]interface{})) error {
	listenKey, ok := s.Api.CreateSpotListenKey()
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
			_, ok := s.Api.CreateSpotListenKey()
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

func (s *Spot) SubscribeMarketData(symbols []string, callback func(map[string]interface{})) error {
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
