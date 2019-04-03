package config

import "database/sql"

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
	Name     string              `json:"app"`
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
}

/****************** App Config Structs ***********************/

type App struct {
	Tables []map[string]string `json:""`
}

type AppConfig struct {
	AppName      string
	Tags         []Tag        `json:"tags"`
	Dependencies []Dependency `json:"dependencies"`
	Ownerships   []Ownership  `json:"ownership"`
	DBConn       *sql.DB
}

type Ownership struct {
	Tag        string              `json:"tag"`
	DependsOn  string              `json:"owned_by"`
	Conditions []map[string]string `json:"conditions"`
}

type Tag struct {
	Name              string              `json:"name"`
	Members           map[string]string   `json:"members"`
	Keys              map[string]string   `json:"keys"`
	InnerDependencies []map[string]string `json:"inner_dependencies"`
	Restrictions      []map[string]string `json:"restrictions"`
}

type Dependency struct {
	Tag       string      `json:"tag"`
	DependsOn []DependsOn `json:"depends_on"`
}

type DependsOn struct {
	Tag        string       `json:"tag"`
	Conditions []DCondition `json:"conditions"`
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
