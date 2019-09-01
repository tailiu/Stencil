package main

type Thread struct {
	Locks []PredicateLock
	Root DependencyNode
}

func (t Thread) KeepAliveLocks () {

	for _, lock := range t.Locks{
		lock.ResetTimer()
	}
}


type DependencyNode struct {
	Children []DependencyNode
}

func (DependencyNode) IsValid() bool {
	return true
}

type MigrationMsg struct {
	Code    int
	Message string
}

type User struct {
}

type PredicateLock struct {
	node   DependencyNode
	owners []User
}

func (PredicateLock) Timer() {

}

func (PredicateLock) ResetTimer() {

}

func (PredicateLock) RemoveOwner() {

}

func (PredicateLock) ReleaseLock() {

}

func MMsg(code int, msg string) *MigrationMsg {
	mmsg := new(MigrationMsg)
	mmsg.Code = code
	mmsg.Message = msg
	return mmsg
}

func getDependencyRootForApp(srcApp string) *DependencyNode {

	return new(DependencyNode)
}

func randomlyGetChild(node *DependencyNode) *DependencyNode {

	return new(DependencyNode)
}

func acquirePredicateLock(node DependencyNode) bool {
	return true
}

func releasePredicateLock(node DependencyNode) bool {
	return true
}

func migrateNode(node DependencyNode) bool {

	// some database functionality to migrate the current node
	return true // or false
}

func refreshNodeInformation(node *DependencyNode) {

}

func (node *DependencyNode) GetChildren() []DependencyNode {

	if node == nil {
		return nil
	}

	var children []DependencyNode
	return children
}

func randomize(nodes []DependencyNode) {

}


func (t Thread) SharedDataMigration(migratingUser, sharingUser int, srcApp, dstApp string) {


}

/*
 * Assumptions:
 * 1. Threads don't communicate with each other. Reasons: Simplicity, performance.
 * 2. If threads die, they just restart
 *
 */

func (t Thread) AggressiveMigration(userid int, srcApp, dstApp string, node *DependencyNode) {
	try:
		if t.Root == node && !checkUserInApp(userid, dstApp) {
			addUserToApplication(userid, dstApp)
		}
		// randomlyGetAdjacentNode(node) not only gets a migrating user's data but also some data shared by other users
		for child := randomlyGetUnvisitedAdjacentNode(node); child != nil; child = randomlyGetUnvisitedAdjacentNode(node) {
			AggressiveMigration(userid, srcApp, dstApp, child)
		}
		acquirePredicateLock(*node)
		for child := randomlyGetAdjacentNode(node); child != nil; child = randomlyGetAdjacentNode(node) {
			AggressiveMigration(userid, srcApp, dstApp, child)
		}
		// Log before migrating, and migrate nodes belonging to the migrating user and shared to the user
		if node.owner == user_id || node.ShareToUser(user_id) {
			migrateNode(*node) 
		} else {
			markAsVisited(*node)
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

func GetMigratedData() {
	return data
}

func DisplayController(migrationID int) {
	for migratedNode := GetMigratedNode(migrationID); 
		!IsMigrationComplete(migrationID);  
		migratedNode = GetMigratedNode(migrationID){
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

func CheckDisplay(node *DependencyNode, finalRound bool) bool {
	try:
		if AlreadyDisplayed(node) {
			return true
		}
		if t.Root == node.GetParent() {
			Display(node)
			return true
		} else {
			if CheckDisplay(node.GetParent(), finalRound) {
				Display(node)
				return true
			}
		}
		if finalRound && node.DisplayFlag {
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
