package schema_mappings

import (
	"stencil/config"
	"log"
)

func addMappingsByPSMThroughOnePath(pairwiseMappings *config.SchemaMappings, 
	mappingsPath []string) {
	
	var procMappings []config.ToTable

	srcApp := mappingsPath[0]
	dstApp := mappingsPath[len(mappingsPath) - 1]
	
	for i := 0; i < len(mappingsPath) - 2; i++ {

		currApp := mappingsPath[i]

		nextApp := mappingsPath[i + 1]

		nextNextApp := mappingsPath[i + 2]
		
		log.Println("************* Process Mappings *********************")
		log.Println(currApp, nextApp, nextNextApp)

		firstMappings, err1 := findFromAppToAppMappings(pairwiseMappings, currApp, nextApp)
		
		// This could happen when there is no mapping defined from currApp to nextApp
		if err1 != nil {
			log.Println(err1)
			continue
		}

		secondMappings, err2 := findFromAppToAppMappings(pairwiseMappings, nextApp, nextNextApp)
		
		// This could happen when there is no mapping defined from nextApp to nextNextApp
		if err2 != nil {
			log.Println(err2)
			continue
		}
		
		procMappings = procMappingsByTables(firstMappings, secondMappings)
		
		log.Println(procMappings)
		log.Println("*****************************************************")

	}

	log.Println("++++++++++++++ Construct Mappings +++++++++++++++++++")
	// if srcApp == "twitter" && dstApp == "gnusocial" {
	constructMappingsUsingProcMappings(pairwiseMappings, procMappings, srcApp, dstApp)
	// }
	// if srcApp == "twitter" && dstApp == "gnusocial" {
		// log.Println(pairwiseMappings)
	// }
	log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++")

}

func DeriveMappingsByPSM() (*config.SchemaMappings, error) {

	pairwiseMappings, err := loadPairwiseSchemaMappings()
	if err != nil {
		log.Fatal(err)
	}

	apps := getApplications(pairwiseMappings)

	log.Println(apps)

	// Get all eligible permutations and combinations from one app to another app
	// One such permutation and combination is one path
	mappingsPaths := getMappingsPaths(apps)

	for _, mappingsPath := range mappingsPaths {
		addMappingsByPSMThroughOnePath(pairwiseMappings, mappingsPath)
	}

	writeMappingsToFile(pairwiseMappings)

	return pairwiseMappings, nil

}