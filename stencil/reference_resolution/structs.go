package reference_resolution

import (
	"database/sql"
	"stencil/config"
)

// app, member, id are all integers corresponding to names
type Identity struct {
	app 	string
	member 	string
	id 		string
}

type RefResolutionConfig struct {
	StencilDBConn 			*sql.DB
	AppDBConn				*sql.DB
	AppID					string
	AppName					string
	MigrationID				int
	AppTableNameIDPairs 	map[string]string
	AppIDNamePairs			map[string]string
	TableIDNamePairs		map[string]string
	AllMappings				*config.SchemaMappings
	MappingsFromSrcToDst	*config.MappedApp
}
