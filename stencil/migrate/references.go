package migrate

import (
	"fmt"
	"log"
	"stencil/db"
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

func (self *MigrationWorkerV2) AddInnerReferences(node *DependencyNode, member string) error {
	log.Printf("@AddInnerReferences | nodeTag: %s | member: %s \n nodeData: %v \n", node.Tag.Name, member, node.Data)
	for _, innerDependency := range node.Tag.InnerDependencies {
		for dependsOn, dependee := range innerDependency {

			depTokens := strings.Split(dependee, ".")
			dependeeMember := node.Tag.Members[depTokens[0]]
			dependeeAttr := depTokens[1]
			dependeeReferencedAttr := fmt.Sprintf("%s.%s", dependeeMember, dependeeAttr)
			dependeeMemberID, err := db.TableID(self.logTxn.DBconn, dependeeMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Fatal("@AddInnerReferences: Unable to resolve id for dependeeMember ", dependeeMember)
			}

			depOnTokens := strings.Split(dependsOn, ".")
			depOnMember := node.Tag.Members[depOnTokens[0]]
			depOnAttr := depOnTokens[1]
			depOnReferencedAttr := fmt.Sprintf("%s.%s", depOnMember, depOnAttr)
			depOnMemberID, err := db.TableID(self.logTxn.DBconn, depOnMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Fatal("@AddInnerReferences: Unable to resolve id for depOnMember ", depOnMember)
			}

			if member != "" {
				if !strings.EqualFold(dependeeMember, member) && !strings.EqualFold(depOnMember, member) {
					continue
				}
			}

			var fromID, toID int64

			if val, ok := node.Data[dependeeMember+".id"]; ok {
				fromID = val.(int64)
			} else {
				fmt.Println(node.Data)
				log.Fatalf("@AddInnerReferences | fromID | '%s.id' doesn't exist in node data? \n", dependeeMember)
			}

			log.Printf("@AddInnerReferences | toID | Checking dependeeReferencedAttr: '%s' \n", dependeeReferencedAttr)

			if val, ok := node.Data[dependeeReferencedAttr]; ok {
				toID = val.(int64)
			} else {
				log.Printf("@AddInnerReferences | toID | dependeeReferencedAttr: '%s' doesn't exist in node data? \n", dependeeReferencedAttr)
				log.Printf("@AddInnerReferences | toID | Checking depOnReferencedAttr: '%s' \n", depOnReferencedAttr)
				if val, ok := node.Data[depOnReferencedAttr]; ok {
					toID = val.(int64)
				} else {
					fmt.Println(node.Data)
					log.Fatalf("@AddInnerReferences | toID | depOnReferencedAttr: '%s' doesn't exist in node data? \n", depOnReferencedAttr)
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

	for _, dep := range self.SrcAppConfig.GetSubDependencies(node.Tag.Name) {
		for _, depOn := range dep.DependsOn {
			if referencedTag, err := self.SrcAppConfig.GetTag(depOn.Tag); err == nil {
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
						fromID = val.(int64)
					} else {
						fmt.Println(node.Data)
						log.Fatal("@AddToReferencesViaDependencies:", fromMember+".id", " doesn't exist in node data? ", node.Tag.Name)
					}

					if val, ok := node.Data[fromReferencedAttr]; ok {
						toID = val.(int64)
					} else {
						fmt.Println(node.Data)
						log.Fatal("@AddToReferencesViaDependencies: '", fromReferencedAttr, "' doesn't exist in node data? ", referencedTag.Name)
					}

					if err := db.CreateNewReference(self.tx.StencilTx, self.SrcAppConfig.AppID, fromMemberID, fromID, toMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), fromReference, toReference); err != nil {
						fmt.Println("#Args: ", self.SrcAppConfig.AppID, fromMemberID, fromID, toMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), fromReference, toReference)
						log.Fatal("@AddToReferencesViaDependencies: Unable to CreateNewReference: ", err)
						return err
					}
				}
			} else {
				log.Fatal("@AddToReferencesViaDependencies: Unable to fetch referencedTag ", depOn.Tag)
			}
		}
	}
	return nil
}

func (self *MigrationWorkerV2) AddToReferences(currentNode *DependencyNode, referencedNode *DependencyNode) error {
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
				fromID = val.(int64)
			} else {
				fmt.Println(currentNode.Data)
				log.Fatal("@AddToReferences:", fromMember+".id", " doesn't exist in node data? ", currentNode.Tag.Name)
			}

			if val, ok := referencedNode.Data[toMember+".id"]; ok {
				toID = val.(int64)
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
