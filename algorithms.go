package main


/*
 * Assumptions:
 * 1. Threads don't communicate with each other. Reasons: Simplicity, performance.
 * 2. If threads die, they just restart
 *
 */


 func (t Thread) AggressiveMigration(uid int, srcApp, dstApp string, node *DependencyNode) {

	if t.Root == node && !checkUserInApp(uid, dstApp) {
		addUserToApplication(uid, dstApp);
	}
	lockAcquired := false;
	for {
		// randomlyGetUnvisitedNextNode(node) not only gets a migrating user's data but also some data shared by other users
		// this part (randomlyGetUnvisitedNextNode) needs to be changed: combine ownership and dependencies for traversal. root goes to ownership and then to dependencies. 
		// root will not be in dependencies anymore; only in ownerships.
		nextNode := randomlyGetUnvisitedNextNode(node); // returns one random adjacent node
		if nextNode != nil {
			AggressiveMigration(uid, srcApp, dstApp, nextNode);
		} else {
			if lockAcquired == true {
				break;
			} else {
				acquirePredicateLock(node);
				lockAcquired = true;
			}
		}
	}
	t1, t2, t3 := srcDB.BeginTransaction(), dstDB.BeginTransaction(), stencilDB.BeginTransaction()
	if node.Owner == uid || node.SharedWith(uid) {
		for _, nextNode := range GetAllNextNodes(node) { // returns a list of all next nodes, ignoring mark as visited flag
			// GetAllPrecedingNodes doesn't consider ownerships, only uses dependencies.
			// We can also add predicate locks to prevent precedingNodes from being added, but since it is very rare in social media applications,
			// We do not do that.
			if precedingNodes := GetAllPrecedingNodes(nextNode); len(precedingNodes) <= 1 && nextNode.Rules.DisplayOnlyIfPrecedingNodeExists {
				sendToBag(nextNode, nextNode.Owner);
			}
			// else {
			addToReferences(node, nextNode);
			// }
		}
		acquireWriteLock(node)
		for _, precedingNode := range GetAllPrecedingNodes(node) {
			if precedingNode.Owner != uid && !precedingNode.SharedWith(uid) {
				addToReferences(precedingNode, node);
			}
		}
		// resolveReferences(node); // Moved To Display
		migrateNode(uid, srcApp, dstApp, node);
		releaseWriteLock(node);
	} else {
		if precedingNodes := GetAllPrecedingNodes(node); len(precedingNodes) <= 1 && node.Rules.DisplayOnlyIfPrecedingNodeExists {
			// GetAllNextNodes doesn't consider ownerships, only uses dependencies.
			for _, nextNode := range GetAllNextNodes(node) { // ignore mark as visited flag while getting next nodes
				if precedingNodes := GetAllPrecedingNodes(nextNode); len(precedingNodes) <= 1 && nextNode.Rules.DisplayOnlyIfPrecedingNodeExists {
					sendToBag(nextNode, nextNode.Owner);
				}
			}
			sendToBag(node, node.Owner);
		} else {
			markAsVisited(node);
		}
	}
	err1, err2, err3 := t1.PrepareTransaction(), t2.PrepareTransaction(), t3.PrepareTransaction()
	if !err1 && !err2 && !err3 {
		t1.commit()
		t2.commit()
		t3.commit()
	} else {
		t1.rollback()
		t2.rollback()
		t3.rollback()
	}
	releasePredicateLock(node);
 }

