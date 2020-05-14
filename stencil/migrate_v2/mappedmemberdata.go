package migrate_v2

import (
	"fmt"
	"log"
	config "stencil/config/v2"
	"stencil/db"
	"stencil/helper"
	"strings"
)

func (mmd MappedMemberData) GetQueryArgs() (string, string, []interface{}) {
	var colStr, phStr string
	var valList []interface{}
	counter := 1

	for mappedAttr, mmv := range mmd.Data {
		colStr += fmt.Sprintf("\"%s\",", mappedAttr)
		phStr += fmt.Sprintf("$%v,", counter)
		valList = append(valList, mmv.Value)
		counter++
	}

	colStr = strings.Trim(colStr, ",")
	phStr = strings.Trim(phStr, ",")

	return colStr, phStr, valList
}

func (mmd MappedMemberData) ValidateMappedData() bool {
	for mappedAttr, mmv := range mmd.Data {
		if mmv.IsInput || mmv.IsExpression || mappedAttr == "id" {
			continue
		}
		if mmv.Value != nil {
			return true
		}
	}
	return false
}

func (mmd *MappedMemberData) SetMember(table string) {
	mmd.ToMember = table
	if mmd.DBConn == nil {
		log.Fatal("@mmd.SetMember: DBConn not set!")
	} else {
		if tableID, err := db.TableID(mmd.DBConn, mmd.ToMember, mmd.ToAppID); err == nil {
			mmd.ToMemberID = tableID
		} else {
			fmt.Println(mmd.ToMember, mmd.ToAppID)
			log.Fatal("@SetMember: ", err)
		}
	}
}

func (mmd MappedMemberData) SrcTables() []Member {

	var srcTables []Member

	added := make(map[string]bool)

	for _, mmv := range mmd.Data {
		if mmv.IsExpression || mmv.IsInput {
			continue
		}
		if _, ok := added[mmv.FromMemberID]; !ok {
			srcTables = append(srcTables, Member{ID: mmv.FromMemberID, Name: mmv.FromMember})
		}
		added[mmv.FromMemberID] = true
	}

	return srcTables
}

func (mmd MappedMemberData) ToCols() []string {
	var toCols []string
	for toCol := range mmd.Data {
		toCols = append(toCols, toCol)
	}
	return toCols
}

func (mmd MappedMemberData) FromCols(table string) []string {
	var fromCols []string
	for _, mmv := range mmd.Data {
		if strings.EqualFold(mmv.FromMember, table) {
			fromCols = append(fromCols, mmv.FromAttr)
		}
	}
	return fromCols
}

func (mmd MappedMemberData) GetSourceAppsAndTables() map[string][]string {

	apptabMap := make(map[string][]string)

	for _, mmv := range mmd.Data {
		if mmv.IsExpression || mmv.IsInput {
			continue
		}
		if _, ok := apptabMap[mmv.AppID]; !ok {
			apptabMap[mmv.AppID] = []string{}
		}
		if !helper.Contains(apptabMap[mmv.AppID], mmv.FromMember) {
			apptabMap[mmv.AppID] = append(apptabMap[mmv.AppID], mmv.FromMember)
		}

	}
	return apptabMap
}

func (mmd MappedMemberData) GetDstDataMap() DataMap {

	data := make(DataMap)

	for col, mmv := range mmd.Data {
		data[mmd.ToMember+"."+col] = mmv.Value
	}

	return data
}

func (mmd MappedMemberData) GetSrcDataMap() DataMap {

	data := make(DataMap)

	for toAttr, mmv := range mmd.Data {
		if strings.EqualFold(toAttr, "id") {
			data[mmv.FromMember+"."+mmv.FromAttr] = mmv.FromID
		} else {
			data[mmv.FromMember+"."+mmv.FromAttr] = mmv.Value
		}
	}

	return data
}

