/*
 * Query Resolver
 */

package qr

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/lib/pq" // postgres driver
	escape "github.com/tj/go-pg-escape"
)

func (self QR) NewRowId() int32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int31n(2147483647) //9223372036854775807
}

func (self QR) GetPhyMappingForLogicalTable(ltable string) map[string][][]string {

	var phyMap = make(map[string][][]string)

	for _, mapping := range append(self.BaseMappings, self.SuppMappings...) {
		// fmt.Println("mapping", mapping)
		// fmt.Println("-------")
		if strings.EqualFold(ltable, mapping["logical_table"]) {
			ptab := mapping["physical_table"]
			pcol := mapping["physical_column"]
			lcol := mapping["logical_column"]
			var pair []string
			pair = append(pair, pcol, lcol)
			// fmt.Println("pair", pair)
			if _, ok := phyMap[ptab]; ok {
				phyMap[ptab] = append(phyMap[ptab], pair)
			} else {
				phyMap[ptab] = [][]string{pair}
			}
		}
	}
	return phyMap
}

func (self QR) GetBaseMappingForLogicalTable(ltable string) map[string][][]string {

	var phyMap = make(map[string][][]string)

	for _, mapping := range self.BaseMappings {
		if strings.EqualFold(ltable, mapping["logical_table"]) {
			ptab := mapping["physical_table"]
			pcol := mapping["physical_column"]
			lcol := mapping["logical_column"]
			var pair []string
			pair = append(pair, pcol, lcol)
			if _, ok := phyMap[ptab]; ok {
				phyMap[ptab] = append(phyMap[ptab], pair)
			} else {
				phyMap[ptab] = [][]string{pair}
			}
		}
	}

	return phyMap
}

func (self QR) GetPhyTabCol(ltabcol string) (string, string) {

	tab := strings.Trim(strings.Split(ltabcol, ".")[0], " ")
	col := strings.Trim(strings.Split(ltabcol, ".")[1], " ")

	return self.GetPhyTabCol_(tab, col)
}

func (self QR) GetPhyTabCol_(tab, col string) (string, string) {

	phyMap := self.GetPhyMappingForLogicalTable(tab)

	for pt, mapping := range phyMap {
		for _, colmap := range mapping {
			if colmap[1] == col {
				return pt, colmap[0]
			}
		}
	}

	return "", ""
}

func (self QR) PhyUpdateAppIDByRowID(new_app_id, ltab string, rowIDs []string) []string {

	var PQs []string

	if len(rowIDs) <= 0 {
		log.Println("Warning: NO ROWIDS!")
	} else {
		phyMap := self.GetBaseMappingForLogicalTable(strings.ToLower(ltab))
		for pt := range phyMap {
			pq := fmt.Sprintf("UPDATE %s SET app_id = %s WHERE app_id = '%s' AND %s_row_id IN (%s);", pt, escape.Literal(new_app_id), self.AppID, pt[0:4], strings.Join(rowIDs[:], ","))
			PQs = append(PQs, pq)
		}
	}
	return PQs
}
