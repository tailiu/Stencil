package SA2_db_populating

import (
	"stencil/db"
	"database/sql"
	"strconv"
	"strings"
	"fmt"
	"log"
	SSHClient "github.com/helloyi/go-sshclient"
)

func existsInSlice(s []int, element int) bool {

	for _, v := range s {
		if element == v {
			return true
		}
	}

	return false

} 

func getAllTablesInDB(dbConn *sql.DB) []map[string]string {

	query := `SELECT tablename FROM pg_catalog.pg_tables WHERE 
		schemaname != 'pg_catalog' AND schemaname != 'information_schema';`

	data := db.GetAllColsOfRows(dbConn, query)

	return data

}

func getAllBaseSupTablesInDB(dbConn *sql.DB) []string {

	var allBaseSupTables []string

	allTables := getAllTablesInDB(dbConn)

	for _, t := range allTables {

		table := t["tablename"]

		if !isBaseOrSupTable(table) {
			continue
		}

		allBaseSupTables = append(allBaseSupTables, table)

	}

	return allBaseSupTables

}

func isSubPartitionTable(subPartitionTableIDs map[int]int, table string) bool {

	for _, subPartitionTableID := range subPartitionTableIDs {

		strID := strconv.Itoa(subPartitionTableID)

		subPartitionTable := "migration_table_" + strID

		if table == subPartitionTable {
			return true
		}

	}

	return false

}

