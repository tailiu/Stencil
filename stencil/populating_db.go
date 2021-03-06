package main

import (
	"stencil/SA2_db_populating"
)

func main() {

	// SA2_db_populating.TruncateSA2Tables()

	// SA2_db_populating.GetTotalRowCountsOfDB()

	// SA2_db_populating.ListRowCountsOfDB()

	// SA2_db_populating.DropPartitions()

	// SA2_db_populating.CreatePartitionedMigrationTable()

	// SA2_db_populating.CreatPartitions(false)

	// SA2_db_populating.DropPartitionedTable()

	// SA2_db_populating.PopulateRangeOfOneTable()

	// SA2_db_populating.AddPrimaryKeysToParitions()

	// SA2_db_populating.TruncateUnrelatedTables()

	// SA2_db_populating.DropPrimaryKeysOfParitions()

	// SA2_db_populating.CreateIndexDataTable()

	// SA2_db_populating.StoreIndexesOfBaseSupTables()

	// SA2_db_populating.DropIndexesConstraintsOfBaseSupTables()

	// SA2_db_populating.CreateIndexesConstraintsOnBaseSupTables()
	
	// SA2_db_populating.DropIndexesConstraintsOfPartitions()

	// SA2_db_populating.DumpAllBaseSupTablesToAnotherDB()

	// SA2_db_populating.CheckpointTruncate()

	// SA2_db_populating.PupulatingControllerForOneTable()

	// SA2_db_populating.CreateConstraintsIndexesOnPartitions()

	// SA2_db_populating.DeleteRowsByDuplicateColumnsInMigrationTablesInTablePartitioning()

	// SA2_db_populating.DeleteRowsByDuplicateColumnsInBaseSupTables()

	// SA2_db_populating.PupulatingControllerWithCheckpointAndTruncate()

	// SA2_db_populating.DropPrimaryKeysOfSA2Tables()

	// SA2_db_populating.DeleteRowsByDuplicateColumnsInMigrationTable()

	// SA2_db_populating.PupulatingControllerForAllTables()

	// SA2_db_populating.PupulatingControllerForAllTablesHandlingPKs(
	// 	"diaspora_1000000", "stencil_sa2_1m",
	// )

	// SA2_db_populating.PupulatingControllerWithCheckpointAndTruncate(
	// 	"diaspora_1000000", "stencil_sa2_1m_inter", "stencil_sa2_1m", "people",
	// )

	SA2_db_populating.PupulatingControllerWithCheckpointAndTruncateForAllTablesHandlingPKs(
		"diaspora_1000000_sa2_1", "stencil_sa2_1m_inter", "stencil_sa2_1m",
	)
}