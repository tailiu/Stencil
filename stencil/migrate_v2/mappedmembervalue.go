package migrate_v2

import (
	"fmt"
	"log"
	config "stencil/config/v2"
	"stencil/db"
	"stencil/helper"
	"strings"
)

func (mmv MappedMemberValue) GetTableID(tableName, appID string) string {
	if len(tableName) == 0 || len(appID) == 0 {
		log.Fatalf("@mmv.GetTableID: Table ID or App ID nil | table: '%s' | App: '%s' ", tableName, appID)
	} else if mmv.DBConn == nil {
		log.Fatal("@mmv.GetTableID: DBConn not set!")
	} else {
		if tableID, err := db.TableID(mmv.DBConn, tableName, appID); err == nil {
			return tableID
		} else {
			fmt.Println(tableName, appID)
			log.Fatal("@mmv.GetTableID: ", err)
		}
	}
	err := fmt.Errorf("This is the end. Shouldn't be here. Table: '%s', App: '%s'", tableName, appID)
	log.Fatal("@mmv.GetTableID: ", err)
	return err.Error()
}

func (mmv MappedMemberValue) GetAttrID(attrName, tableID string) string {
	if len(tableID) == 0 || len(attrName) == 0 {
		log.Fatalf("@mmv.GetAttrID: Table ID or App ID or AttrName nil | Table: '%s' | Attr : '%s'", tableID, attrName)
	} else if mmv.DBConn == nil {
		log.Fatal("@mmv.GetAttrID: DBConn not set!")
	} else {
		if attrID, err := db.AttrID(mmv.DBConn, tableID, attrName); err == nil {
			return attrID
		} else {
			fmt.Printf("@mmv.GetAttrID: Table: '%s' | Attr : '%s'\n", tableID, attrName)
			log.Fatal("@mmv.GetAttrID: ", err)
		}
	}
	err := fmt.Errorf("This is the end. Shouldn't be here. Table: '%s' | Attr : '%s'", tableID, attrName)
	log.Fatal("@mmv.GetAttrID: ", err)
	return err.Error()
}

func (mmv *MappedMemberValue) StoreMemberAndAttr(mappedStmt string) {
	tokens := strings.Split(mappedStmt, ".")
	mmv.FromMember = tokens[0]
	mmv.FromAttr = tokens[1]
	mmv.FromMemberID = mmv.GetTableID(mmv.FromMember, mmv.AppID)
	mmv.FromAttrID = mmv.GetAttrID(mmv.FromAttr, mmv.FromMemberID)
}

func (mmv MappedMemberValue) GetMemberAttr() string {
	return fmt.Sprintf("%s.%s", mmv.FromMember, mmv.FromAttr)
}

func (mmv *MappedMemberValue) CreateReference(fromAttr, toAttr, mappedStmt string, dataMap DataMap) error {

	hardRef := false
	if strings.Contains(mappedStmt, "#REFHARD") {
		hardRef = true
	}

	if toVal, fromVal, fromID, err := dataMap.GetRefValsFromDataMap(fromAttr, toAttr, hardRef); err == nil {
		fromAttrTokens := strings.Split(fromAttr, ".")
		toAttrTokens := strings.Split(toAttr, ".")

		fromMemberID := mmv.GetTableID(fromAttrTokens[0], mmv.AppID)
		fromAttrID := mmv.GetAttrID(fromAttrTokens[1], fromMemberID)

		toMemberID := mmv.GetTableID(toAttrTokens[0], mmv.AppID)
		toAttrID := mmv.GetAttrID(toAttrTokens[1], toMemberID)

		mmv.Ref = &MappingRef{
			appID:        mmv.AppID,
			fromID:       fromID,
			fromVal:      fromVal,
			fromMemberID: fromMemberID,
			fromMember:   fromAttrTokens[0],
			fromAttr:     fromAttrTokens[1],
			fromAttrID:   fromAttrID,
			toVal:        toVal,
			toMemberID:   toMemberID,
			toMember:     toAttrTokens[0],
			toAttr:       toAttrTokens[1],
			toAttrID:     toAttrID,
		}
	} else {
		return err
	}
	return nil
}

func (mmv *MappedMemberValue) SetFromID(dataMap DataMap) error {

	if len(mmv.FromMember) == 0 {
		log.Fatal("@mmv.SetFromID: FromMember is not set!")
	}

	if val, ok := dataMap[mmv.FromMember+".id"]; ok {
		if val != nil {
			mmv.FromID = fmt.Sprint(helper.GetInt64(val))
			return nil
		}
	}
	err := fmt.Errorf("@mmv.SetFromID: Can't find it in: %s", mmv.FromMember)
	fmt.Println(dataMap)
	log.Fatal(err)
	return err
}

