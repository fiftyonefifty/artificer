package util

import (
	"crypto"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

func WaitOnAllChannels(cs ...<-chan bool) <-chan bool {
	out := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(len(cs))

	for _, c := range cs {
		go func(c <-chan bool) {
			<-c
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		out <- true
		close(out)
	}()

	return out
}

// PrintAndLog writes to stdout and to a logger.
func PrintAndLog(message string) {
	log.Println(message)
	fmt.Println(message)
}

func FilterOutStringElement(a *[]string, element string) {
	n := 0
	for _, x := range *a {
		if !strings.EqualFold(x, element) {
			(*a)[n] = x
			n++
		}
	}
	*a = (*a)[:n]
}

func InterfaceArrayToStringArray(source interface{}) (result []string) {
	switch vv := source.(type) {
	case []interface{}:
		result = []string{}
		for _, u := range vv {
			s := u.(string)
			result = append(result, s)
		}
	case []string:
		result = source.([]string)
	default:
		result = nil
	}
	return
}

func Contains(array *[]string, element string, noCase bool) bool {
	for _, e := range *array {
		if noCase {
			if strings.EqualFold(e, element) {
				return true
			}
		} else {
			if e == element {
				return true
			}
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
func ToCanonical(src string) string {
	var replacer = strings.NewReplacer("\\", "/")
	str := replacer.Replace(src)
	return "file:///" + str
}
