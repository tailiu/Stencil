package SA1_display

import (
	"fmt"
	"log"
	"stencil/db"
	"stencil/reference_resolution"
	"strconv"
	"strings"
)

// The only big difference between checkResolveReferenceInGetDataInParentNode 
// and getOneRowBasedOnDependency is that it only needs to check one table and col
func checkResolveReferenceInGetDataInParentNode(displayConfig *displayConfig,
	id, table, col string) (string, error) {

	log.Println("+++++++++++++++++++")
	log.Println(table)
	log.Println(col)
	log.Println("+++++++++++++++++++")
	
	tableID := displayConfig.AppConfig.TableNameIDPairs[table]

	// First, we need to get the attribute that requires reference resolution
	// For example, we have *favourites.status_id*, and we want to get *status*
	// We check whether favourites.status_id needs to be resolved
	if reference_resolution.NeedToResolveReference(displayConfig, table, col) {

		log.Println("Parent Node: before checking reference resolved or not")

		// If favourites.status_id should be resolved (in this case, it should be),
		// we check whether the reference has been resolved or not
		newVal := reference_resolution.ReferenceResolved(displayConfig, tableID, col, id)
		
		// If the reference has been resolved, then use the new reference to get data
		if newVal != "" {

			log.Println("Parent Node: reference has been resolve")

			return newVal, nil
		
		// Otherwise, we try to resolve the reference
		} else {

			hint := CreateHint(table, tableID, id)
			log.Println("Parent Node: Before resolving reference: ", hint)

			ID := hint.TransformHintToIdenity(displayConfig)

			updatedAttrs, _ := reference_resolution.ResolveReference(displayConfig, ID)

			// We check whether the desired attr (col) has been resolved
			foundResolvedAttr := false
			for attr, val := range updatedAttrs {
				if attr == col {
					newVal = val
					foundResolvedAttr = true
					break
				}
			}

			// If we find that col has been resolved, then we can use it to get other data
			if foundResolvedAttr {

				return newVal, nil
			
			// Otherwise we cannot use the unresolved reference to get other data in node
			} else {

				return "", CannotResolveReferencesGetDataInNode
			}
		}
	
	// Normally, there must exist one that needs to be resolved. 
	} else {

		panic("Should not happen since there is always one to solve!")

	}
}

func getHintsInParentNode(displayConfig *displayConfig, 
	hints []*HintStruct, conditions []string, pTag string) (*HintStruct, error) {
	
	log.Println(hints[0])

	var data map[string]interface{}
	var err0, err1 error
	var table string
	var depVal string

	hintID := -1

	for i, condition := range conditions {

		// log.Println(condition)

		tableAttr1 := strings.Split(condition, ":")[0]
		tableAttr2 := strings.Split(condition, ":")[1]

		t1 := strings.Split(tableAttr1, ".")[0]
		a1 := strings.Split(tableAttr1, ".")[1]

		t2 := strings.Split(tableAttr2, ".")[0]
		a2 := strings.Split(tableAttr2, ".")[1]

		// log.Println(t1, a1, t2, a2)

		if i == 0 {

			// There could be mutliple pieces of data in nodes
			// For example:
			// A statuses node contains status, conversation, and status_stats
			for j, hint := range hints {

				if hint.Table == t1 {
					hintID = j
				}

			}

			if hintID == -1 {

				// In this case, since data may be incomplete, 
				// we cannot get the data in the parent node
				return nil, CannotFindAnyDataInParent
			
			} else {

				// This can happen when the data this data depends on is not migrated,
				// e.g., a post does not have correpsonding conversation in Diaspora, 
				// so when it is migrated to Mastodon, and becomes a status, 
				// it does not have conversation_id,  
				// which is actually necessary for each status in Mastodon.
				if hints[hintID].Data[a1] == nil {
					return nil, CannotFindAnyDataInParent
				}

				if displayConfig.ResolveReference {

					depVal, err0 = checkResolveReferenceInGetDataInParentNode(
						displayConfig, 
						fmt.Sprint(hints[hintID].Data["id"]),
						t1, a1)
					
					// If there is an error, it means that the reference has not been resolved
					if err0 != nil {
						return nil, err0
					}

				} else {

					depVal = fmt.Sprint(hints[hintID].Data[a1])

				}

				query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", t2, a2, depVal)

				data, err1 = db.DataCall1(displayConfig.AppConfig.DBConn, query)
				if err1 != nil {
					log.Fatal(err1)
				}
			
				// log.Println(".....first check......")
				// log.Println(data)
				// log.Println("...........")

				if len(data) == 0 {
					return nil, CannotFindAnyDataInParent
				}

				table = t2
			}

		// This is mainly to solve the case in which
		// conversation cannot directly depend on root
		// conversation depends on statuses, which in turn depends on root. 
		// This is now obsolete because there is no dependency between other nodes with root
		// For now, there is always only one condition.
		} else {

			if displayConfig.ResolveReference {

				depVal, err0 = checkResolveReferenceInGetDataInParentNode(
					displayConfig, 
					fmt.Sprint(data["id"]),
					t1, a1)
				
				// If there is an error, it means that the reference has not been resolved
				// so it cannot be used to get data in the parent node
				if err0 != nil {
					return nil, err0
				}

			} else {

				depVal = fmt.Sprint(data[a1])

			}

			query := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", 
				t2, a2, depVal)

			data, err1 = db.DataCall1(displayConfig.AppConfig.DBConn, query)
			if err1 != nil {
				log.Fatal(err1)
			}

			if len(data) == 0 {
				return nil, CannotFindAnyDataInParent
			}

			table = t2
		}
	}

	log.Println("...........")
	log.Println(table)
	log.Println(data)
	log.Println("...........")

	return TransformRowToHint(displayConfig, data, table, pTag), nil

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

	data, err := db.DataCall1(displayConfig.AppConfig.DBConn, query)
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

