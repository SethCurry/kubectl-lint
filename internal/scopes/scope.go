package scopes

import "k8s.io/cli-runtime/pkg/resource"

// TODO add a scope for files vs running instances

type Scope func(*resource.Info) bool
