package exchange

import (
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

//KLine服务
type IKLineService interface {
	SetSymbol(symbol string) IKLineService
	SetInterval(interval string) IKLineService
	SetStartTime(startTime int) IKLineService
	SetEndTime(endTime int) IKLineService
	SetLimit(limit int) IKLineService
	Collect() map[string]interface{}
	Do() ([]KLine, error)
}
type KLineService struct {
	IKLineService
	Api       IApi
	Symbol    string
	Interval  string
	StartTime *int
	EndTime   *int
	Limit     *int
}

func (s *KLineService) SetSymbol(symbol string) IKLineService {
	s.Symbol = symbol
	return s
}
func (s *KLineService) SetInterval(interval string) IKLineService {
	s.Interval = interval
	return s
}
func (s *KLineService) SetStartTime(startTime int) IKLineService {
	s.StartTime = &startTime
	return s
}
func (s *KLineService) SetEndTime(endTime int) IKLineService {
	s.EndTime = &endTime
	return s
}
func (s *KLineService) SetLimit(limit int) IKLineService {
	s.Limit = &limit
	return s
}
func (s *KLineService) Collect() map[string]interface{} {
	params := make(map[string]interface{})
	params["symbol"] = s.Symbol
	params["interval"] = s.Interval
	if s.StartTime != nil {
		params["startTime"] = *s.StartTime
	}
	if s.EndTime != nil {
		params["endTime"] = *s.EndTime
	}
	if s.Limit != nil {
		params["limit"] = *s.Limit
	}
	return params
}
func (s *KLineService) Do() ([]KLine, error) {
	result, err := s.Api.KLines(s.Collect())
	if err != nil {
		log.WithField("err", err).Error("Get KLines error")
		return nil, err
	} else {
		var kLines []KLine
		for _, item := range result {
			open, _ := decimal.NewFromString(item[1].(string))
			high, _ := decimal.NewFromString(item[2].(string))
			low, _ := decimal.NewFromString(item[3].(string))
			close_, _ := decimal.NewFromString(item[4].(string))
			volume, _ := decimal.NewFromString(item[5].(string))
			amount, _ := decimal.NewFromString(item[7].(string))
			kLines = append(kLines, KLine{
				Open:   open,
				Close:  close_,
				High:   high,
				Low:    low,
				Amount: amount,
				Volume: volume,
			})
		}
		return kLines, nil
	}
}
