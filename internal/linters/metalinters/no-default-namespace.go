package metalinters

import (
	"github.com/SethCurry/kubectl-lint/internal/lint"
	"github.com/SethCurry/kubectl-lint/internal/scopes"
	"k8s.io/cli-runtime/pkg/resource"
)

func NoDefaultNamespace() *NoDefaultNamespaceLinter {
	return &NoDefaultNamespaceLinter{}
}

type NoDefaultNamespaceLinter struct {
}

func (n *NoDefaultNamespaceLinter) ErrorCode() lint.ErrorCode {
	return lint.ErrorCode("no-default-namespace")
}

func (n *NoDefaultNamespaceLinter) Severity() lint.Severity {
	return lint.SeverityInfo
}

func (n *NoDefaultNamespaceLinter) ValidScopes() []scopes.Scope {
	return []scopes.Scope{
		scopes.All(),
	}
}

func (n *NoDefaultNamespaceLinter) Lint(info *resource.Info) ([]lint.ErrorMessage, error) {
	var ret []lint.ErrorMessage

	if info.Namespace == "" {
		ret = append(ret, lint.ErrorMessage("in namespace \"\""))
	}

	if info.Namespace == "default" {
		ret = append(ret, lint.ErrorMessage("in namespace \"default\""))
	}

	return ret, nil
}
