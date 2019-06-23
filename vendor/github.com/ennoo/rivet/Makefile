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
	gocyclo -over 18 $(GO_FILES)
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