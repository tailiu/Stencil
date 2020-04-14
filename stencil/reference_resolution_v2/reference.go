package reference_resolution_v2

import (
	"fmt"
	"log"
	"database/sql"
	"stencil/common_funcs"
	"stencil/schema_mappings"
	"stencil/db"

	"github.com/gookit/color"
)

/*
 * Similary, this new reference file aims to generalize reference resolution to different migrations
 * The functions here DO NOT consider migration id.
 */

 func CreateReferenceTableV2(dbConn *sql.DB) {

	var queries []string

	query1 := `CREATE TABLE reference_table_v2 (
		app 			INT8	NOT NULL,
		from_id 		INT8	NOT NULL,
		from_member 	INT8	NOT NULL,
		from_attr 		INT8	NOT NULL,
		from_val 		VARCHAR NOT NULL,
		to_member 		INT8	NOT NULL,
		to_attr 		INT8	NOT NULL,
		to_val 			VARCHAR NOT NULL,
		migration_id	INT8	NOT NULL,
		pk				SERIAL PRIMARY KEY,
		FOREIGN KEY (app) REFERENCES apps (pk),
		FOREIGN KEY (from_member) REFERENCES app_tables (pk),
		FOREIGN KEY (to_member) REFERENCES app_tables (pk),
		FOREIGN KEY (from_attr) REFERENCES app_schemas (pk),
		FOREIGN KEY (to_attr) REFERENCES app_schemas (pk))`
	
	query2 := "CREATE INDEX ON reference_table_v2(app, to_member, to_attr, to_val)"

	query3 := "CREATE INDEX ON reference_table_v2(app, from_member, from_attr, from_val, from_id)"

	queries = append(queries, query1, query2, query3)

	for _, query := range queries {
		err := db.TxnExecute1(dbConn, query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func CreateReferenceTableV2WithoutFromID(dbConn *sql.DB) {

	var queries []string

	query1 := `CREATE TABLE reference_table_v2 (
		app 			INT8	NOT NULL,
		from_member 	INT8	NOT NULL,
		from_attr 		INT8	NOT NULL,
		from_val 		VARCHAR NOT NULL,
		to_member 		INT8	NOT NULL,
		to_attr 		INT8	NOT NULL,
		to_val 			VARCHAR NOT NULL,
		migration_id	INT8	NOT NULL,
		pk				SERIAL PRIMARY KEY,
		FOREIGN KEY (app) REFERENCES apps (pk),
		FOREIGN KEY (from_member) REFERENCES app_tables (pk),
		FOREIGN KEY (to_member) REFERENCES app_tables (pk),
		FOREIGN KEY (from_attr) REFERENCES app_schemas (pk),
		FOREIGN KEY (to_attr) REFERENCES app_schemas (pk))`
	
	query2 := "CREATE INDEX ON reference_table_v2(app, to_member, to_attr, to_val)"

	query3 := "CREATE INDEX ON reference_table_v2(app, from_member, from_attr, from_val)"

	queries = append(queries, query1, query2, query3)

	for _, query := range queries {
		err := db.TxnExecute1(dbConn, query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func CreateResolvedReferencesTable(dbConn *sql.DB) {

	var queries []string

	query1 := `CREATE TABLE resolved_references (
		app 						INT8	NOT NULL,
		id		 					INT8	NOT NULL,
		member		 				INT8	NOT NULL,
		attr	 					INT8	NOT NULL,
		updated_val					varchar	NOT NULL,
		migration_id				INT8	NOT NULL,
		pk							SERIAL PRIMARY KEY,
		FOREIGN KEY (app) REFERENCES apps (pk),
		FOREIGN KEY (member) REFERENCES app_tables (pk),
		FOREIGN KEY (attr) REFERENCES app_schemas (pk))`
	
	query2 := "CREATE INDEX ON resolved_references(app, member, attr, id)"

	queries = append(queries, query1, query2)

	for _, query := range queries {
		err := db.TxnExecute1(dbConn, query)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func (rr *RefResolution) getFromReferencesUsingID(attrRow map[string]string) []map[string]interface{} {

	query := fmt.Sprintf(
		`SELECT * FROM reference_table_v2 WHERE
		app = %s and from_member = %s and from_attr = %s and from_val = '%s' and from_id = %s;`,
		attrRow["from_app"], attrRow["from_member"], attrRow["from_attr"], 
		attrRow["from_val"], attrRow["from_id"],
	)

	log.Println(query)

	data, err := db.DataCall(rr.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data
}

func (rr *RefResolution) getFromReferences(attrRow map[string]string) []map[string]interface{} {

	query := fmt.Sprintf(
		`SELECT * FROM reference_table_v2 WHERE
		app = %s and from_member = %s and from_attr = %s and from_val = '%s';`,
		attrRow["from_app"], attrRow["from_member"], 
		attrRow["from_attr"], attrRow["from_val"],
	)

	log.Println(query)

	data, err := db.DataCall(rr.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data
}

func (rr *RefResolution) getToReferences(attrRow map[string]string) []map[string]interface{} {

	query := fmt.Sprintf(
		`SELECT * FROM reference_table_v2 WHERE
		app = %s and to_member = %s and to_attr = %s and to_val = '%s';`,
		attrRow["from_app"], attrRow["from_member"], 
		attrRow["from_attr"], attrRow["from_val"],
	)

	data, err := db.DataCall(rr.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	return data

}

func (rr *RefResolution) checkDataToUpdateRefExists(member, attr, attrVal string) bool {

	query := fmt.Sprintf(
		"SELECT 1 FROM %s WHERE %s = '%s'",
		member, attr, attrVal,
	)

	log.Println(query)

	data, err := db.DataCall1(rr.appDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) != 0 {
		return true
	} else {
		return false
	}

}

func (rr *RefResolution) getRefByPK(pk string) map[string]interface{} {

	query := fmt.Sprintf("SELECT * FROM reference_table_v2 WHERE pk = %s", pk)

	data1, err1 := db.DataCall1(rr.stencilDBConn, query)
	if err1 != nil {
		log.Fatal(err1)
	}

	return data1

}

func deleteRef(refID string) string {

	return fmt.Sprintf("DELETE FROM reference_table_v2 WHERE pk = %s", refID)

}

func updateDataBasedOnRef(memberToBeUpdated, attrToBeUpdated, attrVal, dataID string) string {

	return fmt.Sprintf(
		`UPDATE "%s" SET %s = %s WHERE id = %s`,
		memberToBeUpdated, attrToBeUpdated, attrVal, dataID,
	)

}

func (rr *RefResolution) referenceExists(attr, attMember, attrToUpdate, attrToUpdateMember string) bool {
	
	app := rr.appName

	log.Println("Checking if the reference exists:", app, attr, attMember, attrToUpdate, attrToUpdateMember)
	
	if rr.dag.ReferenceExists(attr, attMember, attrToUpdate, attrToUpdateMember) {
		return true
	} else if schema_mappings.ReferenceExists(rr.allMappings, app, 
		attr, attMember, attrToUpdate, attrToUpdateMember) {
		return true
	} else {
		return false
	}

}

func (rr *RefResolution) addToResolvedReferences(
	memberToBeUpdated, IDToBeUpdated, attrToBeUpdated, attrVal string) string {

	return fmt.Sprintf(
		`INSERT INTO resolved_references 
		(app, member, id, migration_id, attr, updated_val)
		VALUES (%s, %s, %s, %d, %s, '%s')`,
		rr.appID,
		rr.appTableNameIDPairs[memberToBeUpdated],
		IDToBeUpdated,
		rr.migrationID,
		rr.appAttrNameIDPairs[memberToBeUpdated + ":"+ attrToBeUpdated],
		attrVal,
	)

}

func (rr *RefResolution) ReferenceResolved(member, attr, id string) string {

	query := fmt.Sprintf(
		`select updated_val from resolved_references where 
		app = %s and member = %s and attr = %s and id = %s`,
		rr.appID, 
		member, attr, id,
	)
	
	log.Println(query)
	
	data, err := db.DataCall1(rr.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	if len(data) == 0 {
		return ""
	} else {
		return fmt.Sprint(data["updated_val"])
	}
}

func (rr *RefResolution) updateReferences(refID, member, val, attr, 
	memberToBeUpdated, attrValToBeUpdated, attrToBeUpdated, idToBeUpdated string) (string, error) {

	if attr == "" && attrToBeUpdated == "" {

		return "", notMigrated

	} else if attr != "" && attrToBeUpdated != "" {

		newVal := rr.ReferenceResolved(
			rr.appTableNameIDPairs[memberToBeUpdated],
			rr.appAttrNameIDPairs[memberToBeUpdated + ":"+ attrToBeUpdated],
			idToBeUpdated,
		)

		if newVal != "" {
			return "", alreadySolved
		}

		if rr.checkDataToUpdateRefExists(member, attr, val) {

			log.Println("---------------------------------------------")

			log.Println("Update references:")

			var q0 []string

			// Note that for now when "id" needs to be updated,
			// we don't consider what will happen if the thread crashes
			if attrToBeUpdated == "id" {

				// This is to update id in the display_flags table
				displayFlagsQ0 := rr.getUpdateIDInDisplayFlagsQuery(
					memberToBeUpdated, attrValToBeUpdated, val,
				)

				// This is to update to_id in the identity table
				// There is no need to update id in the reference table
				// since the reference table stores ids in the source app
				updateToAttrQ0 := rr.getUpdateToAttrInAttrChangesTableQuery(
					memberToBeUpdated, attrToBeUpdated, 
					attrValToBeUpdated, val,
				)

				insertIntoIDChanges := rr.getInsertIntoIDChangesTableQuery(
					memberToBeUpdated, attrValToBeUpdated, val,
				)

				// Note that this data must not be used to update other data, because
				// in that case, the updated value of other data is stale value
				// and the data could already be displayed or put into data bags
				// which is incorrect.

				q0 = append(q0, displayFlagsQ0, updateToAttrQ0, insertIntoIDChanges)

			}

			// Even if the thread crashes after executing q1, the crash
			// does not influence the algorithm because the reference record is still there,
			// the thread can still try to update it, which does not change the value actually.
			// As long as q2 and q3 can be executed together, which are in the same transaction,
			// the algorithm is still correct.
			q1 := updateDataBasedOnRef(memberToBeUpdated, attrToBeUpdated, val, idToBeUpdated)

			log.Println(q1)

			err := db.TxnExecute1(rr.appDBConn, q1)
			if err != nil {
				log.Fatal(err)
			}

			var queries []string

			q2 := deleteRef(refID)

			var q3 string

			if attrToBeUpdated == "id" {
				q3 = rr.addToResolvedReferences(
					memberToBeUpdated,
					val,
					attrToBeUpdated,
					val,
				)
			} else {
				q3 = rr.addToResolvedReferences(
					memberToBeUpdated,
					idToBeUpdated,
					attrToBeUpdated,
					val,
				)
			}

			queries = append(queries, q2, q3)
			queries = append(queries, q0...)

			red := color.FgRed.Render

			for j, updateQ := range queries {

				if j != 0 {
					log.Println(updateQ)
				} else {

					refToBeDeleted := rr.getRefByPK(refID)

					if len(refToBeDeleted) == 0 {
						log.Println(red("The reference has already been deleted by other display threads"))
					} else {
						log.Println(red("The reference to be deleted:"))
						profRefToBeDeleted := common_funcs.TransformInterfaceToString(refToBeDeleted)
						refToBeDeletedLog := rr.LogRefRow(profRefToBeDeleted, true)
						log.Println(red(refToBeDeletedLog))
					}
				}
			}

			err1 := db.TxnExecute(rr.stencilDBConn, queries)

			if err1 != nil {
				log.Fatal(err1)
			}

			log.Println("---------------------------------------------")

			return val, nil

		} else {

			return "", dataToUpdateOtherDataNotFound

		}

	} else {

		q1 := deleteRef(refID)

		log.Println("---------------------------------------------")
		log.Println("Already resolved and delete the reference:")
		color.Red.Println(q1)
		log.Println("---------------------------------------------")

		err := db.TxnExecute1(rr.stencilDBConn, q1)
		if err != nil {
			log.Fatal(err)
		}

		return "", alreadySolved

	}

}
