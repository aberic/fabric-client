PKGS_WITH_OUT_EXAMPLES := $(shell go list ./... | grep -v 'examples/')
PKGS_WITH_OUT_EXAMPLES_AND_UTILS := $(shell go list ./... | grep -v 'examples/\|utils/')
GO_FILES := $(shell find . -name "*.go" -not -name "*_test.go" -not -path "./vendor/*" -not -path ".git/*" -print0 | xargs -0)

export GOPROXY=https://goproxy.io
export GO111MODULE=on

checkTravis: overalls vet lint misspell staticcheck cyclo const veralls test

checkLocal: overalls vet lint misspell staticcheck cyclo const test

overalls:
	@echo "overalls"
	overalls -project=github.com/aberic/gnomon -covermode=count -ignore='.git,_vendor'

vet:
	@echo "vet"
	go vet $(PKGS_WITH_OUT_EXAMPLES)

lint:
	@echo "golint"
	golint -set_exit_status $(PKGS_WITH_OUT_EXAMPLES_AND_UTILS)

misspell:
	@echo "misspell"
	misspell -source=text -error $(GO_FILES)

staticcheck:
	@echo "staticcheck"
	staticcheck $(PKGS_WITH_OUT_EXAMPLES)

cyclo:
	@echo "gocyclo"
	gocyclo -over 15 $(GO_FILES)
	gocyclo -top 10 $(GO_FILES)

const:
	@echo "goconst"
	goconst $(PKGS_WITH_OUT_EXAMPLES)

veralls:
	@echo "goveralls"
	goveralls -coverprofile=overalls.coverprofile -service=travis-ci -repotoken $(COVERALLS_TOKEN)

traviscodecovtest:
	@echo "travistest"
	go test -race -coverprofile=coverage.txt -covermode=atomic

test:
	@echo "test"
	go test -v -cover $(PKGS_WITH_OUT_EXAMPLES)