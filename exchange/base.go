package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/long2ice/trader/db"
	"github.com/shopspring/decimal"
	"gopkg.in/resty.v1"
	"strconv"
	"strings"
	"time"
)

type Type string

const (
	BinanceSpot   Type = "BinanceSpot"
	BinanceFuture Type = "BinanceFuture"
	Mock          Type = "Mock"
)

type IExchange interface {
	//订阅市场数据
	SubscribeMarketData(streams []string, callback func(map[string]interface{})) error
	//订阅订单更新
	SubscribeAccount(callback func(map[string]interface{})) error
	//近期K线
	NewKLineService() IKLineService
	//下单
	AddOrder(order db.Order) (map[string]interface{}, error)
	//取消订单
	CancelOrder(symbol string, orderId string) (map[string]interface{}, error)
	//刷新账户信息
	RefreshAccount()
	//获取账户余额
	GetBalance(asset string) Balance
	//获取账户所有余额
	GetBalances() []Balance
	//创建exchange
	NewExchange(apiKey string, apiSecret string) IExchange
}

var exchanges = make(map[Type]IExchange)

func RegisterExchange(exchange Type, iExchange IExchange) {
	exchanges[exchange] = iExchange
}

type BaseExchange struct {
	IExchange
	//余额信息
	Balances []Balance
	Api      BaseApi
}

func (exchange *BaseExchange) GetBalance(asset string) Balance {
	for _, b := range exchange.Balances {
		if b.Asset == asset {
			return b
		}
	}
	return Balance{Asset: asset, Free: decimal.NewFromInt(0), Locked: decimal.NewFromInt(0)}
}

func (exchange *BaseExchange) GetBalances() []Balance {
	return exchange.Balances
}
func (exchange *BaseExchange) NewKLineService() IKLineService {
	var p IKLineService
	p = &KLineService{
		Api: exchange.Api,
	}
	return p
}

type IApi interface {
	CancelOrder(service *CancelOrderService) (map[string]interface{}, error)
	AddOrder(service *CreateOrderService) (map[string]interface{}, error)
	KLines(services *KLineService) ([][]interface{}, error)
}

type BaseApi struct {
	IApi
	ApiKey      string
	ApiSecret   string
	RestyClient *resty.Client
}

// 构建必要参数
func (api *BaseApi) BuildCommonQuery(params map[string]interface{}) string {
	var joins []string
	for key, value := range params {
		joins = append(joins, fmt.Sprintf("%s=%s", key, value))
	}
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	joins = append(joins, fmt.Sprintf("recvWindow=10000&timestamp=%s", timestamp))
	query := strings.Join(joins, "&")
	h := hmac.New(sha256.New, []byte(api.ApiSecret))
	h.Write([]byte(query))
	signature := hex.EncodeToString(h.Sum(nil))
	query += "&signature=" + signature
	return query
}
func NewExchange(exchange Type, apiKey string, apiSecret string) (IExchange, error) {
	if iExchange, ok := exchanges[exchange]; ok {
		e := iExchange.NewExchange(apiKey, apiSecret)
		return e, nil
	} else {
		return nil, errors.New("Unknown exchange:" + string(exchange))
	}
}
