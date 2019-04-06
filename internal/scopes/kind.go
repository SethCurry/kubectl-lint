package scopes

import "k8s.io/cli-runtime/pkg/resource"

func Kind(kindName string) Scope {
	return func(info *resource.Info) bool {
		if info.Object.GetObjectKind().GroupVersionKind().Kind == kindName {
			return true
		}

		return false
	}
}
