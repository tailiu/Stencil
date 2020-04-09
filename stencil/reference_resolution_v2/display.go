package reference_resolution_v2

import (
	"log"
	// "fmt"
	// "stencil/db"
	"stencil/schema_mappings"
)

// func GetUpdatedAttributes(refResolutionConfig *RefResolutionConfig, 
// 	ID *Identity) map[string]string {
	
// 	updatedAttrs := make(map[string]string) 

// 	query := fmt.Sprintf(
// 		`select reference, value from resolved_references where app = %s 
// 		and member = %s and id = %s ORDER BY pk`,
// 		refResolutionConfig.appID, ID.member, ID.id)
	
// 	log.Println(query)

// 	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Even though there could be some duplicate keys here since 
// 	// we don't consider migration_id here,
// 	// the ORDER BY in the query will give us the latest resolved values
// 	for _, data1 := range data {

// 		updatedAttrs[fmt.Sprint(data1["reference"])] = 
// 			fmt.Sprint(data1["value"])

// 	}

// 	return updatedAttrs
	
// }