package migrate_v2

import (
	"fmt"
	"log"
	"stencil/config"
	"stencil/db"
	"stencil/helper"
	"strings"

	"github.com/gookit/color"
)

func (self *MigrationWorkerV2) AddMappedReferencesIfNotExist(refs []MappingRef) error {
	for _, ref := range refs {
		if err := self._CreateMappedReference(ref, true); err != nil {
			return err
		}
	}
	return nil
}

func (self *MigrationWorkerV2) AddMappedReferences(refs []MappingRef) error {

	for _, ref := range refs {
		if err := self._CreateMappedReference(ref, false); err != nil {
			return err
		}
	}
	return nil
}

func (self *MigrationWorkerV2) _CheckReferenceExistsInPreviousMigrations(idRows []IDRow, refAttr string) (bool, error) {
	if len(idRows) > 0 {
		for _, idRow := range idRows {
			if idRow.ToID != 0 {
				if foundAttr, isFound := self.FetchMappedAttribute(idRow.FromAppName, idRow.FromAppID, idRow.ToAppName, idRow.ToAppID, idRow.FromMember, idRow.ToMember, refAttr); isFound {
					self.Logger.Debugf("@_CheckReferenceExistsInPreviousMigrations: Mapped Attr Found | FromAttr: %s | FromApp: %s, FromMember: %s, ToApp: %s, ToMember: %s, ToAttr: %s \n", foundAttr, idRow.FromAppName, idRow.FromMember, idRow.ToAppName, idRow.ToMember, refAttr)
					refAttr = foundAttr
				} else {
					self.Logger.Debugf("@_CheckReferenceExistsInPreviousMigrations: No Mapped Attr found | FromApp: %s, FromMember: %s, ToApp: %s, ToMember: %s, ToAttr: %s \n", idRow.FromAppName, idRow.FromMember, idRow.ToAppName, idRow.ToMember, refAttr)
				}
			}
			fmt.Println("@_CheckReferenceExistsInPreviousMigrations: IDRow | ", idRow)
			fmt.Printf("@_CheckReferenceExistsInPreviousMigrations: Checking Reference for | App: %v, Member: %v, ID: %v, Attr: %v\n", idRow.FromAppID, idRow.FromMemberID, idRow.FromID, refAttr)
			if db.CheckIfReferenceExists(self.logTxn.DBconn, idRow.FromAppID, idRow.FromMemberID, idRow.FromID, refAttr) {
				log.Printf("@_CheckReferenceExistsInPreviousMigrations: Reference Already Exists | App: %v, Member: %v, ID: %v, Attr: %v\n", idRow.FromAppID, idRow.FromMemberID, idRow.FromID, refAttr)
				return true, nil
			}
			fmt.Printf("@_CheckReferenceExistsInPreviousMigrations: Reference doesn't exist | App: %v, Member: %v, ID: %v, Attr: %v\n", idRow.FromAppID, idRow.FromMemberID, idRow.FromID, refAttr)
			fmt.Printf("@_CheckReferenceExistsInPreviousMigrations: Getting IDRows | App: %v, Member: %v, ID: %v, getFrom: %v\n", idRow.FromAppID, idRow.FromMemberID, idRow.FromID, false)
			if newIDRows, err := self.GetRowsFromIDTable(idRow.FromAppID, idRow.FromMemberID, idRow.FromID, false); err == nil {
				if exists, err := self._CheckReferenceExistsInPreviousMigrations(newIDRows, refAttr); err == nil && exists {
					return exists, err
				}
			} else {
				self.Logger.Fatal(err)
			}
		}
	} else {
		log.Println("@_CheckReferenceExistsInPreviousMigrations: No More IDRows | ", idRows)
	}
	return false, nil
}

