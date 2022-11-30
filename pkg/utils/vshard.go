package utils

import "github.com/google/go-cmp/cmp"

func IsVShardRolesEquals(rolesA, rolesB []string, hierarchy map[string][]string) bool {
	return cmp.Equal(VShardRolesToMap(rolesA, hierarchy), VShardRolesToMap(rolesB, hierarchy))
}

func VShardRolesToMap(roles []string, hierarchy map[string][]string) map[string]bool {
	rmap := make(map[string]bool)

	for _, role := range roles {
		rmap[role] = true

		for _, dep := range hierarchy[role] {
			rmap[dep] = true
		}
	}

	return rmap
}
