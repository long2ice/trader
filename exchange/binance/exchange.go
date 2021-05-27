package binance

import (
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/utils"
	"github.com/shopspring/decimal"
	"time"
)

type Exchange struct {
	exchange.BaseExchange
}

func (e *Exchange) ParseTicker(data map[string]interface{}) exchange.Ticker {
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

func (e *Exchange) ParseKLine(data map[string]interface{}) exchange.KLine {
	k, _ := data["k"].(map[string]interface{})
	h, _ := k["h"]
	kh, _ := decimal.NewFromString(h.(string))
	l, _ := k["l"]
	kl, _ := decimal.NewFromString(l.(string))
	o, _ := k["o"]
	ko, _ := decimal.NewFromString(o.(string))
	c, _ := k["c"]
	kc, _ := decimal.NewFromString(c.(string))
	v, _ := k["v"]
	kv, _ := decimal.NewFromString(v.(string))
	q, _ := k["q"]
	kq, _ := decimal.NewFromString(q.(string))
	x, _ := k["x"]
	kx := x.(bool)
	t, _ := k["T"]
	kt := t.(float64)
	return exchange.KLine{
		Open:      ko,
		Close:     kc,
		High:      kh,
		Low:       kl,
		Amount:    kq,
		Volume:    kv,
		Finish:    kx,
		CloseTime: time.Unix(int64(kt/1000), 0),
	}
}
