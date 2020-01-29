package evaluation

import (
	"fmt"
	"log"
	"stencil/db"
	"math/rand"
	"strconv"
	"time"
)

func getMigrationIDBySrcUserID(evalConfig *EvalConfig, userID string) string {

	query := fmt.Sprintf(
		`SELECT migration_id FROM migration_registration 
		WHERE user_id = %s`, userID)
	
	result, err := db.DataCall(evalConfig.DiasporaDBConn, query)
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

func getDanglingDataSizeOfMigration(evalConfig *EvalConfig, migrationID string) int64 {

	query := fmt.Sprintf(`
		SELECT pg_column_size(data) FROM data_bags WHERE
		migration_id = %s`, migrationID)
	
	result, err := db.DataCall(evalConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	var size int64

	for _, data1 := range result {
		size += data1["pg_column_size"].(int64)
	}

	return size

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