/*
 *	Note about #REF in Mappings:
 *
 *	Some samples: 
 *		"statuses.conversation_id": "#REF(posts.id,posts.id)"
 *		"statuses.picture_id": "#REF(pictures.id,pictures.id)"
 *		"statuses.account_id": "#REF(posts.author_id,people.id)"
 *		"media_attachments.status_id": "#REF(#FETCH(posts.id,posts.guid,photos.status_message_guid),posts.id)"
 *
 *	When the migration worker encounters #REF method in the mapped value of the node it's migrating, it assigns a null value to the attr to be migrated. 
 *	Using the #REF method, it stores a reference in the references table. The first argument in the #REF method represents the FROM part of the references table.
 *	The second argument represents the TO part of the references table. We use the first argument to get the value of ToID in the TO part of the references table.
 *	In case where the migrating node doesn't contain this value, we will use #FETCH method, in place of the first argument, to get the value for ToID of the second
 *	argument.
 *
 *	Assumption:
 *	The second argument (ToReference and ToID) will always be ID/PK of the row that is referenced and not any other attr.
 */
func (t Thread) migrateNode(uid int, srcApp, dstApp string, node *DependencyNode) (int, error) {

	idQuery  := "INSERT INTO identity_table (src_app, dst_app, src_table, dst_table, src_id, dst_id, migration_id) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	refQuery := "INSERT INTO references (app, fromMember, fromID, fromReference, toMember, toID, toReference) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	for srcMember := range node.Members {
		mappings := t.Mappings(srcApp, dstApp, srcMember)
		for dstMemberMapping := range mappings {
			// if mapped attr has #REF method, assign NULL value to that attr when migrating.
			insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", dstMemberMapping.MemberName, dstMemberMapping.Columns, srcMember.Values) 
			dstMemberID := t.DstDB.Query(insertQuery)
			t.StencilDB.Query(idQuery, srcApp, dstApp, srcMember, dstMemberMapping.MemberName, node.srcMember.ID, dstMemberID, t.MigrationID)
			for ref := range dstMemberMapping.references {
				var toID, fromMember, fromID, fromAttr 
				if strings.Contains(ref.FirstArgument, "#FETCH"){
					fromID = t.Mapping.Fetch(ref.FirstArgument)
					fromMember = t.Mapping.Fetch(ref.FirstArgument).Args[0].member
					fromAttr =  t.Mapping.Fetch(ref.FirstArgument).Args[0].attr
					toID = fromID
				} else {
					toID = node.Data[ref.FirstArgument.member][ref.FirstArgument.attr] // node.Data[posts][author_id]
					fromMember = ref.FirstArgument.member
					fromID = node.Data[ref.FirstArgument.member][ID] // node.Data[posts][id]
					fromAttr = ref.FirstArgument.attr
				}
				t.StencilDB.Query(refQuery, t.SrcApp, fromMember, fromID, fromAttr, ref.SecondArgument.member, toID, ref.SecondArgument.attr)	
			}
		}
		deleteQuery := fmt.Sprintf("DELETE FROM %s (%s) WHERE id = %s;", srcMember, node.SrcMember.ID) 
		t.SrcDB.Query(deleteQuery)
	}
	for inDep := range node.Tag.InnerDependencies {
		for Dependee, Dependent := range inDep {
			t.StencilDB.Query(refQuery, t.SrcApp, Dependent.Member, Node.Dependent.ID, Dependent.Attr, Dependee.Member, Node.Dependee.ID, Dependee.Attr)
		}
	}
}

func (t Thread) addToReferences(toNode, fromNode) {
	refQuery  := "INSERT INTO references (app, fromMember, toMember, fromID, toID, fromReference, toReference) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	nullQuery := "UPDATE %s SET %s = NULL WHERE id = '%s'"
	for condition := range t.Dependencies.GetDependencyConditions(toNode, fromNode) {
		fromTable := condition.TagAttr.TableName
		toTable   := condition.DependsOnAttr.TableName
		fromID    := fromNode.GetID(fromTable)
		toID 	  := toNode.GetID(toTable)
		reference := condition.TagAttr
		t.StencilDB.Query(refQuery, t.SrcApp, fromTable, toTable, fromID, toID, reference)
		t.SrcDB.Query(nullQuery, fromTable, reference, fromID)
		fromNode.Data.SetData(reference, nil)
	}
}

