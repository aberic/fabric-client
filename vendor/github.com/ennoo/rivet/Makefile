PKGS_WITH_OUT_EXAMPLES := $(shell go list ./... | grep -v 'examples/')
PKGS_WITH_OUT_EXAMPLES_AND_UTILS := $(shell go list ./... | grep -v 'examples/\|utils/')
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" -print0 | xargs -0)

export SQL_Name=circle_test
export DB_USER=root
export DB_NAME=mysql
checkTravis:export DB_PORT=33061
checkTravis:export DB_URL=127.0.0.1:${DB_PORT}
checkTravis:export DB_PASS=secret
checkCircle:export DB_PORT=3306
checkCircle:export DB_URL=127.0.0.1:${DB_PORT}
checkCircle:export DB_PASS=
checkLocal:export DB_PORT=3306
checkLocal:export DB_URL=127.0.0.1:${DB_PORT}
checkLocal:export DB_PASS=secret

checkTravis: start overalls vet lint misspell staticcheck cyclo const veralls test endCommon endConsul endDocker

checkCircle: wright consul overalls vet lint misspell staticcheck cyclo const test endCommon endConsul

checkLocal: wright consul overalls vet lint misspell staticcheck cyclo const endCommon endConsul endDocker

start: wright mysql consul

endCommon:
	@echo "end common"
	rm -rf a.txt

endConsul:
	@echo "end consul"
	consul leave

endDocker:
	@echo "end docker"
	docker rm -f ${SQL_Name}

mysql:
	@echo "docker run mysql"
	docker run --name ${SQL_Name} -e MYSQL_ROOT_PASSWORD=${DB_PASS} -d -p ${DB_PORT}:3306 mysql:5.7.26

consul:
	@echo "consul"
	nohup consul agent -dev &

wright:
	@echo "wright"
	echo "this is my test\n" > a.txt

overalls:
	@echo "overalls"
	overalls -project=github.com/ennoo/rivet -covermode=count -ignore='.git,_vendor'

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
	gocyclo -over 10 $(GO_FILES)
#	gocyclo -top 10 $(GO_FILES)

const:
	@echo "goconst"
	goconst $(PKGS_WITH_OUT_EXAMPLES)

veralls:
	@echo "goveralls"
	goveralls -coverprofile=overalls.coverprofile -service=travis-ci -repotoken $(COVERALLS_TOKEN)

test:
	@echo "test"
	go test -v -cover $(PKGS_WITH_OUT_EXAMPLES)

build_all: build_bow_all build_shunt_all

build_bow_all: build_bow_amd64_all build_bow_386_all

build_bow_amd64_all: build_bow_darwin_amd64 build_bow_linux_amd64 build_bow_windows_amd64 build_bow_freebsd_amd64

build_bow_386_all: build_bow_darwin_386 build_bow_linux_386 build_bow_windows_386 build_bow_freebsd_386

build_shunt_all: build_shunt_amd64_all build_shunt_386_all

build_shunt_amd64_all: build_shunt_darwin_amd64 build_shunt_linux_amd64 build_shunt_windows_amd64 build_shunt_freebsd_amd64

build_shunt_386_all: build_shunt_darwin_386 build_shunt_linux_386 build_shunt_windows_386 build_shunt_freebsd_386


build_bow_darwin_amd64_docker:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./rivet/bow/runner/bow_darwin_amd64 ./rivet/bow/runner

build_bow_darwin_amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(GOPATH)/src/github.com/ennoo/rivet/bow/runner/bow_darwin_amd64 $(GOPATH)/src/github.com/ennoo/rivet/bow/runner

build_bow_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GOPATH)/src/github.com/ennoo/rivet/bow/runner/bow_linux_amd64 $(GOPATH)/src/github.com/ennoo/rivet/bow/runner

build_bow_windows_amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(GOPATH)/src/github.com/ennoo/rivet/bow/runner/bow_windows_amd64 $(GOPATH)/src/github.com/ennoo/rivet/bow/runner

build_bow_freebsd_amd64:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o $(GOPATH)/src/github.com/ennoo/rivet/bow/runner/bow_freebsd_amd64 $(GOPATH)/src/github.com/ennoo/rivet/bow/runner


build_bow_darwin_386:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o $(GOPATH)/src/github.com/ennoo/rivet/bow/runner/bow_darwin_386 $(GOPATH)/src/github.com/ennoo/rivet/bow/runner

build_bow_linux_386:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o $(GOPATH)/src/github.com/ennoo/rivet/bow/runner/bow_linux_386 $(GOPATH)/src/github.com/ennoo/rivet/bow/runner

build_bow_windows_386:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o $(GOPATH)/src/github.com/ennoo/rivet/bow/runner/bow_windows_386 $(GOPATH)/src/github.com/ennoo/rivet/bow/runner

build_bow_freebsd_386:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -o $(GOPATH)/src/github.com/ennoo/rivet/bow/runner/bow_freebsd_386 $(GOPATH)/src/github.com/ennoo/rivet/bow/runner



build_shunt_darwin_amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner/shunt_darwin_amd64 $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner

build_shunt_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner/shunt_linux_amd64 $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner

build_shunt_windows_amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner/shunt_windows_amd64 $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner

build_shunt_freebsd_amd64:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner/shunt_freebsd_amd64 $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner


build_shunt_darwin_386:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner/shunt_darwin_386 $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner

build_shunt_linux_386:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner/shunt_linux_386 $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner

build_shunt_windows_386:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner/shunt_windows_386 $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner

build_shunt_freebsd_386:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -o $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner/shunt_freebsd_386 $(GOPATH)/src/github.com/ennoo/rivet/shunt/runner
