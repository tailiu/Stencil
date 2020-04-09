package migrate_v2

import (
	"fmt"
	"log"
	"stencil/helper"
	"strings"
)

func (dataMap DataMap) GetRefValsFromDataMap(fromAttr, toAttr string, hardRef bool) (string, string, int64, error) {
	var toVal, fromVal string
	var fromID int64

	// log.Printf("@dataMap.GetIDsFromdataMap | Args | fromAttr : %s, toAttr : %s, hardRef: %v", fromAttr, toAttr, hardRef)
	// log.Printf("@dataMap.GetIDsFromdataMap | Args | dataMap : %v ", dataMap)

	fromAttrTokens := strings.Split(fromAttr, ".")

	if val, ok := dataMap[fromAttrTokens[0]+".id"]; ok {
		if val != nil {
			fromID = helper.GetInt64(val)
		}
	} else {
		return toVal, fromVal, fromID, fmt.Errorf("Unable to find fromID ref value in node data: %s.id", fromAttrTokens[0])
	}

	if val, ok := dataMap[fromAttr]; ok {
		if val != nil {
			toVal = fmt.Sprint(val)
			fromVal = fmt.Sprint(val)
		}
	} else {
		return toVal, fromVal, fromID, fmt.Errorf("Unable to find fromAttr ref value in node data: %s", fromAttr)
	}

	if hardRef {
		toVal = ""
		if val, ok := dataMap[toAttr]; ok {
			if val != nil {
				toVal = fmt.Sprint(val)
			}
		} else {
			return toVal, fromVal, fromID, fmt.Errorf("Unable to find toVal ref value in node data: %s", toAttr)
		}
	}

	// log.Printf("@dataMap.GetIDsFromdataMap | Returning | toVal : '%v', fromVal : '%v', fromID : '%v' ", toVal, fromVal, fromID)

	if len(toVal) != 0 && len(fromVal) != 0 && fromID != 0 {
		return toVal, fromVal, fromID, nil
	}

	err := fmt.Errorf("Nil reference(s) | fromAttr : '%s', toAttr : '%s', hardRef : '%v' | fromVal: '%v' | toVal: '%v'", fromAttr, toAttr, hardRef, fromVal, toVal)

	fmt.Println(dataMap)
	log.Println(err)

	return "", "", 0, err
}

func (dataMap DataMap) IsEmptyExcept() bool {
	exceptions := []string{".id", ".display_flag"}
	for key, val := range dataMap {
		if val == nil {
			continue
		}
		isException := false
		for _, exception := range exceptions {
			if strings.Contains(key, exception) {
				isException = true
				break
			}
		}
		if !isException {
			return false
		}
	}
	return true
}
