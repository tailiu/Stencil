package common_funcs

import (
	"stencil/config"
)

type DAG struct {
	Tags         	[]config.Tag        `json:"tags"`
	Dependencies 	[]config.Dependency `json:"dependencies"`
	Ownerships   	[]config.Ownership  `json:"ownership"`
}

type DataInDependencyNode struct {
	Table 	string
	Data	map[string]interface{}
}