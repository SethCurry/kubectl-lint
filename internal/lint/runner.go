package lint

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"go.uber.org/zap"
	apps_v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

func (r *Runner) DisableLinter(ec ErrorCode) error {
	for k, v := range r.linters {
		if v.ErrorCode() == ec {
			r.linters = append(r.linters[:k], r.linters[k+1:]...)
			return nil
		}
	}

	return errors.New("unrecognized error code: " + string(ec))
}

func (r *Runner) AddLinter(linter Linter) {
	r.linters = append(r.linters, linter)
}

func (r *Runner) lintV1Pod(linter Linter, info *resource.Info) ([]ErrorMessage, error) {
	pod := &core_v1.Pod{}

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

	if asPodLinter, ok := linter.(PodLinter); ok {
		return asPodLinter.Lintv1Pod(pod)
	} else if asPodSpecLinter, ok := linter.(PodSpecLinter); ok {
		return asPodSpecLinter.Lintv1PodSpec(pod.Spec)
	}
	return linter.Lint(info)
}

func (r *Runner) lintV1DaemonSet(linter Linter, info *resource.Info) ([]ErrorMessage, error) {
	daemonSet := &apps_v1.DaemonSet{}

	unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(info.Object)
	if err != nil {
		r.logger.Fatal("failed to convert object to unstructured",
			zap.Error(err),
		)
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, daemonSet)
	if err != nil {
		r.logger.Fatal("failed to convert unstructured object to DaemonSet",
			zap.Error(err),
		)
	}

	if asDaemonSetLinter, ok := linter.(DaemonSetLinter); ok {
		return asDaemonSetLinter.Lintv1DaemonSet(daemonSet)
	} else if asPodSpecLinter, ok := linter.(PodSpecLinter); ok {
		return asPodSpecLinter.Lintv1PodSpec(daemonSet.Spec.Template.Spec)
	}
	return linter.Lint(info)
}

// TODO there has to be a way to condense this boilerplate +code-smell
func (r *Runner) lintV1Deployment(linter Linter, info *resource.Info) ([]ErrorMessage, error) {
	deployment := &apps_v1.Deployment{}

	unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(info.Object)
	if err != nil {
		r.logger.Fatal("failed to convert object to unstructured",
			zap.Error(err),
		)
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, deployment)
	if err != nil {
		r.logger.Fatal("failed to convert unstructured object to Deployment",
			zap.Error(err),
		)
	}

	var msgs []ErrorMessage

	if asDeploymentLinter, ok := linter.(DeploymentLinter); ok {
		msgs, err = asDeploymentLinter.Lintv1Deployment(deployment)
	} else if asPodSpecLinter, ok := linter.(PodSpecLinter); ok {
		msgs, err = asPodSpecLinter.Lintv1PodSpec(deployment.Spec.Template.Spec)
	} else {
		msgs, err = linter.Lint(info)
	}

	if err != nil {
		return msgs, nil
	}

	return msgs, nil
}

func (r *Runner) Lint(info *resource.Info) ([]Error, error) {
	var ret []Error

	gvk := info.Object.GetObjectKind().GroupVersionKind()

	r.logger.Debug("linting object",
		zap.String("kind", gvk.Kind),
		zap.String("group", gvk.Group),
		zap.String("version", gvk.Version),
		zap.String("name", info.Name),
		zap.String("namespace", info.Namespace),
	)

	unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(info.Object)
	if err != nil {
		return ret, err
	}

	unstructured := &meta_v1.Unstructured{
		Object: unstructuredObject,
	}

	annotations := unstructured.GetAnnotations()

	// TODO document this
	disableLinters := annotations["lint/disable-linters"]
	bannedLinters := strings.Split(disableLinters, ",")

	for _, linter := range r.linters {
		isLinterBanned := false
		for _, v := range bannedLinters {
			if string(linter.ErrorCode()) == v {
				r.logger.Info("skipping linter because of annotations",
					zap.String("error_code", string(linter.ErrorCode())),
				)
				isLinterBanned = true
				break
			}
		}

		if isLinterBanned {
			continue
		}

		inScope := false
		for _, scope := range linter.ValidScopes() {
			if scope(info) {
				inScope = true
				break
			}
		}

		if !inScope {
			continue
		}

		errorCode := linter.ErrorCode()
		var err error
		var msgs []ErrorMessage

		switch gvk {
		case schema.GroupVersionKind{
			Version: "v1",
			Kind:    "Pod",
		}:
			msgs, err = r.lintV1Pod(linter, info)
		case schema.GroupVersionKind{
			Version: "v1",
			Group:   "apps",
			Kind:    "Deployment",
		}:
			msgs, err = r.lintV1Deployment(linter, info)
		case schema.GroupVersionKind{
			Version: "v1",
			Group:   "apps",
			Kind:    "DaemonSet",
		}:
			msgs, err = r.lintV1DaemonSet(linter, info)
		default:
			return ret, errors.New("no linters for kind")
		}

		if err != nil {
			return ret, errors.Wrap(err, fmt.Sprintf("failed to run linter \"%s\"", string(linter.ErrorCode())))
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
