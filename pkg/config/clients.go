package config

import (
	"artificer/pkg/api/models"

	"sort"

	"github.com/spf13/viper"
)

var (
	ClientsConfig *viper.Viper
	Clients       []models.Client
	ClientMap     = make(map[string]*models.Client)
)

func LoadClientConfig() {

	ClientsConfig = viper.New()
	ClientsConfig.SetConfigFile(`config/clients.json`)
	err := ClientsConfig.ReadInConfig()
	if err != nil {
		panic(err)
	}
	ClientsConfig.UnmarshalKey("clients", &Clients)
	for _, v := range Clients {
		ClientMap[v.ClientID] = &v
		sort.Strings(v.AllowedGrantTypes)
		sort.Strings(v.AllowedScopes)
	}
}
