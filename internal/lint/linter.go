package lint

import (
	"github.com/SethCurry/kubectl-lint/internal/scopes"
	apps_v1 "k8s.io/api/apps/v1"
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

type PodSpecLinter interface {
	Lintv1PodSpec(core_v1.PodSpec) ([]ErrorMessage, error)
}

type DeploymentLinter interface {
	Lintv1Deployment(*apps_v1.Deployment) ([]ErrorMessage, error)
}

type DaemonSetLinter interface {
	Lintv1DaemonSet(*apps_v1.DaemonSet) ([]ErrorMessage, error)
}

// TODO implement this on runner
type ServiceLinter interface {
	Lintv1Service(*core_v1.Service) ([]ErrorMessage, error)
}
