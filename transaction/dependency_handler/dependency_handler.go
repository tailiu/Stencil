package dependency_handler

import (
	"transaction/config"
	"transaction/display"
	"fmt"
)

func GetParent() {

}

func checkDependentDataExists(dependencies config.InnerDependencies, members config.Members, hint display.HintStruct) bool {
	
}

func CheckNodeComplete(innerDependencies []config.Tag, hint display.HintStruct) bool {
	for _, innerDependency := range innerDependencies {
		for _, member := range innerDependency.Members{
			if hint.Table == member {
				if len(innerDependency.Members) == 1 {
					return true
				} else {
					// Note: we assume that one dependency represents one row 
					// 		in one table depends on another row in another table
					// for memberKey, memberVal := range innerDependency.Members {
					// 	fmt.Println(memberKey, memberVal)

					// }
					if !checkDependentDataExists(innerDependency.InnerDependencies, innerDependency.Members, hint) {
						return false
					} else {
						return true
					}
				}
			}
		}
	}
}
