package migrate_v2

func (self *MigrationWorkerV2) IndependentMigration(threadID int) error {

	return self.ConsistentMigration(threadID)
}
