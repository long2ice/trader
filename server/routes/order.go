package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/long2ice/trader/db"
)

func GetOrders(c *gin.Context) {
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
