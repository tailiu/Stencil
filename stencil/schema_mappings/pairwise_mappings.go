package schema_mappings

import (
	"stencil/config"
	"log"
	"fmt"
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

	srcApp := mappingsPath[0]

	for i := 0; i < len(mappingsPath) - 2; i++ {

		currApp := mappingsPath[i]

		nextApp := mappingsPath[i + 1]

		nextNextApp := mappingsPath[i + 2]
		
		log.Println("^^^^^^^^^^ One Round ^^^^^^^^^^")
		log.Println(currApp, nextApp, nextNextApp)

		if i != 0 {
			mappingsSeq = append(mappingsSeq, nextNextApp)
		}

		if !isAlreadyChecked(mappingsSeq, checkedMappingsPaths) {

			if i == 0 {
				firstMappings, err1 = findFromAppToAppMappings(pairwiseMappings, currApp, nextApp)
			
				// This could happen when there is no mapping defined from currApp to nextApp
				if err1 != nil {
					log.Println(err1)
					log.Println(currApp, nextApp)
					break
				}
	
			}
	
			secondMappings, err2 := findFromAppToAppMappings(pairwiseMappings, nextApp, nextNextApp)
			
			// This could happen when there is no mapping defined from nextApp to nextNextApp
			if err2 != nil {
				log.Println(err2)
				log.Println(nextApp, nextNextApp)
				break
			}

			log.Println(secondMappings)
			
			procMappings = procMappingsByTables(firstMappings, secondMappings)
			
			checkedMappingsPaths = append(checkedMappingsPaths, mappingsSeq)
	
			log.Println(procMappings)
	
			log.Println("++++++++ Construct Mappings +++++++++++")
			
			// Note that construct mappings from the source app to the next next app
			firstMappings = constructMappingsByToTables(pairwiseMappings, 
				procMappings, srcApp, nextNextApp)
			
			log.Println(firstMappings)
			
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

	log.Println("******************************************************************\n")

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
	fmt.Println()

	// Get all eligible permutations and combinations from one app to another app
	// One such permutation and combination is one path
	mappingsPaths := getMappingsPaths(apps)

	// checkedMappingsPaths contains checked and constructed paths
	// it does not contain checked but not constructed paths
	// due to missing mappings from one app to another app

	// Checking paths from short to long seems to be better than the other way around
	// because longer paths probably lead to fewer results and using the fewer results
	// to further get more results will lead to even fewer results.
	// Actually checking order matters!!!
	// This can be researched in the future
	// for i := len(mappingsPaths) - 1; i > -1; i-- {
	
	times := 500
	
	for j := 0; j < times; j++ {

		shuffleSlice(mappingsPaths)
		
		log.Println("Shuffle", j + 1, "time(s)")
		log.Println(mappingsPaths)
		fmt.Println()
		
		var checkedMappingsPaths [][]string

		for i := 0; i < len(mappingsPaths); i++ {

			mappingsPath := mappingsPaths[i]

			if !isAlreadyChecked(mappingsPath, checkedMappingsPaths) {

				checkedMappingsPaths = addMappingsByPSMThroughOnePath(pairwiseMappings, 
					mappingsPath, checkedMappingsPaths)

			} else {

				log.Println("Already checked:")
				log.Println(mappingsPath)
			}

		}
	}

	writeMappingsToFile(pairwiseMappings)

	return pairwiseMappings, nil

}