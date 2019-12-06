package main


/*
 * Assumptions:
 * 1. Threads don't communicate with each other. Reasons: Simplicity, performance.
 * 2. If threads die, they just restart
 * 
 * 
 * Independent Migration, Stencil v2:
 * 1. Migration threads neither need to delete data from SrcApp (if data is migrated), nor put data into data bags (if data cannot be migrated).
 * 2. Migration threads only follow ownership relationships to migrate data.
 * 3. Migration threads still need to migrate data from data bags (two phase).
 * 4. If the data cannot be displayed after the display check, delete the data from the destination db, identity table and display_flags. 
 *    Then there will be two cases:
 *    a. If the data is from the source application, delete reference to this data in the references table.
 *    b. If the data is from data bags, put the data back in the bag. 
 * 5. Compared with Stencil v1, there will be "duplicate" data when users migrate back. 
 * 6. We add a column migration_id in the references table and use a unique index involving all columns on this table
 *    to make sure that there are no duplicate rows produced during one migration. So when we deal with concurrent 
 *    independent migrations, we can delete one row after handling that row.
 * 7. Bags are going to be serial. Every migration must acquire some kind of a "right to use bag" or a "lock" and all other subsequent migrations for that user
 *    must wait for the previous migration to finish using the bag before starting to process the bag for itself. This applies to all concurrent migrations for a user.
 *    Migration registration may be used to indicate which migration is using the bag. Probably assign priority numbers and the migration having the lowest number gets to use the bag.
 * 8. Migrated data needs to be marked to avoid being migrated again.
 * 9. Independent migrations can be concurrent?
 *
 *
 * Independent Migration, Stencil v1:
 * 1. Migration threads neither need to delete data from SrcApp (if data is migrated), nor put data into data bags (if data cannot be migrated).
 * 2. Migration threads only follow ownership relationships to migrate data.
 * 3. Migration threads still need to migrate data from data bags.
 * 4. If the data cannot be displayed after the display check, delete the rows pointing to this data in DstApp in Migration Table. 
 *    If the data is from data bags, put the data back in the bag.
 * 5. All migrated data are marked with a copy on write flag. If data is modified, a copy of the data is created in the DstApp.
 * 6. Bags are going to be serial. Every migration must acquire some kind of a "right to use bag" or a "lock" and all other subsequent migrations for that user
 *    must wait for the previous migration to finish using the bag before starting to process the bag for itself. This applies to all concurrent migrations for a user.
 *    Migration registration may be used to indicate which migration is using the bag. Probably assign priority numbers and the migration having the lowest number gets to use the bag.
 * 7. Independent migrations can be concurrent?
 *
 *
 * Consistent Migration, Stencil v1:
 * 1. Migration threads neither need to delete data from SrcApp (if data is migrated), nor put data into data bags (if data cannot be migrated).
 * 2. Migration threads only follow ownership relationships to migrate data.
 * 3. Migration threads still need to migrate data from data bags.
 * 4. If the data cannot be displayed after the display check, delete the rows pointing to this data in DstApp in Migration Table. 
 *    If the data is from data bags, put the data back in the bag.
 * 5. If data is modified in either srcApp or dstApp, the modifications are reflected in both apps.
 * 6. Bags are going to be serial. Every migration must acquire some kind of a "right to use bag" or a "lock" and all other subsequent migrations for that user
 *    must wait for the previous migration to finish using the bag before starting to process the bag for itself. This applies to all concurrent migrations for a user.
 *    Migration registration may be used to indicate which migration is using the bag. Probably assign priority numbers and the migration having the lowest number gets to use the bag.
 * 7. If there is a deletion migration following a consistent migration, Stencil needs to ask user whether to delete consistent data in other applications.
 */

 func (t Thread) OwnershipMigration(uid int, srcApp, dstApp string, root *DependencyNode) {
	 for {
		if node, err := self.GetOwnedNode(); err != nil {
			break; // no owned node found
		} else {
			migrateNode(uid, srcApp, dstApp, node);	
		}
	 }
 }

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
		checkNextNode(node);
		// acquireWriteLock(node)
		for _, precedingNode := range GetAllPrecedingNodes(node) {
			addToReferences(precedingNode, node);
		}
		migrateNode(uid, srcApp, dstApp, node);
		// releaseWriteLock(node);
	} else {
		markAsVisited(node);
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

func (t Thread) checkNextNode(node) {
	// Only through data dependencies
	for _, nextNode := range GetAllNextNodes(node) {
		addToReferences(node, nextNode);
		if precedingNodes := GetAllPrecedingNodes(nextNode); len(precedingNodes) <= 1 && nextNode.Rules.DisplayOnlyIfPrecedingNodeExists {
			acquirePredicateLock(nextNode);
			checkNextNode(nextNode)
			sendNodeToBag(nextNode, nextNode.Owner);
			releasePredicateLock(nextNode);
		}
	}
} 

// This function will run once all threads have finished processing the DAG
// and it will ignore data bags that have been already processed in the same migration
// these bags are referenced in the bags_processed table, which contains the id of the processed bag and the migration_id.
func (t Thread) migrateBags(uid int, srcApp, dstApp string) { 
	idQuery  := "INSERT INTO identity_table (from_app, dst_app, src_table, dst_table, src_id, dst_id, migration_id) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	insertQuery := "INSERT INTO %s (%s) VALUES (%s);"

	for bagRow := range t.GetBagsFromBagTable(uid) { // ignore the processed bags, only consider untouched bags in this migration
		srcMember := bagRow.member
		memberData := bagRow.data
		if bagRow.App == dstApp {
			dstMemberID := t.DstDB.Query(fmt.Sprintf(insertQuery, srcMember, srcMember.Columns, memberData.Values))
			srcMemberIDAttr := node.Tag.FetchMemberAttr(srcMember) // input: statuses, returns: statuses.id
			t.StencilDB.Query(idQuery, srcApp, dstApp, srcMember, srcMember, memberData[srcMemberIDAttr], dstMemberID, t.MigrationID)
			if len(bagRow) == 0 {
				deleteQuery := fmt.Sprintf("DELETE FROM data_bags WHERE id = %s;", memberData.ID) 
				t.SrcDB.Query(deleteQuery)
			}
		} else if mappings := t.Mappings(srcApp, dstApp, srcMember); len(mappings) > 0 {
			for dstMemberMapping := range mappings {
				// if mapped attr has #REF method, assign NULL value to that attr when migrating.
				dstMemberID := t.DstDB.Query(fmt.Sprintf(insertQuery, dstMemberMapping.MemberName, dstMemberMapping.Columns, memberData.Values))
				srcMemberIDAttr := node.Tag.FetchMemberAttr(srcMember) // input: statuses, returns: statuses.id
				t.StencilDB.Query(idQuery, srcApp, dstApp, srcMember, dstMemberMapping.MemberName, memberData[srcMemberIDAttr], dstMemberID, t.MigrationID)
			}
			if len(bagRow) == 0 {
				deleteQuery := fmt.Sprintf("DELETE FROM data_bags WHERE id = %s;", memberData.ID) 
				t.SrcDB.Query(deleteQuery)
			}
		}
	}
}

func (t Thread) sendMemberToBag(member string, data map[string]interface{}, uid int) {
	bagQuery := "INSERT INTO data_bags (user_id, app_id, member, id, data) VALUES ($1,$2,$3,$4,$5); "
	refQuery := "INSERT INTO references (app, fromMember, fromID, fromReference, toMember, toID, toReference) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	db.Query(bagQuery, uid, t.SrcApp, member, data[id], json(data))
	if mappings := t.Mappings(t.srcApp, t.dstApp, member); len(mappings) > 0 {
		for dstMemberMapping := range mappings {
			for ref := range dstMemberMapping.references {
				var toID, fromMember, fromID, fromAttr 
				if strings.Contains(ref.FirstArgument, "#FETCH") {
					fromID = t.Mapping.Fetch(ref.FirstArgument)
					fromMember = t.Mapping.Fetch(ref.FirstArgument).Args[0].member
					fromAttr =  t.Mapping.Fetch(ref.FirstArgument).Args[0].attr
					toID = fromID
				} else {
					toID = data[ref.FirstArgument.member][ref.FirstArgument.attr] // node.Data[posts][author_id]
					fromMember = ref.FirstArgument.member
					fromID = data[ref.FirstArgument.member][ID] // node.Data[posts][id]
					fromAttr = ref.FirstArgument.attr
				}
				t.StencilDB.Query(refQuery, t.SrcApp, fromMember, fromID, fromAttr, ref.SecondArgument.member, toID, ref.SecondArgument.attr)	
			}
		}
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s (%s) WHERE id = %s;", member, data.member.ID) 
	t.SrcDB.Query(deleteQuery)
}

func (t Thread) sendNodeToBag(node *DependencyNode, uid int) {
	refQuery := "INSERT INTO references (app, fromMember, fromID, fromReference, toMember, toID, toReference) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	for member := node.Members {
		t.sendMemberToBag(member.Name, node.member.Data, uid)
	}
	for inDep := range node.Tag.InnerDependencies {
		for Dependee, Dependent := range inDep {
			t.StencilDB.Query(refQuery, t.SrcApp, Dependent.Member, Node.Dependent.ID, Dependent.Attr, Dependee.Member, Node.Dependee.ID, Dependee.Attr)
		}
	}
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

func (t Thread) fetchDataFromBags(uid int, bagData map[string]Interface{}, app, member, id, dstMember) {
	for idRow := range GetRowsFromIDTableByTo(app, member, id) {
		if bagRow := t.GetBagRowFromBagTable(uid, idRow.FromApp, idRow.FromMember, idRow.FromID); bagRow != nil {
			if mappings := t.Mappings(idRow.FromApp, idRow.FromMember, t.DstApp, dstMember); mappings != nil {
				for fromAttr, toAttr := range mappings.mapping {
					if contains(fromAttr, bagRow) {
						if value, exists := bagData[toAttr]; !exists {
							bagData[toAttr] = GetDataFromBagRow(bagRow, fromAttr)
						}
						DeleteDataFromBagRow(bagRow, fromAttr)
						db.AddBagToBagsProcessedTable(t.migrationID, bagRow.id)
					}
				}
				if len(bagRow) == 0 {
					DeleteBagRow(bagRow)
				}
			}
		}
		fetchDataFromBags(uid, bagData, idRow.FromApp, idRow.FromMember, idRow.FromID)
	}
}

func (t Thread) migrateNode(uid int, srcApp, dstApp string, node *DependencyNode) (int,error) {
	refQuery := "INSERT INTO references (app, fromMember, fromID, fromReference, toMember, toID, toReference) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	idQuery  := "INSERT INTO identity_table (src_app, dst_app, src_table, dst_table, src_id, dst_id, migration_id) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	insertQuery := "INSERT INTO %s (%s) VALUES (%s);"
	
	mappings, unmappedMembers := t.FetchMappingsForNode(node)
	var tempNode *DependencyNode = node
	for dstMapping := range mappings {
		var mappedData map[string]Interface{}
		var mappedSrcMembers map[string]bool
		for srcMember := range node.Members {
			for fromAttr, toAttr := range dstMapping.mappings {
				if value, exists := node.Data[fromAttr]; exists {
					mappedData[toAttr] = value
					mappedSrcMembers[srcMember] = true
				}
			}
			fetchDataFromBags(uid, mappedData, t.SrcApp, srcMember, node.Data[fmt.Sprintf("%s.id", srcMember)], dstMapping.MemberName)
		}
		dstMemberID := t.DstDB.Query(fmt.Sprintf(insertQuery, dstMapping.MemberName, mappedData.Columns, mappedData.Values))
		tempNode := DeleteDataFromNode(mappedData.Columns, mappedData.Values)
		for mappedSrcMember := range mappedSrcMembers {
			srcMemberIDAttr := node.Tag.FetchMemberAttr(mappedSrcMember) // input: statuses, returns: statuses.id
			t.StencilDB.Query(idQuery, srcApp, dstApp, mappedSrcMember, dstMapping.MemberName, node.Data[srcMemberIDAttr], dstMemberID, t.MigrationID)
		}

		for ref := range dstMapping.references {
			var toID, fromMember, fromID, fromAttr 
			if strings.Contains(ref.FirstArgument, "#FETCH") {
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
	if t.migrationType == DELETION {
		// partial node sent to bag
		// eg: conversation and statuses are in the same node
		// only statuses has a mapping to dstApp
		// statuses will be migrated and conversation will be sent to bag
		for unmappedMember := range unmappedMembers {
			sendMemberToBag(unmappedMember, unmappedMemberData, uid)
		}
		sendNodeToBag(tempNode)
		for srcMember := range node.Members {
			deleteQuery := fmt.Sprintf("DELETE FROM %s (%s) WHERE id = %s;", srcMember, node.SrcMember.ID) 
			t.SrcDB.Query(deleteQuery)
		}
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
			ForwardTraverseIDTable(ref.App, ref.ToMember, ref.ToMemberID, org_member, org_id, refIdentityRows);
			if len(refIdentityRows) > 0 {
				for refIdentityRow := range refIdentityRows {
					attr := t.GetMappedAttributeFromSchemaMappings(ref.App, ref.ToMember, ref.ToReference, t.DstApp, refIdentityRow.ToMember)
					for AttrToUpdate := range t.GetMappedAttributeFromSchemaMappings(ref.App, ref.FromMember, ref.FromReference, t.DstApp, org_member) {
						updateReferences(ref, refIdentityRow.ToMember, refIdentityRow.ToID, attr, org_member, org_id, AttrToUpdate)
					}
				}
			} else if ref.App == t.DstApp {
				attr := t.GetMappedAttributeFromSchemaMappings(ref.App, ref.ToMember, ref.ToReference, t.DstApp, ref.ToMember) 
				for AttrToUpdate := range t.GetMappedAttributeFromSchemaMappings(ref.App, ref.FromMember, ref.FromReference, t.DstApp, org_member) {
					updateReferences(ref, ref.ToMember, ref.ToID, attr, org_member, org_id, AttrToUpdate)
				}
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
					attr := t.GetMappedAttributeFromSchemaMappings(ref.App, ref.ToMember, ref.ToReference, t.DstApp, org_member) 
					// There could be multiple attributes to update when there are one-to-multple mappings. 
					// For example, when migrating comments to statuses, comments.commentable_id is mapped to both conversation_id and in_reply_to_id.
					// However, when there are multple-to-one mappings, there will be multiple references inserted by migration threads,
					// so we will cope with it in the outer loop "ref" instead of here.
					for AttrToUpdate := range t.GetMappedAttributeFromSchemaMappings(ref.App, ref.FromMember, ref.FromReference, t.DstApp, refIdentityRow.ToMember) {
						updateReferences(ref, org_member, org_id, attr, refIdentityRow.ToMember, refIdentityRow.ToID, AttrToUpdate)
					}
				}
			} else if ref.App == t.DstApp {
				attr := t.GetMappedAttributeFromSchemaMappings(ref.App, ref.ToMember, ref.ToReference, t.DstApp, org_member)
				for AttrToUpdate := range t.GetMappedAttributeFromSchemaMappings(ref.App, ref.FromMember, ref.FromReference, t.DstApp, ref.ToMember) {
					updateReferences(ref, org_member, org_id, attr, refIdentityRow.ToMember, ref.ToID, AttrToUpdate)
				}
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
		if app == t.DstApp && member != org_member && id != org_id {
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

func Display(dstApp string, migrationID int) {
	secondPhase := false

	for migratedData := GetUndisplayedMigratedData(migrationID); 
		!IsMigrationComplete(migrationID);
		migratedData = GetUndisplayedMigratedData(migrationID) {
		
		CheckData(migratedData, secondPhase)

	}

	secondPhase = true
	// Only Executed After The Migration Is Complete
	// Remaning Migration Nodes:
	// -> The Migrated Nodes In The Destination Application That Still Have Their Migration Flags Raised
	for _, migratedNode := range GetUndisplayedMigratedData(migrationID) {
		CheckData(migratedNode, secondPhase)
	}
}

func CheckData(node *DependencyNode, secondPhase bool) {
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
	if secondPhase && node.DisplayFlag {
		Display(node)
		return true
	}
	return  false
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

// func mappingExists(srcApp, dstApp string) bool {
// 	return true
// }

// func main() {

// 	userid := 1
// 	srcApp := "FB"
// 	dstApp := "TW"
// 	threads := 100

// 	if !mappingExists(srcApp, dstApp) {
// 		log.Fatal("No way to migrate between these.")
// 		return
// 	}

// 	root := getDependencyRootForApp(srcApp)

// 	for i:=0; i < threads; i++ {
// 		go AggressiveMigration(userid, srcApp, dstApp, root)
// 	}

// 	go DisplayController(migrationID)

// 	for {

// 	}
// }

// func (t Thread) migrateNode(uid int, srcApp, dstApp string, node *DependencyNode) (int, error) {
// 	refQuery := "INSERT INTO references (app, fromMember, fromID, fromReference, toMember, toID, toReference) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	
// 	for srcMember := range node.Members {
// 		memberData := extractMemberDataFromNodeData(srcMember, node.Data)
// 		idQuery  := "INSERT INTO identity_table (src_app, dst_app, src_table, dst_table, src_id, dst_id, migration_id) VALUES ($1, $2, $3, $4, $5, $6, $7);"
// 		refQuery := "INSERT INTO references (app, fromMember, fromID, fromReference, toMember, toID, toReference) VALUES ($1, $2, $3, $4, $5, $6, $7)"
// 		if mappings := t.Mappings(srcApp, dstApp, srcMember); len(mappings) > 0 {
// 			bagData := fetchDataFromBags(uid, srcMember)
// 			// do something with bag data
// 			for dstMemberMapping := range mappings {
// 				// if mapped attr has #REF method, assign NULL value to that attr when migrating.
// 				insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", dstMemberMapping.MemberName, dstMemberMapping.Columns, memberData.Values) 
// 				dstMemberID := t.DstDB.Query(insertQuery)
// 				srcMemberIDAttr := node.Tag.FetchMemberAttr(srcMember) // input: statuses, returns: statuses.id
// 				t.StencilDB.Query(idQuery, srcApp, dstApp, srcMember, dstMemberMapping.MemberName, node.Data[srcMemberIDAttr], dstMemberID, t.MigrationID)
// 				for ref := range dstMemberMapping.references {
// 					var toID, fromMember, fromID, fromAttr 
// 					if strings.Contains(ref.FirstArgument, "#FETCH") {
// 						fromID = t.Mapping.Fetch(ref.FirstArgument)
// 						fromMember = t.Mapping.Fetch(ref.FirstArgument).Args[0].member
// 						fromAttr =  t.Mapping.Fetch(ref.FirstArgument).Args[0].attr
// 						toID = fromID
// 					} else {
// 						toID = node.Data[ref.FirstArgument.member][ref.FirstArgument.attr] // node.Data[posts][author_id]
// 						fromMember = ref.FirstArgument.member
// 						fromID = node.Data[ref.FirstArgument.member][ID] // node.Data[posts][id]
// 						fromAttr = ref.FirstArgument.attr
// 					}
// 					t.StencilDB.Query(refQuery, t.SrcApp, fromMember, fromID, fromAttr, ref.SecondArgument.member, toID, ref.SecondArgument.attr)	
// 				}
// 			}
// 		} else if t.migrationType == DELETION {
// 			// partial node sent to bag
// 			// eg: conversation and statuses are in the same node
// 			// only statuses has a mapping to dstApp
// 			// statuses will be migrated and conversation will be sent to bag
// 			sendMemberToBag(srcMember, memberData, uid)
// 		}
// 		if t.migrationType == DELETION {
// 			deleteQuery := fmt.Sprintf("DELETE FROM %s (%s) WHERE id = %s;", srcMember, node.SrcMember.ID) 
// 			t.SrcDB.Query(deleteQuery)
// 		}
// 	}
// 	for inDep := range node.Tag.InnerDependencies {
// 		for Dependee, Dependent := range inDep {
// 			t.StencilDB.Query(refQuery, t.SrcApp, Dependent.Member, Node.Dependent.ID, Dependent.Attr, Dependee.Member, Node.Dependee.ID, Dependee.Attr)
// 		}
// 	}
// }