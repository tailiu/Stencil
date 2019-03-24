package config

type DataQuery struct {
	SQL, Table string
}

/****************** Shema Mappings Structs ***********************/

type SchemaMappings struct {
	AllMappings []SchemaMapping `json:"allMappings"`
}

type SchemaMapping struct {
	FromApp   string          `json:"fromApp"`
	VarsFuncs VarsFuncsConfig `json:"varsFuncs"`
	ToApps    []MappedToApp   `json:"toApps"`
}

type VarsFuncsConfig struct {
	Funcs []Func `json:"funcs"`
	Vars  []Vars `json:"vars"`
}

type Func struct {
	Name          string `json:"name"`
	MappingToFunc string `json:"mapping"`
}

type Vars struct {
	Name         string `json:"name"`
	MappingToVar string `json:"mapping"`
}

type MappedToApp struct {
	Name     string    `json:"app"`
	Mappings []Mapping `json:"mappings"`
}

type Mapping struct {
	FromTables string    `json:"fromTables"`
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
