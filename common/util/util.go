package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
)

const (
	FailedToCreateNSQConsumerError = "Failed to create NSQ consumer"
	FailedToConnectToLookupdError  = "Failed to connect NSQ consumer to lookupd"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func ToJSON(a interface{}) string {
	bytes, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		fmt.Println(err)
		return "Couldn't convert inteface to json string"
	}
	return string(bytes)
}

func NewError(baseText string, err error) error {
	return fmt.Errorf("%s: %s", baseText, err)
}
