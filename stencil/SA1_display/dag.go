package SA1_display

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

func getRootMemberAttr(dag *DAG) (string, string, error) {

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
	
	return "", "", CannotFindRootMemberAttr

}

func GetTableByMemberID(dag *DAG, tagName string, checkedMemberID string) (string, error) {

	for _, tag := range dag.Tags {
		if tag.Name == tagName {
			for memberID, memberTable := range tag.Members {
				if memberID == checkedMemberID {
					return memberTable, nil
				}
			}
		}
	}

	return "", NoTableFound
}
