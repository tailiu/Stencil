package migrate

import (
	"errors"
	"fmt"
	"time"
)

func (self DependencyNode) GetValueForKey(key string) (string, error) {

	for _, datum := range self.Data {
		if _, ok := datum[key]; ok {
			switch v := datum[key].(type) {
			case nil:
				return "", nil
			case int, int64:
				val := fmt.Sprintf("%d", v)
				return val, nil
			case string:
				val := fmt.Sprintf("%s", v)
				return val, nil
			case bool:
				val := fmt.Sprintf("%t", v)
				return val, nil
			case time.Time:
				val := v.String()
				return val, nil
			default:
				val := v.(string)
				return val, nil
			}
		}
	}
	return "", errors.New("No value found for " + key)
}
