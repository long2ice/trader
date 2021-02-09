package main

import (
	"container/list"
	"github.com/long2ice/trader/exchange"
	"github.com/shopspring/decimal"
	"sync"
)

var mutex sync.Mutex

//价格窗口，存储近期价格以及快速得到最高价最低价
type priceWindow struct {
	kLines *list.List
}

func newPriceWindow() *priceWindow {
	return &priceWindow{kLines: list.New()}
}
func (p *priceWindow) Prices() (decimal.Decimal, decimal.Decimal) {
	first := p.kLines.Front().Value.(exchange.KLine)
	minPrice := first.Low
	maxPrice := first.High
	for e := p.kLines.Front(); e != nil; e = e.Next() {
		kLine := e.Value.(exchange.KLine)
		if kLine.High.GreaterThan(maxPrice) {
			maxPrice = kLine.High
		}
		if kLine.Low.LessThan(minPrice) {
			minPrice = kLine.Low
		}
	}
	return minPrice, maxPrice
}
func (p *priceWindow) addKline(kLine exchange.KLine) {
	mutex.Lock()
	p.kLines.PushBack(kLine)
	p.kLines.Remove(p.kLines.Front())
	mutex.Unlock()
}
func (p *priceWindow) addKLines(kLines []exchange.KLine) {
	for _, kline := range kLines {
		p.kLines.PushBack(kline)
	}
}

//获取最近K线升降次数
func (p *priceWindow) getUpDown() (bool, decimal.Decimal, decimal.Decimal) {
	upTimes := decimal.Zero
	downTimes := decimal.Zero
	isDown := false
	for e := p.kLines.Front(); e != nil; e = e.Next() {
		kLine := e.Value.(exchange.KLine)
		if kLine.Open.GreaterThan(kLine.Close) {
			downTimes = downTimes.Add(decimal.NewFromInt(1))
		} else {
			upTimes = upTimes.Add(decimal.NewFromInt(1))
		}
	}
	if p.kLines.Back().Value.(exchange.KLine).Close.LessThan(p.kLines.Front().Value.(exchange.KLine).Close) {
		isDown = true
	}
	return isDown, upTimes, downTimes
}
