package SA1_display

import (
	"stencil/config"
	"stencil/reference_resolution_v2"
	"stencil/common_funcs"
	"database/sql"
)

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
	dag									*common_funcs.DAG
	tableNameIDPairs					map[string]string
	colNameIDPairs						map[string]string
	ownershipDisplaySettingsSatisfied 	bool
}

type display struct {
	stencilDBConn 						*sql.DB
	appIDNamePairs						map[string]string
	tableIDNamePairs					map[string]string
	attrIDNamePairs						map[string]string
	appTableNameTableIDPairs			map[string]string
	migrationID							int
	rr									*reference_resolution_v2.RefResolution
	resolveReference					bool
	srcAppConfig						*srcAppConfig
	dstAppConfig						*dstAppConfig
	mappingsFromSrcToDst				*config.MappedApp
	mappingsFromOtherAppsToDst			map[string]*config.MappedApp
	displayInFirstPhase					bool
	markAsDelete						bool
}
