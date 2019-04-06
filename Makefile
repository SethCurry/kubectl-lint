install:
	go install ./cmd/kubectl-lint
lint:
	golangci-lint run
