package evaluation

import (
	"fmt"
	"log"
	"stencil/db"
	"math/rand"
	"strconv"
	"time"
)

func getMigrationIDBySrcUserID(evalConfig *EvalConfig, 
	userID string) string {

	query := fmt.Sprintf(
		`SELECT migration_id FROM migration_registration 
		WHERE user_id = %s`, userID)
	
	result, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(result) != 1 {
		log.Fatal("One user id", userID, "results in more than one migration ids")
	}

	migrationID := fmt.Sprint(result[0]["migration_id"])

	return migrationID

}

func getAllUserIDsInDiaspora(evalConfig *EvalConfig) []string {

	query := fmt.Sprintf(`SELECT id FROM people`)
	
	result, err := db.DataCall(evalConfig.DiasporaDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var userIDs []string

	for _, data1 := range result {
		userIDs = append(userIDs, fmt.Sprint(data1["id"]))
	}

	return userIDs
}

func getDanglingDataSizeOfMigration(evalConfig *EvalConfig, 
	migrationID string) (int64, int64) {

	var size1, size2 int64

	query1 := fmt.Sprintf(`
		SELECT pg_column_size(data), app FROM data_bags WHERE
		migration_id = %s and app = 1`, migrationID)
	
	result1, err := db.DataCall(evalConfig.StencilDBConn, query1)
	if err != nil {
		log.Fatal(err)
	}

	for _, data1 := range result1 {

		appID := fmt.Sprint(data1["app"])

		if appID == "1" {
			size1 += data1["pg_column_size"].(int64)
		} else {
			size2 += data1["pg_column_size"].(int64)
		}
		
	}

	return size1, size2

}

func shuffleSlice(s []string) {
	
	rand.Seed(time.Now().UnixNano())
	
	rand.Shuffle(len(s), func(i, j int) { 
		s[i], s[j] = s[j], s[i] 
	})

}

func getAllMediaSize(evalConfig *EvalConfig) int64 {

	query := fmt.Sprintf(`SELECT id FROM photos`)
	
	result, err := db.DataCall(evalConfig.DiasporaDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var size int64

	for _, data1 := range result {
		
		tmp := fmt.Sprint(data1["id"])
		id, err1 := strconv.Atoi(tmp)
		if err1 != nil {
			log.Fatal(err)
		}

		size += calculateMediaSize(
				evalConfig.DiasporaDBConn, 
				"photos",
				id,
				"1",
			)
	}

	return size

}