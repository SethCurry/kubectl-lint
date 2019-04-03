package podlinters

import (
	"errors"
	"fmt"

	"github.com/SethCurry/kubectl-lint/internal/lint"
	"github.com/SethCurry/kubectl-lint/internal/scopes"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/resource"
)

func HasCPURequest() *HasCPURequestLinter {
	return &HasCPURequestLinter{}
}

type HasCPURequestLinter struct {
}

func (h *HasCPURequestLinter) ErrorCode() lint.ErrorCode {
	return lint.ErrorCode("has-cpu-request")
}

func (h *HasCPURequestLinter) Severity() lint.Severity {
	return lint.SeverityWarn
}

func (h *HasCPURequestLinter) ValidScopes() []scopes.Scope {
	return []scopes.Scope{
		scopes.And(scopes.Kind("Pod"), scopes.Version("v1")),
	}
}

func (h *HasCPURequestLinter) Lint(info *resource.Info) ([]lint.ErrorMessage, error) {
	return []lint.ErrorMessage{}, errors.New("this should be handled by the pod linter")
}

func (h *HasCPURequestLinter) Lintv1Pod(pod *core_v1.Pod) ([]lint.ErrorMessage, error) {
	var ret []lint.ErrorMessage

	containers := append(pod.Spec.Containers, pod.Spec.InitContainers...)

	for _, v := range containers {
		if _, ok := v.Resources.Requests["cpu"]; !ok {
			ret = append(ret, lint.ErrorMessage(fmt.Sprintf("container \"%s\" has no CPU request configured", v.Name)))
		}
	}

	return ret, nil
}
