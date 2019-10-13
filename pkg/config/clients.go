package config

import (
	"artificer/pkg/api/models"

	"github.com/spf13/viper"
)

var (
	ClientsConfig *viper.Viper
	Clients       []models.Client
)

func LoadClientConfig() {

	ClientsConfig = viper.New()
	ClientsConfig.SetConfigFile(`config/clients.json`)
	err := ClientsConfig.ReadInConfig()
	if err != nil {
		panic(err)
	}
	ClientsConfig.UnmarshalKey("clients", &Clients)
}
