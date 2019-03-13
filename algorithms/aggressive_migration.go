func (t Thread) AggressiveMigration(userid int, srcApp, dstApp string, node *DependencyNode) {
	try:
		for child := randomlyGetChild(node); child != nil; child = randomlyGetChild(node) {
			AggressiveMigration(userid, srcApp, dstApp, child)
		}
		acquirePredicateLock(*node)
		for child := randomlyGetChild(node); child != nil; child = randomlyGetChild(node) {
			AggressiveMigration(userid, srcApp, dstApp, child)
		}
		migrateNode(*node)
		releasePredicateLock(*node)
	catch NodeNotFound:
		t.releaseAllLocks()
		if t.Root {
			AggressiveMigration(userid, srcApp, dstApp, t.Root)
		}else{
			log.Println("Congratulations, this migration worker has finished its job!")
		}
}