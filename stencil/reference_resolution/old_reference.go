package reference_resolution

import (
	"stencil/db"
	"fmt"
	"log"
)

func oldGetFromReferences(refResolutionConfig *RefResolution, 
	IDRow map[string]string) []map[string]interface{} {

	query := fmt.Sprintf(`SELECT * FROM reference_table WHERE app = %s and from_member = %s 
		and from_id = %s and migration_id = %d;`, IDRow["from_app"], 
		IDRow["from_member"], IDRow["from_id"], refResolutionConfig.migrationID)
	
	log.Println(query)
	
	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data

}

func oldGetToReferences(refResolutionConfig *RefResolution, 
	IDRow map[string]string) []map[string]interface{} {

	query := fmt.Sprintf(`SELECT * FROM reference_table WHERE app = %s and to_member = %s 
		and to_id = %s and migration_id = %d;`, IDRow["from_app"], 
		IDRow["from_member"], IDRow["from_id"], refResolutionConfig.migrationID)
	
	data, err := db.DataCall(refResolutionConfig.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data

}

func oldGetDataToUpdateRef(refResolutionConfig *RefResolution, member, id, attr string) string {
	
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = %s",
		attr, member, id)
	
	log.Println(query)
	data, err := db.DataCall1(refResolutionConfig.appDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if data[attr] != nil {
		return fmt.Sprint(data[attr]) 
	} else {
		return ""
	}
	
}

func oldDeleteRef(refID string) string {

	return fmt.Sprintf("DELETE FROM reference_table WHERE pk = %s", refID)

	// err := db.TxnExecute1(refResolutionConfig.StencilDBConn, query)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func oldUpdateDataBasedOnRef(memberToBeUpdated, attrToBeUpdated, IDToBeUpdated, data string) string {
	
	return fmt.Sprintf("UPDATE %s SET %s = %s WHERE id = %s",
		memberToBeUpdated, attrToBeUpdated, data, IDToBeUpdated)

	// err := db.TxnExecute1(refResolutionConfig.AppConfig.DBConn, query)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func oldAddToResolvedReferences(refResolutionConfig *RefResolution, 
	memberToBeUpdated, IDToBeUpdated, attrToBeUpdated, data string) string {
	
	return fmt.Sprintf(`INSERT INTO resolved_references 
		(app, member, id, migration_id, reference, value)
		VALUES (%s, %s, %s, %d, '%s', %s)`, 
		refResolutionConfig.appID, 
		memberToBeUpdated,
		IDToBeUpdated,
		refResolutionConfig.migrationID,
		attrToBeUpdated,
		data)

}


func oldUpdateReferences(refResolutionConfig *RefResolution, 
	refID, member, id, attr, memberToBeUpdated, 
	IDToBeUpdated, attrToBeUpdated string) (string, error) {

	if attr == "" && attrToBeUpdated == "" {
		
		return "", notMigrated

	} else if attr != "" && attrToBeUpdated != "" {
		
		data := getDataToUpdateRef(refResolutionConfig, member, id, attr)
		
		log.Println("data to update other data:", data)

		if data != "" {

			log.Println("---------------------------------------------")

			// Even if the thread crashes after executing q1, the crash
			// does not influence the algorithm because the reference record is still there, 
			// the thread can still try to update it, which does not change the value actually.
			// As long as q2 and q3 can be executed together, which are in the same transaction,
			// the algorithm is still correct. 
			q1 := updateDataBasedOnRef(memberToBeUpdated, attrToBeUpdated, IDToBeUpdated, data)
			
			log.Println(q1)

			err := db.TxnExecute1(refResolutionConfig.appDBConn, q1)
			if err != nil {
				log.Fatal(err)
			}

			var queries []string

			q2 := deleteRef(refID)
			
			log.Println(q2)

			q3 := addToResolvedReferences(
				refResolutionConfig, 
				refResolutionConfig.appTableNameIDPairs[memberToBeUpdated],
				IDToBeUpdated, 
				attrToBeUpdated, 
				data)
			
			log.Println(q3)

			queries = append(queries, q2, q3)
			err1 := db.TxnExecute(refResolutionConfig.stencilDBConn, queries)

			if err1 != nil {
				log.Fatal(err1)
			}

			log.Println("---------------------------------------------")
			
			return data, nil

		} else {

			return "", dataToUpdateOtherDataNotFound

		}
	
	} else {
		
		q1 := deleteRef(refID)

		err := db.TxnExecute1(refResolutionConfig.stencilDBConn, q1)
		if err != nil {
			log.Fatal(err)
		}

		return "", alreadySolved

	}

}