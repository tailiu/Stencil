package qr

import "strings"

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func contains(a []string, x string) bool {
	x = strings.Trim(strings.ToLower(x), ", ")
	for _, n := range a {
		if strings.Contains(n, ".") {
			n = strings.Split(n, ".")[1]
		}
		n = strings.Trim(strings.ToLower(n), ", ")
		if x == n {
			return true
		}
	}
	return false
}
