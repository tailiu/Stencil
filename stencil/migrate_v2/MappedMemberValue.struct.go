package migrate_v2

import (
	"fmt"
	"log"
	"stencil/db"
	"stencil/helper"
	"strings"
)

func (mmv MappedMemberValue) GetTableID(tableName, appID string) string {
	if len(tableName) == 0 || len(appID) == 0 {
		log.Fatalf("@mmv.GetTableID: Table ID or App ID nil | table: '%s' | App: '%s' ", tableName, appID)
	} else if mmv.DBConn == nil {
		log.Fatal("@mmv.GetTableID: DBConn not set!")
	} else {
		if tableID, err := db.TableID(mmv.DBConn, tableName, appID); err == nil {
			return tableID
		} else {
			fmt.Println(tableName, appID)
			log.Fatal("@mmv.GetTableID: ", err)
		}
	}
	err := fmt.Errorf("This is the end. Shouldn't be here. Table: '%s', App: '%s'", tableName, appID)
	log.Fatal("@mmv.GetTableID: ", err)
	return err.Error()
}

func (mmv *MappedMemberValue) StoreMemberAndAttr(mappedStmt string) {
	tokens := strings.Split(mappedStmt, ".")
	mmv.FromMember = tokens[0]
	mmv.FromAttr = tokens[1]
	mmv.FromMemberID = mmv.GetTableID(mmv.FromMember, mmv.AppID)
}

func (mmv MappedMemberValue) GetMemberAttr() string {
	return fmt.Sprintf("%s.%s", mmv.FromMember, mmv.FromAttr)
}

func (mmv *MappedMemberValue) CreateReference(fromAttr, toAttr, mappedStmt string, dataMap DataMap) error {

	hardRef := false
	if strings.Contains(mappedStmt, "#REFHARD") {
		hardRef = true
	}

	if toVal, fromVal, fromID, err := dataMap.GetRefValsFromDataMap(fromAttr, toAttr, hardRef); err == nil {
		fromAttrTokens := strings.Split(fromAttr, ".")
		toAttrTokens := strings.Split(toAttr, ".")
		mmv.Ref = &MappingRef{
			appID:        mmv.AppID,
			fromID:       fromID,
			fromVal:      fromVal,
			fromMemberID: mmv.GetTableID(fromAttrTokens[0], mmv.AppID),
			fromMember:   fromAttrTokens[0],
			fromAttr:     fromAttrTokens[1],
			toVal:        toVal,
			toMemberID:   mmv.GetTableID(toAttrTokens[0], mmv.AppID),
			toMember:     toAttrTokens[0],
			toAttr:       toAttrTokens[1],
		}
	} else {
		return err
	}
	return nil
}

func (mmv *MappedMemberValue) SetFromID(fromAttr string, dataMap DataMap) error {
	fromAttrTokens := strings.Split(fromAttr, ".")
	if val, ok := dataMap[fromAttrTokens[0]+".id"]; ok {
		if val != nil {
			mmv.FromID = fmt.Sprint(helper.GetInt64(val))
			return nil
		}
	}
	err := fmt.Errorf("@mmv.SetFromID: Can't find it in: %s", fromAttr)
	fmt.Println(dataMap)
	log.Fatal(err)
	return err
}
