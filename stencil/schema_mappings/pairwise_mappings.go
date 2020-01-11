package schema_mappings

import (
	"stencil/config"
	"log"
)

func addMappingsByPSMThroughOnePath(pairwiseMappings *config.SchemaMappings, 
	mappingsPath []string, checkedMappingsPaths [][]string) [][]string {
	
	// srcApp := mappingsPath[0]
	// dstApp := mappingsPath[len(mappingsPath) - 1]
	
	log.Println("************* Process and Construct Mappings *********************")
	log.Println(mappingsPath)

	var procMappings []config.ToTable

	var firstMappings *config.MappedApp 
	var err1 error
	
	var mappingsSeq = []string {
		mappingsPath[0],
		mappingsPath[1],
		mappingsPath[2],
	}

	// havePath := true

	for i := 0; i < len(mappingsPath) - 2; i++ {

		currApp := mappingsPath[i]

		nextApp := mappingsPath[i + 1]

		nextNextApp := mappingsPath[i + 2]
		
		log.Println("^^^^^^^^^^ One Round ^^^^^^^^^^")
		log.Println(currApp, nextApp, nextNextApp)

		if !isAlreadyChecked(mappingsSeq, checkedMappingsPaths) {

			if i == 0 {
				firstMappings, err1 = findFromAppToAppMappings(pairwiseMappings, currApp, nextApp)
			
				// This could happen when there is no mapping defined from currApp to nextApp
				if err1 != nil {
					log.Println(err1)
					break
				}
	
			}
	
			secondMappings, err2 := findFromAppToAppMappings(pairwiseMappings, nextApp, nextNextApp)
			
			// This could happen when there is no mapping defined from nextApp to nextNextApp
			if err2 != nil {
				log.Println(err2)
				break
			}
			
			procMappings = procMappingsByTables(firstMappings, secondMappings)
			
			if i != 0 {
				mappingsSeq = append(mappingsSeq, nextNextApp)
			}
			
			checkedMappingsPaths = append(checkedMappingsPaths, mappingsSeq)
	
			log.Println(procMappings)
	
			log.Println("++++++++ Construct Mappings +++++++++++")
	
			firstMappings = constructMappingsByToTables(pairwiseMappings, 
				procMappings, currApp, nextNextApp)
			
			log.Println("+++++++++++++++++++++++++++++++++++++++")

		} else {

			log.Println("This round has already been checked:")
			log.Println(mappingsSeq)

			firstMappings, err1 = findFromAppToAppMappings(pairwiseMappings, currApp, nextNextApp)
			
			if err1 != nil {
				log.Fatal(err1)
				break
			}

		}

		log.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	}

	log.Println("******************************************************************")

	// log.Println("++++++++++++++ Construct Mappings +++++++++++++++++++")
	// if srcApp == "twitter" && dstApp == "gnusocial" {
	// constructMappingsUsingProcMappings(pairwiseMappings, procMappings, srcApp, dstApp)
	// }
	// if srcApp == "twitter" && dstApp == "gnusocial" {
		// log.Println(pairwiseMappings)
	// }
	// log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++")

	return checkedMappingsPaths
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

	// checkedMappingsPaths contains checked and constructed paths
	// it does not contain checked but not constructed paths
	// due to missing mappings from one app to another app
	var checkedMappingsPaths [][]string

	for i := len(mappingsPaths) - 1; i > -1; i-- {

		mappingsPath := mappingsPaths[i]

		if !isAlreadyChecked(mappingsPath, checkedMappingsPaths) {

			checkedMappingsPaths = addMappingsByPSMThroughOnePath(pairwiseMappings, 
				mappingsPath, checkedMappingsPaths)

		} else {

			log.Println("Already checked:")
			log.Println(mappingsPath)
		}

	}

	writeMappingsToFile(pairwiseMappings)

	return pairwiseMappings, nil

}