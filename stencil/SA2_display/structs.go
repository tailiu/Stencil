package SA2_display

import (
	"stencil/config"
)

type DataInDependencyNode struct {
	Table 	string
	Data	map[string]interface{}
}

type DAG struct {
	Tags         	[]config.Tag        `json:"tags"`
	Dependencies 	[]config.Dependency `json:"dependencies"`
	Ownerships   	[]config.Ownership  `json:"ownership"`
}

type srcAppConfig struct {
	appID 								string
	appName 							string
	tableNameIDPairs					map[string]string
}

type dstAppConfig struct {
	appID 								string
	appName 							string
	DBConn       						*sql.DB
	dag									*DAG
	tableNameIDPairs					map[string]string
	ownershipDisplaySettingsSatisfied 	bool
}

type displayConfig struct {
	stencilDBConn 			*sql.DB
	appIDNamePairs			map[string]string
	tableIDNamePairs		map[string]string
	migrationID				int
	srcAppConfig			*srcAppConfig
	dstAppConfig			*dstAppConfig
	mappingsFromSrcToDst	*config.MappedApp
	displayInFirstPhase		bool
	userID					string
}