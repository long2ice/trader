package conf

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	BinanceApiKey    string
	BinanceApiSecret string
	DatabaseDsn      string
)

func InitConfig(configFile string) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error reading config file, %s", err)
	}
	BinanceApiKey = viper.GetString("binance.api_key")
	BinanceApiSecret = viper.GetString("binance.api_secret")
	DatabaseDsn = viper.GetString("database.dsn")
}
