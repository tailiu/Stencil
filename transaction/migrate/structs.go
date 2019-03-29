package migrate

type DependencyNode struct {
	Tag  string
	SQL  string
	Data []map[string]interface{}
}

type WaitingNode struct {
	ContainsNodes []DependencyNode
	LookingFor    []string
	Criteria      map[string]interface{}
}

type WaitingList struct {
	Nodes []WaitingNode
}
