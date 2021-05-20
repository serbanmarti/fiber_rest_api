package internal

// StringInSlice checks if a string is in a slice
func StringInSlice(p string, s []string) bool {
	for _, n := range s {
		if p == n {
			return true
		}
	}

	return false
}
