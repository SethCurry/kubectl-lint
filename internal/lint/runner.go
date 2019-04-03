package lint

import (
	"go.uber.org/zap"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
)

type RunnerOption func(*Runner)

func WithLogger(logger *zap.Logger) RunnerOption {
	return func(runner *Runner) {
		runner.logger = logger
	}
}

func WithLinters(linters []Linter) RunnerOption {
	return func(runner *Runner) {
		runner.linters = linters
	}
}

func NewRunner(opts ...RunnerOption) *Runner {
	runner := &Runner{
		logger: zap.NewNop(),
	}

	for _, opt := range opts {
		opt(runner)
	}

	return runner
}

type Runner struct {
	linters []Linter
	logger  *zap.Logger
}

func (r *Runner) AddLinter(linter Linter) {
	r.linters = append(r.linters, linter)
}

func (r *Runner) Lint(info *resource.Info) ([]Error, error) {
	gvk := info.Object.GetObjectKind().GroupVersionKind()

	r.logger.Debug("linting object",
		zap.String("kind", gvk.Kind),
		zap.String("name", info.Name),
		zap.String("namespace", info.Namespace),
	)

	var ret []Error

	for _, linter := range r.linters {
		errorCode := linter.ErrorCode()
		var msgs []ErrorMessage
		var err error

		linted := false

		switch gvk {
		case schema.GroupVersionKind{
			Version: "v1",
			Kind:    "Pod",
		}:
			if asPodLinter, ok := linter.(PodLinter); ok {
				pod := &core_v1.Pod{}

				unstructured := make(map[string]interface{})

				unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(info.Object)
				if err != nil {
					r.logger.Fatal("failed to convert object to unstructured",
						zap.Error(err),
					)
				}

				err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, pod)
				if err != nil {
					r.logger.Fatal("failed to convert unstructured object to Pod",
						zap.Error(err),
					)
				}

				linted = true
				msgs, err = asPodLinter.Lintv1Pod(pod)
			}
		default:
			r.logger.Debug("no more specific handler for linter, using generic",
				zap.String("group", gvk.Group),
				zap.String("kind", gvk.Kind),
				zap.String("version", gvk.Version),
			)
		}

		if !linted {
			msgs, err = linter.Lint(info)
		}

		if err != nil {
			r.logger.Error("linter returned error",
				zap.String("linter", string(errorCode)),
				zap.Error(err),
			)
		}

		for _, msg := range msgs {
			lintErr := Error{
				ErrorCode:       errorCode,
				ErrorMessage:    msg,
				Severity:        linter.Severity(),
				SourceKind:      gvk.Kind,
				SourceNamespace: info.Namespace,
				SourceName:      info.Name,
			}

			ret = append(ret, lintErr)
		}
	}

	return ret, nil
}
