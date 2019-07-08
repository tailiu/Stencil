package migrate

import (
	"stencil/helper"
)

func (self *UnmappedTags) Exists(tag string) bool {
	return helper.Contains(self.tags, tag)
}

func (self *UnmappedTags) Add(tag string) {
	self.Mutex.Lock()
	if !self.Exists(tag) {
		self.tags = append(self.tags, tag)
	}
	self.Mutex.Unlock()
}
