package common_funcs

func ExistsInSlice(s []string, e string) bool {

	for _, s1 := range s {
		if e == s1 {
			return true
		}
	}

	return false
}