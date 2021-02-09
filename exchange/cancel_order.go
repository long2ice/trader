package exchange

import (
	"reflect"
)

//订单取消服务
type ICancelOrderService interface {
	SetSymbol(symbol string) ICancelOrderService
	SetOrderId(orderId string) ICancelOrderService
	SetOthers(map[string]interface{}) ICancelOrderService
	Do() (map[string]interface{}, error)
}
type CancelOrderService struct {
	Symbol            string `json:"Symbol"`
	NewClientOrderId  string `json:"newClientOrderId,omitempty"`
	OrigClientOrderId string `json:"origClientOrderId,omitempty"`
	OrderId           string `json:"orderId,omitempty"`
	Api               IApi
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
func (service *CancelOrderService) Do() (map[string]interface{}, error) {
	ret, err := service.Api.CancelOrder(service)
	if err != nil {
		return ret, err
	}
	return ret, nil
}