func (t Thread) ResolveReferenceByBackTraversal(app, member, id, org_member, org_id) {
	for IDRow := range GetRowsFromIDTableByTo(app, member, id) {
		// You are on the left/from part
		for ref := range t.GetFromReferences(IDRow.FromApp, IDRow.FromMember, IDRow.FromID) {
			// if refIdentityRow := ForwardTraverseIDTable(ref.App, ref.ToMember, ref.ToMemberID, org_member, org_id); refIdentityRow != nil || ref.App == t.DstApp {
			var refIdentityRows [];
			var refID, refMember;
			ForwardTraverseIDTable(ref.App, ref.FromMember, ref.FromMemberID, org_member, org_id, refIdentityRows);
			if len(refIdentityRows) > 0 {
				for refIdentityRow := range refIdentityRows {
					AttrToUpdate:= t.GetMappedAttributeFromSchemaMappings(ref.App, ref.FromMember, ref.FromReference, t.DstApp, org_member) 
					attr := t.GetMappedAttributeFromSchemaMappings(ref.App, ref.ToMember, ref.ToReference, t.DstApp, refIdentityRow.ToMember) 
					updateReferences(ref, refIdentityRow.ToMember, refIdentityRow.ToID, attr, org_member, org_id, AttrToUpdate)
				}
			} else if ref.App == t.DstApp {
				AttrToUpdate:= t.GetMappedAttributeFromSchemaMappings(ref.App, ref.FromMember, ref.FromReference, t.DstApp, org_member) 
				attr := t.GetMappedAttributeFromSchemaMappings(ref.App, ref.ToMember, ref.ToReference, t.DstApp, ref.ToMember) 
				updateReferences(ref, ref.ToMember, ref.ToID, attr, org_member, org_id, AttrToUpdate)
			}
		}
		// You are already on the right/to part
		for ref := range t.GetToReferences(IDRow.FromApp, IDRow.FromMember, IDRow.FromID) {

			// There could be two cases when refIdentityRow is nil: 
			// 1) ref.FromMember is not migrated. 
			// 2) ref.FromMember is migrated, but the currApp of the ref.FromMember is not the t.DstApp		
			var refIdentityRows [];
			var refID, refMember;
			ForwardTraverseIDTable(ref.App, ref.FromMember, ref.FromMemberID, org_member, org_id, refIdentityRows);
			if len(refIdentityRows) > 0 {
				for refIdentityRow := range refIdentityRows {
					AttrToUpdate := t.GetMappedAttributeFromSchemaMappings(ref.App, ref.FromMember, ref.FromReference, t.DstApp, refIdentityRow.ToMember) 
					attr := t.GetMappedAttributeFromSchemaMappings(ref.App, ref.ToMember, ref.ToReference, t.DstApp, org_member) 
					updateReferences(ref, org_member, org_id, attr, refIdentityRow.ToMember, refIdentityRow.ToID, AttrToUpdate)
				}
			} else if ref.App == t.DstApp {
				AttrToUpdate := t.GetMappedAttributeFromSchemaMappings(ref.App, ref.FromMember, ref.FromReference, t.DstApp, ref.ToMember) 
				attr := t.GetMappedAttributeFromSchemaMappings(ref.App, ref.ToMember, ref.ToReference, t.DstApp, org_member) 
				updateReferences(ref, org_member, org_id, attr, refIdentityRow.ToMember, ref.ToID, AttrToUpdate)
			}
		}
		ResolveReferenceByBackTraversal(IDRow.FromApp, IDRow.FromMember, IDRow.FromID, org_member, org_id)
	}
}

func (t Thread) ForwardTraverseIDTable(app, member, id, org_member, org_id, listToBeReturned) {
	IDRows := GetRowsFromIDTableByFrom(app, member, id)
	for _, IDRow := range IDRows {
		ForwardTraverseIDTable(IDRow.ToApp, IDRow.ToMember, IDRow.ToID, org_member, org_id, listToBeReturned)
	}
	if len(IDRows) == 0 {
		if app == t.DstApp && member != org_member {
			listToBeReturned = append(listToBeReturned, []string{tag, member, id})
		}
	}
}

