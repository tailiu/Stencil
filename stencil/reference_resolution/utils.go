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

func combineTwoMaps(m1 map[string]string, m2 map[string]string) map[string]string {

	m3 := make(map[string]string)

	for k1, v1 := range m1 {
		m3[k1] = v1
	}

	for k2, v2 := range m2 {

		if _, ok := m3[k2]; ok {
			log.Println("Found an overlapped key when combing two maps!")
		} else {
			m3[k2] = v2
		}
		
	}

	return m3

}