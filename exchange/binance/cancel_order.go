package binance

import (
	"github.com/long2ice/trader/exchange"
	"reflect"
)

//订单取消服务
type ICancelOrderService interface {
	SetSymbol(symbol string) ICancelOrderService
	SetOrderId(orderId string) ICancelOrderService
	SetOthers(map[string]interface{}) ICancelOrderService
	Collect() map[string]interface{}
	Do() (map[string]interface{}, error)
}
type CancelOrderService struct {
	Symbol            string
	NewClientOrderId  string
	OrigClientOrderId string
	OrderId           string
	Api               exchange.IApi
}

func (service *CancelOrderService) SetSymbol(symbol string) ICancelOrderService {
	service.Symbol = symbol
	return service
}
func (service *CancelOrderService) SetOrderId(orderId string) ICancelOrderService {
	service.OrderId = orderId
	return service
}
func (service *CancelOrderService) SetOthers(params map[string]interface{}) ICancelOrderService {
	v := reflect.ValueOf(service).Elem()
	for key, value := range params {
		v.FieldByName(key).Set(reflect.ValueOf(value))
	}
	return service
}

func (service *CancelOrderService) Collect() map[string]interface{} {
	params := make(map[string]interface{})
	params["symbol"] = service.Symbol
	params["origClientOrderId"] = service.OrigClientOrderId
	params["orderId"] = service.OrderId
	params["newClientOrderId"] = service.NewClientOrderId
	return params
}

func (service *CancelOrderService) Do() (map[string]interface{}, error) {
	ret, err := service.Api.CancelOrder(service.Collect())
	if err != nil {
		return ret, err
	}
	return ret, nil
}
