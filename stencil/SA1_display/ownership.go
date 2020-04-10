package SA1_display

import (
	"stencil/config"
	"stencil/db"
	"strings"
	"log"
	"fmt"
)

func (display *display) getADataInOwner(hints []*HintStruct, ownership *config.Ownership) (*HintStruct, error) {
	
	var hint *HintStruct

	// For now, there is no case in which there are more than one condition
	// so we only need the first condition here
	condition := ownership.Conditions[0]
	
	tableAttr := display.dstAppConfig.dag.ReplaceKey( 
		ownership.Tag, condition.TagAttr,
	)

	dependsOnTableAttr := display.dstAppConfig.dag.ReplaceKey( 
		ownership.OwnedBy, condition.DependsOnAttr,
	)

	table := strings.Split(tableAttr, ".")[0]
	attr := strings.Split(tableAttr, ".")[1]

	dependsOnTable := strings.Split(dependsOnTableAttr, ".")[0]
	dependsOnAttr := strings.Split(dependsOnTableAttr, ".")[1]
	
	var depVal string

	if !display.resolveReference {
		
		for _, hint := range hints {

			if hint.Table == table {
				depVal = fmt.Sprint(hint.Data[attr])
			}

		}

		if depVal == "" {
			return nil, CannotFindDataInOwnership
		}
	
	} else {

		var id string
		var err0 error

		for _, hint1 := range hints {
			if hint1.Table == table {
				id = fmt.Sprint(hint1.KeyVal["id"])
			}
		}

		if id == "" {
			return nil, CannotFindDataInOwnership
		}

		depVal, err0 = display.checkResolveReferenceInGetDataInParentNode(table, attr, id)
		
		// no matter whether this attribute has been resolved before
		// we need to refresh the cached data because this attribute might be
		// resolved by other thread checking other data
		display.refreshCachedDataHints(hints)

		if err0 != nil {
			log.Println(err0)
			if err0 != NoReferenceToResolve {
				return nil, err0
			} else {
				depVal = fmt.Sprint(hint.Data[attr])
			}
		}
	}

	var query string

	if !display.markAsDelete {
		query = fmt.Sprintf(
			`SELECT * FROM "%s" WHERE %s = '%s'`,
			dependsOnTable, dependsOnAttr, depVal,
		)
	} else {
		query = fmt.Sprintf(
			`SELECT * FROM "%s" WHERE %s = '%s' and mark_as_delete = false`,
			dependsOnTable, dependsOnAttr, depVal,
		)
	}
	
	// log.Println("Get a data in the owner node")
	// log.Println(query)

	data, err := db.DataCall1(display.dstAppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("One data in the ownership node:")
	log.Println(data)

	hint = TransformRowToHint(display, data, dependsOnTable, "root")

	return hint, nil

}

func (display *display) getOwner(hints []*HintStruct, ownership *config.Ownership) ([]*HintStruct, error) {
	
	oneDataInOwnerNode, err := display.getADataInOwner(hints, ownership)
	if err != nil {
		return nil, err
	}

	return display.GetDataInNodeBasedOnDisplaySetting(oneDataInOwnerNode)
	
}
