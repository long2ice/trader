package exchange

import (
	"github.com/long2ice/trader/db"
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
	Symbol           string          `json:"symbol"`
	Side             db.Side         `json:"side"`
	Type             db.PriceType    `json:"type"`
	TimeInForce      string          `json:"timeInForce,omitempty"`
	Quantity         decimal.Decimal `json:"quantity,omitempty"`
	Price            decimal.Decimal `json:"price,omitempty"`
	NewClientOrderId string          `json:"newClientOrderId,omitempty"`
	StopPrice        decimal.Decimal `json:"stopPrice,omitempty"`
	NewOrderRespType string          `json:"newOrderRespType,omitempty"`
	Api              IApi
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
	return nil
}
func (service *CreateOrderService) Do() (map[string]interface{}, error) {
	ret, err := service.Api.AddOrder(service.Collect())
	if err != nil {
		return ret, err
	}
	return ret, nil
}
