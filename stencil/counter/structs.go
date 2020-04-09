package counter

import (
	"database/sql"
	config "stencil/config/v2"
	migrate "stencil/migrate_v2"
)

type Counter struct {
	UID           string
	AppConfig     config.AppConfig
	StencilDBConn *sql.DB
	AppDBConn     *sql.DB
	Root          *migrate.DependencyNode
	NodeCount     int
	EdgeCount     int
	VisitedNodes  map[string]map[string]bool
}
