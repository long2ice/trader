package server

import (
	"github.com/gin-gonic/gin"
	"github.com/long2ice/trader/conf"
	"github.com/long2ice/trader/db"
	"github.com/long2ice/trader/engine"
	"github.com/long2ice/trader/exchange"
	"github.com/long2ice/trader/utils"
	"gorm.io/gorm"
)

func getOrders(c *gin.Context) {
	strategy := c.Query("strategy")
	symbol := c.Query("symbol")
	var orders []db.Order
	qs := db.Client
	if strategy != "" {
		qs = qs.Where("strategy = ?", strategy)
	}
	if symbol != "" {
		qs = qs.Where("symbol = ?", symbol)
	}
	qs.Find(&orders)
	c.JSON(200, orders)
}
func getFund(c *gin.Context) {
	strategy := c.Query("strategy")
	var fund db.Fund
	db.Client.Where("strategy = ?", strategy).Find(&fund)
	c.JSON(200, fund)
}
func addFund(c *gin.Context) {
	fund := c.PostForm("fund")
	strategy := c.PostForm("strategy")
	db.Client.Model(&db.Fund{}).Where("strategy = ?", strategy).UpdateColumn("fund", gorm.Expr("fund + ?", fund))
	c.JSON(200, fund)
}
func getStrategy(c *gin.Context) {
	symbol := c.Query("symbol")
	strategy := c.Query("strategy")
	ex := c.Query("exchange")
	eng := (*engine.GetEngine(exchange.Type(ex), conf.BinanceApiKey, conf.BinanceApiSecret)).(*engine.Engine)
	strategies := eng.Base.Strategies
	data := make(map[string]interface{})
	for _, s := range strategies {
		if s.GetSymbol() == symbol && utils.GetTypeName(s) == strategy {
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
