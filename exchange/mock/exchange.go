package mock

import (
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"strings"
	"time"
)

type Mock struct {
	exchange.BaseExchange
	StartTime time.Time
	EndTime   time.Time
	Symbol    string
}

func init() {
	exchange.RegisterExchange(exchange.Mock, &Mock{})
}
func (mock *Mock) SubscribeMarketData(streams []string, callback func(map[string]interface{})) error {
	var symbols []string
	for _, stream := range streams {
		symbols = append(symbols, strings.ToUpper(strings.Split(stream, "@")[0]))
	}
	var kLines []db.KLine
	db.Client.Where("close_time BETWEEN ? AND ?", mock.StartTime, mock.EndTime).Order("close_time").Where("symbol IN ?", symbols).Find(&kLines)
	for _, kline := range kLines {
		callback(map[string]interface{}{
			"h": kline.High,
			"l": kline.Low,
			"o": kline.Open,
			"c": kline.Close,
			"v": kline.Vol,
			"q": kline.Amount,
			"t": kline.CloseTime,
		})
	}
	return nil
}

func (mock *Mock) NewExchange(apiKey string, apiSecret string) exchange.IExchange {
	return &Mock{}
}
func (mock *Mock) NewKLineService() exchange.IKLineService {
	var p exchange.IKLineService
	p = &KLineService{}
	return p
}
