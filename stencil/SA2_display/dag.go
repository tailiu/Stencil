package SA2_display

import (
	"stencil/config"
	"os"
	"strings"
	"fmt"
	"log"
	"errors"
	"encoding/json"
	"io/ioutil"
)


func ReplaceKey(dag *DAG, tag string, key string) string {

	for _, tag1 := range dag.Tags {

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