package scopes

import "k8s.io/cli-runtime/pkg/resource"

// TODO add a scope for files vs running instances

type Scope func(*resource.Info) bool

func Kind(kindName string) Scope {
	return func(info *resource.Info) bool {
		if info.Object.GetObjectKind().GroupVersionKind().Kind == kindName {
			return true
		}

		return false
	}
}

func Group(groupName string) Scope {
	return func(info *resource.Info) bool {
		if info.Object.GetObjectKind().GroupVersionKind().Group == groupName {
			return true
		}

		return false
	}
}

func Version(versionName string) Scope {
	return func(info *resource.Info) bool {
		if info.Object.GetObjectKind().GroupVersionKind().Version == versionName {
			return true
		}

		return false
	}
}

func And(scopes ...Scope) Scope {
	return func(info *resource.Info) bool {
		for _, v := range scopes {
			if v(info) != true {
				return false
			}
		}

		return true
	}
}

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
