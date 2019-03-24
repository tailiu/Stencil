package migrate

type DependencyNode struct {
	Tag  string
	SQL  string
	Data []map[string]interface{}
}
