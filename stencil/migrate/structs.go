package migrate

import "stencil/config"

type DependencyNode struct {
	Tag  config.Tag
	SQL  string
	Data map[string]interface{}
}

type WaitingNode struct {
	ContainedNodes []DependencyNode
	// LookingFor    []config.Tag
	LookingFor map[string]map[string]interface{}
}

type WaitingList struct {
	Nodes []WaitingNode
}

type InvalidList struct {
	Nodes []DependencyNode
}