func (self *MigrationWorkerV2) _CreateMappedReference(ref MappingRef, checkForExistence bool) error {

	fmt.Println("@_CreateMappedReference: Enter >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Printf("@_CreateMappedReference | Ref : %v | checkForExistence : %v \n", ref, checkForExistence)
	defer fmt.Println("@_CreateMappedReference: Exit <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

	dependeeMemberID, err := db.TableID(self.logTxn.DBconn, ref.fromMember, ref.appID)
	if err != nil {
		fmt.Println(err)
		fmt.Println(ref)
		log.Fatal("@_CreateMappedReference: Unable to resolve id for dependeeMember ", ref.fromMember)
		return err
	}

	depOnMemberID, err := db.TableID(self.logTxn.DBconn, ref.toMember, ref.appID)
	if err != nil {
		fmt.Println(err)
		fmt.Println(ref)
		log.Fatal("@_CreateMappedReference: Unable to resolve id for depOnMember ", ref.toMember)
		return err
	}

	if ref.toID == 0 {
		log.Println("@_CreateMappedReference: Unable to CreateNewReference | ", ref.appID, ref.fromMember, dependeeMemberID, ref.fromID, ref.toMember, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
		return nil
	}

	if checkForExistence || ref.mergedFromBag {
		idRows := []IDRow{IDRow{
			FromAppID:    ref.appID,
			FromMemberID: dependeeMemberID,
			FromMember:   ref.fromMember,
			FromID:       ref.fromID}}

		if exists, err := self._CheckReferenceExistsInPreviousMigrations(idRows, ref.fromAttr); err == nil && exists {
			log.Println("@_CreateMappedReference: Reference Indeed Already Exists | ", ref.appID, ref.fromMember, dependeeMemberID, ref.fromID, ref.toMember, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
			return nil
		} else if err != nil {
			log.Fatal("@_CreateMappedReference > _CheckReferenceExistsInPreviousMigrations | Err: ", err)
		}
		log.Println("@_CreateMappedReference: Reference Doesn't Already Exist | ", ref.appID, ref.fromMember, dependeeMemberID, ref.fromID, ref.toMember, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
	}

	if err := db.CreateNewReference(self.tx.StencilTx, ref.appID, dependeeMemberID, ref.fromID, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr); err != nil {
		fmt.Println(ref)
		fmt.Println("#Args: ", ref.appID, ref.fromMember, dependeeMemberID, ref.fromID, ref.toMember, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
		log.Fatal("@_CreateMappedReference: Unable to CreateNewReference: ", err)
		return err
	} else {
		color.Yellow.Printf("New Ref | fromApp: %s, fromMember: %s, fromID: %v, toMember: %s, toID: %v, migrationID: %s, fromAttr: %s, toAttr: %s\n", ref.appID, dependeeMemberID, ref.fromID, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
	}
	return nil
}

func CreateInnerDependencyReferences(appConfig config.AppConfig, tag config.Tag, nodeData map[string]interface{}, attr string) ([]MappingRef, error) {

	log.Printf("@CreateInnerDependencyReferences | nodeTag: %s | nodeData: %v | attr: '%s' \n", tag.Name, nodeData, attr)

	var refs []MappingRef

	for _, innerDependency := range tag.InnerDependencies {
		for dependsOn, dependee := range innerDependency {

			depTokens := strings.Split(dependee, ".")
			dependeeMember := tag.Members[depTokens[0]]
			dependeeAttr := depTokens[1]
			dependeeReferencedAttr := fmt.Sprintf("%s.%s", dependeeMember, dependeeAttr)

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

			var fromID, toID int64

			if val, ok := nodeData[dependeeMember+".id"]; ok {
				fromID = helper.GetInt64(val)
			} else {
				log.Printf("@CreateInnerDependencyReferences | fromID | '%s.id' doesn't exist in node data? \n", dependeeMember)
			}

			log.Printf("@CreateInnerDependencyReferences | toID | Checking dependeeReferencedAttr: '%s' \n", dependeeReferencedAttr)

			if val, ok := nodeData[dependeeReferencedAttr]; ok {
				toID = helper.GetInt64(val)
			} else {
				log.Printf("@CreateInnerDependencyReferences | toID | dependeeReferencedAttr: '%s' doesn't exist in node data? \n", dependeeReferencedAttr)
				log.Printf("@CreateInnerDependencyReferences | toID | Checking depOnReferencedAttr: '%s' \n", depOnReferencedAttr)
				if val, ok := nodeData[depOnReferencedAttr]; ok {
					toID = helper.GetInt64(val)
				} else {
					fmt.Println(nodeData)
					log.Printf("@CreateInnerDependencyReferences | toID | depOnReferencedAttr: '%s' doesn't exist in node data? \n", depOnReferencedAttr)
					continue
				}
			}

			if toID == 0 || fromID == 0 {
				fmt.Println(appConfig.AppID, dependeeMember, fromID, depOnMember, toID, dependeeAttr, depOnAttr, dependeeReferencedAttr, depOnReferencedAttr)
				log.Fatal("@CreateInnerDependencyReferences: Unable to CreateNewReference: toID == 0 || fromID == 0")
			}

			ref := MappingRef{
				appID:      appConfig.AppID,
				fromID:     fromID,
				fromMember: dependeeMember,
				fromAttr:   dependeeAttr,
				toID:       toID,
				toMember:   depOnMember,
				toAttr:     depOnAttr,
			}

			log.Println("@CreateInnerDependencyReferences | Ref Created | ", ref)

			refs = append(refs, ref)
		}
	}

	return refs, nil
}

func CreateReferencesViaDependencies(appConfig config.AppConfig, tag config.Tag, nodeData map[string]interface{}, attr string) ([]MappingRef, error) {

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

					var fromID int64
					var toID string

					if val, ok := nodeData[fromMember+".id"]; ok {
						fromID = helper.GetInt64(val)
					} else {
						log.Println("@CreateReferencesViaDependencies | FromID | ", fromMember+".id", " doesn't exist in node data? NodeTagName: ", tag.Name)
						continue
					}

					if val, ok := nodeData[fromReferencedAttr]; ok {
						toID = fmt.Sprint(val)
					} else {
						log.Println("@CreateReferencesViaDependencies | toID | '", fromReferencedAttr, "' doesn't exist in node data? ReferencedTagName: ", referencedTag.Name)
						continue
					}

					if len(toID) == 0 || fromID == 0 {
						fmt.Println(appConfig.AppID, fromMember, fromID, toMember, toID, fromReference, toReference, fromReferencedAttr)
						log.Fatal("@CreateReferencesViaDependencies: Unable to CreateNewReference: toID == 0 || fromID == 0")
					}

					ref := MappingRef{
						appID:      appConfig.AppID,
						fromID:     fromID,
						fromMember: fromMember,
						fromAttr:   fromReference,
						toID:       toID,
						toMember:   toMember,
						toAttr:     toReference,
					}

					log.Println("@CreateReferencesViaDependencies | Ref Created | ", ref)

					refs = append(refs, ref)
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

func CreateReferencesViaOwnerships(appConfig config.AppConfig, tag config.Tag, nodeData map[string]interface{}, attr string) ([]MappingRef, error) {

	log.Printf("@CreateReferencesViaOwnerships | nodeTag: %s | nodeData: %v | attr: '%s'  \n", tag.Name, nodeData, attr)

	var refs []MappingRef

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

				var fromID, toID int64

				if val, ok := nodeData[fromMember+".id"]; ok {
					fromID = helper.GetInt64(val)
				} else {
					log.Println("@CreateReferencesViaOwnerships | FromID | ", fromMember+".id", " doesn't exist in node data? NodeTagName: ", tag.Name)
					continue
				}

				if val, ok := nodeData[fromReferencedAttr]; ok {
					toID = helper.GetInt64(val)
				} else {
					log.Println("@CreateReferencesViaOwnerships | toID | '", fromReferencedAttr, "' doesn't exist in node data? ReferencedTagName: ", rootTag.Name)
					continue
				}

				if toID == 0 || fromID == 0 {
					fmt.Println(appConfig.AppID, fromMember, fromID, rootMember, toID, fromReference, rootReference, fromReferencedAttr)
					log.Fatal("@CreateReferencesViaOwnerships: Unable to CreateNewReference: toID == 0 || fromID == 0")
				}

				ref := MappingRef{
					appID:      appConfig.AppID,
					fromID:     fromID,
					fromMember: fromMember,
					fromAttr:   fromReference,
					toID:       toID,
					toMember:   rootMember,
					toAttr:     rootReference,
				}

				log.Println("@CreateReferencesViaOwnerships | Ref Created | ", ref)

				refs = append(refs, ref)
			}
		} else {
			log.Fatal("@CreateReferencesViaOwnerships: Unable to fetch referencedTag ", own.Tag)
		}
	} else {
		if tag.Name != "root" {
			log.Fatal("@CreateReferencesViaOwnerships: Unable to fetch ownership ", tag.Name)
		}
	}
	return refs, nil
}

func (self *MigrationWorkerV2) AddInnerReferences(node *DependencyNode) error {
	log.Fatal("Why are you in AddInnerReferences?")
	for _, innerDependency := range node.Tag.InnerDependencies {
		for dependsOn, dependee := range innerDependency {

			depTokens := strings.Split(dependee, ".")
			dependeeMember := node.Tag.Members[depTokens[0]]
			dependeeAttr := depTokens[1]
			dependeeReferencedAttr := fmt.Sprintf("%s.%s", dependeeMember, dependeeAttr)
			dependeeMemberID, err := db.TableID(self.logTxn.DBconn, dependeeMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Println("@AddInnerReferences: Unable to resolve id for dependeeMember ", dependeeMember)
				continue
			}

			depOnTokens := strings.Split(dependsOn, ".")
			depOnMember := node.Tag.Members[depOnTokens[0]]
			depOnAttr := depOnTokens[1]
			depOnReferencedAttr := fmt.Sprintf("%s.%s", depOnMember, depOnAttr)
			depOnMemberID, err := db.TableID(self.logTxn.DBconn, depOnMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Println("@AddInnerReferences: Unable to resolve id for depOnMember ", depOnMember)
				continue
			}

			var fromID, toID int64

			if val, ok := node.Data[dependeeMember+".id"]; ok {
				fromID = helper.GetInt64(val)
			} else {
				log.Printf("@AddInnerReferences | fromID | '%s.id' doesn't exist in node data? \n", dependeeMember)
			}

			log.Printf("@AddInnerReferences | toID | Checking dependeeReferencedAttr: '%s' \n", dependeeReferencedAttr)

			if val, ok := node.Data[dependeeReferencedAttr]; ok {
				toID = helper.GetInt64(val)
			} else {
				log.Printf("@AddInnerReferences | toID | dependeeReferencedAttr: '%s' doesn't exist in node data? \n", dependeeReferencedAttr)
				log.Printf("@AddInnerReferences | toID | Checking depOnReferencedAttr: '%s' \n", depOnReferencedAttr)
				if val, ok := node.Data[depOnReferencedAttr]; ok {
					toID = helper.GetInt64(val)
				} else {
					fmt.Println(node.Data)
					log.Printf("@AddInnerReferences | toID | depOnReferencedAttr: '%s' doesn't exist in node data? \n", depOnReferencedAttr)
					continue
				}
			}

			if toID == 0 || fromID == 0 {
				fmt.Println(self.SrcAppConfig.AppID, dependeeMemberID, fromID, depOnMemberID, toID, dependeeAttr, depOnAttr, dependeeReferencedAttr, depOnReferencedAttr)
				log.Fatal("@AddInnerReferences: Unable to CreateNewReference: toID == 0 || fromID == 0")
			}

			if err := db.CreateNewReference(self.tx.StencilTx, self.SrcAppConfig.AppID, dependeeMemberID, fromID, depOnMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), dependeeAttr, depOnAttr); err != nil {
				fmt.Println("#Args: ", self.SrcAppConfig.AppID, dependeeMemberID, fromID, depOnMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), dependeeAttr, depOnAttr)
				log.Fatal("@AddInnerReferences: Unable to CreateNewReference: ", err)
				return err
			} else {
				color.Yellow.Printf("New Ref | fromApp: %s, fromMember: %s, fromID: %v, toMember: %s, toID: %v, migrationID: %s, fromAttr: %s, toAttr: %s\n", self.SrcAppConfig.AppID, dependeeMemberID, fromID, depOnMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), dependeeAttr, depOnAttr)
			}
		}
	}

	return nil
}

func (self *MigrationWorkerV2) AddToReferencesViaDependencies(node *DependencyNode) error {

	log.Fatal("Why are you in AddToReferencesViaDependencies?")

	log.Printf("@AddToReferencesViaDependencies | nodeTag: %s | nodeData: %v \n", node.Tag.Name, node.Data)

	if dep, err := self.SrcAppConfig.GetDependency(node.Tag.Name); err == nil {
		for _, depOn := range dep.DependsOn {
			if referencedTag, err := self.SrcAppConfig.GetTag(depOn.Tag); err == nil {

				log.Printf("@AddToReferencesViaDependencies | referencedTag: %s  \n", referencedTag.Name)

				for _, condition := range depOn.Conditions {
					tagAttr, err := node.Tag.ResolveTagAttr(condition.TagAttr)
					if err != nil {
						log.Println(err, node.Tag.Name, condition.TagAttr)
						log.Fatal("@AddToReferencesViaDependencies: tagAttr in condition doesn't exist? ", condition.TagAttr)
						break
					}
					tagAttrTokens := strings.Split(tagAttr, ".")
					fromMember := tagAttrTokens[0]
					fromReference := tagAttrTokens[1]
					fromReferencedAttr := fmt.Sprintf("%s.%s", fromMember, fromReference)
					fromMemberID, err := db.TableID(self.logTxn.DBconn, fromMember, self.SrcAppConfig.AppID)
					if err != nil {
						log.Fatal("@AddToReferencesViaDependencies: Unable to resolve id for fromMember ", fromMember)
					}

					depOnAttr, err := referencedTag.ResolveTagAttr(condition.DependsOnAttr)
					if err != nil {
						log.Println(err, referencedTag.Name, condition.DependsOnAttr)
						log.Fatal("@AddToReferencesViaDependencies: depOnAttr in condition doesn't exist? ", condition.DependsOnAttr)
						break
					}
					depOnAttrTokens := strings.Split(depOnAttr, ".")
					toMember := depOnAttrTokens[0]
					toReference := depOnAttrTokens[1]
					toMemberID, err := db.TableID(self.logTxn.DBconn, toMember, self.SrcAppConfig.AppID)
					if err != nil {
						log.Fatal("@AddToReferencesViaDependencies: Unable to resolve id for toMember ", toMember)
					}

					var fromID, toID int64

					if val, ok := node.Data[fromMember+".id"]; ok {
						fromID = helper.GetInt64(val)
					} else {
						log.Println("@AddToReferencesViaDependencies | FromID | ", fromMember+".id", " doesn't exist in node data? NodeTagName: ", node.Tag.Name)
					}

					if val, ok := node.Data[fromReferencedAttr]; ok {
						toID = helper.GetInt64(val)
					} else {
						log.Println("@AddToReferencesViaDependencies | toID | '", fromReferencedAttr, "' doesn't exist in node data? ReferencedTagName: ", referencedTag.Name)
					}

					if err := db.CreateNewReference(self.tx.StencilTx, self.SrcAppConfig.AppID, fromMemberID, fromID, toMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), fromReference, toReference); err != nil {
						fmt.Println("#Args: ", self.SrcAppConfig.AppID, fromMemberID, fromID, toMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), fromReference, toReference)
						log.Fatal("@AddToReferencesViaDependencies: Unable to CreateNewReference: ", err)
						return err
					} else {
						color.Yellow.Printf("New Ref | fromApp: %s, fromMember: %s, fromID: %v, toMember: %s, toID: %v, migrationID: %s, fromAttr: %s, toAttr: %s\n", self.SrcAppConfig.AppID, fromMemberID, fromID, toMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), fromReference, toReference)
					}
				}
			} else {
				log.Fatal("@AddToReferencesViaDependencies: Unable to fetch referencedTag ", depOn.Tag)
			}
		}
	} else {
		log.Fatal("@AddToReferencesViaDependencies: Unable to fetch dependencies ", node.Tag.Name)
	}
	return nil
}

