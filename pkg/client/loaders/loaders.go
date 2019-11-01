package loaders

import (
	"artificer/pkg/client/contracts"
	"artificer/pkg/client/loaders/filesystem"
	"artificer/pkg/client/loaders/keyvault"
	"artificer/pkg/client/models"
	"artificer/pkg/util"
	"context"
	"log"
	"sort"
	"sync"
)

type ClientConfigOptions struct {
	RootFolder  string
	UseKeyVault bool
}
type InMemoryClientStore struct {
	ClientMapA map[string]*models.Client
	ClientMapB map[string]*models.Client
	pCurrent   *map[string]*models.Client
}

var (
	clients              []models.Client
	pInMemoryClientStore *InMemoryClientStore = &InMemoryClientStore{}
	mutex                                     = &sync.Mutex{}
	clientRequest        chan *clientConfigRequestOp
	queueLength          int = 5
	clientConfigOptions  ClientConfigOptions
)

type clientConfigResponse struct {
	Client models.Client
	Found  bool
}
type clientConfigRequestOp struct { // bank operation: deposit or withdraw
	ctx     context.Context
	id      string                    // amount
	confirm chan clientConfigResponse // confirmation channel
}

func InitializeClientConfig(options ClientConfigOptions) {
	// mutex is here to make sure that we only do this once
	mutex.Lock()
	defer mutex.Unlock()
	if clientRequest != nil {
		return
	}
	clientConfigOptions = options
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
					found, client := pInMemoryClientStore.getClientUnsafe(request.id)
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

func (store InMemoryClientStore) getClientUnsafe(id string) (found bool, client models.Client) {

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
func (store InMemoryClientStore) GetClient(ctx context.Context, id string) (found bool, client models.Client) {

	request := &clientConfigRequestOp{ctx: ctx, id: id, confirm: make(chan clientConfigResponse)}
	clientRequest <- request
	response := <-request.confirm
	found = response.Found
	if found {
		client = response.Client
	}
	return
}

func LoadClientConfigFromKeyVault(ctx context.Context) (clients []models.Client, err error) {
	clients, err = keyvault.FetchClientConfigFromKeyVault(clientConfigOptions.RootFolder)
	return
}

func LoadClientConfig(ctx context.Context) (err error) {

	if clientConfigOptions.UseKeyVault {
		clients, err = keyvault.FetchClientConfigFromKeyVault(clientConfigOptions.RootFolder)
	} else {
		clients, err = filesystem.FetchClientConfigFromFileSystem(clientConfigOptions.RootFolder)
	}

	if err != nil {
		log.Fatalf("Failed to load client config: useKeyvault: %t,err %s", clientConfigOptions.UseKeyVault, err.Error())
		return
	}

	a := make(map[string]*models.Client)

	for _, v := range clients {
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
	return
}
