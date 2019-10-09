package config

import (
	"database/sql"
	"math/rand"
	"stencil/qr"
)

type DataQuery struct {
	SQL, Table string
}

/****************** Shema Mappings Structs ***********************/

type SchemaMappings struct {
	AllMappings []SchemaMapping `json:"allMappings"`
}

type SchemaMapping struct {
	FromApp string      `json:"fromApp"`
	ToApps  []MappedApp `json:"toApps"`
}

type MappedApp struct {
	Name     string              `json:"name"`
	Methods  []map[string]string `json:"methods"`
	Inputs   []map[string]string `json:"inputs"`
	Mappings []Mapping           `json:"mappings"`
}

type Mapping struct {
	FromTables []string  `json:"fromTables"`
	ToTables   []ToTable `json:"toTables"`
}

type ToTable struct {
	Table      string            `json:"table"`
	Conditions map[string]string `json:"conditions"`
	Mapping    map[string]string `json:"mapping"`
	Media	   map[string]string `json:"media"`
	TableID    string
}

/****************** App Config Structs ***********************/

type App struct {
	Tables []map[string]string `json:""`
}

type AppConfig struct {
	AppName      string
	AppID        string
	TableIDNamePairs map[string]string
	TableNameIDPairs map[string]string
	Tags         []Tag        `json:"tags"`
	Dependencies []Dependency `json:"dependencies"`
	Ownerships   []Ownership  `json:"ownership"`
	DBConn       *sql.DB
	QR           *qr.QR
	Rand         *rand.Rand
}

type Ownership struct {
	Tag        string       `json:"tag"`
	OwnedBy    string       `json:"owned_by"`
	Conditions []DCondition `json:"conditions"`
}

type Tag struct {
	Name              string              `json:"name"`
	Members           map[string]string   `json:"members"`
	Keys              map[string]string   `json:"keys"`
	InnerDependencies []map[string]string `json:"inner_dependencies"`
	Restrictions      []map[string]string `json:"restrictions"`
	Display_setting   string              `json:"display_setting"`
}

type Dependency struct {
	Tag                    string      `json:"tag"`
	DependsOn              []DependsOn `json:"depends_on"`
	CombinedDisplaySetting string      `json:"combined_display_setting"`
}

type DependsOn struct {
	Tag              string       `json:"tag"`
	As               string       `json:"as"`
	DisplayExistence string       `json:"display_existence"`
	DisplaySetting   string       `json:"display_setting"`
	Conditions       []DCondition `json:"conditions"`
}

type DCondition struct {
	TagAttr       string              `json:"tag_attr"`
	DependsOnAttr string              `json:"depends_on_attr"`
	Restrictions  []map[string]string `json:"restrictions"`
}

type Settings struct {
	UserTable string
	KeyCol    string
	Mappings  Mapping
}
