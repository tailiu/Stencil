package qr

import "database/sql"

// QT: Query Type
const (
	QTInsert = "Insert"
	QTSelect = "Select"
	QTUpdate = "Update"
	QTDelete = "Delete"
	EXISTS   = "EXISTS"
	NEXISTS  = "NOT EXISTS"
	IN       = "IN"
	NOTIN    = "NOT IN"
)

// QR: Query Resolver
type QR struct {
	StencilDB    *sql.DB
	AppName      string
	AppID        string
	BaseMappings []map[string]string
	SuppMappings []map[string]string
	Migration    bool
}

type QI struct {
	TableName        string
	Columns          []string
	Values           []interface{}
	Conditions       string
	ColumnsWithTable map[string][]string
	Type             string
}

type QSold struct {
	QR           *QR
	Columns      []string
	ColumnsWSize []string
	From         string
	Where        string
	Group        string
	Order        string
	Limit        string
	With         string
	TableAliases map[string]map[string]string
	seen         map[string]bool
	vals         []interface{}
	PK           bool
}

type QS struct {
	QR           *QR
	Columns      []string
	ColumnsWSize []string
	From         string
	Where        string
	Group        string
	Order        string
	Limit        string
	With         string
	TableAliases map[string]map[string]string
	seen         map[string]bool
	vals         []interface{}
	PK           bool
}

type QU struct {
	QR              *QR
	Tables          map[string]bool
	Update          map[string]string
	Where           map[string]string
	affected_tables []string
	affected_rows   []string
}

type BulkTrackerTable struct {
	TableName string
	Columns   []string
	Values    [][]interface{}
}

type BulkTracker struct {
	tables map[string]BulkTrackerTable
}
