package migrate

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gookit/color"
)

func (self *MappedData) UpdateData(col, orgCol, fromTable string, ival interface{}) {
	if ival == nil {
		return
	}
	self.ivals = append(self.ivals, ival)
	self.vals += fmt.Sprintf("$%d,", len(self.ivals))
	self.cols += fmt.Sprintf("%s,", col)
	if orgCol != "" {
		self.orgCols += fmt.Sprintf("%s,", orgCol)
	}
	if fromTable != "" {
		if _, ok := self.srcTables[fromTable]; !ok {
			self.srcTables[fromTable] = []string{strings.Split(orgCol, ".")[1]}
		} else {
			self.srcTables[fromTable] = append(self.srcTables[fromTable], strings.Split(orgCol, ".")[1])
		}
	}
}

func (self *MappedData) UpdateRefs(appID, fromID, fromMember, fromAttr, toID, toMember, toAttr interface{}) {

	if toID == nil || fromID == nil {
		fmt.Println(appID, fromID, fromMember, fromAttr, toID, toMember, toAttr)
		color.Yellow.Printf("@UpdateRefs | Returning | toID : '%v' | fromID: '%v' ", toID, fromID)
		return
	}

	self.refs = append(self.refs, MappingRef{
		appID:      fmt.Sprint(appID),
		fromID:     fmt.Sprint(fromID),
		fromMember: fmt.Sprint(fromMember),
		fromAttr:   fmt.Sprint(fromAttr),
		toID:       fmt.Sprint(toID),
		toMember:   fmt.Sprint(toMember),
		toAttr:     fmt.Sprint(toAttr)})
}

func GetIDsFromNodeData(firstMember, secondMember string, nodeData map[string]interface{}, hardRef bool) (interface{}, interface{}, error) {
	var toID, fromID interface{}

	log.Printf("@GetIDsFromNodeData | Args | firstMember : %s, secondMember : %s, hardRef : %v ", firstMember, secondMember, hardRef)
	log.Printf("@GetIDsFromNodeData | Args | nodeData : %v ", nodeData)

	if hardRef {
		if val, ok := nodeData[secondMember]; ok {
			toID = val
		} else {
			return nil, nil, errors.New("Unable to find toID ref value in node data: " + secondMember)
		}
	} else {
		if val, ok := nodeData[firstMember]; ok {
			toID = val
		} else {
			return nil, nil, errors.New("Unable to find toID ref value in node data: " + firstMember)
		}
	}

	firstMemberTokens := strings.Split(firstMember, ".")

	if val, ok := nodeData[firstMemberTokens[0]+".id"]; ok {
		fromID = val
	} else {
		return nil, nil, errors.New("Unable to find fromID ref value in node data: " + firstMemberTokens[0] + ".id")
	}
	log.Printf("@GetIDsFromNodeData | Returning | toID : %v, fromID : %v ", toID, fromID)
	return toID, fromID, nil
}

func (self *MappedData) Trim(chars string) {
	self.vals = strings.Trim(self.vals, chars)
	self.cols = strings.Trim(self.cols, chars)
	self.orgCols = strings.Trim(self.orgCols, chars)
}