// we use toTable and toID to update fromTable, fromID and fromReference
func (t Thread) updateReferences(refID, member, id, attr, memberToBeUpdated, IDToBeUpdated, AttrToBeUpdated) { 
	// If both attr and AttrToBeUpdated are null that means data cannot be not migrated (because there's no mapping for it in the dst app)
	// For example, conversation and status move from Mastodon to posts in Diaspora: There's a reference from status to conversation in Mastodon,
	// but as conversation cannot be mapped to any member in Mastodon, status.conversation_id and conversation.id cannot be found any corresponding mappings in Diaspora.
	if !attr && !AttrToBeUpdated {
		return
	} else if attr && AttrToBeUpdated {
		data := t.DstDB.GetDataForMember(member, id)
		if data != nil {
			// If the memberToBeUpdated is deleted by the application services, then just go ahead 
			t.DstDB.SetDataForTable(memberToBeUpdated, AttrToBeUpdated, IDToBeUpdated, data.Data[attr])
			t.StencilDB.DeleteReference(refID)
		}
	} else {
		// If either attr or AttrToBeUpdated is null that means the reference has probably been implicitly resolved, which is usually the case in multi-to-one mapping.
		// For example, status and status_stats move from Mastodon to posts in Diaspora: There's a reference from status_stats to statuses in Mastodon,
		// but as they combine to create one single Post object in Diaspora, no reference needs to be resolved. So, we delete these references. If this post was 
		// to be moved back to Diaspora, new reference would be created by following normal migration procedures.
		t.StencilDB.DeleteReference(refID)
	}
}

// func (t Thread) ForwardTraverseIDTable(app, tag, member, id, org_member, org_id, listToBeReturned) {
// 	IDRow := GetRowsFromIDTableByFrom(app, tag, member, id)
// 	// Only find the final app is equal to the dst app, will this function return
// 	if IDRow == nil {
// 		if app == t.DstApp && member != org_member { 
// 			return tag, member, id 
// 		}
// 		return nil
// 	}
// 	return ForwardTraverseIDTable(IDRow.ToApp, IDRow.ToTag, IDRow.ToMember, IDRow.ToID, org_member, org_id, listToBeReturned)
// }

// func (t Thread) resolveReferences(idRow, member, memberID) {
	
// 	for ref := range t.GetFromReferences(idRow.FromApp, idRow.FromMember, idRow.FromID) {
// 		for refIdentityRow := range t.GetRowsFromIdentityTable(ref.App, ref.ToMember, ref.ToMemberID) { 
// 			if refIdentityRow.App == t.DstApp  {
// 				updateReferences(toTable, toID, fromTable, fromID, fromReference)
// 			}
// 		}
// 	}
	
// 	for ref := range t.GetToReferences(idRow.FromApp, idRow.FromMember, idRow.FromID) {
// 		updateReferences(member, memberID, ref.FromMembe, ref.FromID, ref.FromReference)
// 	}
// }

// func (t Thread) handleReferences(app, member, memberID) {
// 	for idRow := range t.GetRowsFromIdentityTable(app, member, memberID) {
// 		resolveReferences(idRow, member, memberID)
// 		// if idRow.FromApp != t.DstApp {
// 		// 	handleReferences(idRow.FromApp, idRow.FromMember, idRow.FromID)
// 		// } else {
// 		// 	resolveReferences(idRow, member, memberID)
// 		// }
// 	}
// }



