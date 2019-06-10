package migrate

import (
	"errors"
	"fmt"
	"stencil/config"
	"strings"
)

func (self WaitingList) UpdateIfBeingLookedFor(node DependencyNode) (WaitingNode, error) { // node exists, node is being looked for

	for i, waitingNode := range self.Nodes {
		if lookingFor, ok := waitingNode.LookingFor[node.Tag.Name]; ok {
			satisfied := true
			for lookingForKey, lookingForVal := range lookingFor {
				// for _, datum := range node.Data {
				if val, ok := node.Data[lookingForKey]; ok {
					if val != lookingForVal {
						satisfied = false
						break
					}
				} else {
					satisfied = false
					break
				}
				// }
				// if !satisfied {
				// 	break
				// }
			}
			if satisfied {
				self.Nodes[i].Update(node)
				return self.Nodes[i], nil
			}
		}
	}

	return *new(WaitingNode), errors.New("Not Found In Waiting List")
}

func (self WaitingList) IsAlreadyWaiting(node DependencyNode) bool { // node exists, node is being looked for

	for _, waitingNode := range self.Nodes {
		for _, containedNode := range waitingNode.ContainedNodes {
			if strings.EqualFold(containedNode.Tag.Name, node.Tag.Name) {
				if idAttr1, err := containedNode.Tag.ResolveTagAttr("id"); err == nil {
					if idAttr2, err := node.Tag.ResolveTagAttr("id"); err == nil {
						idAttr1val, err1 := containedNode.GetValueForKey(idAttr1)
						idAttr2val, err2 := node.GetValueForKey(idAttr2)
						if err1 == nil && err2 == nil && strings.EqualFold(idAttr1val, idAttr2val) {
							return true
						}
					}
				}

			}
		}
	}

	return false
}

func (waitingList *WaitingList) AddNewToWaitingList(node DependencyNode, adjTags []config.Tag, srcApp config.AppConfig) error {

	waitingNode := new(WaitingNode)
	waitingNode.ContainedNodes = append(waitingNode.ContainedNodes, node)
	waitingNode.LookingFor = make(map[string]map[string]interface{})

	for _, adjTag := range adjTags {
		if strings.EqualFold(adjTag.Name, node.Tag.Name) {
			continue
		}
		if dependsOn, err := srcApp.CheckDependency(node.Tag.Name, adjTag.Name); err == nil {
			for _, condition := range dependsOn.Conditions {
				for _, restriction := range condition.Restrictions {
					if attr, err := adjTag.ResolveTagAttr(restriction["col"]); err == nil {
						// for _, datum := range node.Data {
						if _, ok := node.Data[attr]; ok {
							if !strings.EqualFold(node.Data[attr].(string), restriction["value"]) {
								return errors.New(fmt.Sprintf("AddToWaitingList: Restriction Failed: %s != %s", node.Data[attr].(string), restriction["value"]))
							}
						} else {
							return errors.New(fmt.Sprintf("AddToWaitingList: Restriction Attr [%s] not found in Node [%s]", restriction["col"], node.Tag))
						}
						// }
					} else {
						return errors.New(fmt.Sprintf("AddToWaitingList: Restriction Attr [%s] not found in tag key list [%s]", restriction["col"], adjTag))
					}
				}
				waitingNode.LookingFor[adjTag.Name] = make(map[string]interface{})

				if tagAttr, err := adjTag.ResolveTagAttr(condition.TagAttr); err == nil {
					if dependsOnAttr, err := node.Tag.ResolveTagAttr(condition.DependsOnAttr); err == nil {
						// for _, datum := range node.Data {
						if _, ok := node.Data[dependsOnAttr]; ok {
							waitingNode.LookingFor[adjTag.Name][tagAttr] = node.Data[dependsOnAttr]
						} else {
							return errors.New(fmt.Sprintf("AddToWaitingList: Criteria failed. Can't find attribute [%s] in data node [%s]", tagAttr, node.Data))
						}
						// }
					} else {
						return errors.New(fmt.Sprintf("AddToWaitingList: Attr can't be resolved [%s], [%s]", condition.DependsOnAttr, adjTag))
					}
				} else {
					return errors.New(fmt.Sprintf("AddToWaitingList: Attr can't be resolved [%s], [%s]", condition.TagAttr, adjTag))
				}
			}
		} else {
			fmt.Println(fmt.Sprintf("No dependency exists between [%s] and [%s].", node.Tag.Name, adjTag.Name))
		}
	}

	if len(waitingNode.LookingFor) > 0 {
		waitingList.Nodes = append(waitingList.Nodes, *waitingNode)
	}

	return nil
}
