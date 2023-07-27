package stringin

// Slice string in slice
func Slice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func slice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Slice2 string in slice 2
func Slice2(a string, list []string) string {
	if slice(a, list) {
		return a
	}

	return ""
}
