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

func (api *Api) CreateSpotListenKey() (string, bool) {
	url := apiAddr + "/fapi/v1/listenKey"
	var result map[string]interface{}
	var respError map[string]interface{}
	_, err := api.RestyClient.R().SetResult(&result).SetError(&respError).Post(url)
	if err != nil || respError != nil {
		log.WithField("respError", respError).WithField("err", err).Error("createSpotListenKey error")
		return "", false
	} else {
		listenKey := result["listenKey"]
		return listenKey.(string), true
	}
}
func (api *Api) AccountInfo() ([]exchange.Balance, error) {
	url := apiAddr + "/fapi/v2/binance?"
	query := api.BuildCommonQuery(map[string]interface{}{}, true)
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
func (api *Api) CancelOrder(params map[string]interface{}) (map[string]interface{}, error) {
	panic("not implemented")
}
func (api *Api) KLines(params map[string]interface{}) ([][]interface{}, error) {
	url := apiAddr + "/fapi/v1/klines"
	var respError map[string]interface{}
	query := api.BuildCommonQuery(params, false)
	resp, err := api.RestyClient.R().SetError(&respError).Post(url + query)
	if err != nil {
		log.WithField("err", err).Error("获取KLine失败")
		return nil, err
	} else if respError != nil {
		log.WithField("respError", respError).Error("获取KLine失败")
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
func (api *Api) AddOrder(params map[string]interface{}) (map[string]interface{}, error) {
	url := apiAddr + "/fapi/v1/order?"
	query := api.BuildCommonQuery(params, true)
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
