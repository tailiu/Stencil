package reference_resolution

import (
	"fmt"
	"log"
)

func transformInterfaceToString(data map[string]interface{}) map[string]string {
	
	res := make(map[string]string)

	for key, val := range data {
		res[key] = fmt.Sprint(val)
	}

	return res
}

func combineTwoMaps(m1 map[string]string, m2 map[string]string) {

	for k, v := range m2 {

		if _, ok := m1[k]; ok {
			log.Println("Found an overlapped key when combing two maps!")
		} else {
			m1[k] = v
		}
		
	}

}