package migrate

func (self *MigrationWorkerV2) IndependentMigration(threadID int) error {

	return self.ConsistentMigration(threadID)
}