func (mmd *MappedMemberData) CreateInnerDependencyReferences(appConfig config.AppConfig, tag config.Tag, nodeData DataMap) ([]MappingRef, error) {

	// log.Printf("@CreateInnerDependencyReferences | nodeTag: %s | nodeData: %v  \n", tag.Name, nodeData)

	var refs []MappingRef

	for _, innerDependency := range tag.InnerDependencies {
		for dependsOn, dependee := range innerDependency {

			depTokens := strings.Split(dependee, ".")
			dependeeMember := tag.Members[depTokens[0]]
			dependeeAttr := depTokens[1]
			dependeeReferencedAttr := fmt.Sprintf("%s.%s", dependeeMember, dependeeAttr)

			var dependeeMemberID, dependeeAttrID string
			if tableID, err := db.TableID(mmd.DBConn, dependeeMember, appConfig.AppID); err == nil {
				dependeeMemberID = tableID
			} else {
				fmt.Println(dependeeMember, appConfig.AppID)
				log.Fatal("@CreateInnerDependencyReferences.GetTableID: ", err)
			}

			if attrID, err := db.AttrID(mmd.DBConn, dependeeMemberID, dependeeAttr); err == nil {
				dependeeAttrID = attrID
			} else {
				fmt.Printf("@CreateInnerDependencyReferences.GetAttrID: Table: '%s' | Attr : '%s'\n", dependeeMemberID, dependeeAttr)
				log.Fatal("@CreateInnerDependencyReferences.GetAttrID: ", err)
			}

			depOnTokens := strings.Split(dependsOn, ".")
			depOnMember := tag.Members[depOnTokens[0]]
			depOnAttr := depOnTokens[1]
			depOnReferencedAttr := fmt.Sprintf("%s.%s", depOnMember, depOnAttr)

			var depOnMemberID, depOnAttrID string
			if tableID, err := db.TableID(mmd.DBConn, depOnMember, appConfig.AppID); err == nil {
				depOnMemberID = tableID
			} else {
				fmt.Println(depOnMember, appConfig.AppID)
				log.Fatal("@CreateInnerDependencyReferences.GetTableID: ", err)
			}

			if attrID, err := db.AttrID(mmd.DBConn, depOnMemberID, depOnAttr); err == nil {
				depOnAttrID = attrID
			} else {
				fmt.Printf("@CreateInnerDependencyReferences.GetAttrID: Table: '%s' | Attr : '%s'\n", depOnMemberID, depOnAttr)
				log.Fatal("@CreateInnerDependencyReferences.GetAttrID: ", err)
			}

			if toVal, fromVal, fromID, err := nodeData.GetRefValsFromDataMap(dependeeReferencedAttr, depOnReferencedAttr, false); err != nil {
				log.Println("@CreateInnerDependencyReferences.GetRefValsFromDataMap: ", err)
			} else {

				ref := MappingRef{
					appID:        appConfig.AppID,
					fromID:       fromID,
					fromVal:      fromVal,
					fromMemberID: dependeeMemberID,
					fromMember:   dependeeMember,
					fromAttr:     dependeeAttr,
					fromAttrID:   dependeeAttrID,
					toVal:        toVal,
					toMemberID:   depOnMemberID,
					toMember:     depOnMember,
					toAttr:       depOnAttr,
					toAttrID:     depOnAttrID,
				}

				log.Println("@CreateInnerDependencyReferences | Ref Created | ", ref)

				refs = append(refs, ref)
			}
		}
	}

	return refs, nil
}

