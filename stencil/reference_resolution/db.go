package reference_resolution

// import (
// 	"log"
// 	"fmt"
// 	"stencil/db"
// 	"stencil/schema_mappings"
// )

// func getDataFromDB1(refResolutionConfig *RefResolutionConfig,
// 	dbConn *sql.DB, query) map[string]interface{} {

// 	data, err := db.DataCall1(dbConn, query)
// 	if err != nil {
// 		if errMsg := err.Error(); strings.Contains(errMsg, "connect: connection timed out") {
// 			log.Println(err)
// 			reconnectToDB(refResolutionConfig)
// 		} else {
// 			log.Fatal(err)
// 		}
// 	}

// 	return data

// }

// func reconnectToDB(refResolutionConfig *RefResolutionConfig) {

// }