func dataFromParentNodeExists(displayConfig *displayConfig,
	hints []*HintStruct, pTag string) (bool, error) {
	
	displayExistenceSetting, _ := hints[0].GetDisplayExistenceSetting(displayConfig, pTag)

	// If display existence setting is not set, 
	// then we have to try to get data in the parent node in any case
	if displayExistenceSetting == "" {

		return true, nil

	} else {

		tableCol := ReplaceKey(displayConfig, hints[0].Tag, displayExistenceSetting)
		table := strings.Split(tableCol, ".")[0]

		for _, hint := range hints {

			if hint.Table == table {

				if hint.Data[tableCol] == nil {

					return false, NotDependsOnAnyData

				} else {

					return true, nil

				}
			}
		}

	}

	// In this case, since data may be incomplete, 
	// we cannot find the existence of the data in a parent node
	// This also implies that it cannot find any data in a parent node
	return false, CannotFindAnyDataInParent

}

// Note: this function may return multiple hints based on dependencies
func GetdataFromParentNode(displayConfig *displayConfig,
	hints []*HintStruct, pTag string) (*HintStruct, error) {

	// Before getting data from a parent node, 
	// we check the existence of the data based on the cols of a child node
	if exists, err := dataFromParentNodeExists(displayConfig, hints, pTag); !exists {
		return nil, err
	}

	tag := hints[0].Tag
	conditions, _ := displayConfig.AppConfig.GetDependsOnConditions(tag, pTag)
	pTag, _ = hints[0].GetOriginalTagNameFromAliasOfParentTagIfExists(displayConfig, pTag)

	var procConditions []string
	var from, to string

	if len(conditions) == 1 {

		condition := conditions[0]
		from = ReplaceKey(displayConfig, tag, condition.TagAttr)
		to = ReplaceKey(displayConfig, pTag, condition.DependsOnAttr)
		procConditions = append(procConditions, from+":"+to)

	} else {

		for i, condition := range conditions {

			if i == 0 {

				from = ReplaceKey(displayConfig, tag, condition.TagAttr)

				to = ReplaceKey(displayConfig,
					strings.Split(condition.DependsOnAttr, ".")[0], 
					strings.Split(condition.DependsOnAttr, ".")[1])

			} else if i == len(conditions)-1 {

				from = ReplaceKey(displayConfig, 
					strings.Split(condition.TagAttr, ".")[0], 
					strings.Split(condition.TagAttr, ".")[1])
				
				to = ReplaceKey(displayConfig, pTag, condition.DependsOnAttr)

			}

			procConditions = append(procConditions, from+":"+to)

		}

	}

	// fmt.Println(procConditions)
	// fmt.Println(hints)

	return getHintsInParentNode(displayConfig, hints, procConditions, pTag)

}
