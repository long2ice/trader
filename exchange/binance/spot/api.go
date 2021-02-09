package spot

import (
	"encoding/json"
	"errors"
	"github.com/long2ice/trader/exchange"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type Api struct {
	exchange.BaseApi
}

const (
	apiAddr = "https://api.binance.com"
)

//创建或延期现货账户listenKey
func (api *Api) CreateSpotListenKey() (string, bool) {
	url := apiAddr + "/api/v3/userDataStream"
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

func (api *Api) CancelOrder(service *exchange.CancelOrderService) (map[string]interface{}, error) {
	url := apiAddr + "/api/v3/order?"
	var params map[string]interface{}
	marshal, _ := json.Marshal(service)
	err := json.Unmarshal(marshal, &params)
	if err != nil {
		log.WithField("err", err).WithField("marshal", marshal).Error("Unmarshal error")
	}
	query := api.BuildCommonQuery(params)
	var respError map[string]interface{}
	resp, err := api.RestyClient.R().SetError(&respError).Delete(url + query)
	if err != nil {
		log.WithField("err", err).Error("撤单失败")
		return nil, err
	} else if respError != nil {
		log.WithField("respError", respError).Error("撤单失败")
		return nil, errors.New(respError["msg"].(string))
	} else {
		var ret map[string]interface{}
		err := json.Unmarshal(resp.Body(), &ret)
		if err != nil {
			log.WithField("err", err).Error("解析撤单返回数据失败")
			return ret, err
		}
		log.WithField("撤单详情", ret).Info("撤单成功")
		return ret, nil
	}
}
func (api *Api) AccountInfo() ([]exchange.Balance, error) {
	url := apiAddr + "/api/v3/account?"
	query := api.BuildCommonQuery(map[string]interface{}{})
	var respError map[string]interface{}
	var result map[string]interface{}
	_, err := api.RestyClient.R().SetResult(&result).SetError(&respError).Get(url + query)
	if err != nil {
		log.WithField("err", err).Error("获取账号信息失败")
		return nil, err
	} else if respError != nil {
		log.WithField("respError", respError).Error("获取账号信息失败")
		return nil, errors.New(respError["msg"].(string))
	} else {
		var balancesRet []exchange.Balance
		balances, _ := result["balances"]
		for _, balance := range balances.([]interface{}) {
			balance := balance.(map[string]interface{})
			asset, _ := balance["asset"]
			free, _ := balance["free"]
			free_, _ := decimal.NewFromString(free.(string))
			locked, _ := balance["locked"]
			locked_, _ := decimal.NewFromString(locked.(string))
			if free_.Add(locked_).GreaterThan(decimal.NewFromInt(0)) {
				balancesRet = append(balancesRet, exchange.Balance{Asset: asset.(string), Free: free_, Locked: locked_})
			}
		}
		log.WithField("balances", balancesRet).Info("Get account balances success")
		return balancesRet, nil
	}
}
func (api *Api) AddOrder(service *exchange.CreateOrderService) (map[string]interface{}, error) {
	url := apiAddr + "/api/v3/order?"

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
func (api *Api) KLines(service *exchange.KLineService) ([][]interface{}, error) {
	url := apiAddr + "/api/v3/klines"
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
