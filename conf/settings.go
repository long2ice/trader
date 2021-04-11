package conf

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	BinanceApiKey    string
	BinanceApiSecret string
	ServerHost       string
	ServerPort       string
	Debug            bool
)

func InitConfig(configFile string) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error reading config file, %s", err)
	}
	BinanceApiKey = viper.GetString("binance.api_key")
	BinanceApiSecret = viper.GetString("binance.api_secret")
	ServerHost = viper.GetString("server.host")
	ServerPort = viper.GetString("server.port")
	Debug = viper.GetBool("Debug")
}
