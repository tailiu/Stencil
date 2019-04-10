package qr

import "database/sql"

type QR struct {
	StencilDB    *sql.DB
	AppName      string
	AppID        string
	BaseMappings []map[string]interface{}
	SuppMappings []map[string]interface{}
}

type QI struct {
	TableName        string
	Columns          []string
	Values           []string
	Conditions       string
	ColumnsWithTable map[string][]string
}
