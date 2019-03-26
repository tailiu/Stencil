package migrate

type DependencyNode struct {
	Tag  string
	SQL  string
	Data []map[string]interface{}
}

type WaitingNode struct {
	members []map[string]bool
	sqls    []map[string]string
}