// func (t Thread) resolveReferences(node) {
// 	for member := range node.Members {
// 		for idRow := range t.GetRowsFromIdentityTable(node.App, member, member.ID) {
// 			if idRow.SrcApp != t.DstApp {continue}
// 			for ref := range t.GetFromReferences(idRow.SrcApp, idRow.SrcMember, idRow.SrcID) {
// 				toData := t.DstDB.GetDataForTable(ref.ToMember, ref.ToID)
// 				if toData != nil {
// 					node.Data[ref.FromReference] = toData.Data[ref.ToReference]
// 					t.StencilDB.DeleteReference(ref)
// 				}
// 			}
// 			break // what if there are multiple
// 		}
// 	}
// }


func (t Thread) AggressiveMigration(userid int, srcApp, dstApp string, node *DependencyNode) {
	try:
		if t.Root == node && !checkUserInApp(userid, dstApp) {
			addUserToApplication(userid, dstApp)
		}
		// randomlyGetAdjacentNode(node) not only gets a migrating user's data but also some data shared by other users
		// this part (randomlyGetUnvisitedAdjacentNode) needs to be changed: combine ownership and dependencies for traversal. root goes to ownership and then to dependencies. 
		// root will not be in dependencies anymore; only in ownerships.
		for child := randomlyGetUnvisitedAdjacentNode(node); child != nil; child = randomlyGetUnvisitedAdjacentNode(node) {
			AggressiveMigration(userid, srcApp, dstApp, child)
		}
		acquirePredicateLock(*node)
		for child := randomlyGetUnvisitedAdjacentNode(node); child != nil; child = randomlyGetUnvisitedAdjacentNode(node) {
			AggressiveMigration(userid, srcApp, dstApp, child)
		}
		// Log before migrating, and migrate nodes belonging to the migrating user and shared to the user
		if node.owner == user_id || node.ShareToUser(user_id) {
			for child := range getDependentNodes(*node) { // ignore mark as visited flag while getting children
				// if parentNodes := getAllParentNodes(child); len(parentNodes) > 1 || node.Rules.DisplayWhatever {
				// 	markAsVisited(*child)
				// }
				if parentNodes := getAllParentNodes(child); len(parentNodes) <= 1 && node.Rules.DisplayOnlyIfParentExists {
					sendToBag(*child, child.Owner)
				}
			}
			migrateNode(*node) 
		} else {
			if parentNodes := getAllParentNodes(*node); len(parentNodes) > 1 || node.Rules.DisplayWhatever {
				markAsVisited(*node)
			} else {
				for child := range GetAllAdjacentNodes(node) {
					if parentNodes := getAllParentNodes(child); len(parentNodes) <= 1 && node.Rules.DisplayOnlyIfParentExists {
						
					}
				}
				sendToBag(*node, node.Owner)
			}
		}
		releasePredicateLock(*node)
	catch NodeNotFound:
		t.releaseAllLocks()
		if t.Root {
			AggressiveMigration(userid, srcApp, dstApp, t.Root)
		}else{
			if checkUserInApp(userid, srcApp){
				removeUserFromApplication(userid, srcApp)
			}
			UpdateMigrationState(userid, srcApp, dstApp)
			log.Println("Congratulations, this migration worker has finished it's job!")
		}
}

//Extending database guarantee to ensure migration correctness (and protect application semantics)
//we are trying not to lock all application service and other users' data and protect application semantics
//Our solution is agnostic to threads or concurrent migration
//we are only interrupting limited things that are going to be migrated and minimizing interruption even not per user basis
//Aggressive and normal: tradeoff between availability and performance(latency) 

