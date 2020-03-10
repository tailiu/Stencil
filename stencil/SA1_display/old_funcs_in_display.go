package SA1_display

import (
	"fmt"
	"log"
	"strconv"
	"stencil/db"
	"stencil/reference_resolution"
	"stencil/schema_mappings"
	"stencil/common_funcs"
	"strings"
)

func oldCheckResolveReferenceInGetDataInNode1(displayConfig *displayConfig,
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
	
	table0ID := displayConfig.dstAppConfig.tableNameIDPairs[table0]
	table1ID := displayConfig.dstAppConfig.tableNameIDPairs[table1]

	// First, we need to get the attribute that requires reference resolution
	// For example, we have *account.id*, and we want to get *users.account_id*
	// We check whether account.id needs to be resolved
	if reference_resolution.NeedToResolveReference(displayConfig.refResolutionConfig, 
		table0, col0) {

		log.Println("Before checking reference1 resolved or not")

		// If account.id should be resolved (in this case, it should not),
		// we check whether the reference has been resolved or not
		newVal := reference_resolution.ReferenceResolved(displayConfig.refResolutionConfig, 
			table0ID, col0, id)
		
		// If the reference has been resolved, then use the new reference to get data
		if newVal != "" {

			log.Println("reference has been resolved")

			return getOneRowBasedOnDependency(displayConfig, table1, col1, newVal)
		
		// Otherwise, we try to resolve the reference
		} else {

			hint0 := CreateHint(table0, table0ID, id)
			log.Println("Before resolving reference1: ", hint0)

			ID0 := hint0.TransformHintToIdenity(displayConfig)

			reference_resolution.ResolveReference(
				displayConfig.refResolutionConfig, ID0)
			
			// Here we check again to get updated attributes and values
			// instead of using the returned values from the ResolveReference
			// because ResolveReference only returns the updated values in that
			// function call. Values could be updated by other threads and in this
			// case, ResolveReference does not return the updated attribute and value
			// Therefore, we check all updated attributes again here by calling 
			// GetUpdatedAttributes
			updatedAttrs := reference_resolution.GetUpdatedAttributes(
				displayConfig.refResolutionConfig,
				ID0,
			)

			log.Println("Updated attributes and values:")
			log.Println(updatedAttrs)
			
			// We check whether the desired attr (col0) has been resolved
			foundResolvedAttr := false
			for attr, val := range updatedAttrs {
				if attr == col0 {
					newVal = val
					foundResolvedAttr = true
					break
				}
			}

			// If we find that col0 has been resolved, then we can use it to get other data
			if foundResolvedAttr {

				return getOneRowBasedOnDependency(displayConfig, table1, col1, newVal)
			
			// Otherwise we cannot use the unresolved reference to get other data in node
			} else {

				return nil, CannotResolveReferencesGetDataInNode
			}
		}

	// We check if users.account_id needs to be resolved (of course, in this case, it should be)
	// However we don't know its id (this is the differece from the above case!!). 
	// Also if the value is the value of "id", this could not be used directly 
	// We decide to make so much efforts to resolve "backwards" because inner-dependencies like
	// "statuses.id":"statuses.status_id"
	// "statuses.id":"mentions.status_id"
	// "statuses.id":"stream_entries.activity_id"
	// force us to do in this way. Otherwise, we cannot get other data in a node through statuses.id
	} else if fromAttrsfirstArg := schema_mappings.GetFirstArgsInREFByToTableToAttr(
		displayConfig.mappingsFromSrcToDst, table1, col1); len(fromAttrsfirstArg) != 0 {

		log.Println("Before checking reference2 resolved or not")
		
		var data1, data2 []map[string]interface{}
		
		var err error

		// First we must assume that it has already been resolved. If it has not been resolved,
		// then we cannot get data. Otherwise we just return the obtained data
		// Note that if we first assume that it has not been resolved and get data using 
		// prevID, then we could get wrong results
		data1, err = getRowsBasedOnDependency(displayConfig, table1, col1, value)
		
		// This could happen when table1 and col1 have been resolved
		if err == nil {

			log.Println("Before checking reference2, the reference seems to have already been resolved")
			
			// Now we have not encountered data1 with more than one piece of data
			for _, data4 := range data1 {

				resolvedVal := reference_resolution.ReferenceResolved(displayConfig.refResolutionConfig, 
					table1ID, col1, fmt.Sprint(data4["id"]))
				
				if resolvedVal == value {
					log.Println("It was indeed resolved")
					return data4, nil 
				
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

		}
		
		// When we reach here, we have to resolve the reference
		log.Println("From attributes:")
		log.Println(fromAttrsfirstArg)

		var prevID string

		fromAttrfirstArgContainID := false

		// Because we only use toTable and toAttr to get the first argument in fromAttrs,
		// there could be multiple results. For example, 1. toTable = status_stats and
		// toAttr = status_id, then fromAttr could be posts.id, comments.id, or messages.id
		// 2. toTable = users and toAttr = accounts.id. then fromAttr is people.id
		// We need to use from attr in the following check because
		// otherwise the fromAttr could be profile, people, and users
		for fromAttrfirstArg, _ := range fromAttrsfirstArg {

			log.Println("Check a from attribute:", fromAttrfirstArg)

			// log.Println("Check a from attribute:")

			// If the first argument of the from attribute contains "id", this indicates 
			// we need to get the *original id* as the value to get the data
			// (the current id is the newly generated one)
			if doesArgAttributeContainID(fromAttrfirstArg) {

				fromAttrfirstArgContainID = true

				// log.Println(displayConfig.dstAppConfig.tableNameIDPairs)
				// log.Println(table0)
				// log.Println(table0ID)
				
				dataID := reference_resolution.CreateIdentity(
					displayConfig.dstAppConfig.appID,
					table0ID,
					id,
				)
				
				log.Println("dataID:")
				log.Println(dataID)

				tableInFirstArg := getTableInArg(fromAttrfirstArg)
				srcTableID := displayConfig.srcAppConfig.tableNameIDPairs[tableInFirstArg]
	
				prevID = reference_resolution.GetPreviousID(displayConfig.refResolutionConfig, 
					dataID, srcTableID)
				
				log.Println("Previous id:", prevID)

				// since there is only one mapping to this toAttr, as long as we find one, 
				// we can set the value as the prevID
				if prevID == "" {
					continue
				} else {
					// value = prevID
					// break

					data2, err = getRowsBasedOnDependency(displayConfig, table1, col1, prevID)
					if err != nil {
						log.Println("The first argument of the from attribute contains id, but")
						log.Println(err)	
					} 
					break
				}
			}
		}

		// If the first argument of the from attribute does not contain "id", 
		// this indicates we can use the current data and the relationship indicated
		// by table0, col0, table1, col1, value to get data
		if !fromAttrfirstArgContainID {

			log.Println(`The from attributes don't contain id`)

			data2, err = getRowsBasedOnDependency(displayConfig, table1, col1, value)
			// This could happen when no data is migrated or there is no mappings.
			// For example, statuses, id, mentions, status_id
			if err != nil {
				return nil, err
			}

		}

		// log.Println("fromAttrfirstArgContainID:", fromAttrfirstArgContainID)
		// log.Println("data:", data)

		// If the first argument of the from attribute contains id and
		// we cannot get data, there could be two cases:
		// 1. The reference has been resolved, so the data contains the up-to-date value
		// 2. The reference has been resolved, but reference resolution crashes before
		// inserting the reference into the resolution resolved table
		// In both cases, try to get data with the new value 
		// For 1, it will be checked afterwards
		// For 2, do the reference resolution again since it does not matter and in the
		// second time, we can remove the reference and add it to the resolved resolution table
		// There is another case when prevID is "" mentioned below
		if fromAttrfirstArgContainID {
			
			// This is the one of the most strange cases found in tests
			// I guess this is because the row is first inserted into display_flags table, 
			// but not inserted into the identity table yet, the display thread can get and check
			// the row in the display_flags table, but cannot find the previous id. 
			if prevID == "" {

				return nil, CannotGetPrevID
			}
			
			// This could happen when another display thread had resolved the ref
			// before we tried to get the data using prevID
			if data2 == nil {

				log.Println(`The from attributes contain id but we cannot get data,
					so we try to get data with the current id value`)

				data2, err = getRowsBasedOnDependency(displayConfig, table1, col1, value)
				// This could happen when the resolved and displayed data is deleted
				if err != nil {
					return nil, err
				}
			}
			
		}

		// Now we have the id of the data, we should check whether it has been resolved before, 
		// but actually if we can get one, it is highly likely that this is the one we want to get because
		// otherwise there will be multiple rows corresponding to one member.
		// There could be the case where ids are not changed. Even if references are not resolved, 
		// we can still get the rows we want, but we need to resolve it.
		// There could be a case where the data we got is from some other unrelated migrations because
		// we use old value (in source app) to get data. In this case, this old value 
		// can be used to get more than one piece of data including the one we want to get
		for _, data3 := range data2 {
 
			newVal := reference_resolution.ReferenceResolved(displayConfig.refResolutionConfig, 
				table1ID, col1, fmt.Sprint(data3["id"]))

			// The reference has been resolved
			if newVal != "" {

				// Theoretically, if it has been resolved, then it should be the value we have 
				// given that one member corresponds to one row
				if newVal == value {

					log.Println("reference has been resolved")
					log.Println(data3)

					return data3, nil

				} else {

					log.Println("Found an unrelated data:")
					log.Println(data3)
				
				}
			} else {
				
				hint1 := CreateHint(table1, table1ID, fmt.Sprint(data3["id"]))
				log.Println("Before resolving reference2:", hint1)

				ID1 := hint1.TransformHintToIdenity(displayConfig)

				reference_resolution.ResolveReference(
					displayConfig.refResolutionConfig, ID1)
				
				updatedAttrs := reference_resolution.GetUpdatedAttributes(
					displayConfig.refResolutionConfig,
					ID1,
				)

				log.Println("Updated attributes and values:")
				log.Println(updatedAttrs)

				// We check whether the desired attr (col1) has been resolved 
				// (until this point, it should be resolved)
				foundResolvedAttr := false
				for attr, val := range updatedAttrs {
					if attr == col1 {
						newVal = val
						foundResolvedAttr = true
						break
					}
				}

				// If we find that col0 has been resolved, then we can use it to get other data
				if foundResolvedAttr {

					if newVal == value {

						return data3, nil
					
					// This can happen when we are trying to resolve the unrelated data
					} else {
						
						log.Println(ID1)
						log.Println("newVal", newVal)
						log.Println("value", value)
						// panic(`Find the resolved attribute, but the value is not what we want. 
						// 	Should not happen given one member corresponds to one row for now!`)
						log.Println(`Find the resolved attribute, but the value is not what we want. 
							This is because we happened to get an unrelated but satisfying data`)
						
					}
				
				// This should not happen
				} else {
					
					// return nil, CannotFindResolvedAttributes
					panic(`Does not find resolved attributes. Should not happen 
						given one member corresponds to one row for now!`)
				}
			}
		}

		panic(`It should never happen since there should be one piece of data which is what we want!`)
	
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

func oldCheckResolveReferenceInGetDataInNode(displayConfig *displayConfig,
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
	
	table0ID := displayConfig.dstAppConfig.tableNameIDPairs[table0]
	table1ID := displayConfig.dstAppConfig.tableNameIDPairs[table1]

	// First, we need to get the attribute that requires reference resolution
	// For example, we have *account.id*, and we want to get *users.account_id*
	// We check whether account.id needs to be resolved
	if reference_resolution.NeedToResolveReference(displayConfig.refResolutionConfig, 
		table0, col0) {

		log.Println("Before checking reference1 resolved or not")

		// If account.id should be resolved (in this case, it should not),
		// we check whether the reference has been resolved or not
		newVal := reference_resolution.ReferenceResolved(displayConfig.refResolutionConfig, 
			table0ID, col0, id)
		
		// If the reference has been resolved, then use the new reference to get data
		if newVal != "" {

			log.Println("reference has been resolved")

			return getOneRowBasedOnDependency(displayConfig, table1, col1, newVal)
		
		// Otherwise, we try to resolve the reference
		} else {

			hint0 := CreateHint(table0, table0ID, id)
			log.Println("Before resolving reference1: ", hint0)

			ID0 := hint0.TransformHintToIdenity(displayConfig)

			updatedAttrs, _ := reference_resolution.ResolveReference(
				displayConfig.refResolutionConfig, ID0)

			// We check whether the desired attr (col0) has been resolved
			foundResolvedAttr := false
			for attr, val := range updatedAttrs {
				if attr == col0 {
					newVal = val
					foundResolvedAttr = true
					break
				}
			}

			// If we find that col0 has been resolved, then we can use it to get other data
			if foundResolvedAttr {

				return getOneRowBasedOnDependency(displayConfig, table1, col1, newVal)
			
			// Otherwise we cannot use the unresolved reference to get other data in node
			} else {

				return nil, CannotResolveReferencesGetDataInNode
			}
		}

	// We check if users.account_id needs be resolved (of course, in this case, it should be)
	// However we don't know its id (this is the differece from the above case!!). 
	// Also if the value is the value of "id", this could not be used directly 
	} else if fromAttrsfirstArg := schema_mappings.GetFirstArgsInREFByToTableToAttr(
		displayConfig.mappingsFromSrcToDst, table1, col1); len(fromAttrsfirstArg) != 0 {

		log.Println("Before checking reference2 resolved or not")
		
		log.Println("From attributes:")
		log.Println(fromAttrsfirstArg)
		
		data := make(map[string]interface{})
		var err error

		fromAttrfirstArgContainID := false

		// Because we only use toTable and toAttr to get the first argument in fromAttrs,
		// there could be multiple results. For example, toTable = status_stats and
		// toAttr = status_id, then fromAttr could be posts.id, comments.id, or messages.id
		for fromAttrfirstArg, _ := range fromAttrsfirstArg {

			// log.Println("Check a from attribute:", fromAttrfirstArg)

			log.Println("Check a from attribute:")

			// If the first argument of the from attribute contains "id", this indicates 
			// we need to get the original id as the value to get the data
			// (the current id is the newly generated one)
			if doesArgAttributeContainID(fromAttrfirstArg) {

				fromAttrfirstArgContainID = true

				// log.Println(displayConfig.dstAppConfig.tableNameIDPairs)
				// log.Println(table0)
				// log.Println(table0ID)
				
				dataID := reference_resolution.CreateIdentity(
					displayConfig.dstAppConfig.appID,
					table0ID,
					id,
				)
				
				tableInFirstArg := getTableInArg(fromAttrfirstArg)
				srcTableID := displayConfig.srcAppConfig.tableNameIDPairs[tableInFirstArg]
	
				prevID := reference_resolution.GetPreviousID(displayConfig.refResolutionConfig, 
					dataID, srcTableID)
				
				log.Println("Previous id:", prevID)

				// since there is only one mapping to this toAttr, as long as we find one, 
				// we can set the value as the prevID
				if prevID == "" {
					continue
				} else {
					// value = prevID
					// break

					data, err = getOneRowBasedOnDependency(displayConfig, table1, col1, prevID)
					if err != nil {
						log.Println("The first argument of the from attribute contains id, but")
						log.Println(err)	
					} 
					break
				}
			}
		} 

		// If the first argument of the from attribute does not contain "id", 
		// this indicates we can use the current data and the relationship indicated
		// by table0, col0, table1, col1, value to get data
		if !fromAttrfirstArgContainID {

			log.Println(`The from attributes don't contain id`)

			data, err = getOneRowBasedOnDependency(displayConfig, table1, col1, value)
			// This could happen when no data is migrated or there is no mappings.
			// For example, statuses, id, mentions, status_id
			if err != nil {
				return nil, err
			}

		}

		// If the first argument of the from attribute contains id and
		// we cannot get data, there could be two cases:
		// 1. The reference has been resolved, so the data contains the up-to-date value
		// 2. The reference has been resolved, but reference resolution crashes before
		// inserting the reference into the resolution resolved table
		// In both cases, try to get data with the new value 
		// For 1, it will be checked afterwards
		// For 2, do the reference resolution again since it does not matter and in the
		// second time, we can remove the reference and add it to the resolved resolution table
		if fromAttrfirstArgContainID && data == nil {

			log.Println(`The from attributes contain id but we cannot get data,
				so we try to get data with the current id value`)

			data, err = getOneRowBasedOnDependency(displayConfig, table1, col1, value)
			// This could happen when the resolved and displayed data is deleted
			if err != nil {
				return nil, err
			}
			
		}

		// Now we have the id of the data, we should check whether it has been resolved before, 
		// but actually if we can get one, it should always be the one we want to get because
		// otherwise there will be multiple rows corresponding to one member.
		// There could be the case where ids are not changed, 
		// so even if references are not resolved, 
		// we can still get the rows we want, but we need to resolve it further.
		newVal := reference_resolution.ReferenceResolved(displayConfig.refResolutionConfig, 
			table1ID, col1, fmt.Sprint(data["id"]))

		// The reference has been resolved
		if newVal != "" {

			// Theoretically, if it has been resolved, then it should be the value we have 
			// given that one member corresponds to one row
			if newVal == value {

				log.Println("reference has been resolved")
				log.Println(data)

				return data, nil

			} else {

				panic(`The reference has been resolved, but the value is not what we want. 
					Should not happen given one member corresponds to one row for now!`)
			
			}
		} else {
			
			hint1 := CreateHint(table1, table1ID, fmt.Sprint(data["id"]))
			log.Println("Before resolving reference2:", hint1)

			ID1 := hint1.TransformHintToIdenity(displayConfig)

			updatedAttrs, _ := reference_resolution.ResolveReference(
				displayConfig.refResolutionConfig, ID1)
			
			log.Println("Updated attributes:")
			log.Println(updatedAttrs)

			// We check whether the desired attr (col1) has been resolved 
			// (until this point, it should be resolved)
			foundResolvedAttr := false
			for attr, val := range updatedAttrs {
				if attr == col1 {
					newVal = val
					foundResolvedAttr = true
					break
				}
			}

			// If we find that col0 has been resolved, then we can use it to get other data
			if foundResolvedAttr {

				if newVal == value {

					return data, nil
	
				} else {

					log.Println("newVal", newVal)
					log.Println("value", value)
					panic(`Find the resolved attribute, but the value is not what we want. 
						Should not happen given one member corresponds to one row for now!`)
				}
			
			// This should not happen
			} else {
				
				// return nil, CannotFindResolvedAttributes
				panic(`Does not find resolved attributes. Should not happen 
					given one member corresponds to one row for now!`)
			}
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

func oldGetHintsInParentNode(displayConfig *displayConfig, 
	hints []*HintStruct, conditions []string, pTag string) (*HintStruct, error) {
	
	query := fmt.Sprintf("SELECT %s.* FROM ", "t"+strconv.Itoa(len(conditions)))
	from := ""
	table := ""
	hintID := -1

	for i, condition := range conditions {

		tableAttr1 := strings.Split(condition, ":")[0]
		tableAttr2 := strings.Split(condition, ":")[1]

		t1 := strings.Split(tableAttr1, ".")[0]
		a1 := strings.Split(tableAttr1, ".")[1]

		t2 := strings.Split(tableAttr2, ".")[0]
		a2 := strings.Split(tableAttr2, ".")[1]

		seq1 := "t" + strconv.Itoa(i)
		seq2 := "t" + strconv.Itoa(i+1)

		if i == 0 {

			// There could be mutliple pieces of data in nodes
			// For example:
			// A statuses node contains status, conversation, and status_stats
			for j, hint := range hints {

				if hint.Table == t1 {
					hintID = j
				}

			}

			// In this case, since data may be incomplete, 
			// we cannot get the data in the parent node
			if hintID == -1 {

				return nil, common_funcs.CannotFindAnyDataInParent

			} else {
				
				// For example:
				// if a condition is [favourites.status_id:statuses.id], 
				// from will be "favourites t0 JOIN statuses t1 ON t0.status_id = t1.id"
				from += fmt.Sprintf("%s %s JOIN %s %s ON %s.%s = %s.%s ",
					t1, seq1, t2, seq2, seq1, a1, seq2, a2)

			}

		// This is mainly to solve the case in which
		// conversation cannot directly depend on root
		// conversation depends on statuses, which in turn depends on root. 
		// This is now obsolete because there is no dependency between other nodes with root
		// For now, there is always only one condition.
		} else {

			from += fmt.Sprintf("JOIN %s %s on %s.%s = %s.%s ",
				t2, seq2, seq1, a1, seq2, a2)
			

		}

		//The last condition
		if i == len(conditions)-1 {

			var depDataKey string
			var depDataValue int

			for k, v := range hints[hintID].KeyVal {

				depDataKey = k
				depDataValue = v

			}

			// Following the above example,
			// the whole query will be:
			// SELECT t1.* 
			// FROM favourites t0 JOIN statuses t1 ON t0.status_id = t1.id
			// WHERE t0.status_id = 80
			where := fmt.Sprintf("WHERE %s.%s = %d", "t0", depDataKey, depDataValue)
			table = t2
			query += from + where

		}
	}
	// fmt.Println(query)

	data, err := db.DataCall1(displayConfig.dstAppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	if len(data) == 0 {

		return nil, common_funcs.CannotFindAnyDataInParent

	} else {

		return TransformRowToHint(displayConfig, data, table, pTag), nil

	}
}