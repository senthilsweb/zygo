package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Setup initialize configuration
func Setup() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AllowEmptyEnv(true)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

}
