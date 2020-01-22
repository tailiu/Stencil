package counter

import (
	"database/sql"
	"stencil/config"
	"stencil/migrate"
)

type Counter struct {
	uid           string
	AppConfig     config.AppConfig
	root          *migrate.DependencyNode
	AppDBConn     *sql.DB
	StencilDBConn *sql.DB
	visitedNodes  map[string]map[string]bool
	NodeCount     int
	EdgeCount     int
}
