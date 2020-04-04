package migrate_v2

import (
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

func (self *MappedData) UpdateRefs(appID, fromID, fromMember, fromAttr, toID, toMember, toAttr interface{}) {

	if toID == nil || fromID == nil {
		fmt.Println(appID, fromID, fromMember, fromAttr, toID, toMember, toAttr)
		color.Yellow.Printf("@UpdateRefs | Returning | toID : '%v' | fromID: '%v' ", toID, fromID)
		return
	}

	self.refs = append(self.refs, MappingRef{
		appID:      fmt.Sprint(appID),
		fromVal:    fmt.Sprint(fromID),
		fromMember: fmt.Sprint(fromMember),
		fromAttr:   fmt.Sprint(fromAttr),
		toVal:      fmt.Sprint(toID),
		toMember:   fmt.Sprint(toMember),
		toAttr:     fmt.Sprint(toAttr)})
}

func (self *MappedData) Trim(chars string) {
	self.vals = strings.Trim(self.vals, chars)
	self.cols = strings.Trim(self.cols, chars)
	self.orgCols = strings.Trim(self.orgCols, chars)
}
