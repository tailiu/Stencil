package SA2_db_populating

import (
	"stencil/db"
	"stencil/apis"
	"strings"
	"log"
	"fmt"
)

func TruncateSA2Tables() {

	db.STENCIL_DB = "stencil_exp_sa2"

	dbConn := db.GetDBConn(db.STENCIL_DB)

	query1 := `TRUNCATE migration_table`

	query3 := "TRUNCATE "
	
	data := getAllTablesInDBs(dbConn)

	for _, data1 := range data {

		tableName := data1["tablename"]

		if strings.Contains(tableName, "base_") {
			query3 += tableName + ", "
			continue
		}

		if strings.Contains(tableName, "supplementary_") &&
			tableName != "supplementary_tables" {
			query3 += tableName + ", "
			continue
		}

	}
	
	query3 = query3[:len(query3) - 2]

	log.Println(query1)
	log.Println(query3)

	queries := []string{query1, query3} 

	err := db.TxnExecute(dbConn, queries)
	if err != nil {
		log.Fatal(err)
	}

	err1 := dbConn.Close()
	if err1 != nil {
		log.Fatal(err1)
	}

}

func GetTotalRowCountsOfDB() {

	dbName := "diaspora_1000000_template"

	dbConn := db.GetDBConn(dbName)

	defer dbConn.Close()

	data := getAllTablesInDBs(dbConn)
	
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

	dbName := "diaspora_1000000_template"

	dbConn := db.GetDBConn(dbName)
	defer dbConn.Close()

	data := getAllTablesInDBs(dbConn)
	
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
	
	db.STENCIL_DB = "stencil_exp_sa2"

	dbConn := db.GetDBConn(db.STENCIL_DB)
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

	db.STENCIL_DB = "stencil_exp_sa2"

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

	db.STENCIL_DB = "stencil_exp_sa2"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()

	data := getAllTablesInDBs(dbConn)
	
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

// When creating a range partition, the lower bound specified with FROM is an inclusive bound,
// whereas the upper bound specified with TO is an exclusive bound.
func CreatPartitions() {

	maxRowID := 2147483647
	ranges1 := [][]int {
		{1, 7}, 		// 1. aspects
		{7, 9},			// 2. comments
		{9, 10},		// 3. contacts
		{10, 11},		// 4. conversations
		{11, 13},		// 5. messages
		{13, 14},		// 6. notification_actors
		{14, 19},		// 7. notifications
		{19, 20},		// 8. people
		{20, 26},		// 9. photos
		{26, 27},		// 10. posts
		{27, 32},		// 11. profiles
		{32, 35},		// 12. aspect_visibilities
		{35, 39},		// 13. users
		{39, 41},		// 14. conversation_visibilities
		{41, 52},		// 15. likes
		{52, 198},		// 16. all other tables
	}

	subPartitionTables := []int {
		13, 14,
	}

	partitionNum1 := len(ranges1)
	partitionNum2 := 5

	db.STENCIL_DB = "stencil_exp_sa2"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()

	var queries []string
	
	subPartitionTableIDs := make(map[int]int)

	for i := 0; i < partitionNum1; i++ {
		
		var query1 string
		
		rangeStart := ranges1[i][0]
		rangeEnd := ranges1[i][1]

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
		step := maxRowID / partitionNum2

		for j := 0; j < partitionNum2; j ++ {

			if j != partitionNum2 - 1 {
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

	tables := getAllTablesInDBs(dbConn)
	
	var queries1 []string

	for _, t := range tables {
		
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

			query18 := fmt.Sprintf(
				`CREATE INDEX ON %s (app_id, table_id, group_id, row_id, mark_as_delete);`,
				table,
			)

			queries1 = append(queries1,
				query19, query20, query3, query4,
				query5, query6, query7, query8,
				query9, query10, query11, query12,
				query13, query14, query15, query16,
				query17, query18,
			)
		}
	}

	err1 := db.TxnExecute(dbConn, queries1)
	if err1 != nil {
		log.Fatal(err1)
	}

}

func DropPrimaryKeysOfParitions() {

	db.STENCIL_DB = "stencil_exp_sa2"

	dbConn := db.GetDBConn(db.STENCIL_DB)
	defer dbConn.Close()

	var queries []string
	
	tables := getAllTablesInDBs(dbConn)

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

func TruncateUnrelatedTables() {

	db.STENCIL_DB = "stencil_exp_sa2"
	
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

// My machine: people, users
// VM: notifications, profiles
// Blade server: notification_actors
func PopulateSA2Tables() {

	var limit int64 

	db.STENCIL_DB = "stencil_exp_sa2"

	table := "users"
	limit = 2000	

	appName := "diaspora_1000000"
	appID := "1"

	apis.Port(appName, appID, table, limit)

}