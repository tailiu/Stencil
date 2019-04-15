package qr

import "database/sql"

// QT: Query Type
const (
	QTInsert = "Insert"
	QTSelect = "Select"
	QTUpdate = "Update"
	QTDelete = "Delete"
)

// QR: Query Resolver
type QR struct {
	StencilDB    *sql.DB
	AppName      string
	AppID        string
	BaseMappings []map[string]interface{}
	SuppMappings []map[string]interface{}
}

// QI: Query Ingredients
type QI struct {
	TableName        string
	Columns          []string
	Values           []interface{}
	Conditions       string
	ColumnsWithTable map[string][]string
	Type             string
}
