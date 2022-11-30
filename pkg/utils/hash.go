package utils

import (
	"encoding/json"
	"fmt"
	"hash/fnv"

	"k8s.io/apimachinery/pkg/util/rand"
)

func HashObject(data interface{}) (string, error) {
	hashFunc := fnv.New32()

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	_, err = hashFunc.Write(jsonBytes)
	if err != nil {
		return "", err
	}

	return rand.SafeEncodeString(fmt.Sprint(hashFunc.Sum32())), nil
}
