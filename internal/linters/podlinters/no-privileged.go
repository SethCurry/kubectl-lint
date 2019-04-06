package podlinters

import (
	"errors"
	"fmt"

	"github.com/SethCurry/kubectl-lint/internal/lint"
	"github.com/SethCurry/kubectl-lint/internal/scopes"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/resource"
)

func NoPrivileged() *NoPrivilegedLinter {
	return &NoPrivilegedLinter{}
}

type NoPrivilegedLinter struct {
}

func (n *NoPrivilegedLinter) ErrorCode() lint.ErrorCode {
	return lint.ErrorCode("no-privileged")
}

func (n *NoPrivilegedLinter) Severity() lint.Severity {
	return lint.SeverityWarn
}

func (n *NoPrivilegedLinter) ValidScopes() []scopes.Scope {
	return []scopes.Scope{
		scopes.Or(
			scopes.And(scopes.Kind("Pod"), scopes.Version("v1")),
			scopes.And(scopes.Kind("Deployment"), scopes.Version("v1")),
			scopes.And(scopes.Kind("DaemonSet"), scopes.Version("v1")),
		),
	}
}

func (n *NoPrivilegedLinter) Lint(info *resource.Info) ([]lint.ErrorMessage, error) {
	return []lint.ErrorMessage{}, errors.New("this should be handled by the pod linter")
}

func (n *NoPrivilegedLinter) Lintv1PodSpec(spec core_v1.PodSpec) ([]lint.ErrorMessage, error) {
	var ret []lint.ErrorMessage

	containers := append(spec.Containers, spec.InitContainers...)

	for _, v := range containers {
		if v.SecurityContext != nil {
			if v.SecurityContext.Privileged != nil && *v.SecurityContext.Privileged == true {
				ret = append(ret, lint.ErrorMessage(fmt.Sprintf("container \"%s\" is running in privileged mode", v.Name)))
			}
		}
	}

	return ret, nil
}
