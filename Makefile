GOCMD = go
GOPATH = $(shell $(GOCMD) env GOPATH)
GOCOV = $(GOPATH)/bin/gocov
GOLINT = $(GOPATH)/bin/golangci-lint
DOCKER = docker
GIT_VERSION = $(shell git rev-list -1 HEAD)
CURRENT_DIR = $(shell pwd)

prereq:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
	@GO111MODULE=on $(GOCMD) install github.com/golang/mock
	@GO111MODULE=on $(GOCMD) install github.com/golang/mock
	@GO111MODULE=on $(GOCMD) install github.com/golang/mock/gomock
	@GO111MODULE=on $(GOCMD) install github.com/golang/mock/mockgen


lint:
	@GO111MODULE=on $(GOLINT) run -v --skip-dirs vendor --timeout=5m

run:
	@GO111MODULE=on $(GOCMD) run main.go http --config=./config/app.conf.yaml

mock:
	$(GOCMD) generate ./...

models:
	@docker run --rm -it -v $(CURRENT_DIR):/opt/mego -w /opt/mego --entrypoint /bin/sh quay.io/goswagger/swagger:latest models.sh