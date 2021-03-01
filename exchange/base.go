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

type IApi interface {
	CancelOrder(params map[string]interface{}) (map[string]interface{}, error)
	AddOrder(params map[string]interface{}) (map[string]interface{}, error)
	KLines(params map[string]interface{}) ([][]interface{}, error)
	CreateSpotListenKey() (string, bool)
}

type BaseApi struct {
	IApi
	ApiKey      string
	ApiSecret   string
	RestyClient *resty.Client
}

// 构建必要参数
func (api *BaseApi) BuildCommonQuery(params map[string]interface{}, withSign bool) string {
	var joins []string
	for key, value := range params {
		switch value.(type) {
		case int:
			joins = append(joins, fmt.Sprintf("%s=%d", key, value))
		case string, interface{}, decimal.Decimal:
			joins = append(joins, fmt.Sprintf("%s=%s", key, value))
		}
	}
	if withSign {
		timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		joins = append(joins, fmt.Sprintf("recvWindow=10000&timestamp=%s", timestamp))
	}
	query := strings.Join(joins, "&")

	if withSign {
		h := hmac.New(sha256.New, []byte(api.ApiSecret))
		h.Write([]byte(query))
		signature := hex.EncodeToString(h.Sum(nil))
		query += "&signature=" + signature
	}
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
