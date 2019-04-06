# kubectl-lint

`kubectl-lint` is a plugin for `kubectl` to allow you to lint your Kubernetes infrastructure.

## Installation

Releases are not currently provided.  You can install `kubectl-lint` by running:

```bash
go get -u github.com/SethCurry/kubectl-lint/cmd/kubectl-lint
```

## Usage

You can run `kubectl-lint` as a `kubectl` plugin with `kubectl lint`.
It takes most of the same arguments as `kubectl get`, i.e.

```bash
kubectl lint pod my-pod-name
```