//Diaspora person and user table <-> Twitter user table
func (t Thread) NormalMigration(userid int, srcApp, dstApp string, node *DependencyNode) {
	try:
		if t.Root == node && !checkUserInApp(userid, dstApp) {
			addUserToApplication(userid, dstApp)
		}
		for child := randomlyGetAdjacentNode(node); child != nil; child = randomlyGetAdjacentNode(node) {
			NormalMigration(userid, srcApp, dstApp, child)
		}
		acquirePredicateLock(*node)
		if child := randomlyGetAdjacentNode(node); child != nil {
			releasePredicateLock(*node)
			NormalMigration(userid, srcApp, dstApp, child)
		} else {
			migrateNode(*node)
			releasePredicateLock(*node)
		}
	catch NodeNotFound:
		t.releaseAllLocks()
		if t.Root {
			NormalMigration(userid, srcApp, dstApp, t.Root)
		}else{
			if checkUserInApp(userid, srcApp){
				removeUserFromApplication(userid, srcApp)
			}
			UpdateMigrationState(userid, srcApp, dstApp)
			log.Println("Congratulations, this migration worker has finished it's job!")
		}
}

func DisplayController(migrationID int) {
	for migratedNode := GetUndisplayedMigratedData(migrationID); !IsMigrationComplete(migrationID); migratedNode = GetUndisplayedMigratedData(migrationID){
		if migratedNode {
			go CheckDisplay(migratedNode. false)
		}
	}
	// Only Executed After The Migration Is Complete
	// Remaning Migration Nodes:
	// -> The Migrated Nodes In The Destination Application That Still Have Their Migration Flags Raised
	for migratedNode := range GetRemainingMigratedNodes(migrationID){
		go CheckDisplay(migratedNode, true)
	}
}

func CheckDisplay(node *DependencyNode, SecondPhase bool) bool {
	try:
		if AlreadyDisplayed(node) {
			return true
		}
		if t.Root == node.GetParent() {
			Display(node)
			return true
		} else {
			if CheckDisplay(node.GetParent(), SecondPhase) {
				Display(node)
				return true
			}
		}
		if SecondPhase && node.DisplayFlag {
			Display(node)
			return true
		}
		return  false
	catch NodeNotFound:
		return false
}

// func CheckDisplay(oneUndisplayedMigratedData dataStruct, finalRound bool) bool {
// 	dataInNode, err := GetDataInNodeBasedOnDisplaySetting(oneUndisplayedMigratedData)
// 	if dataInNode == nil {
// 		return false
// 	} else {
// 		for _, oneDataInNode := range dataInNode {

// 		}
// 	}
// 	if AlreadyDisplayed(node) {
// 		return true
// 	}
// 	if t.Root == node.GetParent() {
// 		Display(node)
// 		return true
// 	} else {
// 		if CheckDisplay(node.GetParent(), finalRound) {
// 			Display(node)
// 			return true
// 		}
// 	}
// 	if finalRound && node.DisplayFlag {
// 		Display(node)
// 		return true
// 	}
// 	return  false
// }

// func DisplayController(migrationID int) {
// 	for undisplayedMigratedData := GetUndisplayedMigratedData(migrationID);
// 		!CheckMigrationComplete(migrationID);
// 		undisplayedMigratedData = GetUndisplayedMigratedData(migrationID){
// 			for _, oneUndisplayedMigratedData := range undisplayedMigratedData {
// 				CheckDisplay(oneUndisplayedMigratedData, false)
// 			}
// 	}
// 	// Only Executed After The Migration Is Complete
// 	// Remaning Migration Nodes:
// 	// -> The Migrated Nodes In The Destination Application That Still Have Their Migration Flags Raised
// 	for _, oneUndisplayedMigratedData := range GetUndisplayedMigratedData(migrationID){
// 		CheckDisplay(oneUndisplayedMigratedData, true)
// 	}
// }

func mappingExists(srcApp, dstApp string) bool {
	return true
}

func main() {

	userid := 1
	srcApp := "FB"
	dstApp := "TW"
	threads := 100

	if !mappingExists(srcApp, dstApp) {
		log.Fatal("No way to migrate between these.")
		return
	}

	root := getDependencyRootForApp(srcApp)

	for i:=0; i < threads; i++ {
		go AggressiveMigration(userid, srcApp, dstApp, root)
	}

	go DisplayController(migrationID)

	for {

	}
}
