
export GO111MODULE=on

TAG := $(shell git describe --tags --always | cut -d '-' -f 1)
VERSION := $(TAG)
LDFLAGS := "-X github.com/cturiel/kubectl-slowdrain/pkg/version.Version=$(VERSION)"

.PHONY: test
test:
	go test ./pkg/... ./cmd/... -coverprofile cover.out

.PHONY: bin
bin: fmt vet
	go build -ldflags $(LDFLAGS) \
					 -o bin/kubectl-slowdrain \
					 github.com/cturiel/kubectl-slowdrain/cmd/plugin

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./cmd/...

.PHONY: kubernetes-deps
kubernetes-deps:
	go mod tidy

.PHONY: setup
setup:
	make -C setup
