package evaluation

import (
	"fmt"
	"log"
	"stencil/db"
	"strconv"
	"database/sql"
	"strings"
)

func getEdgesCounter(evalConfig *EvalConfig,
	counterStart, counterNum, 
	counterInterval int) []map[string]string {
	
	query1 := fmt.Sprintf(
		`SELECT person_id, edges, nodes FROM dag_counter 
		ORDER BY edges ASC`)
	
	data, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}
	
	var res []map[string]string

	prevEdgeNum := -1

	counter := 0

	for i := counterStart; i < len(data); i++ {

		res1 := make(map[string]string)
		res1["person_id"] = fmt.Sprint(data[i]["person_id"])
		res1["edges"] = fmt.Sprint(data[i]["edges"])
		res1["nodes"] = fmt.Sprint(data[i]["nodes"])

		currEdgeNum, err := strconv.Atoi(res1["edges"])
		if err != nil {
			log.Fatal(err)
		}

		// log.Println("prevEdgeNum", prevEdgeNum)
		// log.Println("curr", currEdgeNum)
		// log.Println(prevEdgeNum + counterInterval < currEdgeNum)

		if prevEdgeNum == -1 {
			prevEdgeNum = currEdgeNum
		} else {
			// log.Println("sum", prevEdgeNum + counterInterval)
			// log.Println("curr", currEdgeNum)
			if prevEdgeNum + counterInterval > currEdgeNum {
				continue
			} else {
				prevEdgeNum = currEdgeNum
			}
		}

		// log.Println("added")

		res = append(res, res1)
		counter += 1

		if counter >= counterNum {
			break
		}
	}

	return res
	
}


func getNodesCounter(evalConfig *EvalConfig,
	counterStart, counterNum, 
	counterInterval int) []map[string]string {

	query1 := fmt.Sprintf(
		`SELECT person_id, nodes FROM dag_counter 
		ORDER BY nodes ASC`)
	
	data, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}
	
	var res []map[string]string

	prevNodeNum := -1

	counter := 0

	for i := counterStart; i < len(data); i++ {
		
		res1 := make(map[string]string)
		res1["person_id"] = fmt.Sprint(data[i]["person_id"])
		res1["edges"] = fmt.Sprint(data[i]["edges"])
		res1["nodes"] = fmt.Sprint(data[i]["nodes"])

		currNodeNum, err := strconv.Atoi(res1["nodes"])
		if err != nil {
			log.Fatal(err)
		}

		if prevNodeNum == -1 {
			prevNodeNum = currNodeNum
		} else {
			if prevNodeNum + counterInterval > currNodeNum {
				continue
			} else {
				prevNodeNum = currNodeNum
			}
		}

		// log.Println("added")

		res = append(res, res1)
		counter += 1

		if counter >= counterNum {
			break
		}
	}

	return res
}

func getEdgesCounterByRange(evalConfig *EvalConfig,
	edgeCounterLeft, edgeCounterRight, 
	num int) []map[string]string {
	
	query1 := fmt.Sprintf(
		`SELECT person_id, edges, nodes FROM dag_counter 
		ORDER BY edges ASC`)
	
	data, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}
	
	var res []map[string]string

	count := 0

	for i := 0; i < len(data); i++ {

		currEdgeNum, err := strconv.Atoi(fmt.Sprint(data[i]["edges"]))
		if err != nil {
			log.Fatal(err)
		}

		if currEdgeNum < edgeCounterLeft || currEdgeNum >= edgeCounterRight {
			continue
		}

		res1 := make(map[string]string)
		res1["person_id"] = fmt.Sprint(data[i]["person_id"])
		res1["edges"] = fmt.Sprint(data[i]["edges"])
		res1["nodes"] = fmt.Sprint(data[i]["nodes"])

		res = append(res, res1)
		count += 1

		if count >= num {
			break
		}
	}

	return res
	
}

func getCounter(evalConfig *EvalConfig) []map[string]string {

	query1 := fmt.Sprintf(
		`SELECT person_id, nodes, edges FROM dag_counter 
		ORDER BY nodes ASC`)
	
	data, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	var res []map[string]string
	
	for _, data1 := range data {

		res1 := make(map[string]string)
		res1["nodes"] = fmt.Sprint(data1["nodes"])
		res1["edges"] = fmt.Sprint(data1["edges"])
		res1["person_id"] = fmt.Sprint(data1["person_id"])

		res = append(res, res1)

	}

	return res

}

func isAlreadyCounted(counted []map[string]string, userID string) bool {

	for _, count1 := range counted {
		if count1["person_id"] == userID {
			return true
		}
	}

	return false

}

func insertDataIntoCounterTableIfNotExist(evalConfig *EvalConfig, 
	table string, data Counter) {

	query1 := fmt.Sprintf(
		`SELECT person_id FROM %s WHERE person_id = '%d'`,
		table, data.UserID,
	)

	// log.Println(query1)

	data1, err := db.DataCall1(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	if data1["person_id"] != nil {
		log.Println("UserID", data.UserID, "has already in the table")
		return
	}

	query2 := fmt.Sprintf(
		`INSERT INTO %s (person_id, edges, nodes) 
		VALUES ('%d', %d, %d)`,
		table, data.UserID, data.Edges, data.Nodes,
	)

	// log.Println(query2)

	err1 := db.TxnExecute1(evalConfig.StencilDBConn, query2)
	if err1 != nil {
		log.Fatal(err1)
	}

}

func is1KDatasetForSA2(table string) bool {
	if strings.Contains(table, "sa2") && strings.Contains(table, "1k") {
		return true
	} else {
		return false
	}
}

func is10KDatasetForSA2(table string) bool {
	if strings.Contains(table, "sa2") && strings.Contains(table, "10k") {
		return true
	} else {
		return false
	}
}

func is100KDatasetForSA2(table string) bool {
	if strings.Contains(table, "sa2") && strings.Contains(table, "100k") {
		return true
	} else {
		return false
	}
}

func is1MDatasetForSA2(table string) bool {
	if strings.Contains(table, "sa2") && strings.Contains(table, "1m") {
		return true
	} else {
		return false
	}
}

func getUserIDsWithSameNodesAcrossDatasets(dbConn *sql.DB, 
	databaseName string) []map[string]string {

	counterTables := map[string]string {
		"diaspora_1k_exp13": "a", 
		"diaspora_10k_exp12": "b", 
		"diaspora_100k_exp11": "c",
		"diaspora_1m_exp10": "d",
	}

	tableAlias, ok := counterTables[databaseName]
	if !ok {
		if is1KDatasetForSA2(databaseName) {
			tableAlias = "a"
		} else if is10KDatasetForSA2(databaseName) {
			tableAlias = "b"
		} else if is100KDatasetForSA2(databaseName) {
			tableAlias = "c"
		} else if is1MDatasetForSA2(databaseName) {
			tableAlias = "d"
		} else {
			log.Fatal("Cannot get data by the provided database name:", databaseName)
		}
		
	}

	query := fmt.Sprintf(
		`SELECT DISTINCT %s.person_id, %s.nodes, %s.edges 
		FROM dag_counter_1K a JOIN dag_counter_10K b ON a.nodes = b.nodes 
		JOIN dag_counter_100K c ON b.nodes = c.nodes 
		JOIN dag_counter_1M d ON c.nodes = d.nodes 
		ORDER BY nodes;`,
		tableAlias, tableAlias, tableAlias,
	)

	log.Println(query)

	data, err := db.DataCall(dbConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return procRes1(data)

}