package migrate

import (
	"errors"
	"fmt"
	"strings"
	"transaction/config"
)

func (self WaitingList) SearchWaitingList(node DependencyNode) (WaitingNode, error) {

	for _, waitingNode := range self.Nodes {
		for _, containedNode := range waitingNode.ContainsNodes {
			if strings.EqualFold(containedNode.Tag, node.Tag) {
				for criterionName, criterionValue := range waitingNode.Criteria {
					satisfied := true
					for _, nodeDatum := range node.Data {
						if _, ok := nodeDatum[criterionName]; ok {
							if nodeDatum[criterionName] != criterionValue {
								satisfied = false
							}
						} else {
							satisfied = false
						}
					}
					if satisfied {
						return waitingNode, nil
					}
				}
			}
		}
	}

	return *new(WaitingNode), errors.New("Can't find node in waiting list.")
}

func (self WaitingList) NodeAlreadyExists(node DependencyNode) bool {

	return false
}

func (waitingList WaitingList) AddToWaitingList(node DependencyNode, parentTag config.Tag, dependencyCondition config.DCondition) error {

	for _, restriction := range dependencyCondition.Restrictions {
		attr := fmt.Sprintf("%s.%s", node.Tag, restriction["col"])
		for _, datum := range node.Data {
			if _, ok := datum[attr]; ok {
				if !strings.EqualFold(datum[attr].(string), restriction["value"]) {
					return errors.New(fmt.Sprintf("AddToWaitingList: Restriction Failed: %s != %s", datum[attr].(string), restriction["value"]))
				}
			} else {
				return errors.New(fmt.Sprintf("AddToWaitingList: Restriction Attr [%s] not found in Node [%s]", restriction["col"], node.Tag))
			}
		}
	}

	waitingNode := new(WaitingNode)
	waitingNode.ContainsNodes = append(waitingNode.ContainsNodes, node)
	waitingNode.LookingFor = append(waitingNode.LookingFor, parentTag.Name)

	tagAttr := fmt.Sprintf("%s.%s", node.Tag, dependencyCondition.TagAttr)
	dependsOnAttr := fmt.Sprintf("%s.%s", parentTag.Name, dependencyCondition.DependsOnAttr)

	for _, datum := range node.Data {
		if _, ok := datum[tagAttr]; ok {
			waitingNode.Criteria[dependsOnAttr] = datum[tagAttr]
		} else {
			return errors.New(fmt.Sprintf("AddToWaitingList: Criteria failed. Can't find attribute [%s] in data node [%s]", dependencyCondition.TagAttr, node.Tag))
		}
	}

	waitingList.Nodes = append(waitingList.Nodes, *waitingNode)

	return nil
}
