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
		`SELECT person_id, edges FROM dag_counter 
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

