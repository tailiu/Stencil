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

/*
 * Assumptions:
 * 1. Threads don't communicate with each other. Reasons: Simplicity, performance.
 * 2. If threads die, they just restart
 *
 */

func (t Thread) AggressiveMigration(userid int, srcApp, dstApp string, node *DependencyNode) {
	try:
		for child := randomlyGetChild(node); child != nil; child = randomlyGetChild(node) {
			AggressiveMigration(userid, srcApp, dstApp, child)
		}
		// leaf node ==> len(node.Children) == 0
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
			log.Println("Congratulations, this migration worker has finished it's job!")
		}
}

func (t Thread) NormalMigration(userid int, srcApp, dstApp string, node *DependencyNode) {
	try:
		for child := randomlyGetChild(node); child != nil; child = randomlyGetChild(node) {
			NormalMigration(userid, srcApp, dstApp, child)
		}
		// leaf node ==> len(node.Children) == 0
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
			log.Println("Congratulations, this migration worker has finished it's job!")
		}
}



func mappingExists(srcApp, dstApp string) bool {
	return true
}

func main() {

	userid := 1
	srcApp := "FB"
	dstApp := "TW"

	if !mappingExists(srcApp, dstApp) {
		return
	}

	root := getDependencyRootForApp(srcApp)

	greedyMigration(userid, srcApp, dstApp, root)
}
