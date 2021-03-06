package qr

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"stencil/db"
	"strings"
	"time"
)

func NewQRWithDBConn(app_name, app_id string, dbConn *sql.DB) *QR {
	qr := new(QR)
	qr.AppName = app_name
	qr.AppID = app_id
	qr.StencilDB = dbConn
	// fmt.Println("Fetching Base Mappings...")
	qr.getBaseMappings()
	// fmt.Println("Fetching Supplementary Mappings...")
	qr.getSupplementaryMappings()
	// fmt.Println("QR Created")
	rand.Seed(time.Now().UnixNano())
	return qr
}

func NewQR(app_name, app_id string) *QR {
	qr := new(QR)
	qr.AppName = app_name
	qr.AppID = app_id
	qr.StencilDB = db.GetDBConn(db.STENCIL_DB)
	// fmt.Println("Fetching Base Mappings...")
	qr.getBaseMappings()
	// fmt.Println("Fetching Supplementary Mappings...")
	qr.getSupplementaryMappings()
	// fmt.Println("QR Created")
	rand.Seed(time.Now().UnixNano())
	return qr
}

func NewQRWithAppName(app_name string) *QR {
	qr := new(QR)
	qr.AppName = app_name
	qr.StencilDB = db.GetDBConn(db.STENCIL_DB)
	fmt.Println("Getting App ID...")
	qr.setAppID()
	fmt.Println("Fetching Base Mappings...")
	qr.getBaseMappings()
	fmt.Println("Fetching Supplementary Mappings...")
	qr.getSupplementaryMappings()
	fmt.Println("QR Created")
	rand.Seed(time.Now().UnixNano())
	return qr
}

func NewQRWithAppID(app_id string) *QR {
	qr := new(QR)
	qr.AppID = app_id
	qr.StencilDB = db.GetDBConn(db.STENCIL_DB)
	fmt.Println("Getting App Name...")
	qr.setAppName()
	fmt.Println("Fetching Base Mappings...")
	qr.getBaseMappings()
	fmt.Println("Fetching Supplementary Mappings...")
	qr.getSupplementaryMappings()
	fmt.Println("QR Created")
	rand.Seed(time.Now().UnixNano())
	return qr
}

func (self *QR) setAppID() error {
	sql := fmt.Sprintf("SELECT pk from apps WHERE app_name = '%s'", self.AppName)
	if result, err := db.DataCall1(self.StencilDB, sql); err == nil {
		if val, ok := result["pk"]; ok {
			fmt.Println("App ID:", val)
			self.AppID = fmt.Sprint(val)
			return nil
		} else {
			return errors.New("Can't find pk in Query Results!")
		}
	} else {
		return err
	}
}

func (self *QR) setAppName() error {
	sql := fmt.Sprintf("SELECT app_name from apps WHERE app_id = '%s'", self.AppID)
	if result, err := db.DataCall1(self.StencilDB, sql); err == nil {
		if val, ok := result["app_name"]; ok {
			self.AppName = fmt.Sprint(val)
			return nil
		} else {
			return errors.New("Can't find app_name in Query Results!")
		}
	} else {
		return err
	}
}

func (self *QR) getBaseMappings() error {
	sql := fmt.Sprintf(`SELECT
							LOWER(app_tables.table_name) as logical_table,
							LOWER(app_schemas.column_name) as logical_column,
							LOWER(physical_schema.table_name) as physical_table,
							LOWER(physical_schema.column_name) as physical_column
						FROM
							physical_mappings
							JOIN 	app_schemas ON physical_mappings.logical_attribute = app_schemas.pk
							JOIN 	app_tables ON app_schemas.table_id = app_tables.pk
							JOIN 	physical_schema ON physical_mappings.physical_attribute = physical_schema.pk
						WHERE 	app_tables.app_id  = '%s' `, self.AppID)
	if strings.EqualFold(self.AppID, "1") {
		sql = `SELECT * FROM diaspora_base_mappings`
	}
	// self.BaseMappings = db.DataCall(self.StencilDB, sql)
	if result, err := db.DataCall(self.StencilDB, sql); err == nil {
		for _, mapping := range result {
			mappingStr := make(map[string]string)
			for key, val := range mapping {
				mappingStr[key] = fmt.Sprint(val)
			}
			self.BaseMappings = append(self.BaseMappings, mappingStr)
		}
		return nil
	} else {
		log.Fatal("UNABLE TO FETCH BASE_MAPPINGS:", err)
		return err
	}
}

func (self *QR) getSupplementaryMappings() error {
	sql := fmt.Sprintf(`SELECT
							LOWER(app_tables.table_name) as logical_table,
							LOWER(asm.column_name)  as logical_column,
							CONCAT('supplementary_',st.pk) as physical_table,
							LOWER(asm.column_name)  as physical_column
						FROM 	app_schemas asm
								JOIN app_tables on app_tables.pk = asm.table_id
								JOIN supplementary_tables st ON st.table_id = asm.table_id
						WHERE 	app_tables.app_id  = '%s' AND
								asm.pk NOT IN (
									SELECT logical_attribute FROM physical_mappings
								)`, self.AppID)

	if strings.EqualFold(self.AppID, "1") {
		sql = `SELECT * FROM diaspora_supplementary_mappings`
	}

	// self.SuppMappings = db.DataCall(self.StencilDB, sql)
	if result, err := db.DataCall(self.StencilDB, sql); err == nil {
		for _, mapping := range result {
			mappingStr := make(map[string]string)
			for key, val := range mapping {
				mappingStr[key] = fmt.Sprint(val)
			}
			self.SuppMappings = append(self.SuppMappings, mappingStr)
		}
		return nil
	} else {
		log.Fatal("UNABLE TO FETCH SUPP_MAPPINGS:", err)
		return err
	}
}
