package migrate_v1

func (self *WaitingNode) Update(node *DependencyNode) {
	self.ContainedNodes = append(self.ContainedNodes, node)
	delete(self.LookingFor, node.Tag.Name)
}

func (self WaitingNode) IsComplete() bool {
	if len(self.LookingFor) > 0 {
		return false
	}
	return true
}

func (self WaitingNode) GenDependencyDataNode() *DependencyNode {

	dependencyNode := new(DependencyNode)
	dependencyNode.Data = make(map[string]interface{})
	for _, containedNode := range self.ContainedNodes {
		for key, val := range containedNode.Data {
			dependencyNode.Data[key] = val
		}
	}
	return dependencyNode
}
