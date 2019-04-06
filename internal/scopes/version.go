package scopes

import "k8s.io/cli-runtime/pkg/resource"

func Version(versionName string) Scope {
	return func(info *resource.Info) bool {
		if info.Object.GetObjectKind().GroupVersionKind().Version == versionName {
			return true
		}

		return false
	}
}
