package SA1_display

import (
	"stencil/config"
	"stencil/db"
	"strings"
	"log"
	"fmt"
)


func getADataInOwner(displayConfig *displayConfig, hints []*HintStruct,
	ownership *config.Ownership) (*HintStruct, error) {
	
	// For now, there is no case in which there are more than one condition
	// so we only need the first condition here
	condition := ownership.Conditions[0]
	
	tableAttr := ReplaceKey(displayConfig, 
		ownership.Tag, condition.TagAttr)

	dependsOnTableAttr := ReplaceKey(displayConfig, 
		ownership.OwnedBy, condition.DependsOnAttr)

	table := strings.Split(tableAttr, ".")[0]
	attr := strings.Split(tableAttr, ".")[1]

	dependsOnTable := strings.Split(dependsOnTableAttr, ".")[0]
	dependsOnAttr := strings.Split(dependsOnTableAttr, ".")[1]
	
	var val string
	
	for _, hint := range hints {

		if hint.Table == table {
			val = fmt.Sprint(hint.Data[attr])
		}

	}

	if val == "" {
		return nil, CannotFindDataInOwnership
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s",
		dependsOnTable, dependsOnAttr, val)
	
	data, err := db.DataCall1(displayConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(data)
	// return data, nil

	return nil, nil

}

func getOwner(displayConfig *displayConfig, hints []*HintStruct,
	ownership *config.Ownership) ([]*HintStruct, error) {
	
	oneDataInOwnerNode, err := getADataInOwner(displayConfig, hints, ownership)
	if err != nil {
		return nil, err
	}

	return GetDataInNodeBasedOnDisplaySetting(displayConfig, oneDataInOwnerNode)
	
}