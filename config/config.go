package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type NSO struct {
	Username string
	Password string
	BaseURI  string
	Services []string
}

type ConfServices struct {
	Services []string
}

type HTTP struct {
	Ip   string
	Port string
}

type Log struct {
	Loglevel string
}

type Config struct {
	NSO
	HTTP
	Log
}

func init() {
	viper.SetConfigName("nso_exporter")

	// These 2 imports are for testing only
	viper.AddConfigPath(".")
	viper.AddConfigPath("../.")

	// This is where you are suppose to have the config file in prod
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("/home/nso_exporter/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Could not read config ", err))
	}
}

var ConfigInstance Config

func GetConfig() Config {
	_ = viper.Unmarshal(&ConfigInstance)
	return ConfigInstance
}
