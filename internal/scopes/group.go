package scopes

import "k8s.io/cli-runtime/pkg/resource"

func Group(groupName string) Scope {
	return func(info *resource.Info) bool {
		if info.Object.GetObjectKind().GroupVersionKind().Group == groupName {
			return true
		}

		return false
	}
}
