package SA1_display

import (
	"database/sql"
	"stencil/config"
	"stencil/db"
	"os"
	"strings"
	"fmt"
	"log"
	"errors"
	"encoding/json"
	"io/ioutil"
)

func ReplaceKey(displayConfig *config.DisplayConfig, tag string, key string) string {

	for _, tag1 := range displayConfig.AppConfig.Tags {

		if tag1.Name == tag {
			// fmt.Println(tag)

			for k, v := range tag1.Keys {

				if k == key {

					member := strings.Split(v, ".")[0]
					
					attr := strings.Split(v, ".")[1]
					
					for k1, table := range tag1.Members {

						if k1 == member {

							return table + "." + attr
						}
					}
				}
			}
		}
	}

	return ""

}

func getRootMemberAttr(dag *DAG) (string, string, error) {

	for _, tag1 := range dag.tags {
		
		if tag1.Name == "root" {

			for k, v := range tag1.Keys {

				if k == "root_id" {

					memberNum := strings.Split(v, ".")[0]
					
					attr := strings.Split(v, ".")[1]
					
					for k1, member := range tag1.Members {

						if k1 == memberNum {

							return member, attr, nil
						}
					}
				}
			}
		}
	}
	
	return "", "", CannotFindRootMemberAttr

}

func getDstUserID(stencilDBConn *sql.DB, appID, appName string, migrationID int, dstDAG *DAG) string {

	dstRootMember, _, err2 := getRootMemberAttr(dstDAG)
	if err2 != nil {
		log.Fatal(err2)
	}

	tableID := getTableIDByTableName(stencilDBConn, dstRootMember, appName)

	// Since the in current settings, there is only one row and the root attribute is always id,
	// we only do in the following way. Note that this is not a generic way.
	query := fmt.Sprint(`SELECT id FROM display_flags WHERE app_id = %s 
		and table_id = %s and migration_id = %d`, appID, tableID, migrationID)

	data, err := db.DataCall1(stencilDBConn, query)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprint(data["id"])

}

func loadDAG(app string) (*DAG, error) {
	
	var dag DAG

	var dconfig string

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(dir, "/stencil/") {
		dconfig = "../config/dependencies/" + app + ".json"
	} else {
		dconfig = "./config/dependencies/" + app + ".json"
	}

	jsonFile, err := os.Open(dconfig)

	if err != nil {
		fmt.Println("Some problem with the file: ")
		fmt.Println(err)
		return nil, errors.New("can't open file")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	json.Unmarshal(byteValue, &dag)

	return &dag, nil

}