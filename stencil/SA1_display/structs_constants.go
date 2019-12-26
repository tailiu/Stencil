package SA1_display

import (
	"stencil/config"
	"database/sql"
)

type DAG struct {
	tags         	[]config.Tag        `json:"tags"`
	dependencies 	[]config.Dependency `json:"dependencies"`
	ownerships   	[]config.Ownership  `json:"ownership"`
}

type srcAppConfig struct {
	appID 			string
	appName 		string
	userID			string
	dag				*DAG
}

type dstAppConfig struct {
	appID 			string
	appName 		string
	userID			string
	DBConn       	*sql.DB
	dag				*DAG
}

type displayConfig struct {
	stencilDBConn 		*sql.DB
	appIDNamePairs		map[string]string
	tableIDNamePairs	map[string]string
	attrIDNamePairs		map[string]string
	migrationID			int
	allMappings			*config.SchemaMappings
	mappingsToDst 		*config.MappedApp
	resolveReference	bool
	srcAppConfig		*srcAppConfig
	dstAppConfig		*dstAppConfig
}

type DataInDependencyNode struct {
	Table 	string
	Data	map[string]interface{}
}
