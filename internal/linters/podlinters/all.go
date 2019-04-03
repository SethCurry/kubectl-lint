package podlinters

import "github.com/SethCurry/kubectl-lint/internal/lint"

func All() []lint.Linter {
	return []lint.Linter{
		HasMemoryRequest(),
		HasMemoryLimit(),
		HasCPURequest(),
		HasCPULimit(),
	}
}
