package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/long2ice/trader/db"
	"gorm.io/gorm"
)

func GetFund(c *gin.Context) {
	strategy := c.Query("strategy")
	var fund db.Fund
	db.Client.Where("strategy = ?", strategy).Find(&fund)
	c.JSON(200, fund)
}

func AddFund(c *gin.Context) {
	fund := c.PostForm("fund")
	strategy := c.PostForm("strategy")
	db.Client.Model(&db.Fund{}).Where("strategy = ?", strategy).UpdateColumn("fund", gorm.Expr("fund + ?", fund))
	c.JSON(200, fund)
}
