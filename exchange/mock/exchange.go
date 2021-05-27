package mock

import (
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/utils"
	"github.com/shopspring/decimal"
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
func (mock *Mock) ParseTicker(data map[string]interface{}) exchange.Ticker {
	e_, _ := data["e"]
	E, _ := data["E"]
	s, _ := data["s"]
	o, _ := data["o"]
	h, _ := data["h"]
	l, _ := data["l"]
	c, _ := data["c"]
	v, _ := data["v"]
	q, _ := data["q"]
	tc, _ := decimal.NewFromString(c.(string))
	tv, _ := decimal.NewFromString(v.(string))
	tq, _ := decimal.NewFromString(q.(string))
	to, _ := decimal.NewFromString(o.(string))
	th, _ := decimal.NewFromString(h.(string))
	tl, _ := decimal.NewFromString(l.(string))
	return exchange.Ticker{
		EventType:    e_.(string),
		Time:         utils.TsToTime(E.(float64)),
		Symbol:       s.(string),
		LatestPrice:  tc,
		First24Price: to,
		High24Price:  th,
		Low24Price:   tl,
		Volume:       tv,
		Amount:       tq,
	}
}

func (mock *Mock) ParseKLine(data map[string]interface{}) exchange.KLine {
	h, _ := data["h"]
	h_, _ := h.(decimal.Decimal)
	l, _ := data["l"]
	l_, _ := l.(decimal.Decimal)
	o, _ := data["o"]
	o_, _ := o.(decimal.Decimal)
	c, _ := data["c"]
	c_, _ := c.(decimal.Decimal)
	v, _ := data["v"]
	v_, _ := v.(decimal.Decimal)
	q, _ := data["q"]
	q_, _ := q.(decimal.Decimal)
	t, _ := data["t"]
	return exchange.KLine{
		Open:      o_,
		Close:     c_,
		High:      h_,
		Low:       l_,
		Amount:    q_,
		Volume:    v_,
		Finish:    true,
		CloseTime: t.(time.Time),
	}
}

func (mock *Mock) SubscribeMarketData(streams []string, callback func(map[string]interface{})) error {
	var symbols []string
	for _, stream := range streams {
		symbols = append(symbols, strings.ToUpper(strings.Split(stream, "@")[0]))
	}
	var kLines []db.KLine
	db.Client.Where("close_time BETWEEN ? AND ?", mock.StartTime, mock.EndTime).Order("close_time").Where("symbol IN ?", symbols).Find(&kLines)
	for _, kline := range kLines {
		data := map[string]interface{}{
			"h": kline.High,
			"l": kline.Low,
			"o": kline.Open,
			"c": kline.Close,
			"v": kline.Vol,
			"q": kline.Amount,
			"t": kline.CloseTime,
		}
		callback(map[string]interface{}{
			"stream": strings.ToLower(kline.Symbol) + "@kline_1m",
			"data":   data,
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
	p.SetStartTime(mock.StartTime.Unix())
	return p
}
