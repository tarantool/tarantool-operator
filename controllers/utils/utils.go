package utils

import (
	"fmt"
)

func IsRolesEquals(rolesA, rolesB []string) bool {
	isSubset := func(X, Y []string) bool {
		for _, x := range X {
			match := false
			for _, y := range Y {
				if x == y {
					match = true
					break
				}
			}
			if !match {
				return false
			}
		}
		return true
	}
	return isSubset(rolesA, rolesB) && isSubset(rolesB, rolesA)
}

func MakeStaticPodAddr(pod string, svc string, ns string, domain string, port int) string {
	if port == 0 {
		port = 8081
	}

	if domain == "" {
		domain = "cluster.local"
	}

	return fmt.Sprintf("%s.%s.%s.svc.%s:%d", pod, svc, ns, domain, port)
}
