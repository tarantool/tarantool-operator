package utils

import (
	"fmt"
)

func GetStatefulSetPodName(stsName string, ordinal int32) string {
	return fmt.Sprintf("%s-%d", stsName, ordinal)
}
