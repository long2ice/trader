package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long2ice/trader/conf"
	"github.com/long2ice/trader/server/routes"
	log "github.com/sirupsen/logrus"
)

func Start() {
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.GET("/orders", routes.GetOrders)
	r.GET("/strategy", routes.GetStrategy)
	r.GET("/strategies", routes.GetStrategies)
	r.GET("/fund", routes.GetFund)
	r.POST("/fund", routes.AddFund)

	err := r.Run(fmt.Sprintf("%s:%s", conf.ServerHost, conf.ServerPort))
	if err != nil {
		log.WithField("err", err).Error("Start server failed")
	}
}
