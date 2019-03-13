func (t Thread) NormalMigration(userid int, srcApp, dstApp string, node *DependencyNode) {
	try:
		for child := randomlyGetChild(node); child != nil; child = randomlyGetChild(node) {
			NormalMigration(userid, srcApp, dstApp, child)
		}
		acquirePredicateLock(*node)
		if child := randomlyGetChild(node); child != nil {
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
			log.Println("Congratulations, this migration worker has finished its job!")
		}
}