func getTotalRowCountOfTable(dbConn *sql.DB, table string) int64 {

	query := fmt.Sprintf(
		`SELECT count(*) as num from %s`, table,
	)

	res, err := db.DataCall1(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// tmp := fmt.Sprint(res["num"])

	// num, err1 := strconv.Atoi(tmp)
	// if err1 != nil {
	// 	log.Fatal(err1)
	// }

	return res["num"].(int64)

}

func isBaseOrSupTable(tableName string) bool {

	if strings.Contains(tableName, "base_") {
		return true
	}

	if strings.Contains(tableName, "supplementary_") &&
		tableName != "supplementary_tables" {
		return true
	}

	return false

}

func getIndexesOfBaseSupTables(dbConn *sql.DB) map[string]string {

	data := getAllTablesInDB(dbConn)

	indexData := make(map[string]string)

	for _, data1 := range data {

		table := data1["tablename"]

		if !isBaseOrSupTable(table) {
			continue
		}

		query1 := fmt.Sprintf(
			"SELECT * FROM pg_indexes WHERE tablename = '%s'",
			table,
		)

		indexes := db.GetAllColsOfRows(dbConn, query1)

		for _, index := range indexes {
			
			key := index["tablename"] + ":" + index["indexname"]
			
			indexData[key] = index["indexdef"]
				
		}
	}

	return indexData
}

func getConstraintsOfBaseSupTables(dbConn *sql.DB) map[string]string {

	constraintData := make(map[string]string)

	query2 := `select conrelid::regclass AS table_from, conname, pg_get_constraintdef(c.oid)
				from pg_constraint c join pg_namespace n ON n.oid = c.connamespace
				where contype in ('f', 'p','c','u') order by table_from`
	
	constraints := db.GetAllColsOfRows(dbConn, query2)

	for _, constraint := range constraints {
		
		table := fmt.Sprint(constraint["table_from"])

		if !isBaseOrSupTable(table) {
			continue
		}

		constraintData[table] = fmt.Sprint(constraint["conname"])
	}

	return constraintData
}

func getConstraintsIndexesOfBaseSupTables(
	dbConn *sql.DB) (map[string]string, map[string]string) {

	indexData := getIndexesOfBaseSupTables(dbConn)

	constraintData := getConstraintsOfBaseSupTables(dbConn)

	return indexData, constraintData

}

func createIndexDataTable(dbConn *sql.DB) {

	query := `CREATE TABLE IF NOT EXISTS indexes_of_b_s_table (
				table_name varchar NOT NULL,
				index_name varchar NOT NULL,
				definition varchar NOT NULL
	  		)`

	err := db.TxnExecute1(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

}

func insertIndexDataToTable(dbConn *sql.DB, indexes map[string]string) {

	var queries []string

	for key, definition := range indexes {

		table := strings.Split(key, ":")[0]
		index := strings.Split(key, ":")[1]

		// log.Println(table)
		// log.Println(index)

		query1 := fmt.Sprintf(
			`INSERT INTO indexes_of_b_s_table (table_name, index_name, definition)
			VALUES ('%s', '%s', '%s')`,
			table, index, definition,
		)

		// log.Println(query1)

		queries = append(queries, query1)
	}

	err := db.TxnExecute(dbConn, queries)
	if err != nil {
		log.Fatal(err)
	}

}

func dropIndexes(dbConn *sql.DB, 
	indexes map[string]string) {

	var queries []string

	for key, _ := range indexes {

		index := strings.Split(key, ":")[1]

		// log.Println(table)
		// log.Println(index)

		query1 := fmt.Sprintf(`DROP INDEX IF EXISTS %s`, index)

		// log.Println(query1)

		queries = append(queries, query1)
	}

	err := db.TxnExecute(dbConn, queries)
	if err != nil {
		log.Fatal(err)
	}

}

func dropConstraints(dbConn *sql.DB, 
	constraints map[string]string) {
	
	var queries []string

	for table, constraint := range constraints {

		query1 := fmt.Sprintf(
			`ALTER TABLE %s DROP CONSTRAINT %s;`, 
			table, constraint,
		)

		// log.Println(query1)

		queries = append(queries, query1)
	}

	err := db.TxnExecute(dbConn, queries)
	if err != nil {
		log.Fatal(err)
	}
	
}

func createConstraintsOnBaseSupTables(dbConn *sql.DB) {

	tables := getAllTablesInDB(dbConn)

	var queries []string

	for _, t := range tables {

		table := t["tablename"]

		if !isBaseOrSupTable(table) {
			continue
		}

		query1 := fmt.Sprintf(
			`ALTER TABLE %s ADD CONSTRAINT %s_pk PRIMARY KEY (pk);`,
			table, table,
		)

		queries = append(queries, query1)
	}

	err := db.TxnExecute(dbConn, queries)
	if err != nil {
		log.Fatal(err)
	}

}

func getAllIndexDefsOfBaseSupFromIndexTable(dbConn *sql.DB) []string {

	query := "SELECT definition FROM indexes_of_b_s_table"

	res, err := db.DataCall(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var defs []string
	for _, res1 := range res {
		defs = append(defs, fmt.Sprint(res1["definition"]))
	}

	return defs

}

func createIndexesOfBaseSupTables(dbConn *sql.DB) {

	indexDefs := getAllIndexDefsOfBaseSupFromIndexTable(dbConn)

	var indexDefsWithoutPK []string
	
	// There is no need to create index on pk again since
	// we have already created a unique index on pk when creating
	// primary key constraint
	for _, indexDef := range indexDefs {
		if strings.Contains(indexDef, "(pk)") {
			continue
		}
		indexDefsWithoutPK = append(indexDefsWithoutPK, indexDef)
	}

	err := db.TxnExecute(dbConn, indexDefsWithoutPK)
	if err != nil {
		log.Fatal(err)
	}

}

func isPartitionTable(table string) bool {

	if strings.Contains(table, "migration_table_") {
		return true
	} else {
		return false
	}

}

func getIndexesOfPartitions(dbConn *sql.DB) map[string]string {

	data := getAllTablesInDB(dbConn)

	indexData := make(map[string]string)

	for _, data1 := range data {

		table := data1["tablename"]

		if !isPartitionTable(table) {
			continue
		}

		query1 := fmt.Sprintf(
			"SELECT * FROM pg_indexes WHERE tablename = '%s'",
			table,
		)

		indexes := db.GetAllColsOfRows(dbConn, query1)

		for _, index := range indexes {
			
			key := index["tablename"] + ":" + index["indexname"]
			
			indexData[key] = index["indexdef"]
				
		}
	}

	return indexData

}

func getConstraintsOfPartitions(dbConn *sql.DB) map[string]string {

	constraintData := make(map[string]string)

	query2 := `select conrelid::regclass AS table_from, conname, pg_get_constraintdef(c.oid)
				from pg_constraint c join pg_namespace n ON n.oid = c.connamespace
				where contype in ('f', 'p','c','u') order by table_from`
	
	constraints := db.GetAllColsOfRows(dbConn, query2)

	for _, constraint := range constraints {
		
		table := fmt.Sprint(constraint["table_from"])

		if !isPartitionTable(table) {
			continue
		}

		constraintData[table] = fmt.Sprint(constraint["conname"])
	}

	return constraintData

}

func getConstraintsIndexesOfPartitions(
	dbConn *sql.DB) (map[string]string, map[string]string) {

	indexData := getIndexesOfPartitions(dbConn)

	constraintData := getConstraintsOfPartitions(dbConn)
		
	return indexData, constraintData

}

func SSHMachineExeCommands(host, port, usersname, password string, cmds []string) {

	client, err := SSHClient.DialWithPasswd(host + ":" + port, usersname, password)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	for _, cmd := range cmds {

		out, err1 := client.Cmd(cmd).Output()
		if err1 != nil {
			log.Fatal(err1)
		}
		fmt.Println(string(out))

	}

}

func dumpAllBaseSupTablesToAnotherDB(srcDB, dstDB string, 
	migrationTables []string) [] string {

	log.Println("Src DB:", srcDB)
	
	log.Println("Dst DB:", dstDB)

	var queries []string

	query1 := fmt.Sprintf(
		`pg_dump -U cow -a -t supplementary_* --exclude-table-data='supplementary_tables'  %s | psql -U cow %s`,
		srcDB, dstDB,
	)
	
	query2 := fmt.Sprintf(
		`pg_dump -U cow -a -t base_* %s | psql -U cow %s`,
		srcDB, dstDB,
	)
	
	queries = append(queries, query1, query2)

	var migrationTableQueries []string

	for _, migrationTable := range migrationTables {

		if migrationTable == "migration_table_6" || migrationTable == "migration_table_7" {

			migrationTableNum := migrationTable[len(migrationTable)-1:]

			for i := 0; i < 5; i ++ {

				subMigrationTableNum := strconv.Itoa(i+1)

				subMigrationTable := "migration_table_sub_" + migrationTableNum +
					 "_" +  subMigrationTableNum

				query3 := fmt.Sprintf(
					`pg_dump -U cow -a -t %s %s | psql -U cow %s`,
					subMigrationTable, srcDB, dstDB, 
				)

				migrationTableQueries = append(migrationTableQueries, query3)

			}

		} else {

			query3 := fmt.Sprintf(
				`pg_dump -U cow -a -t %s %s | psql -U cow %s`,
				migrationTable, srcDB, dstDB, 
			)

			migrationTableQueries = append(migrationTableQueries, query3)
		}

	}

	queries = append(queries, migrationTableQueries...)

	return queries

}

func truncateSA2Tables(dbName string) {

	db.STENCIL_DB = dbName

	dbConn := db.GetDBConn(db.STENCIL_DB)

	query1 := `TRUNCATE migration_table`

	query3 := "TRUNCATE "
	
	data := getAllTablesInDB(dbConn)

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