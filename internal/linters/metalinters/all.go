package metalinters

import "github.com/SethCurry/kubectl-lint/internal/lint"

func All() []lint.Linter {
	return []lint.Linter{
		NoDefaultNamespace(),
	}
}
