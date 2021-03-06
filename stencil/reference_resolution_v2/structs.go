package reference_resolution_v2

import (
	"database/sql"
	"stencil/config"
	"stencil/common_funcs"
)

// app, member, id are all integers corresponding to names
type Attribute struct {
	app 		string
	member 		string
	attrName 	string
	val 		string
	id			string
}

type RefResolution struct {
	stencilDBConn 					*sql.DB
	appDBConn						*sql.DB
	appID							string
	appName							string
	migrationID						int
	appTableNameIDPairs 			map[string]string
	appIDNamePairs					map[string]string
	tableIDNamePairs				map[string]string
	attrIDNamePairs					map[string]string
	appAttrNameIDPairs				map[string]string
	allMappings						*config.SchemaMappings
	dag 							*common_funcs.DAG
}
