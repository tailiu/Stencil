package SA2_display

import (
	"database/sql"
	"stencil/qr"
	"stencil/common_funcs"
)

type srcAppConfig struct {
	appID 								string
	appName 							string
	tableNameIDPairs					map[string]string
}

type dstAppConfig struct {
	appID 								string
	appName 							string
	DBConn       						*sql.DB
	dag									*common_funcs.DAG
	tableNameIDPairs					map[string]string
	ownershipDisplaySettingsSatisfied 	bool
	qr									*qr.QR
}

type displayConfig struct {
	stencilDBConn 			*sql.DB
	appIDNamePairs			map[string]string
	tableIDNamePairs		map[string]string
	migrationID				int
	srcAppConfig			*srcAppConfig
	dstAppConfig			*dstAppConfig
	displayInFirstPhase		bool
	userID					string
}