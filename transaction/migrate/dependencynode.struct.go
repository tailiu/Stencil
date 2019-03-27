package migrate

import (
	"fmt"
	"time"
)

func (self DependencyNode) GetValueForKey(key string) *string {

	for _, datum := range self.Data {
		if _, ok := datum[key]; ok {
			switch v := datum[key].(type) {
			case nil:
				return nil
			case int, int64:
				val := fmt.Sprintf("%d", v)
				return &val
			case string:
				val := fmt.Sprintf("%s", v)
				return &val
			case bool:
				val := fmt.Sprintf("%t", v)
				return &val
			case time.Time:
				val := v.String()
				return &val
			default:
				val := v.(string)
				return &val
			}
		}
	}
	return nil
}
