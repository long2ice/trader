package future

import (
	"encoding/json"
	"errors"
	"github.com/long2ice/trader/exchange"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

const (
	apiAddr = "https://fapi.binance.com"
)

type Api struct {
	exchange.BaseApi
}

func (api *Api) AccountInfo() ([]exchange.Balance, error) {
	url := apiAddr + "/fapi/v2/binance?"
	query := api.BuildCommonQuery(map[string]interface{}{})
	var respError map[string]interface{}
	var result []map[string]interface{}
	_, err := api.RestyClient.R().SetResult(&result).SetError(&respError).Get(url + query)
	if err != nil {
		log.WithField("err", err).Error("获取账号信息失败")
		return nil, err
	} else if respError != nil {
		log.WithField("respError", respError).Error("获取账号信息失败")
		return nil, errors.New(respError["msg"].(string))
	} else {
		var balancesRet []exchange.Balance
		for _, balance := range result {
			asset, _ := balance["asset"]
			free, _ := balance["availableBalance"]
			free_, _ := decimal.NewFromString(free.(string))
			all, _ := balance["balance"].(string)
			all_, _ := decimal.NewFromString(all)
			if all_.GreaterThan(decimal.Zero) {
				balancesRet = append(balancesRet, exchange.Balance{Asset: asset.(string), Free: free_, Locked: all_.Sub(free_)})
			}
		}
		log.WithField("balances", balancesRet).Info("Get account balances success")
		return balancesRet, nil
	}
}
func (api *Api) CancelOrder(service *exchange.CancelOrderService) (map[string]interface{}, error) {
	return nil, nil
}
func (api *Api) KLines(service *exchange.KLineService) ([][]interface{}, error) {
	url := apiAddr + "/fapi/v1/klines"
	var params map[string]interface{}
	marshal, _ := json.Marshal(service)
	err := json.Unmarshal(marshal, &params)
	if err != nil {
		log.WithField("err", err).WithField("marshal", marshal).Error("Unmarshal error")
	}
	var respError map[string]interface{}
	query := api.BuildCommonQuery(params)
	resp, err := api.RestyClient.R().SetError(&respError).Post(url + query)
	if err != nil {
		log.WithField("err", err).Error("获取kline失败")
		return nil, err
	} else if respError != nil {
		log.WithField("respError", respError).Error("获取kline失败")
		return nil, errors.New(respError["msg"].(string))
	} else {
		var result [][]interface{}
		err := json.Unmarshal(resp.Body(), &result)
		if err != nil {
			log.WithField("err", err).Error("解析KLine数据失败")
			return result, err
		}
		return result, nil
	}
}
func (api *Api) AddOrder(service *exchange.CreateOrderService) (map[string]interface{}, error) {
	url := apiAddr + "/fapi/v1/order?"
	var params map[string]interface{}
	marshal, _ := json.Marshal(service)
	err := json.Unmarshal(marshal, &params)
	if err != nil {
		log.WithField("err", err).WithField("marshal", marshal).Error("Unmarshal error")
	}
	query := api.BuildCommonQuery(params)
	var respError map[string]interface{}
	resp, err := api.RestyClient.R().SetError(&respError).Post(url + query)
	if err != nil {
		log.WithField("err", err).Error("下单失败")
		return nil, err
	} else if respError != nil {
		log.WithField("respError", respError).Error("下单失败")
		return nil, errors.New(respError["msg"].(string))
	} else {
		var ret map[string]interface{}
		err := json.Unmarshal(resp.Body(), &ret)
		if err != nil {
			log.WithField("err", err).Error("解析下单返回数据失败")
			return ret, err
		}
		log.WithField("订单详情", ret).Info("下单成功")
		return ret, nil
	}
}
