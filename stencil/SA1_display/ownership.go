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
	
	var hint *HintStruct

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
	
	// var val string
	
	// for _, hint := range hints {

	// 	if hint.Table == table {
	// 		val = fmt.Sprint(hint.Data[attr])
	// 	}

	// }

	// if val == "" {
	// 	return nil, CannotFindDataInOwnership
	// }
	
	var id string

	for _, hint1 := range hints {

		if hint1.Table == table {
			id = fmt.Sprint(hint1.KeyVal["id"])
		}

	}

	if id == "" {
		return nil, CannotFindDataInOwnership
	}

	// log.Println(dependsOnTable, dependsOnAttr)
	// log.Println(table, attr, id)

	// Here we must resolve reference first otherwise we cannot get the data in root
	// For example, follows.account_id is still referring to the old id
	depVal, err0 := checkResolveReferenceInGetDataInParentNode(displayConfig, 
		id, table, attr)

	// log.Println(depVal)
	// log.Println(err0)

	if err0 != nil {
		return nil, err0
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s",
		dependsOnTable, dependsOnAttr, depVal)
	
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
	ownership *config.Ownership) ([]*HintStruct, error) {
	
	oneDataInOwnerNode, err := getADataInOwner(displayConfig, hints, ownership)
	if err != nil {
		return nil, err
	}

	return GetDataInNodeBasedOnDisplaySetting(displayConfig, oneDataInOwnerNode)
	
}

func isNodeMigratingUserRootNode(displayConfig *displayConfig, 
	dataInRootNode []*HintStruct) (bool, error) {

	for _, node := range dataInRootNode {

		if node.Table == displayConfig.dstAppConfig.rootTable {

			if node.Data[displayConfig.dstAppConfig.rootAttr] == 
				displayConfig.dstAppConfig.userID {
					
				return true, nil
				
			} else {

				return false, NotMigratingUserRootNode
			}
		}

	}

	return false, CannotFindRootTable

}

func oldHandleRootNode(displayConfig *displayConfig,
	dataInRootNode []*HintStruct) error {

	// If it is the migrating user's root node, the display thread reached this node 
	// by directly picking from migrated data since there are already ownership relationships
	// between normal nodes with the migrating user's root node 
	// and there are no dependencies defined there
	isMigratingUserRootNode, err14 := isNodeMigratingUserRootNode(displayConfig, dataInRootNode)
	if err14 != nil {
		log.Println(err14)
	}

	// If the checked node is the migrating user's root node, then it is not diplayed yet.
	// Since just now the display thread already checked if
	// there is any data displayed in the node and return, diplayed the undisplayed data if
	// there exists some displayed data, and returned the result,
	// Therefore, we need to display the data in the node.
	if isMigratingUserRootNode {

		err15 := Display(displayConfig, dataInRootNode)
		if err15 != nil {
			log.Fatal(err15)
		}
		
		return ReturnResultBasedOnNodeCompleteness(err15)

	// If it is other user's root node, it must be arrived through the dependency relationship
	// since the migrating user root node is only connected with the migrated data with ownership
	// In this case, we do not need to further check the inter-node data dependencies, 
	// ownership, or sharing relationships of the current root node.
	// As the display thread only displays this migrating user's data, even if there is
	// some data not displayed in the root node in this case, it will not display it.
	} else {

		// here the node should have no data able to be displayed
		// since if there is some data already displayed in the checkDisplayConditionsInNode
		// the function returns
		return NoNodeCanBeDisplayed

	}
			
}