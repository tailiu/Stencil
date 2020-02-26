package SA2_display

import (
	"stencil/config"
	"stencil/db"
	"strings"
	"log"
	"fmt"
)

func getADataInOwner(displayConfig *displayConfig, hints []HintStruct,
	ownership *config.Ownership, dag) (HintStruct, error) {
	
	var hint HintStruct

	// For now, there is no case in which there are more than one condition
	// so we only need the first condition here
	condition := ownership.Conditions[0]
	
	tableAttr := ReplaceKey(displayConfig.dstAppConfig.dag, 
		ownership.Tag, condition.TagAttr)

	dependsOnTableAttr := ReplaceKey(displayConfig.dstAppConfig.dag, 
		ownership.OwnedBy, condition.DependsOnAttr)

	table := strings.Split(tableAttr, ".")[0]
	attr := strings.Split(tableAttr, ".")[1]

	dependsOnTable := strings.Split(dependsOnTableAttr, ".")[0]
	dependsOnAttr := strings.Split(dependsOnTableAttr, ".")[1]
	
	var depVal string
		
	for _, hint := range hints {

		if hint.Table == table {
			depVal = fmt.Sprint(hint.Data[attr])
		}

	}

	if depVal == "" {
		return nil, CannotFindDataInOwnership
	}

	var query string

	query = fmt.Sprintf(
		`SELECT * FROM "%s" WHERE %s = '%s'`,
		dependsOnTable, dependsOnAttr, depVal,
	)
	
	// log.Println("Get a data in the owner node")
	// log.Println(query)

	data, err := db.DataCall1(displayConfig.dstAppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("One data in the ownership node:")
	log.Println(data)

	hint = TransformRowToHint(displayConfig, data, dependsOnTable, "root")

	return hint, nil

}

func getOwner(displayConfig *displayConfig, hints []*HintStruct,
	ownership *config.Ownership) ([]HintStruct, error) {
	
	oneDataInOwnerNode, err := getADataInOwner(displayConfig, hints, ownership)
	if err != nil {
		return nil, err
	}

	return GetDataInNodeBasedOnDisplaySetting(displayConfig, oneDataInOwnerNode)
	
}