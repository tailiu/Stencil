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
	
	var depVal string

	if !displayConfig.resolveReference {
		
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

		// log.Println(dependsOnTable, dependsOnAttr)
		// log.Println(table, attr, id)

		// Here we must resolve reference first otherwise we cannot get the data in root
		// For example, follows.account_id is still referring to the old id
		depVal, err0 = checkResolveReferenceInGetDataInParentNode(displayConfig, 
			id, table, attr)

		// log.Println(depVal)
		// log.Println(err0)
		
		// Refresh the cached results which could have changed due to
		// reference resolution 
		refreshCachedDataHints(displayConfig, hints)

		if err0 != nil {
			return nil, err0
		}

	}

	var query string

	if !displayConfig.markAsDelete {
		query = fmt.Sprintf(
			"SELECT * FROM %s WHERE %s = %s",
			dependsOnTable, dependsOnAttr, depVal,
		)
	} else {
		query = fmt.Sprintf(
			"SELECT * FROM %s WHERE %s = %s and mark_as_delete = false",
			dependsOnTable, dependsOnAttr, depVal,
		)
	}
	
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

// This old function does not consider migrating shared data or checked data not belonging to
// the current migrating user
func oldHandleNonRootNode(displayConfig *displayConfig,
	dataInNode []*HintStruct) error {

	// As an optimization, the display thread caches the display result according to
	// ownership display settings.
	// If the cached display result is false, it means that either the display thread has not checked,
	// or the display settings are not satisfied, 
	// so the display thread needs to perform the following checks.
	// We can cache the result because normally after the root node is displayed, it should not be deleted.
	// Users are allowed to do concurrent deletion migraitons and 
	// even if users perform concurrent consistent or independent migrations, 
	// users' root nodes are not affected.
	// In the very rare case in which the migrating user quits the destination application,
	// we can detect such case, and invalide the cache result.
	if !displayConfig.dstAppConfig.ownershipDisplaySettingsSatisfied {

		dataOwnershipSpec, err12 := dataInNode[0].GetOwnershipSpec(displayConfig)
		if err12 != nil {
			log.Fatal(err12)
		}

		log.Println(dataOwnershipSpec)

		dataInOwnerNode, err13 := getOwner(displayConfig, dataInNode, dataOwnershipSpec)

		// The root node could be incomplete
		if err13 != nil {
			log.Println("Get a root node error:")
			log.Println(err13)
		}

		// Display the data not displayed in the root node
		// this root node should be the migrating user's root node
		if len(dataInOwnerNode) != 0 {

			displayedDataInOwnerNode, notDisplayedDataInOwnerNode := checkDisplayConditionsInNode(
				displayConfig, dataInOwnerNode)
			
			if len(displayedDataInOwnerNode) != 0 {

				err6 := Display(displayConfig, notDisplayedDataInOwnerNode)
				if err6 != nil {
					log.Fatal(err6)
				}

			}

		}
		
		// If based on the ownership display settings this node is allowed to be displayed,
		// then continue to check dependencies.
		// Otherwise, no data in the node can be displayed.
		if displayResultBasedOnOwnership := CheckOwnershipCondition(
			dataOwnershipSpec.Display_setting, err13); !displayResultBasedOnOwnership {

			log.Println(`Ownership display settings are not satisfied, 
				so this node cannot be displayed`)

			return NoNodeCanBeDisplayed

		} else {

			// Cache the check result
			displayConfig.dstAppConfig.ownershipDisplaySettingsSatisfied = true

		}
	}

	return nil 

}