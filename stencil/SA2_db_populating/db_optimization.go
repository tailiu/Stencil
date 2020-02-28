package SA2_db_populating

import (
	"stencil/db"
	"strings"
	"log"
	"fmt"
)

func TruncateSA2Tables() {

	dbName := "stencil_exp_sa2_12"

	truncateSA2Tables(dbName)

}

func GetTotalRowCountsOfDB() {

	dbName := "diaspora_100k_exp7"

	dbConn := db.GetDBConn(dbName)

	defer dbConn.Close()

	data := getAllTablesInDB(dbConn)
	
	// log.Println(data)

	var totalRows int64

	for _, data1 := range data {
		
		tableName := data1["tablename"]

		// references table will cause errors
		if tableName == "references" {
			continue
		}

		query2 := fmt.Sprintf(
			`select count(*) as num from %s`, 
			tableName,
		)

		// log.Println(query2)

		res, err := db.DataCall1(dbConn, query2)
		if err != nil {
			log.Fatal(err)
		}

		// log.Println(res)

		totalRows += res["num"].(int64)
		
	}

	log.Println("Total Rows:", totalRows)

}

func ListRowCountsOfDB() {

	dbName := "diaspora_1000000"

	dbConn := db.GetDBConn(dbName)
	defer dbConn.Close()

	data := getAllTablesInDB(dbConn)
	
	// log.Println(data)

	rowCounts := make(map[string]int64)

	for _, data1 := range data {
		
		tableName := data1["tablename"]

		// references table will cause errors
		if tableName == "references" {
			continue
		}

		query2 := fmt.Sprintf(
			`select count(*) as num from %s`, 
			tableName,
		)

		// log.Println(query2)

		res, err := db.DataCall1(dbConn, query2)
		if err != nil {
			log.Fatal(err)
		}

		// log.Println(res)

		rowCounts[tableName] = res["num"].(int64)
		
	}

	for tableName, rowCount := range rowCounts {
		log.Println("Table:", tableName, "rowCount:", rowCount)
	}

}

// Note that Primary keys, index, foreign keys 
// are not supported on a partitioned table
func CreatePartitionedMigrationTable() {
	
	isStencilOnBladeServer := false

	db.STENCIL_DB = "stencil_exp_sa2_1"

	dbConn := db.GetDBConn(db.STENCIL_DB, isStencilOnBladeServer)
	defer dbConn.Close()

	query1 := 
		`CREATE TABLE migration_table (
			row_id	int NOT NULL,
			app_id	int NOT NULL,
			table_id int NOT NULL,
			group_id int NOT NULL,
			user_id int,
			migration_id int,
			mflag int DEFAULT 0,
			mark_as_delete bool DEFAULT false,
			bag bool DEFAULT false,
			copy_on_write bool DEFAULT false,
			created_at timestamp DEFAULT now(),
			updated_at timestamp DEFAULT now()
		) PARTITION BY RANGE (table_id) `
	
	err := db.TxnExecute1(dbConn, query1)
	if err != nil {
		log.Fatal(err)
	}

}

