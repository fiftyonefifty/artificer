package config

import (
	"artificer/pkg/api/contracts"
	"artificer/pkg/api/models"
	"artificer/pkg/util"

	"fmt"
	"sort"
	"strings"

	"path"

	"sync"

	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

type InMemoryClientStore struct {
	ClientMapA map[string]*models.Client
	ClientMapB map[string]*models.Client
	pCurrent   *map[string]*models.Client
}

var (
	ClientsConfig        *viper.Viper
	Clients              []models.Client
	pInMemoryClientStore *InMemoryClientStore = &InMemoryClientStore{}
	mutex                                     = &sync.Mutex{}
	clientRequest        chan *clientConfigRequestOp
	queueLength          int = 5
)

type clientConfigResponse struct {
	Client models.Client
	Found  bool
}
type clientConfigRequestOp struct { // bank operation: deposit or withdraw
	Id      string                    // amount
	confirm chan clientConfigResponse // confirmation channel
}

func InitializeClientConfig() {
	// mutex is here to make sure that we only do this once
	mutex.Lock()
	defer mutex.Unlock()
	if clientRequest != nil {
		return
	}

	// setup the worker
	clientRequest = make(chan *clientConfigRequestOp, queueLength)
	for i := 0; i < queueLength; i++ {
		go func() {
			for {
				/* The select construct is non-blocking:
				-- if there's something to read from a channel, do so
				-- otherwise, fall through to the next case, if any */
				select {
				case request := <-clientRequest:
					found, client := pInMemoryClientStore.GetClientUnsafe(request.Id)
					response := clientConfigResponse{
						Found: found,
					}
					if found {
						response.Client = client
					}

					request.confirm <- response // send back the response
				}
			}
		}()
	}
}

func NewClientStore() contracts.IClientStore {
	return pInMemoryClientStore
}

func (store InMemoryClientStore) GetClientUnsafe(id string) (found bool, client models.Client) {

	currenClientMap := *store.pCurrent

	found = false
	c := currenClientMap[id]
	if c == nil {
		return
	}
	client = *c
	found = true
	return

}
func (store InMemoryClientStore) GetClient(id string) (found bool, client models.Client) {

	request := &clientConfigRequestOp{Id: id, confirm: make(chan clientConfigResponse)}
	clientRequest <- request
	response := <-request.confirm
	found = response.Found
	if found {
		client = response.Client
	}
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
	if pInMemoryClientStore.pCurrent == nil || pInMemoryClientStore.pCurrent == &pInMemoryClientStore.ClientMapB {
		pInMemoryClientStore.pCurrent = &pInMemoryClientStore.ClientMapA
	} else {
		pInMemoryClientStore.pCurrent = &pInMemoryClientStore.ClientMapB
	}

	a := make(map[string]*models.Client)

	for _, v := range Clients {
		a[v.ClientID] = &v
		sort.Strings(v.AllowedGrantTypes)
		sort.Strings(v.AllowedScopes)
		util.FilterOutStringElement(&v.AllowedScopes, "artificer-ns")
		v.AllowedGrantTypesMap = make(map[string]interface{})
		for _, agt := range v.AllowedGrantTypes {
			v.AllowedGrantTypesMap[agt] = ""
		}
	}
	pInMemoryClientStore.pCurrent = &a

}
