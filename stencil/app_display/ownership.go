package app_display

import (
	"stencil/config"
	"strings"
	"log"
	"errors"
)


func getADataInOwner(displayConfig *config.DisplayConfig, hints []*HintStruct,
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
			val = fmt.Sprint(hint[attr])
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

	return data, nil

}

func GetOwner(displayConfig *config.DisplayConfig, hints []*HintStruct,
	ownership *config.Ownership) ([]*HintStruct, error) {
	
	oneDataInOwnership, err := getADataInOwner(displayConfig, hints, ownership)
	if err != nil {
		return nil, err
	}

	
	
}