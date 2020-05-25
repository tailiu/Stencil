package SA1_display

import (
	"fmt"
	"log"
	"stencil/reference_resolution_v2"
)

func (display *display) checkReferenceIndeedResolved(table, col, 
	tableID, colID, value string) (map[string]interface{}, error) {

	// First we must assume that it has already been resolved. If it has not been resolved,
	// then we cannot get data. Otherwise we just return the obtained data
	// Note that if we first assume that it has not been resolved and get data using 
	// prevID, then we could get wrong results
	data, err := display.getRowsBasedOnDependency(table, col, value)

	// This could happen when table1 and col1 have been resolved
	if err == nil {

		log.Println("The reference seems to have already been resolved")
		
		// Now we have not encountered data1 with more than one piece of data
		for _, data1 := range data {

			log.Println("Check one data resolved or not")
			log.Println(data1)

			resolvedVal := display.rr.ReferenceResolved(tableID, colID, fmt.Sprint(data1["id"]))
			
			if resolvedVal == value {
				log.Println("It was indeed resolved")
				return data1, nil 
			
			// This could happen when there was some data happening to satisfy the condition,
			// but that data actually does not have relationships with table0 and col0
			// For example, table0: accounts, col0: id, table1: users, col1: account_id
			// we can get a data with users.account_id pointing to accounts.id, but 
			// the users.account_id of that data is actually old application value and the data
			// may be from other migrations.
			} else {
				log.Println("There happened to be some data satisfying but that data is not what we want")
			}
		}

		log.Println("The reference has not been found in the resolved reference table")

		return nil, DataNotWanted

	} else {
		
		log.Println("The reference has not been found in the resolved reference table")

		return nil, DataNotFound
	} 
}

func (display *display) checkResolveRefWithIDInData(table, col, tableID, colID, id, attrVal string) (string, error) {

	// If favourites.status_id should be resolved (in this case, it should be),
	// we check whether the reference has been resolved or not
	newVal := display.rr.ReferenceResolved(tableID, colID, id)
	
	// If the reference has been resolved, then use the new reference to get data
	// Otherwise, we try to resolve the reference
	if newVal != "" {
		log.Println("The reference has already been resolved")
		return newVal, nil	
	} else {

		attr0 := reference_resolution_v2.CreateAttribute(display.dstAppConfig.appID, tableID, colID, attrVal, id)
		log.Println("Resolving reference:", display.rr.LogAttr(attr0))

		display.rr.ResolveReference(attr0)
		
		// Here we check again to get updated attributes and values
		// instead of using the returned values from the ResolveReference
		// because ResolveReference only returns the updated values in that
		// function call. Values could be updated by other threads and in this
		// case, ResolveReference does not return the updated attribute and value
		// Therefore, we check all updated attributes again here
		updatedAttrs := display.rr.GetUpdatedAttributes(tableID, id)

		log.Println("All updated attributes with values so far:", updatedAttrs)
		
		// We check whether the desired attr (col0) has been resolved
		foundResolvedAttr := false
		for attr, val := range updatedAttrs {
			if attr == col {
				newVal = val
				foundResolvedAttr = true
				break
			}
		}

		// If we find that col has been resolved, then we can use it to get other data
		// Otherwise we cannot use the unresolved reference to get other data in node
		if foundResolvedAttr {
			log.Println("The reference has just been resolved")
			return newVal, nil
		} else {
			log.Println("After trying to resolve the reference, the attribute has not been resolved")
			return "", CannotResolveRefersWithIDInData
		}
	}
}

