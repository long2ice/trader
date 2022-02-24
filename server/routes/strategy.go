package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/long2ice/trader/engine"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/strategy"
	"github.com/long2ice/trader/utils"
)

func getStrategies(ex string) []strategy.IStrategy {
	eng := (*engine.GetEngine(exchange.Type(ex), "", "")).(*engine.Engine)
	return eng.Base.Strategies
}
func GetStrategy(c *gin.Context) {
	symbol := c.Query("symbol")
	strategy_ := c.Query("strategy")
	ex := c.Query("exchange")
	strategies := getStrategies(ex)
	data := make(map[string]interface{})
	for _, s := range strategies {
		if s.GetSymbol() == symbol && utils.GetTypeName(s) == strategy_ {
			data["AvailableFunds"] = s.GetAvailableFunds()
			data["Streams"] = s.GetStreams()
			data["BaseAsset"] = s.GetBaseAsset()
			data["QuoteAsset"] = s.GetQuoteAsset()
			data["FundRatio"] = s.GetFundRatio()
			data["Fund"] = s.GetFund()
			data["StopLoss"] = s.GetStopLoss()
			data["StopProfit"] = s.GetStopProfit()
			data["LatestPrice"] = s.GetLatestPrice()
			c.JSON(200, data)
			break
		}
	}
}

func GetStrategies(c *gin.Context) {
	ex := c.Query("exchange")
	strategies := getStrategies(ex)
	var data []map[string]interface{}
	for _, s := range strategies {
		item := make(map[string]interface{})
		item["AvailableFunds"] = s.GetAvailableFunds()
		item["Streams"] = s.GetStreams()
		item["BaseAsset"] = s.GetBaseAsset()
		item["QuoteAsset"] = s.GetQuoteAsset()
		item["FundRatio"] = s.GetFundRatio()
		item["Fund"] = s.GetFund()
		item["StopLoss"] = s.GetStopLoss()
		item["StopProfit"] = s.GetStopProfit()
		item["LatestPrice"] = s.GetLatestPrice()
		data = append(data, item)
	}
	c.JSON(200, data)
}
