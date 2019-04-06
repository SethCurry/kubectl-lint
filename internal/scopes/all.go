package scopes

import "k8s.io/cli-runtime/pkg/resource"

// All means that the linter applies to all objects.
// This is typically used for things that check metadata,
// since nearly all objects have them.
func All() Scope {
	return func(_ *resource.Info) bool {
		return true
	}
}