func (mmv *MappedMemberValue) CreateInnerDependencyReferences(appConfig config.AppConfig, tag config.Tag, nodeData DataMap, attr string) ([]MappingRef, error) {

	log.Printf("@CreateInnerDependencyReferences | nodeTag: %s | nodeData: %v | attr: '%s' \n", tag.Name, nodeData, attr)

	var refs []MappingRef

	for _, innerDependency := range tag.InnerDependencies {
		for dependsOn, dependee := range innerDependency {

			depTokens := strings.Split(dependee, ".")
			dependeeMember := tag.Members[depTokens[0]]
			dependeeAttr := depTokens[1]
			dependeeReferencedAttr := fmt.Sprintf("%s.%s", dependeeMember, dependeeAttr)

			var dependeeMemberID, dependeeAttrID string
			if tableID, err := db.TableID(mmv.DBConn, dependeeMember, appConfig.AppID); err == nil {
				dependeeMemberID = tableID
			} else {
				fmt.Println(dependeeMember, appConfig.AppID)
				log.Fatal("@CreateInnerDependencyReferences.GetTableID: ", err)
			}

			if attrID, err := db.AttrID(mmv.DBConn, dependeeMemberID, dependeeAttr); err == nil {
				dependeeAttrID = attrID
			} else {
				fmt.Printf("@CreateInnerDependencyReferences.GetAttrID: Table: '%s' | Attr : '%s'\n", dependeeMemberID, dependeeAttr)
				log.Fatal("@CreateInnerDependencyReferences.GetAttrID: ", err)
			}

			if len(attr) > 0 && dependeeReferencedAttr != attr {
				log.Printf("@CreateInnerDependencyReferences | %s != %s \n", dependeeReferencedAttr, attr)
				continue
			} else if dependeeReferencedAttr == attr {
				log.Printf("@CreateInnerDependencyReferences | %s == %s \n", dependeeReferencedAttr, attr)
			}

			depOnTokens := strings.Split(dependsOn, ".")
			depOnMember := tag.Members[depOnTokens[0]]
			depOnAttr := depOnTokens[1]
			depOnReferencedAttr := fmt.Sprintf("%s.%s", depOnMember, depOnAttr)

			var depOnMemberID, depOnAttrID string
			if tableID, err := db.TableID(mmv.DBConn, depOnMember, appConfig.AppID); err == nil {
				depOnMemberID = tableID
			} else {
				fmt.Println(depOnMember, appConfig.AppID)
				log.Fatal("@CreateInnerDependencyReferences.GetTableID: ", err)
			}

			if attrID, err := db.AttrID(mmv.DBConn, depOnMemberID, depOnAttr); err == nil {
				depOnAttrID = attrID
			} else {
				fmt.Printf("@CreateInnerDependencyReferences.GetAttrID: Table: '%s' | Attr : '%s'\n", depOnMemberID, depOnAttr)
				log.Fatal("@CreateInnerDependencyReferences.GetAttrID: ", err)
			}

			if toVal, fromVal, fromID, err := nodeData.GetRefValsFromDataMap(dependeeReferencedAttr, depOnReferencedAttr, false); err != nil {

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

func (mmv *MappedMemberValue) CreateReferencesViaDependencies(appConfig config.AppConfig, tag config.Tag, nodeData DataMap, attr string) ([]MappingRef, error) {

	log.Printf("@CreateReferencesViaDependencies | nodeTag: %s | nodeData: %v | attr: '%s' \n", tag.Name, nodeData, attr)

	var refs []MappingRef

	if dep, err := appConfig.GetDependency(tag.Name); err == nil {
		for _, depOn := range dep.DependsOn {
			if referencedTag, err := appConfig.GetTag(depOn.Tag); err == nil {

				log.Printf("@CreateReferencesViaDependencies | referencedTag: %s  \n", referencedTag.Name)

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
					if tableID, err := db.TableID(mmv.DBConn, fromMember, appConfig.AppID); err == nil {
						fromMemberID = tableID
					} else {
						fmt.Println(fromMember, appConfig.AppID)
						log.Fatal("@CreateReferencesViaDependencies.GetTableID: ", err)
					}

					if attrID, err := db.AttrID(mmv.DBConn, fromMemberID, fromReference); err == nil {
						fromAttrID = attrID
					} else {
						fmt.Printf("@CreateReferencesViaDependencies.GetAttrID: Table: '%s' | Attr : '%s'\n", fromMemberID, fromReference)
						log.Fatal("@CreateReferencesViaDependencies.GetAttrID: ", err)
					}

					if len(attr) > 0 && fromReferencedAttr != attr {
						log.Printf("@CreateReferencesViaDependencies | %s != %s \n", fromReferencedAttr, attr)
						continue
					} else if fromReferencedAttr == attr {
						log.Printf("@CreateReferencesViaDependencies | %s == %s \n", fromReferencedAttr, attr)
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
					if tableID, err := db.TableID(mmv.DBConn, toMember, appConfig.AppID); err == nil {
						toMemberID = tableID
					} else {
						fmt.Println(toMember, appConfig.AppID)
						log.Fatal("@CreateReferencesViaDependencies.GetTableID: ", err)
					}

					if attrID, err := db.AttrID(mmv.DBConn, toMemberID, toReference); err == nil {
						toAttrID = attrID
					} else {
						fmt.Printf("@CreateReferencesViaDependencies.GetAttrID: Table: '%s' | Attr : '%s'\n", toMemberID, toReference)
						log.Fatal("@CreateReferencesViaDependencies.GetAttrID: ", err)
					}

					if toVal, fromVal, fromID, err := nodeData.GetRefValsFromDataMap(fromReferencedAttr, toReferencedAttr, false); err != nil {
						log.Println(err)
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

func (mmv *MappedMemberValue) CreateReferencesViaOwnerships(appConfig config.AppConfig, tag config.Tag, nodeData DataMap, attr string) ([]MappingRef, error) {

	log.Printf("@CreateReferencesViaOwnerships | nodeTag: %s | nodeData: %v | attr: '%s'  \n", tag.Name, nodeData, attr)

	var refs []MappingRef

	if strings.EqualFold(tag.Name, "root") {
		return refs, nil
	}

	if own := appConfig.GetOwnership(tag.Name, "root"); own != nil {
		if rootTag, err := appConfig.GetTag("root"); err == nil {

			log.Printf("@CreateReferencesViaOwnerships | referencedTag: %s  \n", rootTag.Name)

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
				if tableID, err := db.TableID(mmv.DBConn, fromMember, appConfig.AppID); err == nil {
					fromMemberID = tableID
				} else {
					fmt.Println(fromMember, appConfig.AppID)
					log.Fatal("@CreateReferencesViaOwnerships.GetTableID: ", err)
				}

				if attrID, err := db.AttrID(mmv.DBConn, fromMemberID, fromReference); err == nil {
					fromAttrID = attrID
				} else {
					fmt.Printf("@CreateReferencesViaOwnerships.GetAttrID: Table: '%s' | Attr : '%s'\n", fromMemberID, fromReference)
					log.Fatal("@CreateReferencesViaOwnerships.GetAttrID: ", err)
				}

				if len(attr) > 0 && fromReferencedAttr != attr {
					log.Printf("@CreateReferencesViaOwnerships | %s != %s \n", fromReferencedAttr, attr)
					continue
				} else if fromReferencedAttr == attr {
					log.Printf("@CreateReferencesViaOwnerships | %s == %s \n", fromReferencedAttr, attr)
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
				if tableID, err := db.TableID(mmv.DBConn, rootMember, appConfig.AppID); err == nil {
					rootMemberID = tableID
				} else {
					fmt.Println(rootMember, appConfig.AppID)
					log.Fatal("@CreateReferencesViaOwnerships.GetTableID: ", err)
				}

				if attrID, err := db.AttrID(mmv.DBConn, rootMemberID, rootReference); err == nil {
					rootAttrID = attrID
				} else {
					fmt.Printf("@CreateReferencesViaOwnerships.GetAttrID: Table: '%s' | Attr : '%s'\n", rootMemberID, rootReference)
					log.Fatal("@CreateReferencesViaOwnerships.GetAttrID: ", err)
				}

				if toVal, fromVal, fromID, err := nodeData.GetRefValsFromDataMap(fromReferencedAttr, rootReference, false); err != nil {
					log.Println(err)
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

func (mmv *MappedMemberValue) CreateSelfReferences(bagAppConfig config.AppConfig, bagTag config.Tag, bagData DataMap) ([]MappingRef, error) {

	var bagRefs []MappingRef

	if refs, err := mmv.CreateInnerDependencyReferences(bagAppConfig, bagTag, bagData, mmv.FromAttr); err != nil {
		fmt.Println(cleanedFromAttr)
		fmt.Println(bagData)
		log.Fatal(err)
	} else {
		bagRefs = append(bagRefs, refs...)
	}

	if refs, err := mmv.CreateReferencesViaDependencies(bagAppConfig, bagTag, bagData, mmv.FromAttr); err != nil {
		fmt.Println(cleanedFromAttr)
		fmt.Println(bagData)
		log.Fatal(err)
	} else {
		bagRefs = append(bagRefs, refs...)
	}

	if refs, err := mmv.CreateReferencesViaOwnerships(bagAppConfig, bagTag, bagData, mmv.FromAttr); err != nil {
		fmt.Println(cleanedFromAttr)
		fmt.Println(bagData)
		log.Fatal(err)
	} else {
		bagRefs = append(bagRefs, refs...)
	}

	return bagRefs, nil
}
