package evaluation

import (
	"fmt"
	"log"
	"stencil/db"
	"strconv"
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