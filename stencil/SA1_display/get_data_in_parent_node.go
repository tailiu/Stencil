package SA1_display

import (
	"fmt"
	"log"
	"stencil/db"
	"stencil/common_funcs"
	"stencil/reference_resolution"
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
	
	tableID := displayConfig.dstAppConfig.tableNameIDPairs[table]

	// First, we need to get the attribute that requires reference resolution
	// For example, we have *favourites.status_id*, and we want to get *status*
	// We check whether favourites.status_id needs to be resolved
	if reference_resolution.NeedToResolveReference(displayConfig.refResolutionConfig, table, col) {

		log.Println("Parent Node: before checking reference resolved or not")

		// If favourites.status_id should be resolved (in this case, it should be),
		// we check whether the reference has been resolved or not
		newVal := reference_resolution.ReferenceResolved(displayConfig.refResolutionConfig, 
			tableID, col, id)
		
		// If the reference has been resolved, then use the new reference to get data
		if newVal != "" {

			log.Println("Parent Node: reference has been resolved")

			return newVal, nil
		
		// Otherwise, we try to resolve the reference
		} else {

			hint := CreateHint(table, tableID, id)
			log.Println("Parent Node: before resolving reference: ", hint)

			ID := hint.TransformHintToIdenity(displayConfig)

			reference_resolution.ResolveReference(displayConfig.refResolutionConfig, ID)

			updatedAttrs := reference_resolution.GetUpdatedAttributes(
				displayConfig.refResolutionConfig,
				ID,
			)

			log.Println("Updated attributes and values:")
			log.Println(updatedAttrs)

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

				return "", CannotResolveReferencesGetDataInParentNode
			}
		}
	
	// Normally, there must exist one that needs to be resolved. 
	// But this could happen for example, in Diaspora, posts.id depends on aspects.shareable_id
	// There is no need to resolve id here
	} else {

		// panic("Should not happen since there is always one to solve!")

		return "", NoReferenceToResolve

	}
}

func getHintInParentNode(displayConfig *displayConfig, 
	hints []*HintStruct, conditions []string, 
	pTag string) (*HintStruct, error) {
	
	// log.Println(hints[0])

	var data map[string]interface{}
	var err0, err1 error
	var table string
	var depVal string

	hintID := -1

	for i, condition := range conditions {

		// log.Println(condition)

		tableAttr1 := strings.Split(condition, ":")[0]
		tableAttr2 := strings.Split(condition, ":")[1]

		// log.Println("processing conditions")
		// log.Println(tableAttr1)
		// log.Println(tableAttr2)

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

				if displayConfig.resolveReference {

					depVal, err0 = checkResolveReferenceInGetDataInParentNode(
						displayConfig, 
						fmt.Sprint(hints[hintID].Data["id"]),
						t1, a1)
					
					// no matter whether this attribute has been resolved before
					// we need to refresh the cached data because this attribute might be
					// resolved by other thread checking other data
					refreshCachedDataHints(displayConfig, hints)

					// If there is an error, it means that the reference has not been resolved
					if err0 != nil {
						
						log.Println(err0)

						if err0 != NoReferenceToResolve {
							return nil, CannotFindAnyDataInParent
						} else {
							depVal = fmt.Sprint(hints[hintID].Data[a1])
						}
					}

				} else {

					depVal = fmt.Sprint(hints[hintID].Data[a1])

				}

				var query string

				if !displayConfig.markAsDelete {
					query = fmt.Sprintf(
						`SELECT * FROM "%s" WHERE %s = '%s'`, 
						t2, a2, depVal,
					)
				} else {
					query = fmt.Sprintf(
						`SELECT * FROM "%s" WHERE %s = '%s' and mark_as_delete = false`, 
						t2, a2, depVal,
					)
				}

				data, err1 = db.DataCall1(displayConfig.dstAppConfig.DBConn, query)
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

			if displayConfig.resolveReference {

				depVal, err0 = checkResolveReferenceInGetDataInParentNode(
					displayConfig, 
					fmt.Sprint(data["id"]),
					t1, a1)
				
				// no matter whether this attribute has been resolved before
				// we need to refresh the cached data because this attribute might be
				// resolved by other thread checking other data
				refreshCachedDataHints(displayConfig, hints)

				// If there is an error, it means that the reference has not been resolved
				// so it cannot be used to get data in the parent node
				if err0 != nil {
					
					log.Println(err0)

					if err0 != NoReferenceToResolve {
						return nil, CannotFindAnyDataInParent
					} else {
						depVal = fmt.Sprint(data[a1])
					}
				}

			} else {

				depVal = fmt.Sprint(data[a1])

			}

			var query string

			if !displayConfig.markAsDelete {
				query = fmt.Sprintf(
					`SELECT * FROM "%s" WHERE %s = '%s'`, 
					t2, a2, depVal,
				)
			} else {
				query = fmt.Sprintf(
					`SELECT * FROM "%s" WHERE %s = '%s' and mark_as_delete = false`, 
					t2, a2, depVal,
				)
			}

			data, err1 = db.DataCall1(displayConfig.dstAppConfig.DBConn, query)
			if err1 != nil {
				log.Fatal(err1)
			}

			if len(data) == 0 {
				return nil, CannotFindAnyDataInParent
			}

			table = t2
		}
	}

	// log.Println("...........")
	// log.Println(table)
	// log.Println(data)
	// log.Println("...........")

	return TransformRowToHint(displayConfig, data, table, pTag), nil

}

