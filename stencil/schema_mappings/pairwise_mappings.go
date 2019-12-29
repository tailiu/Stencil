package schema_mappings

import (
	"stencil/config"
	// "stencil/db"
	"os"
	"log"
	"io/ioutil"
	"encoding/json"

)

func getApplications(pairwiseMappings *config.SchemaMappings) []string {

	var apps []string

	for _, mapping := range pairwiseMappings.AllMappings {

		apps = append(apps, mapping.FromApp)
	}

	return apps

}

func designPairwiseMappings(apps []string)  {



}

func loadPairwiseSchemaMappings() (*config.SchemaMappings, error) {

	var SchemaMappingsObj config.SchemaMappings

	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(dir)

	pairwiseSchemaMappingFile := "./config/app_settings/pairwise_mappings.json"

	jsonFile, err := os.Open(pairwiseSchemaMappingFile)
	if err != nil {
		log.Println(err)
		return &SchemaMappingsObj, CannotOpenPSMFile
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	jsonAsBytes, err1 := ioutil.ReadAll(jsonFile)
	if err1 != nil {
		log.Fatal(err1)
	}

	json.Unmarshal(jsonAsBytes, &SchemaMappingsObj)

	// dbConn := db.GetDBConn(db.STENCIL_DB)

	// log.Println(SchemaMappingsObj)

	// for i, mapping := range SchemaMappingsObj.AllMappings {
	// 	for j, toApp := range mapping.ToApps {
	// 		appID := db.GetAppIDByAppName(dbConn, toApp.Name)
	// 		for k, toAppMapping := range toApp.Mappings {
	// 			for l, toTable := range toAppMapping.ToTables {
	// 				ToTableID, err := db.TableID(dbConn, toTable.Table, appID);
	// 				if  err != nil{
	// 					log.Println("LoadSchemaMappings: Unable to resolve ToTableID for table: ", 
	// 						toTable.Table, toApp.Name, appID)
	// 					log.Fatal(err)
	// 				}
	// 				SchemaMappingsObj.AllMappings[i].ToApps[j].Mappings[k].ToTables[l].TableID = ToTableID
	// 				// fmt.Println(toTable.Table, toApp.Name, appID, ToTableID)
	// 			}
	// 		}
	// 	}
	// }
	// fmt.Println(SchemaMappingsObj.AllMappings[0].ToApps[0].Mappings[0].ToTables)
	// log.Fatal()

	return &SchemaMappingsObj, nil

}

func DeriveMappingsByPSM() (*config.SchemaMappings, error) {

	pairwiseMappings, err := loadPairwiseSchemaMappings()
	if err != nil {
		log.Fatal(err)
	}

	apps := getApplications(pairwiseMappings)

	log.Println(apps)

	designPairwiseMappings(apps)

	return nil, nil

}