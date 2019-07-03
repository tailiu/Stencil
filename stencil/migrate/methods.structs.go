package migrate

import (
	"fmt"
	"stencil/helper"
)

func (self *UnmappedTags) Exists(tag string) bool {
	return helper.Contains(self.tags, tag)
}

func (self *UnmappedTags) Add(tag string) {
	if self.Exists(tag) {
		return
	}
	self.Mutex.Lock()
	self.tags = append(self.tags, tag)
	self.Mutex.Unlock()
	fmt.Println("==>> ADDED NEW UNMAPPED TAG:", tag, self.tags)
}
