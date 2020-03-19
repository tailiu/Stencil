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
			if idRow.ToID != nil {
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
			if newIDRows, err := self.GetRowsFromIDTable(idRow.FromAppID, idRow.FromMemberID, idRow.FromID, false); err == nil {
				if exists, err := self._CheckReferenceExistsInPreviousMigrations(newIDRows, refAttr); err == nil && exists {
					return exists, err
				}
			}
		}
	} else {
		log.Println("@_CheckReferenceExistsInPreviousMigrations: IDRows | null | ", idRows)
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

	if ref.toID == nil {
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
			log.Println("@_CreateMappedReference: Reference Does Already Exist | ", ref.appID, ref.fromMember, dependeeMemberID, ref.fromID, ref.toMember, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
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
		color.Yellow.Printf("New Ref | fromApp: %s, fromMember: %s, fromID: %s, toMember: %s, toID: %s, migrationID: %s, fromAttr: %s, toAttr: %s\n", ref.appID, dependeeMemberID, ref.fromID, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
	}
	return nil
}

func (self *MigrationWorkerV2) AddInnerReferences(node *DependencyNode, member string) error {
	return nil
	for _, innerDependency := range node.Tag.InnerDependencies {
		for dependee, dependsOn := range innerDependency {

			depTokens := strings.Split(dependee, ".")
			dependeeMember := node.Tag.Members[depTokens[0]]
			dependeeMemberID, err := db.TableID(self.logTxn.DBconn, dependeeMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Fatal("@AddInnerReferences: Unable to resolve id for dependeeMember ", dependeeMember)
			}

			depOnTokens := strings.Split(dependsOn, ".")
			depOnMember := node.Tag.Members[depOnTokens[0]]
			depOnMemberID, err := db.TableID(self.logTxn.DBconn, depOnMember, self.SrcAppConfig.AppID)
			if err != nil {
				log.Fatal("@AddInnerReferences: Unable to resolve id for depOnMember ", depOnMember)
			}

			if member != "" {
				if !strings.EqualFold(dependeeMember, member) && !strings.EqualFold(depOnMember, member) {
					continue
				}
			}

			var fromID, toID string

			if val, ok := node.Data[dependeeMember+".id"]; ok {
				fromID = fmt.Sprint(val)
			} else {
				fmt.Println(node.Data)
				log.Fatal("@AddInnerReferences:", dependeeMember+".id", " doesn't exist in node data? ", node.Tag.Name)
			}

			if val, ok := node.Data[depOnMember+".id"]; ok {
				toID = fmt.Sprint(val)
			} else {
				fmt.Println(node.Data)
				log.Fatal("@AddInnerReferences:", depOnMember+".id", " doesn't exist in node data? ", node.Tag.Name)
			}

			if err := db.CreateNewReference(self.tx.StencilTx, self.SrcAppConfig.AppID, dependeeMemberID, fromID, depOnMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), depTokens[1], depOnTokens[1]); err != nil {
				fmt.Println("#Args: ", self.SrcAppConfig.AppID, dependeeMemberID, fromID, depOnMemberID, toID, fmt.Sprint(self.logTxn.Txn_id), depTokens[1], depOnTokens[1])
				log.Fatal("@AddInnerReferences: Unable to CreateNewReference: ", err)
				return err
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

			var fromID, toID string

			if val, ok := currentNode.Data[fromMember+".id"]; ok {
				fromID = fmt.Sprint(val)
			} else {
				fmt.Println(currentNode.Data)
				log.Fatal("@AddToReferences:", fromMember+".id", " doesn't exist in node data? ", currentNode.Tag.Name)
			}

			if val, ok := referencedNode.Data[toMember+".id"]; ok {
				toID = fmt.Sprint(val)
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
