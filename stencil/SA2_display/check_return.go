package SA2_display

import (
	"log"
	"strings"
	"strconv"
)

func CheckCombinedDisplayConditions(displayConfig *displayConfig, 
	pTagConditions map[string]bool, oneMigratedData *HintStruct) bool {	
	
	if len(pTagConditions) == 1 {

		for _, result := range pTagConditions{
			return result
		}

	}

	combinedSettings, err := oneMigratedData.GetCombinedDisplaySettings(displayConfig)
	if err != nil {
		log.Fatal(err)
	}
	
	for key, val := range pTagConditions {
		combinedSettings = strings.Replace(combinedSettings, key, strconv.FormatBool(val), 1)
	}
	strs := strings.Split(combinedSettings, " ")

	var combinedResults bool 
	var operator string

	for i, val := range strs {

		if i == 0 {

			result, err := strconv.ParseBool(val)
			if err != nil {
				log.Fatal(err)
			}
			
			combinedResults = result
		
		} else if i % 2 == 1 {
			
			operator = val
		
		} else {

			result, err := strconv.ParseBool(val)
			if err != nil {
				log.Fatal(err)
			}

			if operator == "or" {

				combinedResults = combinedResults || result

			} else {
				
				combinedResults = combinedResults && result
			}
		}
	}

	return combinedResults

}