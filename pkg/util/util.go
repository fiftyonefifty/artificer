package util

import (
	"crypto"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// PrintAndLog writes to stdout and to a logger.
func PrintAndLog(message string) {
	log.Println(message)
	fmt.Println(message)
}

func Contains(array []string, element string) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}
	return false
}
func StringSha256Encode64(value string) string {
	return ByteArraySha256Encode64([]byte(value))
}

func ByteArraySha256Encode64(value []byte) string {
	digest := crypto.SHA256.New()
	digest.Write(value)
	h := digest.Sum(nil)
	sEnc := b64.StdEncoding.EncodeToString(h)
	return sEnc
}

// ReadJSON reads a json file, and unmashals it.
// Very useful for template deployments.
func ReadJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read template file: %v\n", err)
	}
	contents := make(map[string]interface{})
	json.Unmarshal(data, &contents)
	return &contents, nil
}
