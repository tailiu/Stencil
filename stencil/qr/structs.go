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
	BaseMappings []map[string]string
	SuppMappings []map[string]string
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

type QS struct {
	QR      *QR
	Columns []string
	From    string
	Where   string
	Group   string
	Order   string
	Limit   string
	seen    map[string]bool
}
