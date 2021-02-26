package binance

import (
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"github.com/shopspring/decimal"
)

//订单创建服务
type ICreateOrderService interface {
	Collect() map[string]interface{}
	Do() (map[string]interface{}, error)
}

type CreateOrderService struct {
	Symbol           string
	Side             db.Side
	Type             db.PriceType
	TimeInForce      *string
	Quantity         *decimal.Decimal
	Price            *decimal.Decimal
	QuoteOrderQty    *decimal.Decimal
	NewClientOrderId *string
	StopPrice        *decimal.Decimal
	NewOrderRespType *string
	Api              exchange.IApi
}

func (service *CreateOrderService) Collect() map[string]interface{} {
	params := make(map[string]interface{})
	params["symbol"] = service.Symbol
	params["side"] = service.Side
	params["type"] = service.Type
	params["timeInForce"] = "GTC"
	if service.TimeInForce != nil {
		params["timeInForce"] = *service.TimeInForce
	}
	if service.Price != nil {
		params["price"] = *service.Price
	}
	if service.Quantity != nil {
		params["quantity"] = *service.Quantity
	}
	if service.QuoteOrderQty != nil {
		params["quoteOrderQty"] = *service.QuoteOrderQty
	}
	if service.NewClientOrderId != nil {
		params["newClientOrderId"] = *service.NewClientOrderId
	}
	if service.StopPrice != nil {
		params["stopPrice"] = *service.StopPrice
	}
	if service.NewOrderRespType != nil {
		params["newOrderRespType"] = *service.NewOrderRespType
	}

	return params
}
func (service *CreateOrderService) Do() (map[string]interface{}, error) {
	params := service.Collect()
	ret, err := service.Api.AddOrder(params)
	if err != nil {
		return ret, err
	}
	return ret, nil
}
