package SA1_display

import (
	"stencil/config"
	"stencil/reference_resolution"
	"database/sql"
	"sync"
)

type DAG struct {
	Tags         	[]config.Tag        `json:"tags"`
	Dependencies 	[]config.Dependency `json:"dependencies"`
	Ownerships   	[]config.Ownership  `json:"ownership"`
}

type srcAppConfig struct {
	appID 								string
	appName 							string
	userID								string
	tableNameIDPairs					map[string]string
}

type dstAppConfig struct {
	appID 								string
	appName 							string
	rootTable							string
	rootAttr							string
	userID								string
	DBConn       						*sql.DB
	dag									*DAG
	tableNameIDPairs					map[string]string
	ownershipDisplaySettingsSatisfied 	bool
}

type displayConfig struct {
	stencilDBConn 			*sql.DB
	appIDNamePairs			map[string]string
	tableIDNamePairs		map[string]string
	attrIDNamePairs			map[string]string
	migrationID				int
	refResolutionConfig		*reference_resolution.RefResolutionConfig
	resolveReference		bool
	srcAppConfig			*srcAppConfig
	dstAppConfig			*dstAppConfig
	mappingsFromSrcToDst	*config.MappedApp
	wg 						*sync.WaitGroup
}

type DataInDependencyNode struct {
	Table 	string
	Data	map[string]interface{}
}