func dataFromParentNodeExists(displayConfig *displayConfig,
	hints []*HintStruct, pTag string) (bool, error) {
	
	log.Println("check dataFromParentNodeExists")

	displayExistenceSetting, _ := hints[0].GetDisplayExistenceSetting(displayConfig, pTag)

	// If display existence setting is not set, 
	// then we have to try to get data in the parent node in any case
	if displayExistenceSetting == "" {

		return true, nil

	} else {

		tableCol := common_funcs.ReplaceKey(displayConfig.dstAppConfig.dag, 
			hints[0].Tag, displayExistenceSetting)
		table := strings.Split(tableCol, ".")[0]
		col := strings.Split(tableCol, ".")[1]

		// log.Println(tableCol)

		for _, hint := range hints {

			if hint.Table == table {
				
				log.Println(hint.Data)
				log.Println(tableCol)
				log.Println(hint.Data[col])

				if hint.Data[col] == nil {

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
	conditions, _ := common_funcs.GetDependsOnConditionsInDeps(displayConfig.dstAppConfig.dag, tag, pTag)
	pTag, _ = hints[0].GetOriginalTagNameFromAliasOfParentTagIfExists(displayConfig, pTag)

	// log.Println("conditions")
	// log.Println(conditions)

	var procConditions []string
	var from, to string

	if len(conditions) == 1 {

		condition := conditions[0]
		from = common_funcs.ReplaceKey(displayConfig.dstAppConfig.dag, tag, condition.TagAttr)
		to = common_funcs.ReplaceKey(displayConfig.dstAppConfig.dag, pTag, condition.DependsOnAttr)
		procConditions = append(procConditions, from+":"+to)

	} else {

		for i, condition := range conditions {

			if i == 0 {

				from = common_funcs.ReplaceKey(displayConfig.dstAppConfig.dag, tag, condition.TagAttr)

				to = common_funcs.ReplaceKey(displayConfig.dstAppConfig.dag,
					strings.Split(condition.DependsOnAttr, ".")[0], 
					strings.Split(condition.DependsOnAttr, ".")[1])

			} else if i == len(conditions)-1 {

				from = common_funcs.ReplaceKey(displayConfig.dstAppConfig.dag, 
					strings.Split(condition.TagAttr, ".")[0], 
					strings.Split(condition.TagAttr, ".")[1])
				
				to = common_funcs.ReplaceKey(displayConfig.dstAppConfig.dag, pTag, condition.DependsOnAttr)

			}

			procConditions = append(procConditions, from+":"+to)

		}

	}

	// fmt.Println(procConditions)
	// fmt.Println(hints)

	return getHintInParentNode(displayConfig, hints, procConditions, pTag)

}
