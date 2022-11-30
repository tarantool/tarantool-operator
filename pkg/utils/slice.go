package utils

func SliceContains[Type comparable](a []Type, x Type) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}

	return false
}