// If you drop the partitioned table, then all partitions of the table 
// will also be dropped.
func DropPartitionedTable() {

	db.STENCIL_DB = "stencil_exp_sa2_1"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()
	
	query := "DROP TABLE migration_table"

	log.Println(query)

	err := db.TxnExecute1(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

}

func DropPartitions() {

	db.STENCIL_DB = "stencil_exp_sa2_1"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()

	data := getAllTablesInDB(dbConn)
	
	query := "DROP TABLE "

	for _, data1 := range data {

		tableName := data1["tablename"]

		if strings.Contains(tableName, "migration_table_") {
			query += tableName + ", "
			continue
		}

	}

	query = query[:len(query) - 2]

	log.Println(query)

	err := db.TxnExecute1(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

}

func CreateConstraintsIndexesOnPartitions() {

	subPartitionTableIDs := map[int]int{
		0: 6,
		1: 7,
	}

	db.STENCIL_DB = "stencil_exp_sa2_100k"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()

	tables := getAllTablesInDB(dbConn)

	for _, t := range tables {
		
		var queries1 []string

		table := t["tablename"]

		if strings.Contains(table, "migration_table_") &&
			table != "migration_table_backup" && 
			!isSubPartitionTable(subPartitionTableIDs, table) {
			
			query19 := fmt.Sprintf(
				`ALTER TABLE %s ADD CONSTRAINT %s_pk 
				PRIMARY KEY (app_id, table_id, group_id, row_id, mark_as_delete);`,
				table, table,
			)

			query20 := fmt.Sprintf(
				`ALTER TABLE %s ADD CONSTRAINT %s_app_tables_fkey
				FOREIGN KEY (table_id) REFERENCES app_tables (pk);`,
				table, table,
			)

			query3 := fmt.Sprintf(
				`ALTER TABLE %s ADD CONSTRAINT %s_apps_fkey
				FOREIGN KEY (app_id) REFERENCES apps (pk);`,
				table, table,
			)

			query4 := fmt.Sprintf(
				`CREATE INDEX ON %s (app_id, table_id, mflag, mark_as_delete);`,
				table,
			)

			query5 := fmt.Sprintf(
				`CREATE INDEX ON %s (app_id, table_id, row_id);`,
				table,
			)

			query6 := fmt.Sprintf(
				`CREATE INDEX ON %s (bag);`,
				table,
			)

			query7 := fmt.Sprintf(
				`CREATE INDEX ON %s (app_id);`,
				table,
			)

			query8 := fmt.Sprintf(
				`CREATE INDEX ON %s (app_id, group_id, row_id);`,
				table,
			)
			
			query9 := fmt.Sprintf(
				`CREATE INDEX ON %s (row_id);`,
				table,
			)

			query10 := fmt.Sprintf(
				`CREATE INDEX ON %s (row_id, group_id);`,
				table,
			)

			query11 := fmt.Sprintf(
				`CREATE INDEX ON %s (table_id);`,
				table,
			)
			
			query12 := fmt.Sprintf(
				`CREATE INDEX ON %s (app_id, group_id, row_id, table_id);`,
				table,
			)
			
			query13 := fmt.Sprintf(
				`CREATE INDEX ON %s (mark_as_delete);`,
				table,
			)

			query14 := fmt.Sprintf(
				`CREATE INDEX ON %s (mflag);`,
				table,
			)

			query15 := fmt.Sprintf(
				`CREATE INDEX ON %s (migration_id);`,
				table,
			)

			query16 := fmt.Sprintf(
				`CREATE INDEX ON %s (group_id);`,
				table,
			)

			query17 := fmt.Sprintf(
				`CREATE INDEX ON %s (user_id);`,
				table,
			)

			// Primary key has already covered this index:
			// query18 := fmt.Sprintf(
			// 	`CREATE INDEX ON %s (app_id, table_id, group_id, row_id, mark_as_delete);`,
			// 	table,
			// )

			queries1 = append(queries1,
				query19, query20, query3, query4,
				query5, query6, query7, query8,
				query9, query10, query11, query12,
				query13, query14, query15, query16,
				query17,
			)

			log.Println("Creating indexes and constraints for table:", table)

			for _, q1 := range queries1 {

				log.Println(q1)

				err1 := db.TxnExecute1(dbConn, q1)
				if err1 != nil {
					log.Fatal(err1)
				}

			}

		}
	}

}

// When creating a range partition, the lower bound specified with FROM is an inclusive bound,
// whereas the upper bound specified with TO is an exclusive bound.
// Creating constraints is optional
// We can add constraints after populating to increase populating speed
func CreatPartitions(createConstrainsts ...bool) {

	partitionNum1 := len(ranges)

	isStencilOnBladeServer := false

	db.STENCIL_DB = "stencil_exp_sa2_1"

	dbConn := db.GetDBConn(db.STENCIL_DB, isStencilOnBladeServer)
	defer dbConn.Close()

	var queries []string
	
	subPartitionTableIDs := make(map[int]int)

	for i := 0; i < partitionNum1; i++ {
		
		var query1 string
		
		rangeStart := ranges[i][0]
		rangeEnd := ranges[i][1]

		if ok := existsInSlice(subPartitionTables, rangeStart); !ok {

			query1 = fmt.Sprintf(
				`CREATE TABLE migration_table_%d PARTITION OF migration_table 
				FOR VALUES FROM ('%d') TO ('%d') `,
				i + 1,
				rangeStart, 
				rangeEnd,
			)

		} else {

			query1 = fmt.Sprintf(
				`CREATE TABLE migration_table_%d PARTITION OF migration_table 
				FOR VALUES FROM ('%d') TO ('%d') 
				PARTITION BY RANGE (row_id) `,
				i + 1,
				rangeStart, 
				rangeEnd,
			)

			subPartitionTableIDs[rangeStart] = i + 1

		}

		// log.Println(query1)

		queries = append(queries, query1)
	
	}

	for _, subPartitionTableID := range subPartitionTableIDs {

		var rangeEnd1 int

		rangeStart1 := 0
		step := maxRowID / subPartionNum

		for j := 0; j < subPartionNum; j ++ {

			if j != subPartionNum - 1 {
				rangeEnd1 = rangeStart1 + step
			} else {
				rangeEnd1 = maxRowID
			}

			query2 := fmt.Sprintf(
				`CREATE TABLE migration_table_sub_%d_%d PARTITION OF migration_table_%d 
				FOR VALUES FROM ('%d') TO ('%d')`,
				subPartitionTableID,
				j + 1,
				subPartitionTableID,
				rangeStart1,
				rangeEnd1,
			)
			
			rangeStart1 = rangeEnd1

			// log.Println(query2)
	
			queries = append(queries, query2)

		}

	}

	err := db.TxnExecute(dbConn, queries)
	if err != nil {
		log.Fatal(err)
	}

	if len(createConstrainsts) == 0 || createConstrainsts[0] {
		CreateConstraintsIndexesOnPartitions()
	}

}

func DropPrimaryKeysOfParitions() {

	db.STENCIL_DB = "stencil_exp_sa2_1"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()

	var queries []string
	
	tables := getAllTablesInDB(dbConn)

	for _, t := range tables {
		
		table := t["tablename"]

		if strings.Contains(table, "migration_table_sub_") {

			query := fmt.Sprintf(
				`ALTER TABLE %s DROP CONSTRAINT %s_pk;`,
				table, table,
			)

			log.Println(query)

			queries = append(queries, query)
		}

	} 

	err := db.TxnExecute(dbConn, queries)
	if err != nil {
		log.Fatal(err)
	}

}

func AddPrimaryKeysToParitions() {

	db.STENCIL_DB = "stencil_exp_sa2_100k"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()

	var queries []string
	
	tables := getAllTablesInDB(dbConn)

	for _, t := range tables {
		
		table := t["tablename"]

		if strings.Contains(table, "migration_table_sub_") {

			query := fmt.Sprintf(
				`ALTER TABLE %s ADD CONSTRAINT %s_pk 
				PRIMARY KEY (app_id, table_id, group_id, row_id, mark_as_delete);`,
				table, table,
			)

			log.Println(query)

			queries = append(queries, query)
		}

	} 

	err := db.TxnExecute(dbConn, queries)
	if err != nil {
		log.Fatal(err)
	}

}

func TruncateUnrelatedTables() {

	db.STENCIL_DB = "stencil"
	
	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()

	query1 := `TRUNCATE identity_table, migration_registration, 
		reference_table, resolved_references, txn_logs,
		evaluation, data_bags, display_flags, display_registration`
	
	err1 := db.TxnExecute1(dbConn, query1)
	if err1 != nil {
		log.Fatal(err1)
	} 

}

func CreateIndexDataTable() {

	db.STENCIL_DB = "stencil_exp_sa2_3"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	
	defer dbConn.Close()

	createIndexDataTable(dbConn)
	
}

func StoreIndexesOfBaseSupTables() {

	db.STENCIL_DB = "stencil_exp_sa2_3"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	
	defer dbConn.Close()

	indexData := getIndexesOfBaseSupTables(dbConn)

	insertIndexDataToTable(dbConn, indexData)

}

func DropIndexesConstraintsOfPartitions() {
	
	db.STENCIL_DB = "stencil_exp_sa2_3"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	
	defer dbConn.Close()

	indexData, constraintData := getConstraintsIndexesOfPartitions(dbConn)

	dropConstraints(dbConn, constraintData)

	dropIndexes(dbConn, indexData)

}

func DropIndexesConstraintsOfBaseSupTables() {

	isStencilOnBladeServer := false

	db.STENCIL_DB = "stencil_exp_sa2_3"

	dbConn := db.GetDBConn(db.STENCIL_DB, isStencilOnBladeServer)
	
	defer dbConn.Close()

	indexData, constraintData := getConstraintsIndexesOfBaseSupTables(dbConn)

	dropConstraints(dbConn, constraintData)

	dropIndexes(dbConn, indexData)

}

func CreateIndexesConstraintsOnBaseSupTables() {
	
	db.STENCIL_DB = "stencil_exp_sa2_100k"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	
	defer dbConn.Close()

	createConstraintsOnBaseSupTables(dbConn)

	createIndexesOfBaseSupTables(dbConn)

}

func DeleteRowsByDuplicateColumnsInMigrationTables() {

	db.STENCIL_DB = "stencil_exp_sa2_100k"

	uniqueCols := []string {
		"app_id", "table_id", "group_id", "row_id", "mark_as_delete",
	}
	
	dbConn := db.GetDBConn(db.STENCIL_DB)
	
	defer dbConn.Close()

	deleteRowsByDuplicateColumnsInMigrationTables(dbConn, uniqueCols)

}

func DeleteRowsByDuplicateColumnsInMigrationTable() {

	db.STENCIL_DB = "stencil_exp_sa2_10k"

	uniqueCols := []string {
		"app_id", "table_id", "group_id", "row_id", "mark_as_delete",
	}
	
	dbConn := db.GetDBConn(db.STENCIL_DB)
	
	defer dbConn.Close()

	deleteRowsByDuplicateColumnsInMigrationTable(dbConn, uniqueCols)

}

func DeleteRowsByDuplicateColumnsInBaseSupTables() {

	db.STENCIL_DB = "stencil_exp_sa2_10k"

	uniqueCols := []string {
		"pk",
	}
	
	dbConn := db.GetDBConn(db.STENCIL_DB)
	
	defer dbConn.Close()

	deleteRowsByDuplicateColumnsInBaseSupTables(dbConn, uniqueCols)

}

func OldDumpAllBaseSupTablesToAnotherDB() {

	srcDB := "stencil_exp_sa2_7"

	dstDB := "stencil_exp_sa2_8"

	db.STENCIL_DB = srcDB

	dbConn := db.GetDBConn(db.STENCIL_DB)
	
	defer dbConn.Close()

	baseSupTables := getAllBaseSupTablesInDB(dbConn)

	query := "pg_dump -U cow -t "

	for i, table := range baseSupTables {

		if i == len(baseSupTables) - 1 {
			query += table + " "
		} else {
			query += table + ","
		}

	}

	query += srcDB + " | " + "psql -U cow " + dstDB

	log.Println(query)

}

func DumpAllBaseSupTablesToAnotherDB() {

	srcDB := "stencil_exp_sa2_10"

	dstDB := "stencil_exp_sa2_100k" 

	query1 := fmt.Sprintf(
		`pg_dump -U cow -a -t supplementary_* --exclude-table-data='supplementary_tables'  %s | psql -U cow %s`,
		srcDB, dstDB,
	)
	
	query2 := fmt.Sprintf(
		`pg_dump -U cow -a -t base_* %s | psql -U cow %s`,
		srcDB, dstDB,
	)
	
	query3 := fmt.Sprintf(
		`pg_dump -U cow -a -t table_name_to_be_replaced %s | psql -U cow %s`,
		srcDB, dstDB, 
	)

	log.Println(query1)

	log.Println(query2)

	log.Println(query3)

}

func CheckpointTruncate() {

	srcDB := "stencil_exp_sa2_10"

	migrationTable := "migration_table_13"

	dstDB := "stencil_exp_sa2_100k" 

	checkpointTruncate(srcDB, dstDB, migrationTable)

}

func DropPrimaryKeysOfSA2TablesWithoutPartitions() {

	db.STENCIL_DB = "stencil_exp_sa2_10k"

	dbConn := db.GetDBConn(db.STENCIL_DB)

	defer dbConn.Close()

	dropPrimaryKeysOfSA2TablesWithoutPartitions(dbConn)

}

func AddPrimaryKeysToSA2TablesWithoutPartitions() {

	db.STENCIL_DB = "stencil_exp_sa2_10k"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()
	
	tables := getAllTablesInDB(dbConn)

	for _, t := range tables {
		
		table := t["tablename"]

		if isMigrationTable(table) {

			query := fmt.Sprintf(
				`ALTER TABLE %s ADD CONSTRAINT %s_pk 
				PRIMARY KEY (app_id, table_id, group_id, row_id, mark_as_delete);`,
				table, table,
			)

			log.Println(query)

			err := db.TxnExecute1(dbConn, query)
			if err != nil {
				log.Fatal(err)
			}

		}

	} 
}

func AddPrimaryKeysToBaseSupTables() {

	db.STENCIL_DB = "stencil_exp_sa2_10k"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()
	
	tables := getAllTablesInDB(dbConn)

	for _, t := range tables {
		
		table := t["tablename"]

		if isBaseOrSupTable(table) {

			query := fmt.Sprintf(
				`ALTER TABLE %s ADD CONSTRAINT %s_pk PRIMARY KEY (pk);`,
				table, table,
			)

			log.Println(query)

			err := db.TxnExecute1(dbConn, query)
			if err != nil {
				log.Fatal(err)
			}

		}

	} 
}