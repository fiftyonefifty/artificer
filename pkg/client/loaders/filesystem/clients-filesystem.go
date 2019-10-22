package filesystem

import (
	"artificer/pkg/client/models"
	"artificer/pkg/util"

	"errors"
	"fmt"
	"path"

	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

func FetchClientConfigFromFileSystem(rootFolder string) (clients []models.Client, err error) {

	schemaPath := util.ToCanonical(path.Join(rootFolder, "config/clients.schema.json"))
	documentPath := util.ToCanonical(path.Join(rootFolder, "config/clients.json"))

	schemaLoader := gojsonschema.NewReferenceLoader(schemaPath)
	documentLoader := gojsonschema.NewReferenceLoader(documentPath)
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
		err = errors.New(fmt.Sprintf("document:[%s] did not pass schema:[%s] validation errors:[%s]",
			documentPath, schemaPath, sError))
		return
	}

	clientsConfig := viper.New()
	clientsConfig.SetConfigFile("config/clients.json")
	err = clientsConfig.ReadInConfig()
	if err != nil {
		return
	}

	err = clientsConfig.UnmarshalKey("clients", &clients)
	if err != nil {
		return
	}
	return
}
