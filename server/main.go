package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long2ice/trader/conf"
	log "github.com/sirupsen/logrus"
)

func Start() {
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.GET("/orders", getOrders)
	r.GET("/strategy", getStrategy)
	r.GET("/fund", getFund)
	r.POST("/fund", addFund)

	err := r.Run(fmt.Sprintf("%s:%s", conf.ServerHost, conf.ServerPort))
	if err != nil {
		log.WithField("err", err).Error("Start server failed")
	}
}
