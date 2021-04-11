package db

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var Client *gorm.DB

func Init(client *gorm.DB) {
	var err error
	Client = client
	err = Client.AutoMigrate(&Order{}, &KLine{}, &Fund{})
	if err != nil {
		log.WithField("err", err).Error("AutoMigrate db fail")
	}
}
