package migrate_v2

import "fmt"

func (pBags *ProcessedBags) Init() {
	pBags.Bags = make(map[string]bool)
}

func (pBags *ProcessedBags) Add(id string) {
	pBags.Bags[id] = true
}

func (pBags *ProcessedBags) Exists(id string) bool {
	_, ok := pBags.Bags[id]
	return ok
}

func (pBags *ProcessedBags) Update(bagNode *DependencyNode) {
	for _, pk := range bagNode.PKs {
		pBags.Add(fmt.Sprint(pk))
	}
}
