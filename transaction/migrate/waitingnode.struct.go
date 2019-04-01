package migrate

func (self *WaitingNode) Update(node DependencyNode) {
	self.ContainsNodes = append(self.ContainsNodes, node)
	delete(self.LookingFor, node.Tag.Name)
}

func (self WaitingNode) IsComplete() bool {
	if len(self.LookingFor) > 0 {
		return false
	}
	return true
}

func (self WaitingNode) GenDependencyDataNode() DependencyNode {

	var dependencyNode DependencyNode
	dependencyNode.Data = make(map[string]interface{})
	for _, containedNode := range self.ContainsNodes {
		for key, val := range containedNode.Data {
			dependencyNode.Data[key] = val
		}
	}
	return dependencyNode
}
