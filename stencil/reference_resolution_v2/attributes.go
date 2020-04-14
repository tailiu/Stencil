package reference_resolution_v2

import (
	"database/sql"
	"fmt"
	"log"
	"stencil/db"
	"stencil/common_funcs"
)

/*
 * This new identity file aims to generalize reference resolution to different migrations
 * The functions here DO NOT consider migration id and GetPreviousID DOES NOT restrict
 * from app and from member
 */

func CreateAttribute(app, member, attrName, val string, id ...string) *Attribute {

	var idVal string

	// "-2" indicates id is missing in an attribute
	if len(id) == 0 {
		idVal = "-2"
	} else if len(id) == 1 {
		idVal = id[0]
	} else {
		log.Fatal("Pass multiple ids to an attribute")
	}

	attr := &Attribute{
		app:    	  app,
		member: 	  member,
		attrName:     attrName,
		val:		  val,
		id:			  idVal,
	}

	return attr
}

func createAttributeWithoutID(app, member, attrName, val string) *Attribute {

	attr := &Attribute{
		app:    	  app,
		member: 	  member,
		attrName:     attrName,
		val:		  val,
	}

	return attr
}

func CreateAttributeChangesTable(dbConn *sql.DB) {
	
	var queries []string

	query1 := `CREATE TABLE attribute_changes (
			from_app 		INT8	NOT NULL,
			from_member 	INT8	NOT NULL,
			from_attr 		INT8	NOT NULL,
			from_val 		VARCHAR NOT NULL,
			from_id 		INT8	NOT NULL,
			to_app 			INT8	NOT NULL,
			to_member 		INT8	NOT NULL,
			to_attr 		INT8	NOT NULL,
			to_val 			VARCHAR NOT NULL,
			to_id	 		INT8	NOT NULL,
			migration_id	INT8	NOT NULL,
			pk				SERIAL PRIMARY KEY,
			FOREIGN KEY (from_member) REFERENCES app_tables (pk),
			FOREIGN KEY (to_member) REFERENCES app_tables (pk),
			FOREIGN KEY (from_app) REFERENCES apps (pk),
			FOREIGN KEY (to_app) REFERENCES apps (pk),
			FOREIGN KEY (from_attr) REFERENCES app_schemas (pk),
			FOREIGN KEY (to_attr) REFERENCES app_schemas (pk))`

	query2 := "CREATE INDEX ON attribute_changes(to_app, to_member, to_attr, to_val)"

	query3 := "CREATE INDEX ON attribute_changes(from_app, from_member, from_attr, from_val)"

	query4 := "CREATE INDEX ON attribute_changes(from_app, from_member, from_id)"

	query5 := "CREATE INDEX ON attribute_changes(to_app, to_member, to_id)"

	queries = append(queries, query1, query2, query3, query4, query5)

	for _, query := range queries {
		err := db.TxnExecute1(dbConn, query)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func createAttributeByRefRowToPart(procAttrRow map[string]string) *Attribute {

	attr := CreateAttribute(
		procAttrRow["to_app"],
		procAttrRow["to_member"],
		procAttrRow["to_attr"],
		procAttrRow["to_val"],
		procAttrRow["to_id"],
	)

	return attr
}

func (rr *RefResolution) getRowsFromAttrChangesTableByTo(attr *Attribute) []map[string]interface{} {

	query := fmt.Sprintf(
		`SELECT * FROM attribute_changes WHERE 
		to_app = %s and to_member = %s and to_attr = %s and to_id = %s`,
		attr.app, attr.member, attr.attrName, attr.id,
	)

	log.Println(query)

	data, err := db.DataCall(rr.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data
}

func (rr *RefResolution) getRowsFromAttrChangesTableByFromByAttrVal(attr *Attribute) []map[string]interface{} {

	query := fmt.Sprintf(
		`SELECT * FROM attribute_changes WHERE
		from_app = %s and from_member = %s and from_attr = %s and from_val = '%s'`,
		attr.app, attr.member, attr.attrName, attr.val,
	)

	data, err := db.DataCall(rr.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}

func (rr *RefResolution) getRowsFromAttrChangesTableByFromUsingID(attr *Attribute) []map[string]interface{} {

	query := fmt.Sprintf(
		`SELECT * FROM attribute_changes WHERE
		from_app = %s and from_member = %s and from_attr = %s and from_id = '%s'`,
		attr.app, attr.member, attr.attrName, attr.id,
	)

	data, err := db.DataCall(rr.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(data)

	return data

}
func (rr *RefResolution) forwardTraverseAttrChangesTable(attr, orgAttr *Attribute, inRecurrsion bool) []*Attribute {

	var res []*Attribute
	var attrRows []map[string]interface{}

	if attr.id != "-2" {
		attrRows = rr.getRowsFromAttrChangesTableByFromUsingID(attr)
	} else {
		attrRows = rr.getRowsFromAttrChangesTableByFromByAttrVal(attr)
	}
	
	// log.Println(IDRows)

	for _, attrRow := range attrRows {

		procAttrRow := common_funcs.TransformInterfaceToString(attrRow)

		nextData := createAttributeByRefRowToPart(procAttrRow)

		res = append(res, rr.forwardTraverseAttrChangesTable(nextData, orgAttr, true)...)

	}

	// If recurrsion has not started yet and we cannot find IDRows, then
	// we should directly return null result and 
	// not execute the code in the if block since this could directly return
	// the provided ID to us which is wrong!
	if len(attrRows) == 0 && inRecurrsion {

		// We don't need to test ID.id != orginalID.id becaseu as long as
		// ID.member != orginalID.member, this means that
		// this is different from the original row.
		// ID.id may be the same as orginalID.id in the scenario in which
		// migration does not change ids.
		// We don't find the cases in which ID.member == orginalID.member but
		// ID.id != orginalID.id, however this may happen..
		// Before changing:
		// if ID.app == rr.AppConfig.AppID &&
		// 	ID.member != orginalID.member && ID.id != orginalID.id {
		// if ID.app == rr.appID &&
		// 	(ID.member != orginalID.member || ID.id != orginalID.id) {
		// We remove the conditions: ID.member != orginalID.member || ID.id != orginalID.id
		// because there could be data referencing the same data
		// For example, when resolving tweets.guid in the mapping Twitter.tweets -> Diaspora.posts, 
		// tweets.guid is referring to the id of its own tweets.
		// In this case, we cannot use the condition to filter its own data. 
		// If there are more than one mappings to different tables, 
		// for example, when resolving Diaspora.notification_actors.notification_id, there are two ids:
		// Twitter.notifications.id -> Diaspora.notifications.id and Diaspora.notification_actors.id,
		// we use the third argument in the #REF to point out which table we are referring to 
		// (excluding Twitter.notifications.id -> Diaspora.notifications.id (same as the original data under check)
		// in that example). 
		if attr.app == rr.appID {
			res = append(res, attr)
		}
	}

	return res
}

func (rr *RefResolution) getUpdateToAttrInAttrChangesTableQuery(
	memberName, attrToBeUpdated, attrValToBeUpdated, newAttrVal string) string {
		
	query := fmt.Sprintf(
		`UPDATE attribute_changes SET to_val = %s 
		WHERE to_app = %s and to_member = %s 
		and to_attr = %s and to_val = %s`,
		newAttrVal, 
		rr.appID, 
		rr.appTableNameIDPairs[memberName],
		rr.appAttrNameIDPairs[memberName + ":" + attrToBeUpdated],
		attrValToBeUpdated,
	)

	return query

}

func (rr *RefResolution) GetUpdatedAttributes(member, id string) map[string]string {
	
	updatedAttrs := make(map[string]string) 

	query := fmt.Sprintf(
		`select attr, updated_val from resolved_references where
		app = %s and member = %s and id = %s ORDER BY pk`,
		rr.appID, member, id,
	)
	
	log.Println(query)

	data, err := db.DataCall(rr.stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	// Even though there could be some duplicate keys here since 
	// we don't consider migration_id here,
	// the ORDER BY in the query will give us the latest resolved values
	for _, data1 := range data {
		updatedAttrs[rr.attrIDNamePairs[fmt.Sprint(data1["attr"])]] = fmt.Sprint(data1["updated_val"])
	}

	return updatedAttrs
}