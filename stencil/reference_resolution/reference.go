package reference_resolution

import (
	"stencil/db"
	"stencil/config"
	"fmt"
	"log"
)

func getFromReferences(displayConfig *config.DisplayConfig, 
	IDRow map[string]string) []map[string]interface{} {

	query := fmt.Sprintf(`SELECT * FROM reference_table WHERE app = %s and from_member = %s 
		and from_id = %s and migration_id = %d;`, IDRow["from_app"], 
		IDRow["from_member"], IDRow["from_id"], displayConfig.MigrationID)
	
	log.Println(query)
	
	data, err := db.DataCall(displayConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data

}

func getToReferences(displayConfig *config.DisplayConfig, 
	IDRow map[string]string) []map[string]interface{} {

	query := fmt.Sprintf(`SELECT * FROM reference_table WHERE app = %s and to_member = %s 
		and to_id = %s and migration_id = %d;`, IDRow["from_app"], 
		IDRow["from_member"], IDRow["from_id"], displayConfig.MigrationID)
	
	data, err := db.DataCall(displayConfig.StencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data

}

func getDataToUpdateRef(displayConfig *config.DisplayConfig, member, id, attr string) string {
	
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = %s",
		attr, member, id)
	
	log.Println(query)
	data, err := db.DataCall1(displayConfig.AppConfig.DBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if data[attr] != nil {
		return fmt.Sprint(data[attr]) 
	} else {
		return ""
	}
	
}

func deleteRef(refID string) string {

	return fmt.Sprintf("DELETE FROM reference_table WHERE pk = %s", refID)

	// err := db.TxnExecute1(displayConfig.StencilDBConn, query)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func updateDataBasedOnRef(memberToBeUpdated, attrToBeUpdated, IDToBeUpdated, data string) string {
	
	return fmt.Sprintf("UPDATE %s SET %s = %s WHERE id = %s",
		memberToBeUpdated, attrToBeUpdated, data, IDToBeUpdated)

	// err := db.TxnExecute1(displayConfig.AppConfig.DBConn, query)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func addToResolvedReferences(displayConfig *config.DisplayConfig, 
	memberToBeUpdated, IDToBeUpdated, attrToBeUpdated, data string) string {
	
	return fmt.Sprintf(`INSERT INTO resolved_references 
		(app, member, id, migration_id, reference, value)
		VALUES (%s, %s, %s, %d, '%s', %s)`, 
		displayConfig.AppConfig.AppID, 
		memberToBeUpdated,
		IDToBeUpdated,
		displayConfig.MigrationID,
		attrToBeUpdated,
		data)

}


func updateReferences(displayConfig *config.DisplayConfig, 
	refID, member, id, attr, memberToBeUpdated, 
	IDToBeUpdated, attrToBeUpdated string) (string, error) {

	if attr == "" && attrToBeUpdated == "" {
		
		return "", notMigrated

	} else if attr != "" && attrToBeUpdated != "" {
		
		data := getDataToUpdateRef(displayConfig, member, id, attr)
		
		log.Println(data)

		if data != "" {

			// Even if the thread crashes after executing q1, the crash
			// does not influence the algorithm because the reference record is still there, 
			// the thread can still try to update it, which does not change the value actually.
			// As long as q2 and q3 can be executed together, which are in the same transaction,
			// the algorithm is still correct. 
			q1 := updateDataBasedOnRef(memberToBeUpdated, attrToBeUpdated, IDToBeUpdated, data)
			// log.Println(q1)

			err := db.TxnExecute1(displayConfig.AppConfig.DBConn, q1)
			if err != nil {
				log.Fatal(err)
			}

			var queries []string

			q2 := deleteRef(refID)
			
			log.Println(q2)

			q3 := addToResolvedReferences(
				displayConfig, 
				displayConfig.AppConfig.TableNameIDPairs[memberToBeUpdated],
				IDToBeUpdated, 
				attrToBeUpdated, 
				data)
			
			log.Println(q3)

			queries = append(queries, q2, q3)
			err1 := db.TxnExecute(displayConfig.StencilDBConn, queries)

			if err1 != nil {
				log.Fatal(err1)
			}

			return data, nil

		} else {

			return "", dataToUpdateOtherDataNotFound

		}
	
	} else {
		
		q1 := deleteRef(refID)

		err := db.TxnExecute1(displayConfig.StencilDBConn, q1)
		if err != nil {
			log.Fatal(err)
		}

		return "", alreadySolved

	}

}