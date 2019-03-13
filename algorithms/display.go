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
