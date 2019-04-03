package scopes

import "k8s.io/cli-runtime/pkg/resource"

func Not(scope Scope) Scope {
	return func(info *resource.Info) bool {
		return !scope(info)
	}
}
