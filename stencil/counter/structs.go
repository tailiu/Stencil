package counter

import (
	"database/sql"
	"stencil/config"
	"stencil/migrate"
)

type Counter struct {
	AppConfig     config.AppConfig
	StencilDBConn *sql.DB
	AppDBConn     *sql.DB
	root          *migrate.DependencyNode
	NodeCount     int
	EdgeCount     int
}
