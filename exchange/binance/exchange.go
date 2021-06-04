package binance

import (
	"encoding/json"
	"github.com/long2ice/trader/exchange"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"time"
)

type Exchange struct {
	exchange.BaseExchange
}

func (e *Exchange) ParseTicker(data map[string]interface{}) exchange.Ticker {
	dbByte, _ := json.Marshal(data)
	var ticker exchange.Ticker
	err := json.Unmarshal(dbByte, &ticker)
	if err != nil {
		log.WithField("err", err).Error("ParseTicker failed")
	}
	return ticker
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
