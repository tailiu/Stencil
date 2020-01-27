package counter

import (
	"database/sql"
	"stencil/config"
	"stencil/migrate"
)

type Counter struct {
	uid           string
	AppConfig     config.AppConfig
	StencilDBConn *sql.DB
	AppDBConn     *sql.DB
	root          *migrate.DependencyNode
	NodeCount     int
	EdgeCount     int
	visitedNodes  map[string]map[string]bool
}
