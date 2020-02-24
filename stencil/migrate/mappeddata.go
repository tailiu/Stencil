package migrate

import (
	"errors"
	"fmt"
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

func (self *MappedData) UpdateRefs(fromID, fromMember, fromAttr, toID, toMember, toAttr interface{}) {

	if toID == nil || fromID == nil {
		return
	}

	self.refs = append(self.refs, MappingRef{
		fromID:     fmt.Sprint(fromID),
		fromMember: fmt.Sprint(fromMember),
		fromAttr:   fmt.Sprint(fromAttr),
		toID:       fmt.Sprint(toID),
		toMember:   fmt.Sprint(toMember),
		toAttr:     fmt.Sprint(toAttr)})
}

func GetIDsFromNodeData(firstMember string, secondMember string, nodeData map[string]interface{}) (interface{}, interface{}, error) {
	var toID, fromID interface{}

	color.Yellow.Println("GetIDsFromNodeData: ", firstMember, secondMember)

	if val, ok := nodeData[firstMember]; ok {
		toID = val
	} else {
		return nil, nil, errors.New("Unable to find toID ref value in node data: " + firstMember)
	}

	if val, ok := nodeData[secondMember]; ok {
		toID = val
	}

	firstMemberTokens := strings.Split(firstMember, ".")

	if val, ok := nodeData[firstMemberTokens[0]+".id"]; ok {
		fromID = val
	} else {
		return nil, nil, errors.New("Unable to find fromID ref value in node data: " + firstMemberTokens[0] + ".id")
	}

	return toID, fromID, nil
}

func (self *MappedData) Trim(chars string) {
	self.vals = strings.Trim(self.vals, chars)
	self.cols = strings.Trim(self.cols, chars)
	self.orgCols = strings.Trim(self.orgCols, chars)
}
