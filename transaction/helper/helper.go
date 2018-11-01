/*
 * Helping Functions, because Go is shit and Python is <3
 */

package helper

import (
	"fmt"
	"strings"
)

func Linebreak(ch string, times ...int) {

	t := 1

	if len(times) > 0 {
		t = times[0]
	}

	for i := 0; i < t; i++ {
		fmt.Printf(ch)
	}
	fmt.Println()
}

func Contains(list []string, str string) bool {
	for _, v := range list {
		// fmt.Println(v, "==", str)
		if strings.ToLower(v) == strings.ToLower(str) {
			return true
		}
	}
	return false
}