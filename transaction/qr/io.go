package qr

import (
	"errors"
	"fmt"
	"transaction/db"
)

func NewQRWithAppName(app_name string) *QR {
	qr := new(QR)
	qr.AppName = app_name
	qr.StencilDB = db.GetDBConn("stencil")
	qr.setAppID()
	qr.getBaseMappings()
	qr.getSupplementaryMappings()
	return qr
}

func NewQRWithAppID(app_id string) *QR {
	qr := new(QR)
	qr.AppID = app_id
	qr.StencilDB = db.GetDBConn("stencil")
	qr.setAppName()
	qr.getBaseMappings()
	qr.getSupplementaryMappings()
	return qr
}

func (self *QR) setAppID() error {
	sql := fmt.Sprintf("SELECT rowid from apps WHERE app_name = '%s'", self.AppName)
	result := db.DataCall1(self.StencilDB, sql)
	if val, ok := result["rowid"]; ok {
		self.AppID = fmt.Sprint(val)
		return nil
	}
	return errors.New("can't set app name")
}

func (self *QR) setAppName() error {
	sql := fmt.Sprintf("SELECT app_name from apps WHERE app_id = '%s'", self.AppID)
	result := db.DataCall1(self.StencilDB, sql)
	if val, ok := result["app_name"]; ok {
		self.AppName = fmt.Sprint(val)
		return nil
	}
	return errors.New("can't set app name")
}

func (self *QR) getBaseMappings() {
	sql := fmt.Sprintf(`SELECT
							LOWER(app_schemas.table_name) as logical_table, 
							LOWER(app_schemas.column_name) as logical_column, 
							LOWER(physical_schema.table_name) as physical_table,  
							LOWER(physical_schema.column_name) as physical_column
						FROM 	
							physical_mappings 
							JOIN 	app_schemas ON physical_mappings.logical_attribute = app_schemas.rowid
							JOIN 	app_tables ON app_schemas.table_id = app_tables.rowid
							JOIN 	physical_schema ON physical_mappings.physical_attribute = physical_schema.rowid
						WHERE 	app_tables.app_id  = '%s' `, self.AppID)

	self.BaseMappings = db.DataCall(self.StencilDB, sql)
}

func (self *QR) getSupplementaryMappings() {
	sql := fmt.Sprintf(`SELECT  LOWER(asm.table_name) as logical_table,
							LOWER(asm.column_name)  as logical_column,
							CONCAT('supplementary_',st.rowid::string) as physical_table,
							LOWER(asm.column_name)  as physical_column
						FROM 	app_schemas asm JOIN
						supplementary_tables st ON 
						st.table_name = asm.table_name AND 
						st.app_id = asm.app_id
						WHERE 	asm.app_id  = '%s' AND
						asm.rowid NOT IN (
							SELECT logical_attribute FROM physical_mappings
						)`, self.AppID)

	self.SuppMappings = db.DataCall(self.StencilDB, sql)
}
