package qr

import (
	"errors"
	"fmt"
)

func (self QI) valueOfColumn(col string) (string, error) {

	for i, c := range self.Columns {
		if col == c {
			return self.Values[i], nil
		}
	}
	return "", errors.New("No column: " + col)
}

func (self QI) Print() {

	fmt.Println(self)

	// for i, c := range self.Columns {
	// 	fmt.Println("i:", i, "Col:", c, " || Val:", self.Values[i])
	// }
}
