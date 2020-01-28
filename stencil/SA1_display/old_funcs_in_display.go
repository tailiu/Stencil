package SA1_display

import (
	"fmt"
	"log"
	"strconv"
	"stencil/db"
	"stencil/reference_resolution"
	"stencil/schema_mappings"
	"strings"
)

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

				return nil, CannotFindAnyDataInParent

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

		return nil, CannotFindAnyDataInParent

	} else {

		return TransformRowToHint(displayConfig, data, table, pTag), nil

	}
}