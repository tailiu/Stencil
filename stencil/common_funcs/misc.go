package common_funcs

import (
	"fmt"
)

func ExistsInSlice(s []string, e string) bool {

	for _, s1 := range s {
		if e == s1 {
			return true
		}
	}

	return false
}

func TransformInterfaceToString(data map[string]interface{}) map[string]string {
	
	res := make(map[string]string)

	for key, val := range data {
		res[key] = fmt.Sprint(val)
	}

	return res
}