package keyvault

import (
	"artificer/pkg/client/models"
	"artificer/pkg/keyvault"
	"artificer/pkg/util"
	"encoding/json"
	"path"

	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

func FetchClientConfigFromKeyVault(rootFolder string) (clients []models.Client, err error) {

	schemaPath := util.ToCanonical(path.Join(rootFolder, "config/clients.schema.json"))
	schemaLoader := gojsonschema.NewReferenceLoader(schemaPath)

	keyVaultSecretName := viper.GetString("clientConfig.KeyVaultSecretName") // "artificer-clients-int"

	bundle, er := keyvault.GetSecret(keyVaultSecretName)
	if er != nil {
		err = er
		return
	}
	documentLoader := gojsonschema.NewStringLoader(*bundle.Value)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return
	}
	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		sError := ""
		for _, desc := range result.Errors() {
			sError += fmt.Sprintf("- %s\n", desc)
		}
		fmt.Println(sError)
		err = errors.New(fmt.Sprintf("document:did not pass schema:[%s] validation errors:[%s]",
			schemaPath, sError))
		return
	}

	var clientContainer struct {
		Clients []models.Client
	}

	err = json.Unmarshal([]byte(*bundle.Value), &clientContainer)
	clients = clientContainer.Clients
	return
}