func (display *display) checkResolveReferenceInGetDataInNode(
	id, table0, col0, table1, col1, value string) (map[string]interface{}, error) {

	// We use table0 and col0 to get table1 and col1
	log.Println("+++++++++++++++++++")
	log.Println("id:", id)
	log.Println(table0)
	log.Println(col0)
	log.Println("value:", value)
	log.Println(table1)
	log.Println(col1)
	log.Println("+++++++++++++++++++")
	
	table0ID := display.dstAppConfig.tableNameIDPairs[table0]
	table1ID := display.dstAppConfig.tableNameIDPairs[table1]

	col0ID := display.dstAppConfig.colNameIDPairs[table0 + ":" + col0]
	col1ID := display.dstAppConfig.colNameIDPairs[table1 + ":" + col1]

	// First, we need to get the attribute that requires reference resolution
	// For example, we have *account.id*, and we want to get *users.account_id*
	// We check whether account.id needs to be resolved
	if display.needToResolveReference(table0, col0) {

		log.Println("Checking reference1 resolved or not")

		newVal, err := display.checkResolveRefWithIDInData(table0, col0, table0ID, col0ID, id, value)

		// If account.id should be resolved (in this case, it should not),
		// we check whether the reference has been resolved or not
		// newVal := display.rr.ReferenceResolved(table0ID, col0ID, id)
		
		// If the reference has been resolved, then use the new reference to get data
		if newVal != "" {
			return display.getOneRowBasedOnDependency(table1, col1, newVal)
		} else {
			display.logUnresolvedRefAndData(table0, table0ID, col0, id)
			return nil, err
		}

	// We check if users.account_id needs to be resolved (of course, in this case, it should be)
	// However we don't know its id (this is the differece from the above case!!). 
	// Also if the value is the value of "id", this could not be used directly 
	// We decide to make so much efforts to resolve "backwards" because inner-dependencies like
	// "statuses.id":"statuses.status_id"
	// "statuses.id":"mentions.status_id"
	// "statuses.id":"stream_entries.activity_id"
	// force us to do in this way. Otherwise, we cannot get other data in a node through statuses.id
	} else if display.needToResolveReference(table1, col1) {

		log.Println("Checking reference2 resolved or not")
		
		data1, err1 := display.checkReferenceIndeedResolved(table1, col1, table1ID, col1ID, value)
		if err1 != nil {
			log.Println(err1)
		} else {
			return data1, nil
		}
		
		attr1 := reference_resolution_v2.CreateAttribute(display.dstAppConfig.appID, table0ID, col0ID, id, id)
		
		log.Println("Resolving reference2:", display.rr.LogAttr(attr1))

		display.rr.ResolveReference(attr1)
		
		log.Println("After trying to resolving reference2:", display.rr.LogAttr(attr1))

		data2, err2 := display.checkReferenceIndeedResolved(table1, col1, table1ID, col1ID, value)
		if err2 != nil {
			display.logUnresolvedRefAndData(table1, table1ID, col1)
			return nil, CannotGetDataAfterResolvingRef2
		} else {
			return data2, nil
		}
	
	// Normally, there must exist one that needs to be resolved. 
	// However, up to now there is a case breaking the above rule. 
	// When there is no mapping
	// For example: 
	// When migrating from Diaspora to Mastodon, 
	// there is no mapping to stream_entries.activity_id.
	} else {
		return nil, NoReferenceToResolve
	}
}

func (display *display) checkResolveReferenceInGetDataInParentNode(table, col, attrVal, id string) (string, error) {
	
	log.Println("+++++++++++++++++++")
	log.Println(table)
	log.Println(col)
	log.Println("+++++++++++++++++++")

	log.Println("Parent Node: before checking reference resolved or not")

	tableID := display.dstAppConfig.tableNameIDPairs[table]
	colID := display.dstAppConfig.colNameIDPairs[table + ":" + col]

	// Normally, there must exist one that needs to be resolved. 
	// But this could happen for example, in Diaspora, posts.id depends on aspects.shareable_id
	// There is no need to resolve id here in the else case
	if display.needToResolveReference(table, col) {
		if newVal, err := display.checkResolveRefWithIDInData(table, col, tableID, colID, id, attrVal); err != nil {
			display.logUnresolvedRefAndData(table, tableID, col, id)
			return "", err
		} else {
			return newVal, nil
		}
	} else {
		return "", NoReferenceToResolve
	}
}