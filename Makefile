PKGS_WITH_OUT_EXAMPLES := $(shell go list ./... | grep -v 'examples/')
PKGS_WITH_OUT_EXAMPLES_AND_UTILS := $(shell go list ./... | grep -v 'examples/\|utils/')
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" -print0 | xargs -0)
BASE_VERSION1 = 1.4
BASEIMAGE_RELEASE=0.4.14

checkLocal: vet misspell cyclo const

checkTest: images

images:
	@echo "docker pull images"
	docker pull hyperledger/fabric-zookeeper:$(BASEIMAGE_RELEASE)
	docker pull hyperledger/fabric-kafka:$(BASEIMAGE_RELEASE)
	docker pull hyperledger/fabric-baseos:$(BASEIMAGE_RELEASE)
    docker pull hyperledger/fabric-orderer:$(BASE_VERSION)
    docker pull hyperledger/fabric-peer:$(BASE_VERSION)
    docker pull hyperledger/fabric-ccenv:$(BASE_VERSION)

overalls:
	@echo "overalls"
	overalls -project=github.com/aberic/fabric-client -covermode=count -ignore='.git,_vendor'
	go tool cover -func=overalls.coverprofile

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

const:
	@echo "goconst"
	goconst $(PKGS_WITH_OUT_EXAMPLES)

veralls:
	@echo "goveralls"
	goveralls -coverprofile=overalls.coverprofile -service=travis-ci -repotoken $(COVERALLS_TOKEN)

test:
	@echo "test"
	go test -v -cover $(PKGS_WITH_OUT_EXAMPLES)