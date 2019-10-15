package handlers

import (
	"artificer/pkg/api/models"
	"artificer/pkg/config"
	"artificer/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func validateArbitraryNoSubjectRequest(req *ArbitraryNoSubjectRequest) (err error) {

	schemaLoader := gojsonschema.NewStringLoader(`{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"type": "object",
		"additionalProperties": {
		  "type": "array",
		  "items": {
			"type": "string"
		  }
		}
	  }`)

	if !json.Valid([]byte(req.ArbitraryClaims)) {
		err = errors.New("arbitrary_claims: is not a valid json")
		fmt.Println(err.Error())
		return
	}
	documentLoader := gojsonschema.NewStringLoader(req.ArbitraryClaims)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !result.Valid() {
		err = errors.New("arbitrary_claims: did not pass schema validation")
		fmt.Println(err.Error())
		return
	}

	schemaLoader = gojsonschema.NewStringLoader(`{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"type": "array",
		 "items": {
			"type": "string"
		  }
	  }`)

	ts := strings.TrimSpace(req.ArbitraryAmrs)
	if len(ts) > 0 {
		if !json.Valid([]byte(ts)) {
			err = errors.New("arbitrary_amrs: is not a valid json")
			fmt.Println(err.Error())
			return
		}
		documentLoader = gojsonschema.NewStringLoader(req.ArbitraryAmrs)
		result, err = gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if !result.Valid() {
			err = errors.New("arbitrary_amrs: did not pass schema validation")
			fmt.Println(err.Error())
			return
		}
	}

	ts = strings.TrimSpace(req.ArbitraryAudiences)
	if len(ts) > 0 {
		if !json.Valid([]byte(ts)) {
			err = errors.New("arbitrary_audiences: is not a valid json")
			fmt.Println(err.Error())
			return
		}
		documentLoader = gojsonschema.NewStringLoader(req.ArbitraryAudiences)
		result, err = gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if !result.Valid() {
			err = errors.New("arbitrary_audiences: did not pass schema validation")
			fmt.Println(err.Error())
			return
		}
	}
	return
}

func validateClient(req *TokenRequest) (err error) {

	if len(req.ClientID) == 0 || len(req.ClientSecret) == 0 {
		err = errors.New("client_id or client_secret is not present")
		fmt.Println(err.Error())
		return
	}

	sEnc := util.StringSha256Encode64(req.ClientSecret)

	var client *models.Client
	client = config.ClientMap[req.ClientID]
	if client == nil {
		err = errors.New(fmt.Sprintf("client_id: %s does not exist", req.ClientID))
		fmt.Println(err.Error())
		return
	}

	foundSecret := false
	for _, element := range client.ClientSecrets {
		foundSecret = (sEnc == element.Value)
		if foundSecret {
			break
		}
	}
	if !foundSecret {
		err = errors.New(fmt.Sprintf("client_id: %s does not have a match for client_secret: %s", req.ClientID, req.ClientSecret))
		fmt.Println(err.Error())
		return
	}

	agt := client.AllowedGrantTypesMap[req.GrantType]

	if agt == nil {
		err = errors.New(fmt.Sprintf("client_id: %s is not authorized for grant_type: %s", req.ClientID, req.GrantType))
		fmt.Println(err.Error())
		return
	}

	return
}
