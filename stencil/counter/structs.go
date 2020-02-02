package counter

import (
	"database/sql"
	"stencil/config"
	"stencil/migrate"
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
