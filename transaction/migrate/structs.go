package migrate

import "transaction/config"

type DependencyNode struct {
	Tag  config.Tag
	SQL  string
	Data map[string]interface{}
}

type WaitingNode struct {
	ContainsNodes []DependencyNode
	// LookingFor    []config.Tag
	LookingFor map[string]map[string]interface{}
}

type WaitingList struct {
	Nodes []WaitingNode
}
