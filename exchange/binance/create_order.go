package binance

import (
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/exchange"
	"github.com/shopspring/decimal"
	"reflect"
)

//订单创建服务
type ICreateOrderService interface {
	SetSymbol(symbol string) ICreateOrderService
	SetPrice(price decimal.Decimal) ICreateOrderService
	SetVol(vol decimal.Decimal) ICreateOrderService
	SetSide(side db.Side) ICreateOrderService
	SetType(type_ db.PriceType) ICreateOrderService
	SetOthers(map[string]interface{}) ICreateOrderService
	Collect() map[string]interface{}
	Do() (map[string]interface{}, error)
}

type CreateOrderService struct {
	Symbol           string
	Side             db.Side
	Type             db.PriceType
	TimeInForce      string
	Quantity         decimal.Decimal
	Price            decimal.Decimal
	QuoteOrderQty    decimal.Decimal
	NewClientOrderId string
	StopPrice        decimal.Decimal
	NewOrderRespType string
	Api              exchange.IApi
}

func (service *CreateOrderService) SetSymbol(symbol string) ICreateOrderService {
	service.Symbol = symbol
	return service
}
func (service *CreateOrderService) SetPrice(price decimal.Decimal) ICreateOrderService {
	service.Price = price
	return service
}
func (service *CreateOrderService) SetSide(side db.Side) ICreateOrderService {
	service.Side = side
	return service
}
func (service *CreateOrderService) SetVol(vol decimal.Decimal) ICreateOrderService {
	service.Quantity = vol
	return service
}

func (service *CreateOrderService) SetType(type_ db.PriceType) ICreateOrderService {
	service.Type = type_
	return service
}
func (service *CreateOrderService) SetOthers(params map[string]interface{}) ICreateOrderService {
	v := reflect.ValueOf(service).Elem()
	for key, value := range params {
		v.FieldByName(key).Set(reflect.ValueOf(value))
	}
	return service
}
func (service *CreateOrderService) Collect() map[string]interface{} {
	params := make(map[string]interface{})
	params["symbol"] = service.Symbol
	params["side"] = service.Side
	params["type"] = service.Type
	params["timeInForce"] = service.TimeInForce
	params["quantity"] = service.Quantity
	params["quoteOrderQty"] = service.QuoteOrderQty
	params["newClientOrderId"] = service.NewClientOrderId
	params["stopPrice"] = service.StopPrice
	params["newOrderRespType"] = service.NewOrderRespType
	return params
}
func (service *CreateOrderService) Do() (map[string]interface{}, error) {
	ret, err := service.Api.AddOrder(service.Collect())
	if err != nil {
		return ret, err
	}
	return ret, nil
}
