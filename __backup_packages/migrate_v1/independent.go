package migrate_v1

func (self *MigrationWorkerV2) IndependentMigration(threadID int) error {

	return self.ConsistentMigration(threadID)
}
