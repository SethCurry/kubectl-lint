package lint

import (
	"github.com/SethCurry/kubectl-lint/internal/scopes"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/resource"
)

type Linter interface {
	ErrorCode() ErrorCode
	ValidScopes() []scopes.Scope
	Lint(*resource.Info) ([]ErrorMessage, error)
	Severity() Severity
}

type PodLinter interface {
	Lintv1Pod(*core_v1.Pod) ([]ErrorMessage, error)
}
