package SA1_display

import (
	"stencil/common_funcs"
	"strings"
)

func getRootMemberAttr(dag *common_funcs.DAG) (string, string, error) {

	for _, tag1 := range dag.Tags {
		
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
	
	return "", "", common_funcs.CannotFindRootMemberAttr

}
