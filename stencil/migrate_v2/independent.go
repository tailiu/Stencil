package migrate_v2

func (self *MigrationWorker) IndependentMigration(threadID int) error {

	return self.ConsistentMigration(threadID)
}
