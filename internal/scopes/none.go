package scopes

import "k8s.io/cli-runtime/pkg/resource"

// None applies to nothing.  It is used for testing,
// though I'm not sure it has a useful purpose beyond
// that.
func None() Scope {
	return func(_ *resource.Info) bool {
		return false
	}
}
