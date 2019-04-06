package scopes

import "k8s.io/cli-runtime/pkg/resource"

func Or(scopes ...Scope) Scope {
	return func(info *resource.Info) bool {
		for _, v := range scopes {
			if v(info) == true {
				return true
			}
		}

		return false
	}
}
