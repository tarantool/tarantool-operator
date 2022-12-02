package utils

import "github.com/google/go-cmp/cmp"

func IsMapSubset(set, subset map[string]any) bool {
	if len(subset) > len(set) {
		return false
	}

	for k, subsetValue := range subset {
		setValue, found := set[k]
		if !found {
			return false
		}

		if !cmp.Equal(setValue, subsetValue) {
			return false
		}
	}

	return true
}
