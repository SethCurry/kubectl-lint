package metalinters

import (
	"github.com/SethCurry/kubectl-lint/internal/lint"
	"github.com/SethCurry/kubectl-lint/internal/scopes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/resource"
)

func HasLabel() *HasLabelLinter {
	return &HasLabelLinter{}
}

type HasLabelLinter struct {
}

func (h *HasLabelLinter) ErrorCode() lint.ErrorCode {
	return lint.ErrorCode("has-label")
}

func (n *HasLabelLinter) Severity() lint.Severity {
	return lint.SeverityInfo
}

func (n *HasLabelLinter) ValidScopes() []scopes.Scope {
	return []scopes.Scope{
		scopes.All(),
	}
}

func (n *HasLabelLinter) Lint(info *resource.Info) ([]lint.ErrorMessage, error) {
	var ret []lint.ErrorMessage

	unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(info.Object)
	if err != nil {
		return ret, err
	}

	unstr := &unstructured.Unstructured{
		Object: unstructuredObject,
	}

	labels := unstr.GetLabels()

	if len(labels) == 0 {
		ret = append(ret, lint.ErrorMessage("has no labels, consider adding some for organization or tagging"))
	}

	return ret, nil
}
