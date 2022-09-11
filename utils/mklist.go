package utils

func MkStringSet(x ...string) map[string]bool {
	n := make(map[string]bool, len(x))
	for _, i := range x {
		n[i] = true
	}

	return n
}
