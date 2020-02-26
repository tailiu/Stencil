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
