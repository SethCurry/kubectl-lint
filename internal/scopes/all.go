package scopes

import "k8s.io/cli-runtime/pkg/resource"

func All() Scope {
	return func(_ *resource.Info) bool {
		return true
	}
}
