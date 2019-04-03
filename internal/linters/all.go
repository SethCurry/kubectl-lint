package linters

import (
	"github.com/SethCurry/kubectl-lint/internal/lint"
	"github.com/SethCurry/kubectl-lint/internal/linters/metalinters"
	"github.com/SethCurry/kubectl-lint/internal/linters/podlinters"
)

func All() []lint.Linter {
	return append(metalinters.All(), podlinters.All()...)
}
