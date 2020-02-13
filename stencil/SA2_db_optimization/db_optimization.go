package SA2_db_optimization

import (
	"stencil/db"
	"strings"
	"log"
	"fmt"
)

func TruncateSA2Tables() {

	db.STENCIL_DB = "stencil"

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
// Nor can one truncate a partitioned table 
// since all data is stored in partitions
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

func CreatPartitions() {

	maxRowID := 2147483647
	ranges1 := [][]int {
		{1, 7}, 		// aspects
		{7, 9},			// comments
		{9, 10},		// contacts
		{10, 11},		// conversations
		{11, 13},		// messages
		{13, 14},		// notification_actors
		{14, 19},		// notifications
		{19, 20},		// people
		{20, 26},		// photos
		{26, 27},		// posts
		{27, 32},		// profiles
		{32, 35},		// aspect_visibilities
		{35, 39},		// users
		{39, 41},		// conversation_visibilities
		{41, 52},		// likes
		{52, 198},		// all other tables
	}

	subPartitionTables := []int {
		13, 14,
	}

	partitionNum1 := len(ranges1)
	partitionNum2 := 5

	db.STENCIL_DB = "stencil_exp_sa2"

	dbConn := db.GetDBConn(db.STENCIL_DB)

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

		log.Println(query1)

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

			log.Println(query2)
	
			queries = append(queries, query2)

		}

	}

	err := db.TxnExecute(dbConn, queries)
	if err != nil {
		log.Fatal(err)
	}
	
}