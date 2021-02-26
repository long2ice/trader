package binance

import (
	"github.com/long2ice/trader/exchange"
)

//订单取消服务
type ICancelOrderService interface {
	Collect() map[string]interface{}
	Do() (map[string]interface{}, error)
}
type CancelOrderService struct {
	Symbol            string
	NewClientOrderId  *string
	OrigClientOrderId *string
	OrderId           *string
	Api               exchange.IApi
}

func (service *CancelOrderService) Collect() map[string]interface{} {
	params := make(map[string]interface{})
	params["symbol"] = service.Symbol
	if service.OrigClientOrderId != nil {
		params["origClientOrderId"] = *service.OrigClientOrderId
	}
	if service.OrderId != nil {
		params["orderId"] = *service.OrderId
	}
	if service.NewClientOrderId != nil {
		params["newClientOrderId"] = *service.NewClientOrderId
	}
	return params
}

func (service *CancelOrderService) Do() (map[string]interface{}, error) {
	ret, err := service.Api.CancelOrder(service.Collect())
	if err != nil {
		return ret, err
	}
	return ret, nil
}
