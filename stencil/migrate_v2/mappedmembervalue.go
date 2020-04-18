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

func (mmv MappedMemberValue) GetAttrID(attrName, tableID string) string {
	if len(tableID) == 0 || len(attrName) == 0 {
		log.Fatalf("@mmv.GetAttrID: Table ID or App ID or AttrName nil | Table: '%s' | Attr : '%s'", tableID, attrName)
	} else if mmv.DBConn == nil {
		log.Fatal("@mmv.GetAttrID: DBConn not set!")
	} else {
		if attrID, err := db.AttrID(mmv.DBConn, tableID, attrName); err == nil {
			return attrID
		} else {
			fmt.Printf("@mmv.GetAttrID: Table: '%s' | Attr : '%s'\n", tableID, attrName)
			log.Fatal("@mmv.GetAttrID: ", err)
		}
	}
	err := fmt.Errorf("This is the end. Shouldn't be here. Table: '%s' | Attr : '%s'", tableID, attrName)
	log.Fatal("@mmv.GetAttrID: ", err)
	return err.Error()
}

func (mmv *MappedMemberValue) StoreMemberAndAttr(mappedStmt string) {
	tokens := strings.Split(mappedStmt, ".")
	mmv.FromMember = tokens[0]
	mmv.FromAttr = tokens[1]
	mmv.FromMemberID = mmv.GetTableID(mmv.FromMember, mmv.AppID)
	mmv.FromAttrID = mmv.GetAttrID(mmv.FromAttr, mmv.FromMemberID)
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

		fromMemberID := mmv.GetTableID(fromAttrTokens[0], mmv.AppID)
		fromAttrID := mmv.GetAttrID(fromAttrTokens[1], fromMemberID)

		toMemberID := mmv.GetTableID(toAttrTokens[0], mmv.AppID)
		toAttrID := mmv.GetAttrID(toAttrTokens[1], toMemberID)

		mmv.Ref = &MappingRef{
			appID:        mmv.AppID,
			fromID:       fromID,
			fromVal:      fromVal,
			fromMemberID: fromMemberID,
			fromMember:   fromAttrTokens[0],
			fromAttr:     fromAttrTokens[1],
			fromAttrID:   fromAttrID,
			toVal:        toVal,
			toMemberID:   toMemberID,
			toMember:     toAttrTokens[0],
			toAttr:       toAttrTokens[1],
			toAttrID:     toAttrID,
		}
	} else {
		return err
	}
	return nil
}

func (mmv *MappedMemberValue) SetFromID(dataMap DataMap) error {

	if len(mmv.FromMember) == 0 {
		log.Fatal("@mmv.SetFromID: FromMember is not set!")
	}

	if val, ok := dataMap[mmv.FromMember+".id"]; ok {
		if val != nil {
			mmv.FromID = fmt.Sprint(helper.GetInt64(val))
			return nil
		}
	}
	err := fmt.Errorf("@mmv.SetFromID: Can't find from ID in: %s", mmv.FromMember)
	fmt.Println(dataMap)
	// log.Fatal(err)
	return err
}