func (self *MigrationWorkerV2) AddToReferences(currentNode *DependencyNode, referencedNode *DependencyNode) error {
	log.Fatal("Why are you in AddToReferences?")
	return nil
	if dep, err := self.SrcAppConfig.CheckDependency(currentNode.Tag.Name, referencedNode.Tag.Name); err != nil {
		fmt.Println(err)
		log.Fatal("@AddToReferences: CheckDependency can't find dependency!")
	} else {
		for _, condition := range dep.Conditions {
			tagAttr, err := currentNode.Tag.ResolveTagAttr(condition.TagAttr)
			if err != nil {
				log.Println(err, currentNode.Tag.Name, condition.TagAttr)
				log.Fatal("@AddToReferences: tagAttr in condition doesn't exist? ", condition.TagAttr)
				break
			}
			tagAttrTokens := strings.Split(tagAttr, ".")
			fromMember := tagAttrTokens[0]
			fromReference := tagAttrTokens[1]
			fromMemberID, err := db.TableID(self.logTxn.DBconn, fromMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Fatal("@AddToReferences: Unable to resolve id for fromMember ", fromMember)
			}

			depOnAttr, err := referencedNode.Tag.ResolveTagAttr(condition.DependsOnAttr)
			if err != nil {
				log.Println(err, referencedNode.Tag.Name, condition.DependsOnAttr)
				log.Fatal("@AddToReferences: depOnAttr in condition doesn't exist? ", condition.DependsOnAttr)
				break
			}
			depOnAttrTokens := strings.Split(depOnAttr, ".")
			toMember := depOnAttrTokens[0]
			toReference := depOnAttrTokens[1]
			toMemberID, err := db.TableID(self.logTxn.DBconn, toMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Fatal("@AddToReferences: Unable to resolve id for toMember ", toMember)
			}

			var fromID, toID int64

			if val, ok := currentNode.Data[fromMember+".id"]; ok {
				fromID = helper.GetInt64(val)
			} else {
				fmt.Println(currentNode.Data)
				log.Fatal("@AddToReferences:", fromMember+".id", " doesn't exist in node data? ", currentNode.Tag.Name)
			}

			if val, ok := referencedNode.Data[toMember+".id"]; ok {
				toID = helper.GetInt64(val)
			} else {
				fmt.Println(referencedNode.Data)
				log.Fatal("@AddToReferences:", toMember+".id", " doesn't exist in node data? ", referencedNode.Tag.Name)
			}

			if err := db.CreateNewReference(self.tx.StencilTx, self.SrcAppConfig.AppID, fromMemberID, fromID, toMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), fromReference, toReference); err != nil {
				fmt.Println("#Args: ", self.SrcAppConfig.AppID, fromMemberID, fromID, toMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), fromReference, toReference)
				log.Fatal("@AddToReferences: Unable to CreateNewReference: ", err)
				return err
			}
		}
	}
	return nil
}