func (mmd *MappedMemberData) CreateReferencesViaDependencies(appConfig config.AppConfig, tag config.Tag, nodeData DataMap) ([]MappingRef, error) {

	// log.Printf("@CreateReferencesViaDependencies | nodeTag: %s | nodeData: %v  \n", tag.Name, nodeData)

	var refs []MappingRef

	if dep, err := appConfig.GetDependency(tag.Name); err == nil {
		for _, depOn := range dep.DependsOn {
			if referencedTag, err := appConfig.GetTag(depOn.Tag); err == nil {

				// log.Printf("@CreateReferencesViaDependencies | referencedTag: %s  \n", referencedTag.Name)

				for _, condition := range depOn.Conditions {
					tagAttr, err := tag.ResolveTagAttr(condition.TagAttr)
					if err != nil {
						log.Println(err, tag.Name, condition.TagAttr)
						log.Fatal("@CreateReferencesViaDependencies: tagAttr in condition doesn't exist? ", condition.TagAttr)
						break
					}
					tagAttrTokens := strings.Split(tagAttr, ".")
					fromMember := tagAttrTokens[0]
					fromReference := tagAttrTokens[1]
					fromReferencedAttr := fmt.Sprintf("%s.%s", fromMember, fromReference)

					var fromMemberID, fromAttrID string
					if tableID, err := db.TableID(mmd.DBConn, fromMember, appConfig.AppID); err == nil {
						fromMemberID = tableID
					} else {
						fmt.Println(fromMember, appConfig.AppID)
						log.Fatal("@CreateReferencesViaDependencies.GetTableID: ", err)
					}

					if attrID, err := db.AttrID(mmd.DBConn, fromMemberID, fromReference); err == nil {
						fromAttrID = attrID
					} else {
						fmt.Printf("@CreateReferencesViaDependencies.GetAttrID: Table: '%s' | Attr : '%s'\n", fromMemberID, fromReference)
						log.Fatal("@CreateReferencesViaDependencies.GetAttrID: ", err)
					}

					depOnAttr, err := referencedTag.ResolveTagAttr(condition.DependsOnAttr)
					if err != nil {
						log.Println(err, referencedTag.Name, condition.DependsOnAttr)
						log.Fatal("@CreateReferencesViaDependencies: depOnAttr in condition doesn't exist? ", condition.DependsOnAttr)
						break
					}

					depOnAttrTokens := strings.Split(depOnAttr, ".")
					toMember := depOnAttrTokens[0]
					toReference := depOnAttrTokens[1]
					toReferencedAttr := fmt.Sprintf("%s.%s", toMember, toReference)

					var toMemberID, toAttrID string
					if tableID, err := db.TableID(mmd.DBConn, toMember, appConfig.AppID); err == nil {
						toMemberID = tableID
					} else {
						fmt.Println(toMember, appConfig.AppID)
						log.Fatal("@CreateReferencesViaDependencies.GetTableID: ", err)
					}

					if attrID, err := db.AttrID(mmd.DBConn, toMemberID, toReference); err == nil {
						toAttrID = attrID
					} else {
						fmt.Printf("@CreateReferencesViaDependencies.GetAttrID: Table: '%s' | Attr : '%s'\n", toMemberID, toReference)
						log.Fatal("@CreateReferencesViaDependencies.GetAttrID: ", err)
					}

					if toVal, fromVal, fromID, err := nodeData.GetRefValsFromDataMap(fromReferencedAttr, toReferencedAttr, false); err != nil {
						log.Println("@CreateReferencesViaDependencies.GetRefValsFromDataMap: ", err)
					} else {
						ref := MappingRef{
							appID:        appConfig.AppID,
							fromID:       fromID,
							fromVal:      fromVal,
							fromMemberID: fromMemberID,
							fromMember:   fromMember,
							fromAttr:     fromReference,
							fromAttrID:   fromAttrID,
							toVal:        toVal,
							toMemberID:   toMemberID,
							toMember:     toMember,
							toAttr:       toReference,
							toAttrID:     toAttrID,
						}

						log.Println("@CreateReferencesViaDependencies | Ref Created | ", ref)

						refs = append(refs, ref)
					}
				}
			} else {
				log.Fatal("@CreateReferencesViaDependencies: Unable to fetch referencedTag ", depOn.Tag)
			}
		}
	} else {
		log.Fatal("@CreateReferencesViaDependencies: Unable to fetch dependencies ", tag.Name)
	}
	return refs, nil
}

