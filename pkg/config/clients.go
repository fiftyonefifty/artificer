package config

import (
	"artificer/pkg/api/contracts"
	"artificer/pkg/api/models"
	"artificer/pkg/util"

	"fmt"
	"sort"
	"strings"

	"path"

	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

var (
	ClientsConfig *viper.Viper
	Clients       []models.Client
	ClientMap     = make(map[string]*models.Client)
)

type inMemoryClientStore struct {
	clientMap map[string]*models.Client
}

func NewClientStore() contracts.IClientStore {
	store := inMemoryClientStore{
		clientMap: ClientMap,
	}
	return store
}
func (store inMemoryClientStore) GetClient(id string) (found bool, client models.Client) {

	found = false
	c := store.clientMap[id]
	if c == nil {
		return
	}
	client = *c
	return

}

func ToCanonical(src string) string {
	var replacer = strings.NewReplacer("\\", "/")
	str := replacer.Replace(src)
	return "file:///" + str
}

func LoadClientConfig(processDirectory string) {

	schemaPath := ToCanonical(path.Join(processDirectory, "config/clients.schema.json"))
	documentPath := ToCanonical(path.Join(processDirectory, "config/clients.json"))

	schemaLoader := gojsonschema.NewReferenceLoader(schemaPath)
	documentLoader := gojsonschema.NewReferenceLoader(documentPath)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}
	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	ClientsConfig = viper.New()
	ClientsConfig.SetConfigFile("config/clients.json")
	err = ClientsConfig.ReadInConfig()
	if err != nil {
		panic(err)
	}
	ClientsConfig.UnmarshalKey("clients", &Clients)
	for _, v := range Clients {
		ClientMap[v.ClientID] = &v
		sort.Strings(v.AllowedGrantTypes)
		sort.Strings(v.AllowedScopes)
		util.FilterOutStringElement(&v.AllowedScopes, "artificer-ns")
		v.AllowedGrantTypesMap = make(map[string]interface{})
		for _, agt := range v.AllowedGrantTypes {
			v.AllowedGrantTypesMap[agt] = ""
		}
	}
}
