package utils

func MergeMaps[KeyType comparable, ValueType any](ms ...map[KeyType]ValueType) map[KeyType]ValueType {
	res := map[KeyType]ValueType{}

	for _, m := range ms {
		if ms == nil {
			continue
		}

		for k, v := range m {
			res[k] = v
		}
	}

	return res
}
