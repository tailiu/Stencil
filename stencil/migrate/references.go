package migrate

import (
	"fmt"
	"log"
	"stencil/db"
	"strings"
)

func (self *MigrationWorkerV2) AddMappedReferences(refs []MappingRef) error {

	for _, ref := range refs {

		dependeeMemberID, err := db.TableID(self.logTxn.DBconn, ref.fromMember, self.SrcAppConfig.AppID)
		if err != nil {
			log.Fatal("@AddMappedReferences: Unable to resolve id for dependeeMember ", ref.fromMember)
			return err
		}

		depOnMemberID, err := db.TableID(self.logTxn.DBconn, ref.toMember, self.SrcAppConfig.AppID)
		if err != nil {
			log.Fatal("@AddMappedReferences: Unable to resolve id for depOnMember ", ref.toMember)
			return err
		}

		if len(ref.toID) < 1 {
			log.Println("@AddMappedReferences: Unable to CreateNewReference | ", self.SrcAppConfig.AppID, ref.fromMember, dependeeMemberID, ref.fromID, ref.toMember, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
			continue
		}

		if err := db.CreateNewReference(self.tx.StencilTx, self.SrcAppConfig.AppID, dependeeMemberID, ref.fromID, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr); err != nil {
			fmt.Println(refs)
			fmt.Println("#Args: ", self.SrcAppConfig.AppID, ref.fromMember, dependeeMemberID, ref.fromID, ref.toMember, depOnMemberID, ref.toID, fmt.Sprint(self.logTxn.Txn_id), ref.fromAttr, ref.toAttr)
			log.Fatal("@AddMappedReferences: Unable to CreateNewReference: ", err)
			return err
		}
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