func (mmd *MappedMemberData) CreateReferencesViaOwnerships(appConfig config.AppConfig, tag config.Tag, nodeData DataMap) ([]MappingRef, error) {

	// log.Printf("@CreateReferencesViaOwnerships | nodeTag: %s | nodeData: %v  \n", tag.Name, nodeData)

	var refs []MappingRef

	if strings.EqualFold(tag.Name, "root") {
		return refs, nil
	}

	if own := appConfig.GetOwnership(tag.Name, "root"); own != nil {
		if rootTag, err := appConfig.GetTag("root"); err == nil {

			// log.Printf("@CreateReferencesViaOwnerships | referencedTag: %s  \n", rootTag.Name)

			for _, condition := range own.Conditions {

				tagAttr, err := tag.ResolveTagAttr(condition.TagAttr)
				if err != nil {
					log.Println(err, tag.Name, condition.TagAttr)
					log.Fatal("@CreateReferencesViaOwnerships: tagAttr in condition doesn't exist? ", condition.TagAttr)
					break
				}
				tagAttrTokens := strings.Split(tagAttr, ".")
				fromMember := tagAttrTokens[0]
				fromReference := tagAttrTokens[1]
				fromReferencedAttr := fmt.Sprintf("%s.%s", fromMember, fromReference)

				var fromMemberID, fromAttrID string
				if tableID, err := db.TableID(mmd.DBConn, fromMember, appConfig.AppID); err == nil {
					fromMemberID = tableID
				} else {
					fmt.Println(fromMember, appConfig.AppID)
					log.Fatal("@CreateReferencesViaOwnerships.GetTableID: ", err)
				}

				if attrID, err := db.AttrID(mmd.DBConn, fromMemberID, fromReference); err == nil {
					fromAttrID = attrID
				} else {
					fmt.Printf("@CreateReferencesViaOwnerships.GetAttrID: Table: '%s' | Attr : '%s'\n", fromMemberID, fromReference)
					log.Fatal("@CreateReferencesViaOwnerships.GetAttrID: ", err)
				}

				rootAttr, err := rootTag.ResolveTagAttr(condition.DependsOnAttr)
				if err != nil {
					log.Println(err, rootTag.Name, condition.DependsOnAttr)
					log.Fatal("@CreateReferencesViaOwnerships: depOnAttr in condition doesn't exist? ", condition.DependsOnAttr)
					break
				}
				rootAttrTokens := strings.Split(rootAttr, ".")
				rootMember := rootAttrTokens[0]
				rootReference := rootAttrTokens[1]

				var rootMemberID, rootAttrID string
				if tableID, err := db.TableID(mmd.DBConn, rootMember, appConfig.AppID); err == nil {
					rootMemberID = tableID
				} else {
					fmt.Println(rootMember, appConfig.AppID)
					log.Fatal("@CreateReferencesViaOwnerships.GetTableID: ", err)
				}

				if attrID, err := db.AttrID(mmd.DBConn, rootMemberID, rootReference); err == nil {
					rootAttrID = attrID
				} else {
					fmt.Printf("@CreateReferencesViaOwnerships.GetAttrID: Table: '%s' | Attr : '%s'\n", rootMemberID, rootReference)
					log.Fatal("@CreateReferencesViaOwnerships.GetAttrID: ", err)
				}

				if toVal, fromVal, fromID, err := nodeData.GetRefValsFromDataMap(fromReferencedAttr, rootReference, false); err != nil {
					log.Println("@CreateReferencesViaOwnerships.GetRefValsFromDataMap: ", err)
				} else {
					ref := MappingRef{
						appID:        appConfig.AppID,
						fromID:       fromID,
						fromVal:      fromVal,
						fromMemberID: fromMemberID,
						fromMember:   fromMember,
						fromAttr:     fromReference,
						fromAttrID:   fromAttrID,
						toVal:        toVal,
						toMemberID:   rootMemberID,
						toMember:     rootMember,
						toAttr:       rootReference,
						toAttrID:     rootAttrID,
					}

					log.Println("@CreateReferencesViaOwnerships | Ref Created | ", ref)

					refs = append(refs, ref)
				}
			}
		} else {
			log.Fatal("@CreateReferencesViaOwnerships: Unable to fetch referencedTag ", own.Tag)
		}
	} else {
		log.Fatal("@CreateReferencesViaOwnerships: Unable to fetch ownership ", tag.Name)
	}
	return refs, nil
}

func (mmd *MappedMemberData) CreateSelfReferences(bagAppConfig config.AppConfig, bagTag config.Tag, bagData DataMap) ([]MappingRef, error) {

	var bagRefs []MappingRef

	if refs, err := mmd.CreateInnerDependencyReferences(bagAppConfig, bagTag, bagData); err != nil {
		fmt.Println(bagData)
		log.Fatal(err)
	} else {
		bagRefs = append(bagRefs, refs...)
	}

	if refs, err := mmd.CreateReferencesViaDependencies(bagAppConfig, bagTag, bagData); err != nil {
		fmt.Println(bagData)
		log.Fatal(err)
	} else {
		bagRefs = append(bagRefs, refs...)
	}

	if refs, err := mmd.CreateReferencesViaOwnerships(bagAppConfig, bagTag, bagData); err != nil {
		fmt.Println(bagData)
		log.Fatal(err)
	} else {
		bagRefs = append(bagRefs, refs...)
	}

	return bagRefs, nil
}

func (mmd *MappedMemberData) FindMMV(appID, fromID, fromMemberID, fromAttrID string) (string, *MappedMemberValue) {
	for toAttr, mmv := range mmd.Data {
		if mmv.AppID == appID &&
			mmv.FromMemberID == fromMemberID &&
			mmv.FromID == fromID &&
			mmv.FromAttrID == fromAttrID {
			return toAttr, &mmv
		}
	}
	return "", nil
}
