package reference_resolution

import (
	"fmt"
)

func transformInterfaceToString(data map[string]interface{}) map[string]string {
	
	res := make(map[string]string)

	for key, val := range data {
		res[key] = fmt.Sprint(val)
	}

	return res
}