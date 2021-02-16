package db

import (
	"github.com/long2ice/trader/conf"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Client *gorm.DB

func init() {
	var err error
	Client, err = gorm.Open(mysql.Open(conf.DatabaseDsn), &gorm.Config{})
	if err != nil {
		log.WithField("err", err).Fatal("Fail to connect db")
	}
	err = Client.AutoMigrate(&Order{}, &KLine{}, &Fund{})
	if err != nil {
		log.WithField("err", err).Fatal("AutoMigrate db fail")
	}
}
